package commands

import tea "github.com/charmbracelet/bubbletea"

func SendMessageCmd(msg any) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return msg
	})
}
