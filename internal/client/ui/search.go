package ui

import (
	"log/slog"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/gen-4/gorrent/internal/client/commands"
	"github.com/gen-4/gorrent/internal/client/models"
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
	spinner   spinner.Model
	loading   bool
	torrents  []string
}

func SearchInitialModel() SearchView {
	input := textinput.New()
	input.Placeholder = "Search"
	input.Width = 30
	input.CharLimit = 32
	input.Focus()

	loadingSpinner := spinner.New()
	loadingSpinner.Spinner = spinner.Dot
	loadingSpinner.Style = SpinnerStyle

	return SearchView{
		help:      help.New(),
		keys:      searchKeys,
		textInput: input,
		err:       nil,
		spinner:   loadingSpinner,
		loading:   false,
		torrents:  []string{},
	}
}

func (s SearchView) Init() tea.Cmd {
	return textinput.Blink
}

func (s SearchView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, s.keys.Help):
			s.help.ShowAll = !s.help.ShowAll

		case key.Matches(msg, s.keys.Enter):
			s.textInput.Blur()
			s.loading = true
			s.torrents = []string{}
			cmds = append(cmds, s.spinner.Tick)
			cmds = append(cmds, commands.GetSuperserverTorrents(s.textInput.Value()))
		case key.Matches(msg, s.keys.Edit):
			if s.textInput.Focused() {
				s.textInput.Blur()
			} else {
				s.textInput.Focus()
				cmds = append(cmds, textinput.Blink)
			}
		}

	case models.SuperserverTorrents:
		s.torrents = msg.Torrents
		s.loading = false

	case spinner.TickMsg:
		if s.loading {
			s.spinner, cmd = s.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case error:
		slog.Error("Error in search view", "error", msg.Error())
		s.err = msg
		return s, nil
	}

	s.textInput, cmd = s.textInput.Update(msg)
	cmds = append(cmds, cmd)
	return s, tea.Batch(cmds...)
}

func (s SearchView) View() string {
	helpView := lipgloss.NewStyle().MarginTop(1).Render(s.help.View(s.keys))
	title := TitleStyle.Render("Search for a torrent")
	spin := ""
	if s.loading {
		spin = lipgloss.NewStyle().MarginTop(1).MarginLeft(1).Render(s.spinner.View())
	}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		title,
		lipgloss.NewStyle().Render(s.textInput.View()),
		spin,
		strings.Join(s.torrents, "\n"),
		helpView,
	)

	return view
}
