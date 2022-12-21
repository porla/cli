package ui

import (
	"fmt"
	"strconv"
	"time"

	"osprey/config"
	"osprey/data/torrents"
	"osprey/http"
	"osprey/ui/components"
	"osprey/ui/styling"
	"osprey/utils"

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
	LoadingViewIota = iota
	TorrentListIota
	AddTorrentIota
	RemoveTorrentIota
	MoveTorrentIota
	TorrentSettingsIota
	QuittingIota
)

const (
	AddTorrentMagnetLinkInput = iota
	AddTorrentSavePathInput
)

const (
	TorrentSettingsDownloadLimitInput = iota
	TorrentSettingsMaxConnectionsInput
	TorrentSettingsMaxUploadsInput
	TorrentSettingsUploadLimitInput
)

func InitialModel() Model {
	// Add torrent text inputs
	var addTorrentTextInputs []textinput.Model = make([]textinput.Model, 2)
	for i := range addTorrentTextInputs {
		addTorrentTextInputs[i] = textinput.New()
		addTorrentTextInputs[i].CharLimit = -1
		addTorrentTextInputs[i].Width = 50
		addTorrentTextInputs[i].Prompt = ""
	}
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.MagnetLinkPlaceHolder
	addTorrentTextInputs[AddTorrentMagnetLinkInput].Focus()
	addTorrentTextInputs[AddTorrentSavePathInput].Placeholder = config.Currenti18n.SaveDirPlaceHolder

	// Move torrent text input
	moveTorrentPathTextInput := textinput.New()
	moveTorrentPathTextInput.Placeholder = config.Currenti18n.NewSaveDirPlaceHolder
	moveTorrentPathTextInput.Focus()
	moveTorrentPathTextInput.CharLimit = -1
	moveTorrentPathTextInput.Width = 50
	moveTorrentPathTextInput.Prompt = ""

	// TorrentSettings text inputs
	var torrentSettingsTextInputs []textinput.Model = make([]textinput.Model, 4)
	for i := range torrentSettingsTextInputs {
		torrentSettingsTextInputs[i] = textinput.New()
		torrentSettingsTextInputs[i].CharLimit = -1
		torrentSettingsTextInputs[i].Width = 50
		torrentSettingsTextInputs[i].Prompt = ""
	}
	torrentSettingsTextInputs[TorrentSettingsDownloadLimitInput].Placeholder = "-1"
	torrentSettingsTextInputs[TorrentSettingsMaxConnectionsInput].Placeholder = "16777215"
	torrentSettingsTextInputs[TorrentSettingsMaxUploadsInput].Placeholder = "1000"
	torrentSettingsTextInputs[TorrentSettingsUploadLimitInput].Placeholder = "-1"

	return Model{
		Page:           0,
		Cursor:         0,
		SubMenuCursor:  0,
		SubMenuEntries: 0,
		CurrentView:    TorrentListIota,
		Progress:       0.0,
		TorrentList:    torrents.TorrentList{},
		AddTorrentSubMenuState: AddTorrentSubMenuState{
			AddTorrentTextInputs: addTorrentTextInputs,
			AddingMagnetLink:     true,
		},
		MoveTorrentSubMenuState: MoveTorrentSubMenuState{
			MoveTorrentPathTextInput: moveTorrentPathTextInput,
		},
		TorrentSettingsSubMenuState: TorrentSettingsSubMenuState{
			TorrentIsAutomaticallyManaged:    false,
			TorrentIsSequenciallyDownloading: false,
			TorrentSettingsTextInputs:        torrentSettingsTextInputs,
		},
		NinjaMode: false,
	}
}

type TorrentSettingsSubMenuState struct {
	TorrentIsAutomaticallyManaged    bool
	TorrentIsSequenciallyDownloading bool
	TorrentSettingsTextInputs        []textinput.Model
}

type AddTorrentSubMenuState struct {
	AddTorrentTextInputs []textinput.Model
	AddingMagnetLink     bool
}

type MoveTorrentSubMenuState struct {
	MoveTorrentPathTextInput textinput.Model
}

