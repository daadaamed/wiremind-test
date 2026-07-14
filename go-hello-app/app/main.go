package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	shutdownTimeout time.Duration = 10 * time.Second
	serverPort      string        = ":8080"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	firstname := r.URL.Query().Get("firstname")
	lastname := r.URL.Query().Get("lastname")
	if firstname == "" || lastname == "" {
		http.Error(w, "missing required query parameters: firstname, lastname", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s %s \n", firstname, lastname)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)

	mux.HandleFunc("GET /healthz", healthzHandler)

	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}
	// Listen for SIGTERM and SIGINT
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		slog.Info("Starting server on " + serverPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Block until signal received.
	<-ctx.Done()
	stop()
	slog.Info("shutdown signal received, draining connections")

	// Give app time to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown gracefully
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}
