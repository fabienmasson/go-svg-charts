package charts

import (
	"fmt"
	"io"
)

type LineChart struct {
	Dimension
	xaxis           []string
	series          []string
	data            [][]float64
	horizontalLines int
	numberFormat    string
	colors          *ColorScheme
}

func NewLineChart(
	width int,
	height int,
	xaxis []string,
	series []string,
	data [][]float64,
) *LineChart {
	return &LineChart{
		Dimension: Dimension{
			width:  width,
			height: height,
		},
		colors: &ColorScheme{
			Foreground:   "black",
			Background:   "white",
			ColorPalette: DefaultPalette,
		},
		horizontalLines: 8,
		xaxis:           xaxis,
		series:          series,
		data:            data,
	}
}

func (l *LineChart) SetColorDcheme(colorScheme *ColorScheme) *LineChart {
	l.colors = colorScheme
	return l
}

func (l *LineChart) SetNumberFormat(numberFormat string) *LineChart {
	l.numberFormat = numberFormat
	return l
}

func (l *LineChart) SetHorizontalLines(horizontalLines int) *LineChart {
	l.horizontalLines = horizontalLines
	fmt.Println("coucou")
	return l
}

func (l *LineChart) RenderSVG(w io.Writer) error {

	headerHeight := 70
	xaxisHeight := 50
	yaxisWidth := 50
	gap := 10
	rightMargin := 20

	// start svg
	fmt.Fprintf(
		w,
		"<svg varsion='1.1' xmlns='http://www.w3.org/2000/svg' width='%d' height='%d'>",
		l.width,
		l.height,
	)

	// background
	fmt.Fprintf(
		w,
		"<rect width=\"100%%\" height=\"100%%\" fill=\"%s\" />",
		l.colors.Background,
	)

	var color string

	// horizontal lines and labels
	labels, hlines, convy := yAxisFit(headerHeight, l.height-xaxisHeight-gap, l.data)
	color = "#eee"
	for i, hline := range hlines {
		fmt.Fprintf(
			w,
			"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
			yaxisWidth,
			l.width-rightMargin,
			convy(hline),
			convy(hline),
			color,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f'>%s</text>",
			float64(gap),
			convy(hline),
			labels[i],
		)
	}

	// vertical lines
	dw := float64(l.width-yaxisWidth-gap*2-rightMargin) / float64(len(l.xaxis)-1)
	for i := 0; i < len(l.xaxis); i++ {
		if i == 0 {
			color = "#777"
		} else {
			color = "#eee"
		}
		fmt.Fprintf(
			w,
			"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
			float64(yaxisWidth+gap)+dw*float64(i),
			float64(yaxisWidth+gap)+dw*float64(i),
			headerHeight,
			l.height-xaxisHeight,
			color,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			float64(yaxisWidth+gap)+dw*float64(i),
			float64(l.height-xaxisHeight/2),
			l.xaxis[i],
		)
	}

	// xaxis
	color = "#777"
	fmt.Fprintf(
		w,
		"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
		yaxisWidth,
		l.width,
		float64(l.height-xaxisHeight-gap),
		float64(l.height-xaxisHeight-gap),
		color,
	)

	// series
	for s, serie := range l.data {
		for i := 0; i < len(serie); i++ {
			if i < len(serie)-1 {
				fmt.Fprintf(
					w,
					"<line x1='%f' x2='%f' y1='%f' y2='%f' stroke='%s' stroke-width='2'/>",
					float64(yaxisWidth+gap)+dw*float64(i),
					float64(yaxisWidth+gap)+dw*float64(i+1),
					convy(serie[i]),
					convy(serie[i+1]),
					l.colors.ColorPalette[s%len(l.colors.ColorPalette)],
				)
			}
			fmt.Fprintf(
				w,
				"<circle cx='%f' cy='%f' r='3' fill='%s' />",
				float64(yaxisWidth+gap)+dw*float64(i),
				convy(serie[i]),
				l.colors.ColorPalette[s%len(l.colors.ColorPalette)],
			)
		}
	}

	seriesLegend(w, headerHeight, l.width, l.series, l.colors)

	// start svg
	fmt.Fprintf(w, "</svg>")

	return nil
}
