package ui

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SearchKeyMap struct {
	Quit  key.Binding
	Back  key.Binding
	Enter key.Binding
	Help  key.Binding
	Edit  key.Binding
}

func (k SearchKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Help, k.Back, k.Edit, k.Enter}
}

func (k SearchKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Quit, k.Help, k.Back, k.Enter},
		{k.Edit},
	}
}

var searchKeys = SearchKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "go back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Search"),
	),
	Edit: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle write mode"),
	),
}

type SearchView struct {
	textInput textinput.Model
	err       error
	help      help.Model
	keys      SearchKeyMap
}

func SearchInitialModel() SearchView {
	input := textinput.New()
	input.Placeholder = "Search"
	input.Width = 30
	input.CharLimit = 32
	input.Focus()

	return SearchView{
		help:      help.New(),
		keys:      searchKeys,
		textInput: input,
		err:       nil,
	}
}

func (s SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (s SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Help):
			s.help.ShowAll = !s.help.ShowAll

		case key.Matches(msg, s.keys.Enter):
			s.textInput.Blur()
			// TODO: actually search in the superserver

		case key.Matches(msg, s.keys.Edit):
			if s.textInput.Focused() {
				s.textInput.Blur()
			} else {
				s.textInput.Focus()
			}
		}

	case error:
		slog.Error("Error in search view", "error", msg.Error())
		s.err = msg
		return s, nil
	}

	s.textInput, cmd = s.textInput.Update(msg)
	return s, cmd
}

func (s SearchView) View() string {
	helpView := s.help.View(s.keys)
	title := Title.Render("Search for a torrent")
	view := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		lipgloss.NewStyle().MarginBottom(1).Render(s.textInput.View()),
		helpView,
	)

	return view
}
