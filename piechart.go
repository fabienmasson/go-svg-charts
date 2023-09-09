/* Cf. https://marian-caikovski.medium.com/drawing-sectors-and-pie-charts-with-svg-paths-b99b5b6bf7bd */

package charts

import (
	"fmt"
	"io"
	"math"
	"sort"
)

type PieChart struct {
	Dimension
	series        []string
	data          []float64
	numberFormat  string
	colors        *ColorScheme
	showValues    bool
	isInteractive bool
}

func NewPieChart(
	width int,
	height int,
	series []string,
	data []float64,
) *PieChart {
	return &PieChart{
		Dimension: Dimension{
			width:  width,
			height: height,
		},
		colors: &ColorScheme{
			Foreground:   "black",
			Background:   "white",
			ColorPalette: DefaultPalette,
		},
		series:        series,
		data:          data,
		showValues:    false,
		isInteractive: false,
	}
}

func (pc *PieChart) SetColorDcheme(colorScheme *ColorScheme) *PieChart {
	pc.colors = colorScheme
	return pc
}

func (pc *PieChart) SetNumberFormat(numberFormat string) *PieChart {
	pc.numberFormat = numberFormat
	return pc
}

func (pc *PieChart) SetInteractive(interactive bool) *PieChart {
	pc.isInteractive = interactive
	return pc
}
func (pc *PieChart) SetShowValue(showValues bool) *PieChart {
	pc.showValues = showValues
	return pc
}

func (pc *PieChart) RenderSVG(w io.Writer) error {

	startSVG(w, pc.width, pc.height, pc.colors)
	writeFontStyle(w, pc.isInteractive)

	type pieSlice struct {
		value          float64
		label          string
		startX, startY float64
		endX, endY     float64
		large          int
		labelX, labelY float64
	}
	pieSlices := make([]pieSlice, len(pc.data))
	total := 0.0
	for i, v := range pc.data {
		pieSlices[i] = pieSlice{
			value: v,
			label: pc.series[i],
		}
		total += v
	}
	sort.SliceStable(pieSlices, func(i, j int) bool {
		return pieSlices[i].value > pieSlices[j].value
	})
	sortSeries := make([]string, len(pieSlices))
	for i, v := range pieSlices {
		sortSeries[i] = v.label
	}

	// series
	legendfHeight := writeBarSeriesLegend(w, pc.width, sortSeries, pc.colors)
	centerX := float64(pc.width / 2)
	centerY := float64(pc.height-legendfHeight)/2.0 + float64(legendfHeight)
	var radius float64
	if pc.width > (pc.height - legendfHeight) {
		radius = float64(pc.height-legendfHeight) / 2.0
	} else {
		radius = float64(pc.width) / 2.0
	}

	curAlpha := math.Pi / 2
	for i, v := range pieSlices {
		percent := v.value / total
		alpha := percent*math.Pi*2 + curAlpha
		pieSlices[i].startX = math.Cos(curAlpha) * radius
		pieSlices[i].startY = math.Sin(curAlpha) * radius
		pieSlices[i].endX = math.Cos(alpha) * radius
		pieSlices[i].endY = math.Sin(alpha) * radius
		pieSlices[i].labelX = math.Cos((alpha+curAlpha)/2.0) * radius * 0.8
		pieSlices[i].labelY = math.Sin((alpha+curAlpha)/2.0) * radius * 0.8
		if alpha-curAlpha > math.Pi {
			pieSlices[i].large = 1
		}
		curAlpha = alpha
	}

	// pie
	for i, _ := range pieSlices {
		fmt.Fprintf(
			w,
			"<path d='M %f %f A %f %f 0 %d 1 %f %f L %f %f L %f %f Z' fill='%s' stroke='%s' />",
			centerX-pieSlices[i].startX,
			centerY-pieSlices[i].startY,
			radius, radius,
			pieSlices[i].large,
			centerX-pieSlices[i].endX,
			centerY-pieSlices[i].endY,
			centerX,
			centerY,
			centerX-pieSlices[i].startX,
			centerY-pieSlices[i].startY,
			pc.colors.ColorPalette[i%len(pc.colors.ColorPalette)],
			pc.colors.Background,
		)
	}
	if pc.showValues {
		for i, _ := range pieSlices {
			fmt.Fprintf(
				w,
				"<text x='%f' y='%f' text-anchor='middle' alignment-baseline='middle' fill='#fff'>%g</text>",
				centerX-pieSlices[i].labelX,
				centerY-pieSlices[i].labelY,
				pieSlices[i].value,
			)
		}
	}

	endSVG(w)

	return nil
}
