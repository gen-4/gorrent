package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/utils"
)

func HasTorrent(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")

	for _, torrent := range config.Configuration.Torrents {
		if torrent.File == file {
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

func DownloadChunk(w http.ResponseWriter, req *http.Request) {
	file := req.URL.Query().Get("file")
	chunk := req.URL.Query().Get("chunk")
	chunkLength := req.URL.Query().Get("chunk_length")

	path := ""
	found := false

	for _, torrent := range config.Configuration.Torrents {
		if torrent.File == file {
			path = fmt.Sprintf("%s%s", torrent.DownloadDir, file)
			if _, err := os.Stat(path); err != nil {
				slog.Error("File not found", "error", err.Error())
				break
			}

			found = true
		}
	}

	if !found {
		slog.Error("File not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	f, err := utils.OpenFile(utils.READ, path)
	if err != nil {
		slog.Error("Error opening file", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		data := map[string]string{
			"error": "Internal server error trying to download data chunk",
		}
		json.NewEncoder(w).Encode(data)
		return
	}
	defer f.Close()

	l, err := strconv.ParseInt(chunkLength, 10, 0)
	if err != nil {
		slog.Error("Error parsing chunk length", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		data := map[string]string{
			"error": "Internal server error trying to download data chunk",
		}
		json.NewEncoder(w).Encode(data)
		return

	}
	ch, err := strconv.ParseInt(chunk, 10, 0)
	if err != nil {
		slog.Error("Error parsing chunk", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		data := map[string]string{
			"error": "Internal server error trying to download data chunk",
		}
		json.NewEncoder(w).Encode(data)
		return

	}

	dataBytes := make([]byte, l)
	if _, err := f.ReadAt(dataBytes, (ch-1)*l); err != nil && err != io.EOF {
		slog.Error("Error reading data file")
		w.WriteHeader(http.StatusInternalServerError)
		data := map[string]string{
			"error": "Internal server error trying to download data chunk",
		}
		json.NewEncoder(w).Encode(data)
		return
	}

	w.Write(dataBytes)
}
