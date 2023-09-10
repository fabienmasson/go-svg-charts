package charts

import (
	"fmt"
	"math"
)

var DefaultColorScheme = ColorScheme{
	Foreground:      "#000",
	Background:      "#fff",
	LightAxisColor:  "#eee",
	DarkerAxisColor: "#777",
	ColorPalette:    defaultColorPalette,
}

type ColorScheme struct {
	Foreground      string
	Background      string
	LightAxisColor  string
	DarkerAxisColor string
	ColorPalette    ColorPalette
}

type ColorPalette func(i int) string

func defaultColorPalette(i int) string {
	s := 0.5
	l := 0.5
	h, _, _ := rGBToHSL(0, 0, 255)
	h = float64((int(h) + int(i)*69) % 360)
	r, g, b := hSLToRGB(h, s, l)
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func rGBToHSL(r, g, b uint8) (h, s, l float64) {

	Rnot := float64(r) / 255
	Gnot := float64(g) / 255
	Bnot := float64(b) / 255
	Cmax, Cmin := getMaxMin(Rnot, Gnot, Bnot)
	delta := Cmax - Cmin
	// Lightness
	l = (Cmax + Cmin) / 2
	// Hue and Saturation
	if delta == 0 {
		h = 0
		s = 0
	} else {
		switch Cmax {
		case Rnot:
			h = 60 * (math.Mod((Gnot-Bnot)/delta, 6))
		case Gnot:
			h = 60 * (((Bnot - Rnot) / delta) + 2)
		case Bnot:
			h = 60 * (((Rnot - Gnot) / delta) + 4)
		}
		if h < 0 {
			h += 360
		}

		s = delta / (1 - math.Abs((2*l)-1))
	}

	return h, round(s), round(l)
}

func hSLToRGB(h, s, l float64) (r, g, b uint8) {

	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64

	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}
	r = uint8(math.Round((Rnot + m) * 255))
	g = uint8(math.Round((Gnot + m) * 255))
	b = uint8(math.Round((Bnot + m) * 255))
	return r, g, b
}

func getMaxMin(a, b, c float64) (max, min float64) {
	if a > b {
		max = a
		min = b
	} else {
		max = b
		min = a
	}
	if c > max {
		max = c
	} else if c < min {
		min = c
	}
	return max, min
}

func round(x float64) float64 {
	return math.Round(x*1000) / 1000
}
