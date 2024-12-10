package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/elazarl/goproxy"
	"go.uber.org/zap"
)

func main() {
	mkcertDir := path.Join(os.Getenv("HOME"), "Library", "Application Support", "mkcert")

	port := flag.Uint64("port", 8000, "port to listen on")
	caCertPath := flag.String("ca-cert", path.Join(mkcertDir, "rootCA.pem"), "path to a CA certificate")
	caKeyPath := flag.String("ca-key", path.Join(mkcertDir, "rootCA-key.pem"), "path to a CA private key")
	flag.Parse()

	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	cert, err := tls.LoadX509KeyPair(*caCertPath, *caKeyPath)
	if err != nil {
		logger.Panic("Failed loading CA certificate", zap.Error(err))
	}

	goproxy.GoproxyCa = cert

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest(goproxy.ReqConditionFunc(
		func(req *http.Request, _ *goproxy.ProxyCtx) bool {
			return req.Method != http.MethodConnect
		},
	)).DoFunc(
		func(req *http.Request, _ *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			logger.Info("received a request", zap.String("method", req.Method), zap.String("url", req.URL.String()))

			return req, nil
		},
	)

	server := &http.Server{
		Addr:              ":" + strconv.FormatUint(*port, 10),
		Handler:           proxy,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go waitForSignal(server, logger)

	logger.Info("Starting HTTP proxy server...", zap.Uint64("port", *port))

	if err := server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			logger.Warn("HTTP proxy server was stopped", zap.Error(err))
		} else {
			logger.Panic("Failed running HTTP proxy server", zap.Error(err))
		}
	}
}

func waitForSignal(server *http.Server, logger *zap.Logger) {
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("Failed shutting down HTTPS proxy server", zap.Error(err))
	}
}
