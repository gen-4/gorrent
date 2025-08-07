package utils

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"reflect"

	"github.com/gen-4/gorrent/internal/client/models"
)

type mode string

const (
	READ       mode = "read"
	READ_WRITE      = "read_write"
)

func OpenTorrentsFile(mode mode) (*os.File, error) {
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

func GetTorrentsData() []models.Torrent {
	var torrents []models.Torrent = []models.Torrent{}
	torrentsData := map[string]any{}

	f, _ := OpenTorrentsFile(READ)
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
