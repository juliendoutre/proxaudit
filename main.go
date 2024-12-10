package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	port := flag.Uint64("port", 8000, "port to listen on")
	flag.Parse()

	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Panic(err)
	}

	server := &http.Server{
		Addr:              ":" + strconv.FormatUint(*port, 10),
		Handler:           &handler{logger},
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

type handler struct {
	logger *zap.Logger
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.logger.Info("received request", zap.String("method", req.Method), zap.String("url", req.URL.String()))

	if req.Method == http.MethodConnect {
		h.handleConnect(rw, req)
	} else {
		httputil.NewSingleHostReverseProxy(req.URL).ServeHTTP(rw, req)
	}
}

func (h *handler) handleConnect(rw http.ResponseWriter, req *http.Request) {
	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		http.Error(rw, "Hijacking not supported", http.StatusInternalServerError)

		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)

		return
	}
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", req.Host)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusServiceUnavailable)

		return
	}
	defer serverConn.Close()

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}
