package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	parser "github.com/j-muller/go-torrent-parser"

	config "github.com/gen-4/gorrent/config/superserver"
)

func readTorrentFile(path string) (string, []string) {
	content, _ := parser.ParseFromFile(path)
	fileName := content.Files[0].Path[0]
	return fileName, content.Announce
}

func GetStoredTorrents(w http.ResponseWriter, req *http.Request) {
	torrents := []string{}
	if config.Configuration.TorrentsFolder == "" {
		slog.Info("Torrents folder is not configured")
		return
	}

	criteriaName := req.URL.Query().Get("name")

	err := filepath.Walk(config.Configuration.TorrentsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			name, _ := readTorrentFile(path)
			if strings.Contains(strings.ToLower(name), strings.ToLower(criteriaName)) {
				torrents = append(torrents, name)
			}
		}

		return nil
	})
	if err != nil {
		slog.Error("Error walking torrents dir", "error", err.Error())
	}

	w.Header().Add("Content-Type", "application/json")

	data := map[string]any{
		"torrents": torrents,
	}
	json.NewEncoder(w).Encode(data)
}
