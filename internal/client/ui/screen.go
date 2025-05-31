package ui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

type Model string

const (
	WELCOME_VIEW Model = "welcome_view"
)

type Screen struct {
	Fragments map[Model]tea.Model
	Selected  Model
}

func InitialModel() Screen {
	return Screen{
		Fragments: map[Model]tea.Model{
			WELCOME_VIEW: WelcomeInitialModel(),
		},
		Selected: WELCOME_VIEW,
	}
}

func (s Screen) Init() tea.Cmd {
	return s.Fragments[s.Selected].Init()
}

func (s Screen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		slog.Debug("KeyMsg received by Screen")
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		}

	case Model:
		slog.Debug("Case Modelll")
		switch msg {
		case WELCOME_VIEW:
			s.Selected = WELCOME_VIEW
		}
	}
	updatedModel, cmd := s.Fragments[s.Selected].Update(msg)
	s.Fragments[s.Selected] = updatedModel
	return s, cmd
}

func (s Screen) View() string {
	return s.Fragments[s.Selected].View()
}
