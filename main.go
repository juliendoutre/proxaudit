package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/elazarl/goproxy"
	"go.uber.org/zap"
)

var (
	version = "unknown"
	commit  = "unknown" //nolint:gochecknoglobals
	date    = "unknown" //nolint:gochecknoglobals
)

func main() {
	mkcertDir := defaultMkcertDir()

	port := flag.Uint64("port", 8000, "port to listen on")
	caCertPath := flag.String("ca-cert", path.Join(mkcertDir, "rootCA.pem"), "path to a CA certificate")
	caKeyPath := flag.String("ca-key", path.Join(mkcertDir, "rootCA-key.pem"), "path to a CA private key")
	outputPath := flag.String("output", "stderr", "Path to a file to write logs to")
	showVersion := flag.Bool("version", false, "Show this program's version and exit")
	logHeader := flag.Bool("log-header", true, "Log request's header")
	logBody := flag.Bool("log-body", false, "Log request's body")
	serverMode := flag.Bool("server", false, "Run proxaudit as a server")
	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout, "proxaudit v%s, commit %s, built at %s\n", version, commit, date)

		return
	}

	os.Exit(run(*port, *outputPath, *caCertPath, *caKeyPath, *serverMode, *logHeader, *logBody))
}

func defaultMkcertDir() string {
	switch runtime.GOOS {
	case "darwin":
		return path.Join(os.Getenv("HOME"), "Library", "Application Support", "mkcert")
	default:
		return path.Join(os.Getenv("HOME"), ".local", "share", "mkcert")
	}
}

func run(port uint64, outputPath, caCertPath, caKeyPath string, serverMode, logHeader, logBody bool) int {
	logger, err := getLogger(outputPath)
	if err != nil {
		log.Println("Failed creating logger")

		return 1
	}

	defer func() { _ = logger.Sync() }()

	goproxy.GoproxyCa, err = tls.LoadX509KeyPair(caCertPath, caKeyPath)
	if err != nil {
		logger.Error("Failed loading CA certificate", zap.Error(err))

		return 1
	}

	server := newProxyServer(port, logger, logHeader, logBody)
	defer func() {
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Failed shutting down HTTP proxy server", zap.Error(err))
		}
	}()

	go runServer(server, logger, port)

	if serverMode {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals

		return 0
	}

	command, err := getCommand(logger)
	if err != nil {
		logger.Error("Failed reading command", zap.Error(err))

		return 1
	}

	instrumentedCommand := newInstrumentedCommand(context.Background(), command, port, caCertPath)

	status, err := runCommand(instrumentedCommand)
	if err != nil {
		logger.Error("Failed running command", zap.Error(err), zap.String("command", instrumentedCommand.String()))

		return 1
	}

	return status
}

func getLogger(outputPath string) (*zap.Logger, error) {
	if outputPath == "" {
		outputPath = "stderr"
	}

	config := zap.NewProductionConfig()
	config.OutputPaths = []string{outputPath}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("building logger from config: %w", err)
	}

	return logger, nil
}

func getCommand(logger *zap.Logger) (string, error) {
	if flag.NArg() == 0 {
		logger.Info("No command passed as argument, reading from stdin...")

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			return scanner.Text(), nil
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("reading from stdin: %w", err)
		}

		return "", nil
	}

	return strings.Join(flag.Args(), " "), nil
}

func newProxyServer(port uint64, logger *zap.Logger, logHeader, logBody bool) *http.Server {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.ReqConditionFunc(
		func(req *http.Request, _ *goproxy.ProxyCtx) bool {
			return req.Method != http.MethodConnect
		},
	)).DoFunc(
		func(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("url", req.URL.String()),
			}

			if logHeader {
				fields = append(fields, zap.Any("header", req.Header))
			}

			if logBody {
				body, err := io.ReadAll(req.Body)
				if err == nil {
					req.Body.Close()
					req.Body = io.NopCloser(bytes.NewReader(body))
				}

				fields = append(fields, zap.String("body", string(body)))
			}

			logger.Info("received a request", fields...)

			return req, nil
		},
	)

	return &http.Server{
		Addr:              ":" + strconv.FormatUint(port, 10),
		Handler:           proxy,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

func runServer(server *http.Server, logger *zap.Logger, port uint64) {
	logger.Info("Starting HTTP proxy server...", zap.Uint64("port", port))

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Warn("HTTP proxy server was stopped", zap.Error(err))
		} else {
			logger.Error("Failed running HTTP proxy server", zap.Error(err))
		}
	}
}

func newInstrumentedCommand(ctx context.Context, command string, port uint64, caCertPath string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("HTTP_PROXY=http://localhost:%d", port),
		fmt.Sprintf("HTTPS_PROXY=http://localhost:%d", port),
		fmt.Sprintf("http_proxy=http://localhost:%d", port),
		fmt.Sprintf("https_proxy=http://localhost:%d", port),
		"NODE_EXTRA_CA_CERTS="+caCertPath,
	)

	return cmd
}

func runCommand(cmd *exec.Cmd) (int, error) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	if err := cmd.Start(); err != nil {
		return 0, fmt.Errorf("starting command: %w", err)
	}

	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case sig := <-signals:
				_ = cmd.Process.Signal(sig)
			case <-done:
				return
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return 0, fmt.Errorf("waiting for command: %w", err)
	}

	return cmd.ProcessState.ExitCode(), nil
}
