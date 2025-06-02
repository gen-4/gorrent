package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeKeyMap struct {
	Quit  key.Binding
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Help  key.Binding
}

func (k WelcomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Up, k.Down, k.Enter}
}

func (k WelcomeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help, k.Up, k.Down, k.Enter},
		{},
	}
}

var welcomeKeys = WelcomeKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "h"),
		key.WithHelp("?/h", "help"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Select"),
	),
}

type ModelItem struct {
	name  string
	model Model
}

type Welcome struct {
	models []ModelItem
	cursor int8
	help   help.Model
	keys   WelcomeKeyMap
}

func WelcomeInitialModel() Welcome {
	return Welcome{
		models: []ModelItem{
			ModelItem{
				"Search torrent to download",
				SEARCH_VIEW,
			},
			ModelItem{
				"My torrents",
				TORRENT_VIEW,
			},
		},
		cursor: 0,
		help:   help.New(),
		keys:   welcomeKeys,
	}
}

func (m Welcome) Init() tea.Cmd {
	return nil
}

func (m Welcome) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor -= 1
			}

		case key.Matches(msg, m.keys.Down):
			if m.cursor < int8(len(m.models)-1) {
				m.cursor += 1
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Enter):
			return m, func() tea.Msg { return m.models[m.cursor].model }
		}
	}

	return m, nil
}

func (m Welcome) View() string {
	title := lipgloss.JoinVertical(lipgloss.Top, Title.Render("GORRENT"))

	list := ""
	for i, model := range m.models {
		row := model.name
		if int8(i) == m.cursor {
			row = lipgloss.NewStyle().Background(PURPLE).Render(model.name)
		}

		list += fmt.Sprintf("%s\n", row)
	}

	helpView := m.help.View(m.keys)
	view := lipgloss.JoinVertical(lipgloss.Top, title, list, helpView)
	return view
}
