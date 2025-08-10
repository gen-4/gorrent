package ui

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/commands"
	"github.com/gen-4/gorrent/internal/client/models"
	"github.com/gen-4/gorrent/internal/client/sharing"
)

type TorrentsKeyMap struct {
	Quit    key.Binding
	Back    key.Binding
	Help    key.Binding
	Up      key.Binding
	Down    key.Binding
	DirBack key.Binding
	DirNext key.Binding
	Add     key.Binding
	Enter   key.Binding
	Stop    key.Binding
}

func (k TorrentsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Back}
}

func (k TorrentsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help, k.Back},
		{k.Stop, k.Up, k.Down, k.DirBack, k.DirNext, k.Add, k.Enter},
	}
}

var torrentsKeys = TorrentsKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "go back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("k"),
		key.WithHelp("k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j"),
		key.WithHelp("j", "down"),
	),
	DirBack: key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "previous directory"),
	),
	DirNext: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "navigate to directory"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add new torrent"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "navigate to directory/select torrent"),
	),
	Stop: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit selecting torrent"),
	),
}

type Torrents struct {
	filepicker filepicker.Model
	picking    bool
	err        error
	help       help.Model
	keys       TorrentsKeyMap
	torrents   []models.Torrent
	table      table.Model
	ch         chan models.Torrent
}

func updateTableRows(t *Torrents) {
	rows := []table.Row{}
	for _, torrentRow := range t.torrents {
		rows = append(rows, table.Row{
			torrentRow.File,
			string(torrentRow.Status),
			strconv.Itoa(int(torrentRow.Peers)),
			fmt.Sprintf("%d%%", torrentRow.Progress),
		})
	}
	t.table.SetRows(rows)
}

func TorrentsInitialModel() Torrents {
	var err error
	tableHeaders := []table.Column{
		{Title: "File", Width: 40},
		{Title: "Status", Width: 12},
		{Title: "Peers", Width: 5},
		{Title: "Progress", Width: 15},
	}
	torrentsTable := table.New(
		table.WithColumns(tableHeaders),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	torrentsTable.SetStyles(SetTableRowStyles())

	fp := filepicker.New()
	fp.AllowedTypes = []string{".torrent"}
	fp.CurrentDirectory, err = os.UserHomeDir()
	fp.SetHeight(10)
	fp.Styles.Selected = fp.Styles.Selected.Foreground(SECONDARY)
	fp.Styles.Directory = fp.Styles.Directory.Foreground(PRIMARY)
	fp.Styles.File = fp.Styles.File.Foreground(lipgloss.Color("#FFFFFF"))
	if err != nil {
		slog.Error("Error setting filepicker current directory", "error", err.Error())
	}

	var torrents []models.Torrent = config.Configuration.Torrents

	model := Torrents{
		filepicker: fp,
		picking:    false,
		help:       help.New(),
		keys:       torrentsKeys,
		err:        nil,
		torrents:   torrents,
		table:      torrentsTable,
		ch:         make(chan models.Torrent, 10),
	}
	updateTableRows(&model)

	for i, _ := range model.torrents {
		sharing.DownloadTorrent(&model.torrents[i], model.ch)
	}

	return model
}

func (t Torrents) Init() tea.Cmd {
	return nil
}

func (t Torrents) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	var torrentUpdate models.Torrent

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keys.Help):
			t.help.ShowAll = !t.help.ShowAll

		case key.Matches(msg, t.keys.Stop) && t.picking:
			t.picking = false

		case key.Matches(msg, t.keys.Add):
			t.picking = true
			cmds = append(cmds, t.filepicker.Init())

		default:
			t.table, cmd = t.table.Update(msg)
			cmds = append(cmds, cmd)
		}

	case models.NewTorrentRequest:
		t.torrents = append(t.torrents, models.Torrent(msg))
		updateTableRows(&t)
		sharing.DownloadTorrent(&t.torrents[len(t.torrents)-1], t.ch)

	case models.TorrentUpdate:
		updateTableRows(&t)

	case error:
		slog.Error("Error in torrents view", "error", msg.Error())
		t.err = msg
		return t, nil
	}

	select {
	case torrentUpdate = <-t.ch:
		cmds = append(cmds, commands.SendMessageCmd(models.TorrentUpdate(torrentUpdate)))

	default:
	}

	if t.picking {
		t.filepicker, cmd = t.filepicker.Update(msg)
		cmds = append(cmds, cmd)
		if didSelect, path := t.filepicker.DidSelectFile(msg); didSelect {
			cmds = append(cmds, commands.CreateTorrent(path))
			t.picking = false
		}
	}

	return t, tea.Batch(cmds...)
}

func (t Torrents) View() string {
	helpView := lipgloss.NewStyle().MarginTop(1).Render(t.help.View(t.keys))
	title := TitleStyle.Render("Manage your torrents")
	torrents := ""
	filepickerView := ""
	if t.picking {
		filepickerView = t.filepicker.View()
	} else {
		torrents = TableStyle.Render(t.table.View())
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		filepickerView,
		torrents,
		helpView,
	)

	return view
}
