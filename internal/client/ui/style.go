package ui

import "github.com/charmbracelet/lipgloss"

const (
	PURPLE = lipgloss.Color("#521296")
	YELLOW = lipgloss.Color("#FCED77")
)

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(PURPLE).
	Padding(0, 1).
	MarginBottom(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(YELLOW)

var SpinnerStyle = lipgloss.NewStyle().
	Foreground(YELLOW)
