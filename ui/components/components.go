package components

import (
	"fmt"
	"math"
	"osprey/config"
	"osprey/data/torrents"
	"osprey/ui/styling"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/muesli/termenv"
)

func VersionNumber() string {
	return "osprey " + styling.ColorFg(config.Osprey_version, styling.HighlightedColor)
}

func Checkbox(label string, checked bool, selected bool) string {
	s := fmt.Sprintf("[ ] %s", label)
	if checked {
		s = "[x] " + label
	}
	if selected {
		return styling.ColorFg(s, styling.HighlightedColor)
	}
	return s
}

func Progressbar(width int, percent float64) string {
	w := float64(styling.ProgressBarWidth)

	fullSize := int(math.Round(w * percent))
	var fullCells string
	for i := 0; i < fullSize; i++ {
		fullCells += termenv.String(styling.ProgressFullChar).Foreground(styling.Term.Color(styling.Ramp[i])).String()
	}

	emptySize := int(w) - fullSize
	emptyCells := strings.Repeat(styling.ProgressEmpty, emptySize)

	return fmt.Sprintf("%s%s %3.0f", fullCells, emptyCells, math.Round(percent*100))
}

func Torrent(torrent torrents.Torrent, selected bool) string {
	s := ""
	torrentNameString := fmt.Sprintf("- %-9s %s\n", fmt.Sprintf("[%s]", torrents.StateString(torrent)), torrent.Name)
	if selected {
		s += styling.ColorFg(torrentNameString, styling.HighlightedColor)
	} else {
		s += styling.ColorFg(torrentNameString, torrents.StateColor(torrent))
	}
	s += fmt.Sprintf("↓ %-9s  ↑ %-9s  ↔ %-9s  P %-6d  S %-6d\n", humanize.Bytes(torrent.DownloadRate)+"/s", humanize.Bytes(torrent.UploadRate)+"/s", humanize.Bytes(torrent.Size), torrent.NumPeers, torrent.NumSeeds)
	s += Progressbar(20, torrent.Progress) + "\n"

	s += "\n"
	return s
}

func KeybindsHints(keybinds []string) string {
	s := ""
	for index, keybind := range keybinds {
		if index != 0 {
			s += styling.Dot
		}
		s += styling.Subtle(keybind)
	}
	return s
}
