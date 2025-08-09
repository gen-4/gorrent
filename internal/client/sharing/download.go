package sharing

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/models"
	"github.com/gen-4/gorrent/internal/client/utils"
)

type order uint8

type response struct {
	chunk uint8
	id    int
	data  []byte
	error error
}

type peer struct {
	id     int
	orders chan order
	cancel context.CancelFunc
}

func getPeersWithFile(torrentFile string) []string {
	peers := []string{}
	var peersRes models.GetPeersDto
	var data []byte

	for _, ss := range config.Configuration.Superservers {
		response, err := http.Get(fmt.Sprintf(config.Configuration.SuperserverUrlTemplate, ss, fmt.Sprintf("torrent/?file=%s", torrentFile)))
		if err != nil {
			slog.Error("Error asking for torrent", "superserver", ss, "torrent", torrentFile)
		}

		if _, err := response.Body.Read(data); err != nil {
			slog.Error("Error reading response body", "error", err.Error())
		}
		json.Unmarshal(data, &peersRes)
		peers = append(peers, peersRes.Peers...)
	}

	return peers
}

func peerProcess(
	wg *sync.WaitGroup,
	ctx context.Context,
	orders <-chan order,
	responses chan<- response,
	peer string,
	file string,
	length uint64,
	id int,
) {
	defer wg.Done()
	byteReponse := make([]byte, length)

	for {
		select {
		case <-ctx.Done():
			return

		case order := <-orders:
			resp, err := http.Get(fmt.Sprintf(config.Configuration.PeerUrlTemplate, peer, fmt.Sprintf("chunk?file=%s&chunk=%d&chunk_length=%d", file, order, length)))
			if err != nil {
				slog.Error("Error requesting peer chunk", "error", err.Error())
				responses <- response{
					id:    id,
					chunk: uint8(order),
					error: err,
				}
			}
			if _, err := resp.Body.Read(byteReponse); err != nil {
				slog.Error("Error reading peer chunk", "error", err.Error())
				responses <- response{
					id:    id,
					chunk: uint8(order),
					error: err,
				}
			}

			responses <- response{
				id:    id,
				chunk: uint8(order),
				data:  byteReponse,
				error: nil,
			}
		}
	}
}

func getNextChunk(torrent models.Torrent, chunks uint8) uint8 {
	nChunk := uint8(1)
	for nChunk <= chunks {
		if !slices.Contains(torrent.ChunksDownloaded, nChunk) {
			return nChunk
		}

		nChunk++
	}

	return 0
}

func manageChunkResponse(torrent *models.Torrent, resp response) error {
	f, err := utils.OpenFile(utils.WRITE, fmt.Sprintf("%s%d_%s", torrent.DownloadDir, resp.chunk, torrent.File))
	if err != nil {
		slog.Error("Error opening chunk file", "error", err.Error())
		return err
	}

	if _, err := f.Write(resp.data); err != nil {
		slog.Error("Error writing chunk file", "error", err.Error())
		return err
	}

	torrent.ChunksDownloaded = append(torrent.ChunksDownloaded, resp.chunk)
	torrent.Progress = utils.CalculateTorrentProgress(torrent.ChunksDownloaded, torrent.ChunkLength)
	torrent.Status = utils.CalculateTorrentStatus(torrent.Length, torrent.ChunkLength, torrent.ChunksDownloaded, models.IN_PROGRESS)

	return nil
}

func aggregateChunks(torrent models.Torrent) error {
	data := make([]byte, torrent.ChunkLength)
	wf, err := utils.OpenFile(utils.WRITE, fmt.Sprintf("%s%s", torrent.DownloadDir, torrent.File))
	if err != nil {
		slog.Error("Error opening final file", "error", err.Error())
		return err
	}
	defer wf.Close()

	for chunk := range torrent.ChunksDownloaded {
		rf, err := utils.OpenFile(utils.READ, fmt.Sprintf("%s%d_%s", torrent.DownloadDir, chunk, torrent.File))
		if err != nil {
			slog.Error("Error opening chunk file", "chunk", chunk, "error", err.Error())
			return err
		}
		defer rf.Close()

		if _, err := wf.Seek(0, 2); err != nil {
			slog.Error("Error seeking last position in final file", "error", err.Error())
			return err
		}
		if _, err := rf.Read(data); err != nil {
			slog.Error("Error reading chunk file", "chunk", chunk, "error", err.Error())
			return err

		}
		if _, err := wf.Write(data); err != nil {
			slog.Error("Error writing final file", "error", err.Error())
			return err
		}
	}

	return nil
}

func DownloadTorrent(torrent *models.Torrent, ch chan<- models.Torrent) {
	// TODO: Add context to finalize process
	var wg *sync.WaitGroup
	var resp response
	peerAddrs := getPeersWithFile(torrent.File)
	peers := []peer{}
	responses := make(chan response, 5)
	done := false
	chunks := utils.CalculateChunksNumber(torrent.Length, torrent.ChunkLength)
	nextChunk := getNextChunk(*torrent, chunks)

	// TODO: Make this a goroutine and check periodcallly for new peers
	for i, p := range peerAddrs {
		if nextChunk == 0 {
			break
		}

		ctx, cancel := context.WithCancel(context.Background())
		peers = append(peers, peer{
			id:     i,
			orders: make(chan order),
			cancel: cancel,
		})

		wg.Add(1)
		go peerProcess(wg, ctx, peers[i].orders, responses, p, torrent.File, torrent.ChunkLength, i)
	}

	for !done {
		resp = <-responses
		if resp.error != nil {
			nextChunk = resp.chunk
			// TODO: Better retry strategy
		} else if err := manageChunkResponse(torrent, resp); err != nil {
			nextChunk = resp.chunk
			// TODO: Better retry strategy
		}

		targetPeer := peer{}
		index := -1
		for i, p := range peers {
			if p.id == resp.id {
				targetPeer = p
				index = i
				break
			}
		}
		if nextChunk == 0 {
			targetPeer.cancel()
			close(targetPeer.orders)
			slices.Delete(peers, index, index+1)

			if int(chunks) >= len(torrent.ChunksDownloaded) {
				done = true
				break
			}
		}

		targetPeer.orders <- order(nextChunk)
		nextChunk = getNextChunk(*torrent, chunks)
		ch <- *torrent
	}

	wg.Wait()
	if err := aggregateChunks(*torrent); err != nil {
		return
	}

	torrent.Status = models.DOWNLOADED
	ch <- *torrent

}
