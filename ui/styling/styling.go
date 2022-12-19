package styling

import (
	"fmt"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
)

const (
	HighlightedColor = "212"
	SecondaryColor   = "225"
)

const (
	ProgressBarWidth  = 71
	ProgressFullChar  = "█"
	ProgressEmptyChar = "░"
)

var (
	Term          = termenv.EnvColorProfile()
	Keyword       = MakeFgStyle("211")
	Subtle        = MakeFgStyle("241")
	ProgressEmpty = Subtle(ProgressEmptyChar)
	Dot           = ColorFg(" • ", "236")

	// Gradient colors we'll use for the progress bar
	Ramp = MakeRamp("#B14FFF", "#00FFA3", ProgressBarWidth)
)

// Utils

// Color a string's foreground with the given value.
func ColorFg(val, color string) string {
	return termenv.String(val).Foreground(Term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func MakeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(Term.Color(color)).Styled
}

// Color a string's foreground and background with the given value.
func MakeFgBgStyle(fg, bg string) func(string) string {
	return termenv.Style{}.
		Foreground(Term.Color(fg)).
		Background(Term.Color(bg)).
		Styled
}

// Generate a blend of colors.
func MakeRamp(colorA, colorB string, steps float64) (s []string) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, ColorToHex(c))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format compatible with termenv.
func ColorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", ColorFloatToHex(c.R), ColorFloatToHex(c.G), ColorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func ColorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}