type Model struct {
	Page                        int
	Cursor                      int
	SubMenuCursor               int
	SubMenuEntries              int
	CurrentView                 int
	Progress                    float64
	TorrentList                 torrents.TorrentList
	AddTorrentSubMenuState      AddTorrentSubMenuState
	MoveTorrentSubMenuState     MoveTorrentSubMenuState
	TorrentSettingsSubMenuState TorrentSettingsSubMenuState
	NinjaMode                   bool
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
	var cmd tea.Cmd
	cmd = tea.EnterAltScreen
	return tea.Batch(tick(), cmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		// Check if is a submenu
		if utils.Contains([]int{AddTorrentIota, MoveTorrentIota, RemoveTorrentIota, TorrentSettingsIota}, m.CurrentView) {
			switch k {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.CurrentView = TorrentListIota
			case "up":
				if m.SubMenuCursor > 0 {
					m.SubMenuCursor--
				}
			case "down":
				if m.SubMenuCursor < m.SubMenuEntries-1 {
					m.SubMenuCursor++
				}
			}
		} else {
			if k == "q" || k == "ctrl+c" {
				return m, tea.Quit
			}
		}
	}

	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		newPageSize := (msg.Height - 11) / 4 // The torrent elements are 4 lines high and there are 11 lines not used for displaying torrents
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
	case TorrentSettingsIota:
		return updateTorrentSettingsView(msg, m)
	}
	return m, nil
}

func updateRemoveTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
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
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.AddTorrentSubMenuState.AddTorrentTextInputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.AddTorrentSubMenuState.AddingMagnetLink = !m.AddTorrentSubMenuState.AddingMagnetLink
			if m.AddTorrentSubMenuState.AddingMagnetLink {
				m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.MagnetLinkPlaceHolder
			} else {
				m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Placeholder = config.Currenti18n.TorrentFilePath
			}

		case "enter":
			if (m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Value() != "") && (m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentSavePathInput].Value() != "") {
				http.AddTorrent(m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Value(), m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentSavePathInput].Value(), m.AddTorrentSubMenuState.AddingMagnetLink)
				m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].Reset()
				m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentSavePathInput].Reset()
				m.CurrentView = TorrentListIota
			}

		}

		for i := range m.AddTorrentSubMenuState.AddTorrentTextInputs {
			m.AddTorrentSubMenuState.AddTorrentTextInputs[i].Blur()
		}
		m.AddTorrentSubMenuState.AddTorrentTextInputs[m.SubMenuCursor].Focus()

	case tickMsg:
		return m, tick()
	}

	for i := range m.AddTorrentSubMenuState.AddTorrentTextInputs {
		m.AddTorrentSubMenuState.AddTorrentTextInputs[i], cmds[i] = m.AddTorrentSubMenuState.AddTorrentTextInputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func updateMoveTorrentView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			http.MoveTorrent(m.TorrentList.Torrents[m.Cursor], m.MoveTorrentSubMenuState.MoveTorrentPathTextInput.Value())
			m.MoveTorrentSubMenuState.MoveTorrentPathTextInput.Reset()
			m.CurrentView = TorrentListIota
		}
	case tickMsg:
		return m, tick()
	}
	m.MoveTorrentSubMenuState.MoveTorrentPathTextInput, cmd = m.MoveTorrentSubMenuState.MoveTorrentPathTextInput.Update(msg)
	return m, cmd
}

