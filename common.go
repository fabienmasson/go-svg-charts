package charts

import (
	"fmt"
	"io"
	"math"
)

type Dimension struct {
	width, height int
}

type BezierPoint struct {
	x, y                                         float64
	beforeCtlx, beforeCtly, afterCtlx, afterCtly float64
}

func writeBackground(w io.Writer, width, height int, colorScheme *ColorScheme) {
	fmt.Fprintf(w, "<rect x='0' y='0' width='%d' height='%d' fill='%s' />", width, height, colorScheme.Background)
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
func yAxisFit(start int, end int, data [][]float64, showZero bool) ([]string, []float64, func(float64) float64) {
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
	if min > 0 && showZero {
		min = 0
	}

	if max < 0 && showZero {
		max = 0
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

	conv := func(val float64) float64 {
		return float64(start) + bottom - (bottom-top)*(val-min)/(max-min)
	}

	k := 0
	labels := make([]string, 0)
	lines := make([]float64, 0)
	for {
		val := float64(int(min/interval)+1+k) * interval
		if val < max {
			labels = append(labels, fmt.Sprintf("%.8g", val))
			lines = append(lines, val)
		} else {
			break
		}
		k++
	}

	return labels, lines, conv
}

func writeLineSeriesLegend(
	w io.Writer,
	width int,
	markerModulo int,
	series []string,
	colorScheme *ColorScheme) int {
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
			colorScheme.ColorPalette(s),
			s%markerModulo,
		)
		x += samplewidth + gap
		fmt.Fprintf(
			w,
			"<text x='%d' y='%d' alignment-baseline='middle'>%s</text>",
			x, y+2.0,
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

func writeBarSeriesLegend(
	w io.Writer,
	width int,
	series []string,
	colorScheme *ColorScheme) int {
	const samplewidth = 30
	const sampleHeight = 15
	const labelwidth = 70
	const gap = 5

	x := 10
	y := 10

	for s, serie := range series {

		fmt.Fprintf(
			w,
			"<rect x='%d' y='%d' width='%d' height='%d' fill='%s' />",
			x, y,
			samplewidth, sampleHeight,
			colorScheme.ColorPalette(s),
		)
		x += samplewidth + gap
		fmt.Fprintf(
			w,
			"<text x='%d' y='%d' alignment-baseline='middle'>%s</text>",
			x, y+sampleHeight/2+2.0,
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

func writeFontStyle(w io.Writer, isInteractive bool) {
	fmt.Fprintf(w, "<style>")
	fmt.Fprintf(w, "text { font-size: 8pt; font-family: sans-serif }  ")
	fmt.Fprintf(w, ".axislegend { font-size: 12pt; font-weight: bold } ")
	fmt.Fprintf(w, ".hovercircle {z-index:0; cursor:pointer; fill:'none'; stroke:'none'; } ")
	if isInteractive {
		fmt.Fprintf(w, ".value {z-index: 1; display:none; } ")
		fmt.Fprintf(w, ".hovercircle:hover + .value, .value:hover { display:block; }")
	} else {
		fmt.Fprintf(w, ".value {z-index: 1; } ")
	}
	fmt.Fprintf(w, "</style>")
}

func startSVG(w io.Writer, width, height int, colorScheme *ColorScheme) {
	fmt.Fprintf(
		w,
		"<svg version='1.1' xmlns='http://www.w3.org/2000/svg' viewBox='0 0 %d %d'>",
		width,
		height,
	)
}

func writeDefsTxtBg(w io.Writer, colorScheme *ColorScheme) {
	fmt.Fprintf(w, "<defs>")
	fmt.Fprintf(w, `<filter x='0' y='0' width='1' height='1' id='textbg'>
						<feFlood flood-color='%s' result='bg' />
						<feMerge>
							<feMergeNode in='bg'/>
							<feMergeNode in='SourceGraphic'/>
						</feMerge>
					</filter>`, colorScheme.Background)
	fmt.Fprintf(w, "</defs>")
}
func writeDefsMarkers(w io.Writer, size float64, n int, colorScheme *ColorScheme) int {

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
			fmt.Fprintf(w, "<circle cx='%f' cy='%f' r='%f' fill='%s' />", halfsize, halfsize, halfsize, colorScheme.ColorPalette(i))
		case 1:
			fmt.Fprintf(w, "<rect x='0' y='0' width='%f' height='10' fill='%s' />", fullsize, colorScheme.ColorPalette(i))
		case 2:
			fmt.Fprintf(w, "<polygon points='0,%f %f,0 %f,%f' fill='%s' />", fullsize, halfsize, fullsize, fullsize, colorScheme.ColorPalette(i))
		case 3:
			fmt.Fprintf(w, "<line x1='0' y1='0' x2='%f' y2='%f' stroke='%s' stroke-width='1.5'/>", fullsize, fullsize, colorScheme.ColorPalette(i))
			fmt.Fprintf(w, "<line x1='0' y1='%f' x2='%f' y2='0' stroke='%s' stroke-width='1.5'/>", fullsize, fullsize, colorScheme.ColorPalette(i))
		case 4:
			fmt.Fprintf(w, "<circle cx='%f' cy='%f' r='%f' stroke='%s' stroke-width='1.5' fill='none'/>", halfsize, halfsize, halfsize, colorScheme.ColorPalette(i))
		case 5:
			fmt.Fprintf(w, "<rect x='0' y='0' width='%f' height='%f' stroke='%s' fill='none' />", fullsize, fullsize, colorScheme.ColorPalette(i))
		case 6:
			fmt.Fprintf(w, "<polygon points='0,%f %f,0 %f,%f' stroke='%s' fill='none' />", fullsize, halfsize, fullsize, fullsize, colorScheme.ColorPalette(i))
		}

		fmt.Fprintf(w, "</marker>")
	}
	fmt.Fprintf(w, "</defs>")

	return maxMarkers

}

func endSVG(w io.Writer) {
	fmt.Fprintf(w, "</svg>")
}
