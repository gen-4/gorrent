package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SearchView struct{}

func SearchInitialModel() SearchView {
	return SearchView{}
}

func (s SearchView) Init() tea.Cmd {
	return nil
}

func (s SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}

	return s, nil
}

func (s SearchView) View() string {
	return "Hello from Search View"
}
