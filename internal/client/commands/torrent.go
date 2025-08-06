package commands

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"reflect"

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
		flags = os.O_RDWR | os.O_CREATE
		permissions = 0600

	} else {
		slog.Error("Invalid file mode provided")
		return nil, errors.New("Invalid file mode provided")
	}

	f, err := os.OpenFile("torrents.json", flags, permissions)
	if err != nil {
		slog.Error("Error opening torrents file", "error", err.Error())
	}
	return f, err
}

func readTorrentFile(path string) (string, uint64, []string, error) {
	content, err := parser.ParseFromFile(path)
	if err != nil {
		slog.Error("Unable to read .torrent file", "error", err.Error())
		return "", 0, []string{}, err
	}

	return content.Files[0].Path[1], uint64(content.Files[0].Length), content.Announce, nil
}

func calculateChunkLength(length uint64) uint64 {
	lengths := []uint16{512, 1024, 2048, 4096}
	for l := range lengths {
		if int(length)/l > 500 && int(length)/l < 2000 {
			return uint64(l)
		}
	}

	return uint64(2048)
}

func CreateTorrent(path string) tea.Cmd {
	return func() tea.Msg {
		f, err := openTorrentsFile(READ_WRITE)
		if err != nil {
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
		isTorrentsFileEmpty := stat.Size() == 0

		torrentsData := map[string]any{}
		if !isTorrentsFileEmpty {
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

			if _, err := f.Seek(-2, 2); err != nil {
				slog.Error("Error setting file offset", "error", err.Error())
				return err
			}
			if _, err := f.WriteString(",\n"); err != nil {
				slog.Error("Unable to write torrents file", "error", err.Error())
				return err
			}
		}

		_, found := torrentsData[file]
		if found {
			slog.Warn("Torrent is already present")
			return nil
		}

		var chunkLength uint64 = calculateChunkLength(length)

		tData := map[string]any{
			file: map[string]any{
				"download_directory": "~/Downloads/",
				"length":             length,
				"chunk_length":       chunkLength,
				"chunks_downloaded":  []int{},
			},
		}

		byteTData, err := json.Marshal(tData)
		if err != nil {
			slog.Error("Unable to marshal torrent data", "error", err.Error())
			return err
		}

		if !isTorrentsFileEmpty {
			byteTData = byteTData[1:]
		}

		if _, err = f.Write(byteTData); err != nil {
			slog.Error("Unable to write torrents file", "error", err.Error())
			return err
		}

		return models.NewTorrentRequest{
			Name:             file,
			Peers:            uint8(0),
			Progress:         uint8(0),
			Status:           models.STOPPED,
			Superservers:     superservers,
			ChunkLength:      chunkLength,
			ChunksDownloaded: []uint8{},
			Length:           length,
			DownloadDir:      "~/Downloads/",
		}
	}
}

func GetTorrentsData() []models.Torrent {
	var torrents []models.Torrent = []models.Torrent{}
	torrentsData := map[string]any{}

	f, _ := openTorrentsFile(READ)
	stat, err := f.Stat()
	if err != nil {
		slog.Error("Error reading torrents file stats", "error", err.Error())
		return torrents
	}

	if stat.Size() != 0 {
		jsonData := make([]byte, stat.Size())
		if _, err := f.Read(jsonData); err != nil {
			slog.Error("Unable to read torrents file", "erro", err.Error())
		}

		json.Unmarshal(jsonData, &torrentsData)
	}

	for file, torrentData := range torrentsData {
		var peers uint8 = 0
		var progress uint8 = 0
		var status models.Status = models.STOPPED
		var superservers []string = []string{}
		var chunkLength uint64 = 0
		var length uint64 = 0
		var chunksDownloaded []uint8 = []uint8{}
		var downloadDir string = "~/Downloads/"

		tData := torrentData.(map[string]any)
		if v, found := tData["download_directory"]; found {
			downloadDir = v.(string)
		}
		if v, found := tData["chunk_length"]; found {
			chunkLength = uint64(v.(float64))
		}
		if v, found := tData["chunks_downloaded"]; found {
			cDownloadedAny, ok := v.([]any)
			if !ok {
				slog.Error("Wrong type assertion, expected []any", "type", reflect.TypeOf(cDownloadedAny))
			}
			for _, ch := range cDownloadedAny {
				chunksDownloaded = append(chunksDownloaded, uint8(ch.(float64)))
			}
		}
		if v, found := tData["length"]; found {
			length = uint64(v.(float64))
		}

		chunks := length / chunkLength
		if length%chunkLength != 0 {
			chunks += 1
		}
		if len(chunksDownloaded) == int(chunks) {
			status = models.DOWNLOADED
		}

		torrents = append(torrents, models.Torrent{
			Name:             file,
			Peers:            peers,
			Progress:         progress,
			Status:           status,
			Superservers:     superservers,
			ChunkLength:      chunkLength,
			Length:           length,
			DownloadDir:      downloadDir,
			ChunksDownloaded: chunksDownloaded,
		})
	}

	return torrents
}
