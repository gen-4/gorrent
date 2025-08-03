package models

type Status string

const (
	IN_PROGRESS Status = "in_progress"
	DOWNLOADED         = "downloaded"
	STOPPED            = "stopped"
	SEEDING            = "seeding"
)

type Torrent struct {
	Name     string
	Peers    int
	Progress int
	Status   Status
}
