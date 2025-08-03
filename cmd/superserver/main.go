package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	config "github.com/gen-4/gorrent/config/superserver"
	"github.com/gen-4/gorrent/internal/superserver/handlers"
)

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(recorder, r)
		slog.Info(fmt.Sprintf("[%s] %s %s %d %s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL, recorder.statusCode, r.UserAgent()))
	})
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig

		fmt.Print("\nCtrl+C noticed. Exiting...\n")
		os.Exit(1)
	}()
	config.Config()
	defer config.CloseConfig()

	mux := http.NewServeMux()
	gorrentMux := http.NewServeMux()
	gorrentMux.HandleFunc("/{$}", func(w http.ResponseWriter, req *http.Request) { fmt.Print("root hehe\n") })
	gorrentMux.HandleFunc("/healthcheck", func(w http.ResponseWriter, req *http.Request) { fmt.Print("hehe\n") })
	gorrentMux.HandleFunc("GET /get-torrents/", handlers.GetStoredTorrents)
	gorrentMux.HandleFunc("GET /get-torrents", handlers.GetStoredTorrents)
	mux.Handle("/gorrent/", http.StripPrefix("/gorrent", gorrentMux))
	appliedMiddlewareRouter := loggingMiddleware(mux)

	err := http.ListenAndServe("localhost:8000", appliedMiddlewareRouter)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Print("Server closed\n")
	} else if err != nil {
		slog.Error("Server closed with an error", "error", err.Error())
	}
}
