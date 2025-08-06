package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/handlers"
	"github.com/gen-4/gorrent/internal/client/ui"
	"github.com/gen-4/gorrent/internal/middleware"
)

func main() {
	config.Config()
	defer config.CloseConfig()
	mux := http.NewServeMux()
	gorrentMux := http.NewServeMux()
	gorrentMux.HandleFunc("/{$}", func(w http.ResponseWriter, req *http.Request) { fmt.Print("root hehe\n") })
	gorrentMux.HandleFunc("/healthcheck", func(w http.ResponseWriter, req *http.Request) { fmt.Print("hehe\n") })
	gorrentMux.HandleFunc("GET /torrent/", handlers.HasTorrent)
	mux.Handle("/gorrent/", http.StripPrefix("/gorrent", gorrentMux))
	appliedMiddlewareRouter := middleware.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: appliedMiddlewareRouter,
	}

	go func() {
		slog.Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server stopped unexpectedly", "error", err.Error())
		}
	}()

	ui.Run()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", "error", err.Error())
	} else {
		slog.Info("Server shut down gracefully")
	}
}
