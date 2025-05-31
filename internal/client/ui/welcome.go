package ui

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
)

type Welcome struct{}

func WelcomeInitialModel() Welcome {
	return Welcome{}
}

func (m Welcome) Init() tea.Cmd {
	return nil
}

func (m Welcome) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		slog.Debug("It was a Keymsgggg in Welcome")
		switch msg.String() {
		case "a":
			slog.Debug("He pressed aaaa")
			return m, func() tea.Msg { return WELCOME_VIEW }
		}
	}

	return m, nil
}

func (m Welcome) View() string {
	return "Hello, helloooo!!"
}
