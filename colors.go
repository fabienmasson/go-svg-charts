package charts

var DefaultPalette = []string{"0f7b6d", "6940a6", "0b6e98", "dfab00", "ac1b72", "df3e3e"}
var PopPalette = []string{"1be7ff", "6eeb83", "e4ff1a", "ffb800", "ff5714", "ffbe0b", "fb5607", "ff006e", "8338ec", "3a86ff"}
var PastelPalette = []string{"ff99c8", "fcf6bd", "d0f4de", "a9def9", "e4c1f9", "d8e2dc", "ffe5d9", "ffcad4", "f4acb7", "9d8189"}

type ColorScheme struct {
	Foreground   string
	Background   string
	ColorPalette []string
}
