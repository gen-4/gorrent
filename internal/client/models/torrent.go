package models

type Status string

const (
	IN_PROGRESS Status = "in_progress"
	DOWNLOADED         = "downloaded"
	STOPPED            = "stopped"
	SEEDING            = "seeding"
)

type Torrent struct {
	Name             string
	Peers            uint8
	Progress         uint8
	Status           Status
	Superservers     []string
	ChunkLength      uint64
	Length           uint64
	DownloadDir      string
	ChunksDownloaded []uint8
}
