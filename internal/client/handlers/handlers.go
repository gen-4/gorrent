package handlers

import (
	"net/http"

	"github.com/gen-4/gorrent/config/client"
)

func HasTorrent(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")

	for _, torrent := range config.Configuration.Torrents {
		if torrent.Name == file {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
