package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	PRIMARY   = lipgloss.Color("#521296")
	SECONDARY = lipgloss.Color("#FCED77")
)

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(PRIMARY).
	Padding(0, 1).
	MarginBottom(1).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(SECONDARY)

var SpinnerStyle = lipgloss.NewStyle().
	Foreground(SECONDARY)

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(PRIMARY)

func SetTableRowStyles() table.Styles {
	var tableRowStyles table.Styles = table.DefaultStyles()
	tableRowStyles.Header = tableRowStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(SECONDARY).
		BorderBottom(true).
		Bold(false)
	tableRowStyles.Selected = tableRowStyles.Selected.
		Foreground(PRIMARY).
		Background(SECONDARY).
		Bold(false)

	return tableRowStyles
}
