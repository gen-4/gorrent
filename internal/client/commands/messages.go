package commands

import tea "github.com/charmbracelet/bubbletea"

func SendMessageCmd(torrent any) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return torrent
	})
}
