package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	Quit key.Binding
	Back key.Binding
}

var keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b", "backspace"),
		key.WithHelp("<-", "go back"),
	),
}

type Model string

const (
	WELCOME_VIEW Model = "welcome_view"
	SEARCH_VIEW  Model = "search_view"
	TORRENT_VIEW Model = "torrent_view"
)

type Screen struct {
	Fragments map[Model]tea.Model
	Selected  Model
	Help      help.Model
	Keys      KeyMap
}

func InitialModel() Screen {
	return Screen{
		Fragments: map[Model]tea.Model{
			WELCOME_VIEW: WelcomeInitialModel(),
			SEARCH_VIEW:  SearchInitialModel(),
			TORRENT_VIEW: TorrentsInitialModel(),
		},
		Selected: WELCOME_VIEW,
		Help:     help.New(),
		Keys:     keys,
	}
}

func (s Screen) Init() tea.Cmd {
	return s.Fragments[s.Selected].Init()
}

func (s Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.Keys.Quit):
			return s, tea.Quit

		case key.Matches(msg, s.Keys.Back):
			s.Selected = WELCOME_VIEW
		}

	case Model:
		switch msg {
		case WELCOME_VIEW:
			s.Selected = WELCOME_VIEW

		case SEARCH_VIEW:
			s.Selected = SEARCH_VIEW

		case TORRENT_VIEW:
			s.Selected = TORRENT_VIEW
		}
	}
	updatedModel, cmd := s.Fragments[s.Selected].Update(msg)
	s.Fragments[s.Selected] = updatedModel
	return s, cmd
}

func (s Screen) View() string {
	return s.Fragments[s.Selected].View()
}
