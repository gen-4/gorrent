package commands

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	parser "github.com/j-muller/go-torrent-parser"

	"github.com/gen-4/gorrent/internal/client/models"
)

type mode string

const (
	READ       mode = "read"
	READ_WRITE      = "read_write"
)

func openTorrentsFile(mode mode) (*os.File, error) {
	var flags int
	var permissions fs.FileMode

	if mode == READ {
		flags = os.O_RDONLY | os.O_CREATE
		permissions = 0400

	} else if mode == READ_WRITE {
		flags = os.O_APPEND | os.O_RDWR | os.O_CREATE
		permissions = 0600

	} else {
		return nil, errors.New("Invalid file mode provided")
	}

	return os.OpenFile("torrents.json", flags, permissions)
}

func readTorrentFile(path string) (string, uint64, []string, error) {
	content, err := parser.ParseFromFile(path)
	if err != nil {
		slog.Error("Unable to read .torrent file", "error", err.Error())
		return "", 0, []string{}, err
	}

	return content.Files[0].Path[1], uint64(content.Files[0].Length), content.Announce, nil
}

func CreateTorrent(path string) tea.Cmd {
	return func() tea.Msg {
		f, err := openTorrentsFile(READ_WRITE)
		if err != nil {
			slog.Error("Unable to open torrents file", "errors", err.Error())
			return err
		}
		defer f.Close()

		file, length, superservers, err := readTorrentFile(path)
		if err != nil {
			slog.Error("Unable to read .torrent file", "errors", err.Error())
			return err
		}

		stat, err := f.Stat()
		if err != nil {
			slog.Error("Error reading torrents file stats", "error", err.Error())
			return err
		}

		torrentsData := map[string]any{}
		if stat.Size() != 0 {
			jsonData := make([]byte, stat.Size())
			_, err := f.Read(jsonData)
			if err != nil {
				slog.Error("Error reading torrents file", "error", err.Error())
				return err
			}
			if err := json.Unmarshal(jsonData, &torrentsData); err != nil {
				slog.Error("Error unmarshaling torrents file", "error", err.Error())
				return err
			}
		}

		_, found := torrentsData[file]
		if found {
			slog.Warn("Torrent is already present")
			return nil
		}

		// TODO: Calculate chunks and chunk length, which is basically the same

		torrentsData[file] = map[string]any{
			"download_directory": "~/Downloads/",
			"length":             length,
			"chunks":             1,
			"chunks_downloaded":  0,
		}

		byteTData, err := json.Marshal(torrentsData)
		if err != nil {
			slog.Error("Unable to marshal torrent data", "error", err.Error())
			return err
		}

		if _, err = f.WriteString(string(byteTData)); err != nil {
			slog.Error("Unable to write torrents file", "errors", err.Error())
			return err
		}

		return models.NewTorrentRequest{
			Name:         file,
			Peers:        uint8(0),
			Progress:     uint8(0),
			Status:       models.STOPPED,
			Superservers: superservers,
			ChunkLength:  length,
			Length:       length,
		}
	}
}
