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
		Cursor:      0,
		CurrentView: "TorrentList",
		Progress:    0.0,
		TorrentList: torrents.TorrentList{},
	}
}

type Model struct {
	Cursor      int
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
	case "RemoveTorrent":
		return updateRemoveTorrentView(msg, m)
	}
	return m, nil
}

func updateRemoveTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentView = "TorrentList"
		case "y":
			http.DeleteTorrent(m.TorrentList.Torrents[m.Cursor], true)
			m.CurrentView = "TorrentList"
		case "n":
			http.DeleteTorrent(m.TorrentList.Torrents[m.Cursor], false)
			m.CurrentView = "TorrentList"
		}
	case tickMsg:
		m.TorrentList = http.UpdateTorrentList()
		return m, tick()
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
	case tickMsg:
		m.TorrentList = http.UpdateTorrentList()
		return m, tick()
	}
	return m, nil
}

func updateListView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.CurrentView = "AddTorrent"
		case "r":
			if len(m.TorrentList.Torrents) != 0 {
				m.CurrentView = "RemoveTorrent"
			}
		case "p":
			http.PauseResumeTorrent(m.TorrentList.Torrents[m.Cursor])
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.TorrentList.Torrents)-1 {
				m.Cursor++
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
	case "RemoveTorrent":
		s = removeTorrentView(m)
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

func removeTorrentView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]

	tpl := fmt.Sprintf("Deleting %s\n\n", selectedTorrent.Name)
	tpl += "Keep data?\n\n"
	tpl += components.KeybindsHints([]string{"y: yes", "n: no", "esc: back"})

	return fmt.Sprintf(tpl)
}

func listView(m Model) string {
	tpl := "%s active:\n"
	for index, torrent := range m.TorrentList.Torrents {
		tpl += components.Torrent(torrent, index == m.Cursor)
	}
	tpl += components.KeybindsHints([]string{"j/k, up/down: select", "p: pause/resume torrent", "a: add new torrent", "r: remove torrent", "q: quit"})

	return fmt.Sprintf(tpl, english.Plural(m.TorrentList.TorrentsTotal, "torrent", ""))
}