func updateTorrentSettingsView(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if m.SubMenuCursor == 0 {
				m.TorrentSettingsSubMenuState.TorrentIsAutomaticallyManaged = !m.TorrentSettingsSubMenuState.TorrentIsAutomaticallyManaged
			}
			if m.SubMenuCursor == 1 {
				m.TorrentSettingsSubMenuState.TorrentIsSequenciallyDownloading = !m.TorrentSettingsSubMenuState.TorrentIsSequenciallyDownloading
			}
		case "enter":
			http.SetTorrentProperties(m.TorrentList.Torrents[m.Cursor], torrents.TorrentPropertiesSetData{
				IsAutomaticallyManaged:    m.TorrentSettingsSubMenuState.TorrentIsAutomaticallyManaged,
				IsSequenciallyDownloading: m.TorrentSettingsSubMenuState.TorrentIsSequenciallyDownloading,
				DownloadLimit:             m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsDownloadLimitInput].Value(),
				MaxConnections:            m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxConnectionsInput].Value(),
				MaxUploads:                m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxUploadsInput].Value(),
				UploadLimit:               m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsUploadLimitInput].Value(),
			})
			m.CurrentView = TorrentListIota
		}

		for i := range m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs {
			m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[i].Blur()
		}
		if m.SubMenuCursor > 1 {
			m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[m.SubMenuCursor-2].Focus()
		}

	case tickMsg:
		return m, tick()
	}

	for i := range m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs {
		m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[i], cmds[i] = m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
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
			m.SubMenuEntries = len(m.AddTorrentSubMenuState.AddTorrentTextInputs)
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
			if len(m.TorrentList.Torrents) != 0 {
				m.MoveTorrentSubMenuState.MoveTorrentPathTextInput.SetValue(m.TorrentList.Torrents[m.Cursor].SavePath)
				m.CurrentView = MoveTorrentIota
			}
		case "s":
			if len(m.TorrentList.Torrents) != 0 {
				torrentProperties := http.GetTorrentProperties(m.TorrentList.Torrents[m.Cursor])
				m.TorrentSettingsSubMenuState.TorrentIsAutomaticallyManaged = torrents.IsAutoManaged(uint64(torrentProperties.Flags))
				m.TorrentSettingsSubMenuState.TorrentIsSequenciallyDownloading = torrents.IsSequenciallyDownloading(uint64(torrentProperties.Flags))
				m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsDownloadLimitInput].SetValue(strconv.Itoa(torrentProperties.DownloadLimit))
				m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsUploadLimitInput].SetValue(strconv.Itoa(torrentProperties.UploadLimit))
				m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxConnectionsInput].SetValue(strconv.Itoa(torrentProperties.MaxConnections))
				m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxUploadsInput].SetValue(strconv.Itoa(torrentProperties.MaxUploads))
				for i := range m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs {
					m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[i].Blur()
				}
				m.SubMenuCursor = 0
				m.SubMenuEntries = 2 + len(m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs)
				m.CurrentView = TorrentSettingsIota
			}
		case "n":
			m.NinjaMode = !m.NinjaMode
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
	case TorrentListIota:
		s = listView(m)
	case AddTorrentIota:
		s = addTorrentView(m)
	case RemoveTorrentIota:
		s = removeTorrentView(m)
	case MoveTorrentIota:
		s = moveTorrentView(m)
	case TorrentSettingsIota:
		s = torrentSettingsView(m)
	case QuittingIota:
		return "\n  " + config.Currenti18n.SeeYouLater + "\n\n"
	default:
		return "\n  " + config.Currenti18n.ErrorNonExistantView + "\n\n"
	}
	return indent.String("\n"+s+"\n\n", 2)
}

func loadingView(m Model) string {
	tpl := components.VersionNumber() + "\n\n"
	tpl += "%s\n\n"
	tpl += config.Currenti18n.ConnectingToPorlaBackend + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.QKeybind})

	return fmt.Sprintf(tpl, components.Progressbar(80, m.Progress))
}

func addTorrentView(m Model) string {
	tpl := styling.ColorFg(config.Currenti18n.AddTorrent, styling.SecondaryColor) + "\n\n"
	if m.AddTorrentSubMenuState.AddingMagnetLink {
		tpl += styling.ColorFg(config.Currenti18n.MagnetLink, styling.SecondaryColor) + "\n"
	} else {
		tpl += styling.ColorFg(config.Currenti18n.PathToTorrentFile, styling.SecondaryColor) + "\n"
	}
	tpl += m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentMagnetLinkInput].View() + "\n\n"
	tpl += styling.ColorFg(config.Currenti18n.SavePath, styling.SecondaryColor) + "\n"
	tpl += m.AddTorrentSubMenuState.AddTorrentTextInputs[AddTorrentSavePathInput].View() + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.ToggleMagnetTorrentKeybind, config.Currenti18n.Keybinds.SelectReducedKeybind, config.Currenti18n.Keybinds.DoneKeybind, config.Currenti18n.Keybinds.EscKeybind})

	return fmt.Sprintf(tpl)
}

