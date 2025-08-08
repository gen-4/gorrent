package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gen-4/gorrent/config/client"
)

func HasTorrent(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")

	for _, torrent := range config.Configuration.Torrents {
		if torrent.Name == file {
			if _, err := os.Stat(fmt.Sprintf("%s%s", torrent.DownloadDir, file)); err != nil {
				slog.Error("File not found", "error", err.Error())
				break
			}

			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
