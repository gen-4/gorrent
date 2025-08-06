package commands

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	config "github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/models"
)

func GetSuperserverTorrents(torrentName string) tea.Cmd {
	return func() tea.Msg {
		torrents := []string{}

		for _, superserver := range config.Configuration.Superservers {
			var superserverTorrents models.SuperserverTorrentsDto

			response, err := http.Get(fmt.Sprintf(
				config.Configuration.SuperserverUrlTemplate,
				superserver,
				fmt.Sprintf("torrents?name=%s", torrentName),
			))
			if err != nil {
				slog.Error("Error requesting superserver torrents", "superserver", superserver, "error", err.Error())
				continue
			}

			if response.StatusCode != http.StatusOK {
				slog.Error("Wrong response code requesting superserver torrents", "superserver", superserver, "status", response.StatusCode)
				continue
			}

			if err := json.NewDecoder(response.Body).Decode(&superserverTorrents); err != nil {
				slog.Error("Error decoding superserver torrents", "superserver", superserver, "error", err.Error())
				continue
			}
			torrents = append(torrents, superserverTorrents.Torrents...)
		}

		return models.SuperserverTorrents{Torrents: torrents}
	}
}
