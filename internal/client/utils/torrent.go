package utils

import (
	"encoding/json"
	"errors"
	"fmt"
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
	WRITE           = "write"
)

func OpenFile(mode mode, path string) (*os.File, error) {
	var flags int
	var permissions fs.FileMode

	if mode == READ {
		flags = os.O_RDONLY | os.O_CREATE
		permissions = 0400

	} else if mode == READ_WRITE {
		flags = os.O_RDWR | os.O_CREATE
		permissions = 0600

	} else if mode == WRITE {
		flags = os.O_WRONLY | os.O_CREATE
		permissions = 0600

	} else {
		slog.Error("Invalid file mode provided")
		return nil, errors.New("Invalid file mode provided")
	}

	f, err := os.OpenFile(path, flags, permissions)
	if err != nil {
		slog.Error("Error opening torrents file", "error", err.Error())
	}
	return f, err
}

func CalculateChunksNumber(length uint64, chunkLength uint64) uint8 {
	chunks := length / chunkLength
	if length%chunkLength != 0 {
		chunks += 1
	}

	return uint8(chunks)
}

func CalculateTorrentProgress(chunksDownloaded []uint8, chunkLength uint64) uint8 {
	return uint8(len(chunksDownloaded) * int(chunkLength))
}

func CalculateTorrentStatus(length uint64, chunkLength uint64, chunksDownloaded []uint8, defaultMode models.Status) models.Status {
	chunks := CalculateChunksNumber(length, chunkLength)
	if len(chunksDownloaded) == int(chunks) {
		return models.DOWNLOADED
	}

	return defaultMode
}

func GetTorrentsData() []models.Torrent {
	var torrents []models.Torrent = []models.Torrent{}
	torrentsData := map[string]any{}

	f, _ := OpenFile(READ, "torrents.json")
	defer f.Close()
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

		if homeDir, err := os.UserHomeDir(); err == nil {
			downloadDir = fmt.Sprintf("%s%s", homeDir, "Downloads/")
		}

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

		status = CalculateTorrentStatus(length, chunkLength, chunksDownloaded, models.STOPPED)

		progress = CalculateTorrentProgress(chunksDownloaded, chunkLength)

		torrents = append(torrents, models.Torrent{
			File:             file,
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
