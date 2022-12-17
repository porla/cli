package ui

import (
	"fmt"
	"time"

	"osprey/config"
	"osprey/data/torrents"
	"osprey/http"
	"osprey/ui/components"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dustin/go-humanize/english"
	"github.com/muesli/reflow/indent"
)

type (
	tickMsg  struct{}
	frameMsg struct{}
)

const (
	AddTorrentMagnetLinkInput = iota
	AddTorrentSavePathInput
)

const (
	LoadingViewIota = iota
	TorrentListIota
	AddTorrentIota
	RemoveTorrentIota
	MoveTorrentIota
	QuittingIota
)

func InitialModel() Model {
	var addTorrentTextInputs []textinput.Model = make([]textinput.Model, 2)
	addTorrentTextInputs[AddTorrentMagnetLinkInput] = textinput.New()
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.MagnetLinkPlaceHolder
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Focus()
	addTorrentTextInputs[AddTorrentMagnetLinkInput].CharLimit = -1
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Width = 50
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Prompt = ""

	addTorrentTextInputs[AddTorrentSavePathInput] = textinput.New()
	addTorrentTextInputs[AddTorrentSavePathInput].Placeholder = config.Currenti18n.SaveDirPlaceHolder
	addTorrentTextInputs[AddTorrentSavePathInput].CharLimit = -1
	addTorrentTextInputs[AddTorrentSavePathInput].Width = 50
	addTorrentTextInputs[AddTorrentSavePathInput].Prompt = ""

	moveTorrentPathTextInput := textinput.New()
	moveTorrentPathTextInput.Placeholder = config.Currenti18n.NewSaveDirPlaceHolder
	moveTorrentPathTextInput.Focus()
	moveTorrentPathTextInput.CharLimit = -1
	moveTorrentPathTextInput.Width = 50
	moveTorrentPathTextInput.Prompt = ""

	return Model{
		Page:                     0,
		Cursor:                   0,
		SubMenuCursor:            0,
		CurrentView:              TorrentListIota,
		Progress:                 0.0,
		TorrentList:              torrents.TorrentList{},
		AddTorrentTextInputs:     addTorrentTextInputs,
		AddingMagnetLink:         true,
		MoveTorrentPathTextInput: moveTorrentPathTextInput,
	}
}

type Model struct {
	Page                     int
	Cursor                   int
	SubMenuCursor            int
	CurrentView              int
	Progress                 float64
	TorrentList              torrents.TorrentList
	AddTorrentTextInputs     []textinput.Model
	AddingMagnetLink         bool
	MoveTorrentPathTextInput textinput.Model
}

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

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		newPageSize := (msg.Height - 8) / 4
		if newPageSize != 0 {
			config.Config.PageSize = newPageSize
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.CurrentView {
	case TorrentListIota:
		return updateListView(msg, m)
	case AddTorrentIota:
		return updateAddTorrentView(msg, m)
	case RemoveTorrentIota:
		return updateRemoveTorrentView(msg, m)
	case MoveTorrentIota:
		return updateMoveTorrentView(msg, m)
	}
	return m, nil
}

func updateRemoveTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentView = TorrentListIota
		case "y":
			http.DeleteTorrent(m.TorrentList.Torrents[m.Cursor], true)
			m.CurrentView = TorrentListIota
		case "n":
			http.DeleteTorrent(m.TorrentList.Torrents[m.Cursor], false)
			m.CurrentView = TorrentListIota
		}
	case tickMsg:
		return m, tick()
	}
	return m, nil
}

func updateAddTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.AddTorrentTextInputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentView = TorrentListIota

		case "up":
			if m.SubMenuCursor > 0 {
				m.SubMenuCursor--
			}
		case "down":
			if m.SubMenuCursor < len(m.AddTorrentTextInputs)-1 {
				m.SubMenuCursor++
			}

		case "tab":
			m.AddingMagnetLink = !m.AddingMagnetLink
			if m.AddingMagnetLink {
				m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.MagnetLinkPlaceHolder
			} else {
				m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.TorrentFilePath
			}

		case "enter":
			if (m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Value() != "") && (m.AddTorrentTextInputs[AddTorrentSavePathInput].Value() != "") {
				http.AddTorrent(m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Value(), m.AddTorrentTextInputs[AddTorrentSavePathInput].Value(), m.AddingMagnetLink)
				m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Reset()
				m.AddTorrentTextInputs[AddTorrentSavePathInput].Reset()
				m.CurrentView = TorrentListIota
			}

		}

		for i := range m.AddTorrentTextInputs {
			m.AddTorrentTextInputs[i].Blur()
		}
		m.AddTorrentTextInputs[m.SubMenuCursor].Focus()

	case tickMsg:
		return m, tick()
	}

	for i := range m.AddTorrentTextInputs {
		m.AddTorrentTextInputs[i], cmds[i] = m.AddTorrentTextInputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func updateMoveTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.CurrentView = TorrentListIota

		case "enter":
			http.MoveTorrent(m.TorrentList.Torrents[m.Cursor], m.MoveTorrentPathTextInput.Value())
			m.MoveTorrentPathTextInput.Reset()
			m.CurrentView = TorrentListIota
		}
	case tickMsg:
		return m, tick()
	}
	m.MoveTorrentPathTextInput, cmd = m.MoveTorrentPathTextInput.Update(msg)
	return m, cmd
}

func decrementPage(m *Model) {
	if m.Page > 0 {
		m.Page--
		m.TorrentList, m.Page = http.UpdateTorrentList(m.Page)
		m.Cursor = len(m.TorrentList.Torrents) - 1
	}
}

