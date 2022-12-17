package ui

import (
	"fmt"
	"time"

	"osprey/config"
	"osprey/data/torrents"
	"osprey/http"
	"osprey/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize/english"
	"github.com/muesli/reflow/indent"
)

func InitialModel() Model {
	return Model{
		Choices:     []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		Cursor:      0,
		Selected:    make(map[int]struct{}),
		CurrentView: "TorrentList",
		Progress:    0.8,
		TorrentList: torrents.TorrentList{},
	}
}

type Model struct {
	Choices     []string
	Cursor      int
	Selected    map[int]struct{}
	CurrentView string
	Progress    float64
	TorrentList torrents.TorrentList
}

type (
	tickMsg  struct{}
	frameMsg struct{}
)

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func frame() tea.Cmd {
	return tea.Tick(time.Second/60, func(time.Time) tea.Msg {
		return frameMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.CurrentView {
	case "TorrentList":
		return updateListView(msg, m)
	case "AddTorrent":
		return updateAddTorrentView(msg, m)
	}
	return m, nil
}

func updateAddTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentView = "TorrentList"
		}
	}
	return m, nil
}

func updateListView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.CurrentView = "AddTorrent"

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.TorrentList.Torrents)-1 {
				m.Cursor++
			}

		case "enter", " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = struct{}{}
			}
		}

	// Get updated info
	case tickMsg:
		m.TorrentList = http.UpdateTorrentList()
		return m, tick()
	}
	return m, nil
}

func (m Model) View() string {
	var s string
	switch m.CurrentView {

	case "Loading":
		s = loadingView(m)
		break
	case "TorrentList":
		s = listView(m)
		break
	case "AddTorrent":
		s = addTorrentView(m)
		break
	case "Quitting":
		return "\n  See you later!\n\n"
	default:
		return "\n  Error: Non existant view called.\n\n"
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func loadingView(m Model) string {
	tpl := "osprey %s\n\n"
	tpl += "%s\n\n"
	tpl += "establishing connection to Porla backend.\n\n"
	tpl += components.KeybindsHints([]string{"q: quit"})

	return fmt.Sprintf(tpl, config.Osprey_version, components.Progressbar(80, m.Progress))
}

func addTorrentView(m Model) string {
	tpl := "Add torrent\n\n"
	tpl += components.KeybindsHints([]string{"esc: back", "q: quit"})

	return fmt.Sprintf(tpl)
}

func listView(m Model) string {
	c := m.Cursor

	tpl := "%s active:\n"
	for index, torrent := range m.TorrentList.Torrents {
		tpl += components.Torrent(torrent, index == c)
	}
	tpl += components.KeybindsHints([]string{"a: add new torrent", "j/k, up/down: select", "enter: choose", "q: quit"})

	return fmt.Sprintf(tpl, english.Plural(m.TorrentList.TorrentsTotal, "torrent", ""))
}
