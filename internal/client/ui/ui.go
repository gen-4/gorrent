package ui

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		slog.Error(fmt.Sprintf("Error running TUI"), "error", err.Error())
		os.Exit(1)
	}
}
