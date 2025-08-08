package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	config "github.com/gen-4/gorrent/config/superserver"
	gUtils "github.com/gen-4/gorrent/internal/utils"
)

func namesMatch(a string, b string) bool {
	matchedChars := 0
	for _, char := range a {
		if matchedChars >= len(b) {
			break
		}
		if char == rune(b[matchedChars]) {
			matchedChars += 1
		}
	}

	return matchedChars == len(b)
}

func GetStoredTorrents(w http.ResponseWriter, req *http.Request) {
	torrents := []string{}
	if config.Configuration.TorrentsFolder == "" {
		slog.Warn("Torrents folder is not configured")
		return
	}

	criteriaName := req.URL.Query().Get("name")

	err := filepath.Walk(config.Configuration.TorrentsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Error("Error reading file info", "error", err.Error())
			return err
		}

		if !info.IsDir() {
			name, _, _, err := gUtils.ReadTorrentFile(path)
			if err != nil {
				slog.Error("Unable to read .torrent file", "error", err.Error())
				return err
			}
			if namesMatch(strings.ToLower(name), strings.ToLower(criteriaName)) {
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

func SubscribePeer(w http.ResponseWriter, req *http.Request) {
	addr := strings.Split(req.RemoteAddr, ":")[0]
	for _, p := range config.Configuration.Peers {
		if p == addr {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	config.Configuration.Peers = append(config.Configuration.Peers, addr)
	slog.Info("New Peer subscribed", "peer", addr)
	w.WriteHeader(http.StatusOK)
}

func GetPeersWithFile(w http.ResponseWriter, req *http.Request) {
	var wg sync.WaitGroup
	peersWithFile := []string{}
	file := req.URL.Query().Get("file")

	wg.Add(len(config.Configuration.Peers))

	for _, peer := range config.Configuration.Peers {
		go func() {
			defer wg.Done()

			if config.Configuration.Environment == config.PRO && peer == strings.Split(req.RemoteAddr, ":")[0] {
				return
			}

			response, err := http.Get(fmt.Sprintf(config.Configuration.PeerUrlTemplate, peer, fmt.Sprintf("torrent/?file=%s", file)))
			if err != nil {
				return
			}

			if response.StatusCode == http.StatusOK {
				peersWithFile = append(peersWithFile, peer)
			}

		}()
	}

	wg.Wait()

	w.Header().Add("Content-Type", "application/json")
	data := map[string]any{
		"peers": peersWithFile,
	}
	json.NewEncoder(w).Encode(data)

}
