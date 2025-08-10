package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/gen-4/gorrent/internal/client/models"
	"github.com/gen-4/gorrent/internal/client/utils"
	gUtils "github.com/gen-4/gorrent/internal/utils"
)

func calculateChunkLength(length uint64) uint64 {
	lengths := []uint16{512, 1024, 2048, 4096}
	for _, l := range lengths {
		if l == 0 {
			slog.Error(fmt.Sprintf("Division by zero. %d in %s list", l, lengths))
			continue
		}
		nChunks := length / uint64(l)
		if nChunks > 500 && nChunks < 2000 {
			return uint64(l)
		}
	}

	return uint64(2048)
}

func CreateTorrent(path string) tea.Cmd {
	return func() tea.Msg {
		var downloadDir string = "~/Downloads/"
		if homeDir, err := os.UserHomeDir(); err == nil {
			downloadDir = fmt.Sprintf("%s%s", homeDir, "Downloads/")
		}

		f, err := utils.OpenFile(utils.READ_WRITE, "torrents.json")
		if err != nil {
			return err
		}
		defer f.Close()

		file, length, superservers, err := gUtils.ReadTorrentFile(path)
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
				"superservers":       superservers,
				"download_directory": downloadDir,
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
			File:             file,
			Peers:            uint8(0),
			Progress:         uint8(0),
			Status:           models.STOPPED,
			Superservers:     superservers,
			ChunkLength:      chunkLength,
			ChunksDownloaded: []uint8{},
			Length:           length,
			DownloadDir:      downloadDir,
		}
	}
}