func incrementPage(m *Model) {
	if m.Page < getPageCount(*m)-1 {
		m.Page++
		m.TorrentList, m.Page = http.UpdateTorrentList(m.Page)
		m.Cursor = 0
	}
}

func updateListView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.SubMenuCursor = 0
			m.CurrentView = AddTorrentIota
		case "r":
			if len(m.TorrentList.Torrents) != 0 {
				m.CurrentView = RemoveTorrentIota
			}
		case "p":
			http.PauseResumeTorrent(m.TorrentList.Torrents[m.Cursor])
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				decrementPage(&m)
			}
		case "down", "j":
			if m.Cursor < len(m.TorrentList.Torrents)-1 {
				m.Cursor++
			} else {
				incrementPage(&m)
			}
		case "left", "g":
			decrementPage(&m)
		case "right", "h":
			incrementPage(&m)
		case "m":
			m.MoveTorrentPathTextInput.SetValue(m.TorrentList.Torrents[m.Cursor].SavePath)
			m.CurrentView = MoveTorrentIota
		}
	// Get updated info
	case tickMsg:
		m.TorrentList, m.Page = http.UpdateTorrentList(m.Page)
		if (getPageCount(m)-1 != -1) && (m.Page > getPageCount(m)-1) {
			m.Page = getPageCount(m) - 1
			m.TorrentList, m.Page = http.UpdateTorrentList(m.Page)
			m.Cursor = len(m.TorrentList.Torrents) - 1
		}
		if m.Page < 0 {
			m.Page = 0
			m.TorrentList, m.Page = http.UpdateTorrentList(m.Page)
			m.Cursor = 0
		}
		if m.Cursor > len(m.TorrentList.Torrents)-1 {
			m.Cursor = len(m.TorrentList.Torrents) - 1
		}
		if m.Cursor < 0 {
			m.Cursor = 0
		}
		return m, tick()
	}
	return m, nil
}

func (m Model) View() string {
	var s string
	switch m.CurrentView {

	case LoadingViewIota:
		s = loadingView(m)
		break
	case TorrentListIota:
		s = listView(m)
		break
	case AddTorrentIota:
		s = addTorrentView(m)
		break
	case RemoveTorrentIota:
		s = removeTorrentView(m)
		break
	case MoveTorrentIota:
		s = moveTorrentView(m)
		break
	case QuittingIota:
		return "\n  " + config.Currenti18n.SeeYouLater + "\n\n"
	default:
		return "\n  " + config.Currenti18n.ErrorNonExistantView + "\n\n"
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func loadingView(m Model) string {
	tpl := "osprey %s\n\n"
	tpl += "%s\n\n"
	tpl += config.Currenti18n.ConnectingToPorlaBackend + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.QKeybind})

	return fmt.Sprintf(tpl, config.Osprey_version, components.Progressbar(80, m.Progress))
}

func addTorrentView(m Model) string {
	tpl := config.Currenti18n.AddTorrent + "\n\n"
	if m.AddingMagnetLink {
		tpl += config.Currenti18n.MagnetLink + "\n"
	} else {
		tpl += config.Currenti18n.PathToTorrentFile + "\n"
	}
	tpl += m.AddTorrentTextInputs[AddTorrentMagnetLinkInput].View() + "\n\n"
	tpl += config.Currenti18n.SavePath + "\n"
	tpl += m.AddTorrentTextInputs[AddTorrentSavePathInput].View() + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.ToggleMagnetTorrentKeybind, config.Currenti18n.Keybinds.SelectReducedKeybind, config.Currenti18n.Keybinds.DoneKeybind, config.Currenti18n.Keybinds.EscKeybind, config.Currenti18n.Keybinds.QKeybind})

	return fmt.Sprintf(tpl)
}

func removeTorrentView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]

	tpl := fmt.Sprintf(config.Currenti18n.DeletingTorrentName+"\n\n", selectedTorrent.Name)
	tpl += config.Currenti18n.KeepDataQuestion + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.YesKeybind, config.Currenti18n.Keybinds.NoKeybind, config.Currenti18n.Keybinds.EscKeybind})

	return fmt.Sprintf(tpl)
}

func moveTorrentView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]
	tpl := fmt.Sprintf(config.Currenti18n.MovingTorrentName+"\n\n", selectedTorrent.Name)
	tpl += config.Currenti18n.NewSavePath + "\n"
	tpl += m.MoveTorrentPathTextInput.View() + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.DoneKeybind, config.Currenti18n.Keybinds.EscKeybind, config.Currenti18n.Keybinds.QKeybind})

	return fmt.Sprintf(tpl)
}

func listView(m Model) string {
	tpl := config.Currenti18n.TorrentsActive + "\n"
	for index, torrent := range m.TorrentList.Torrents {
		tpl += components.Torrent(torrent, index == m.Cursor)
	}
	tpl += config.Currenti18n.PageInfo + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.SelectKeybind, config.Currenti18n.Keybinds.ChangePageKeybind, config.Currenti18n.Keybinds.PauseResumeKeybind, config.Currenti18n.Keybinds.AddTorrentKeybind, config.Currenti18n.Keybinds.RemoveTorrentKeybind, config.Currenti18n.Keybinds.MoveTorrentKeybind, config.Currenti18n.Keybinds.QKeybind})

	return fmt.Sprintf(tpl, english.Plural(m.TorrentList.TorrentsTotal, config.Currenti18n.Torrent, ""), m.Page+1, getPageCount(m), config.Config.PageSize)
}

func getPageCount(m Model) int {
	return (m.TorrentList.TorrentsTotal-1)/config.Config.PageSize + 1
}
