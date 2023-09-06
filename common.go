package charts

import (
	"fmt"
	"io"
	"math"
)

type Dimension struct {
	width, height int
}

// yAxisDimensions calculates the y-axis dimensions for a given height and data.
//
// Parameters:
// - height: the height of the y-axis.
// - data: a 2D slice of float64 values representing the data.
//
// Returns:
// - slice of line labels
// - slice of lines y
// - slice of converted data y
func yAxisFit(start int, end int, data [][]float64) ([]string, []float64, func(float64) float64) {
	min, max := data[0][0], data[0][0]
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] < min {
				min = data[i][j]
			}
			if data[i][j] > max {
				max = data[i][j]
			}
		}
	}

	diff := max - min
	log10 := math.Log10(diff)
	i, f := math.Modf(log10)
	var interval float64
	if f < 0.3 {
		interval = math.Pow10(int(i)) / 5
	} else if f < 0.7 {
		interval = math.Pow10(int(i)) / 2
	} else {
		interval = math.Pow10(int(i))
	}

	height := float64(end - start)
	top := 0.0       // where max value goes
	bottom := height // where min value goes

	fmt.Printf("top:%f; bottom:%f; start:%d; end: %d; min: %f, max: %f, interval: %f\n", top, bottom, start, end, min, max, interval)

	conv := func(val float64) float64 {
		return float64(start) + bottom - (bottom-top)*(val-min)/(max-min)
	}

	k := 0
	labels := make([]string, 0)
	lines := make([]float64, 0)
	for {
		val := float64(int(min/interval)+1+k) * interval
		if val < max {
			if val > 1 || val <= -1 {
				labels = append(labels, fmt.Sprintf("%d", int(val)))
			} else {
				labels = append(labels, fmt.Sprintf("%f", val))
			}
			lines = append(lines, val)
		} else {
			break
		}
		k++
	}

	return labels, lines, conv
}

func seriesLegend(
	w io.Writer,
	width int,
	markerModulo int,
	series []string,
	colors *ColorScheme) int {
	const samplewidth = 30
	const sampleHeight = 15
	const labelwidth = 70
	const gap = 5

	x := 10
	y := 10

	for s, serie := range series {

		fmt.Fprintf(
			w,
			"<polyline points='%d,%d %d,%d %d,%d' fill='none' stroke='%s' stroke-width='2' marker-mid='url(#dot%d)' />",
			x, y,
			x+samplewidth/2, y,
			x+samplewidth, y,
			colors.ColorPalette[s%len(colors.ColorPalette)],
			s%markerModulo,
		)
		x += samplewidth + gap
		fmt.Fprintf(
			w,
			"<text x='%d' y='%d'>%s</text>",
			x, y,
			serie,
		)
		x += labelwidth + gap

		if x+samplewidth+labelwidth > width {
			x = 10
			y += sampleHeight + gap
		}
	}
	legendHeight := y + sampleHeight + gap
	return legendHeight
}

func writeFontStyle(w io.Writer) {
	fmt.Fprintf(
		w,
		"<style>text { font-size: 10px; font-family: sans-serif }</style>",
	)
}

func startSVG(w io.Writer, width, height int, colorScheme *ColorScheme) {
	fmt.Fprintf(
		w,
		"<svg varsion='1.1' xmlns='http://www.w3.org/2000/svg' width='%d' height='%d'>",
		width,
		height,
	)
}

func writeDefsMarkers(w io.Writer, size float64, n int, colors *ColorScheme) int {

	const maxMarkers = 7

	fullsize := float64(size)
	halfsize := float64(size / 2)

	// marker
	fmt.Fprintf(w, "<defs>")
	for i := 0; i < n && i < 7; i++ {
		fmt.Fprintf(
			w,
			"<marker id='dot%d' viewBox='0 0 %f %f' refX='%f' refY='%f'  markerWidth='%f' markerHeight='%f'>",
			i,
			fullsize, fullsize,
			halfsize, halfsize,
			halfsize, halfsize,
		)
		switch i {
		case 0:
			fmt.Fprintf(w, "<circle cx='%f' cy='%f' r='%f' fill='%s' />", halfsize, halfsize, halfsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 1:
			fmt.Fprintf(w, "<rect x='0' y='0' width='%f' height='10' fill='%s' />", fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 2:
			fmt.Fprintf(w, "<polygon points='0,%f %f,0 %f,%f' fill='%s' />", fullsize, halfsize, fullsize, fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 3:
			fmt.Fprintf(w, "<line x1='0' y1='0' x2='%f' y2='%f' stroke='%s' stroke-width='1.5'/>", fullsize, fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
			fmt.Fprintf(w, "<line x1='0' y1='%f' x2='%f' y2='0' stroke='%s' stroke-width='1.5'/>", fullsize, fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 4:
			fmt.Fprintf(w, "<circle cx='%f' cy='%f' r='%f' stroke='%s' stroke-width='1.5' fill='none'/>", halfsize, halfsize, halfsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 5:
			fmt.Fprintf(w, "<rect x='0' y='0' width='%f' height='%f' stroke='%s' fill='none' />", fullsize, fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		case 6:
			fmt.Fprintf(w, "<polygon points='0,%f %f,0 %f,%f' stroke='%s' fill='none' />", fullsize, halfsize, fullsize, fullsize, colors.ColorPalette[i%len(colors.ColorPalette)])
		}

		fmt.Fprintf(w, "</marker>")
	}
	fmt.Fprintf(w, "</defs>")

	return maxMarkers

}

func endSVG(w io.Writer) {
	fmt.Fprintf(w, "</svg>")
}