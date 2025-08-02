package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "github.com/gen-4/gorrent/config/superserver"
	"github.com/gen-4/gorrent/internal/superserver/handlers"
)

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

	http.HandleFunc("/{$}", func(w http.ResponseWriter, req *http.Request) { fmt.Print("hehe") })
	http.HandleFunc("GET /get-torrents/", handlers.GetStoredTorrents)

	err := http.ListenAndServe("localhost:8000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Print("Server closed\n")
	} else if err != nil {
		slog.Error("Server closed with an error", "error", err.Error())
	}
}
