package utils

import (
	"log/slog"

	parser "github.com/j-muller/go-torrent-parser"
)

func ReadTorrentFile(path string) (string, uint64, []string, error) {
	content, err := parser.ParseFromFile(path)
	if err != nil {
		slog.Error("Unable to read .torrent file", "error", err.Error())
		return "", 0, []string{}, err
	}

	return content.Files[0].Path[1], uint64(content.Files[0].Length), content.Announce, nil
}
