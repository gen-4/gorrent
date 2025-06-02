package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Torrents struct{}

func TorrentsInitialModel() Torrents {
	return Torrents{}
}

func (t Torrents) Init() tea.Cmd {
	return nil
}

func (t Torrents) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	return t, nil
}

func (w Torrents) View() string {
	return "Hello from Torrents View"
}