func removeTorrentView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]

	tpl := fmt.Sprintf(styling.ColorFg(config.Currenti18n.DeletingTorrentName, styling.SecondaryColor)+"\n\n", selectedTorrent.Name)
	tpl += config.Currenti18n.KeepDataQuestion + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.YesKeybind, config.Currenti18n.Keybinds.NoKeybind, config.Currenti18n.Keybinds.EscKeybind})

	return fmt.Sprintf(tpl)
}

func moveTorrentView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]
	tpl := fmt.Sprintf(styling.ColorFg(config.Currenti18n.MovingTorrentName, styling.SecondaryColor)+"\n\n", selectedTorrent.Name)
	tpl += styling.ColorFg(config.Currenti18n.NewSavePath, styling.SecondaryColor) + "\n"
	tpl += m.MoveTorrentSubMenuState.MoveTorrentPathTextInput.View() + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.DoneKeybind, config.Currenti18n.Keybinds.EscKeybind})

	return fmt.Sprintf(tpl)
}

func torrentSettingsView(m Model) string {
	selectedTorrent := m.TorrentList.Torrents[m.Cursor]
	tpl := fmt.Sprintf(styling.ColorFg(config.Currenti18n.TorrentSettingsForTorrentName, styling.SecondaryColor)+"\n\n", selectedTorrent.Name)
	tpl += components.Checkbox(config.Currenti18n.AutomaticallyManaged, m.TorrentSettingsSubMenuState.TorrentIsAutomaticallyManaged, m.SubMenuCursor == 0) + "\n"
	tpl += components.Checkbox(config.Currenti18n.SequentialDownload, m.TorrentSettingsSubMenuState.TorrentIsSequenciallyDownloading, m.SubMenuCursor == 1) + "\n\n"
	tpl += styling.ColorFg(config.Currenti18n.DownloadLimit, styling.SecondaryColor) + "\n"
	tpl += m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsDownloadLimitInput].View() + "\n"
	tpl += styling.Subtle(config.Currenti18n.DownloadLimitHint + "\n\n")
	tpl += styling.ColorFg(config.Currenti18n.MaxConnections, styling.SecondaryColor) + "\n"
	tpl += m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxConnectionsInput].View() + "\n\n"
	tpl += styling.ColorFg(config.Currenti18n.MaxUploads, styling.SecondaryColor) + "\n"
	tpl += m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsMaxUploadsInput].View() + "\n\n"
	tpl += styling.ColorFg(config.Currenti18n.UploadLimit, styling.SecondaryColor) + "\n"
	tpl += m.TorrentSettingsSubMenuState.TorrentSettingsTextInputs[TorrentSettingsUploadLimitInput].View() + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.SelectReducedKeybind, config.Currenti18n.Keybinds.ToggleOptionKeybind, config.Currenti18n.Keybinds.DoneKeybind, config.Currenti18n.Keybinds.EscKeybind})

	return fmt.Sprintf(tpl)
}

func listView(m Model) string {
	tpl := components.VersionNumber() + "\n\n"
	tpl += config.Currenti18n.TorrentsActive + "\n"
	for index, torrent := range m.TorrentList.Torrents {
		tpl += components.Torrent(torrent, index, m.NinjaMode, index == m.Cursor)
	}
	tpl += styling.Subtle(config.Currenti18n.PageInfo) + "\n\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.SelectKeybind, config.Currenti18n.Keybinds.ChangePageKeybind, config.Currenti18n.Keybinds.PauseResumeKeybind, config.Currenti18n.Keybinds.AddTorrentKeybind}) + "\n"
	tpl += components.KeybindsHints([]string{config.Currenti18n.Keybinds.RemoveTorrentKeybind, config.Currenti18n.Keybinds.MoveTorrentKeybind, config.Currenti18n.Keybinds.TorrentSettingsKeybind, config.Currenti18n.Keybinds.NinjaModeKeybind, config.Currenti18n.Keybinds.QKeybind})
	return fmt.Sprintf(tpl, english.Plural(m.TorrentList.TorrentsTotal, config.Currenti18n.Torrent, ""), m.Page+1, getPageCount(m), config.Config.PageSize)
}

func getPageCount(m Model) int {
	return (m.TorrentList.TorrentsTotal-1)/config.Config.PageSize + 1
}
