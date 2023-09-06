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
	fmt.Println(l.colors.Background)
	return l
}

func (l *LineChart) SetHorizontalLines(horizontalLines int) *LineChart {
	l.horizontalLines = horizontalLines
	return l
}

func (l *LineChart) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const rightMargin = 20
	const lightAxisColor = "#eee"
	const darkerAxisColor = "#777"

	startSVG(w, l.width, l.height, l.colors)
	writeFontStyle(w)
	markerModulo := writeDefsMarkers(w, 8.0, len(l.series), l.colors)
	headerHeight := seriesLegend(w, l.width, markerModulo, l.series, l.colors)

	// horizontal lines and labels
	labels, hlines, convy := yAxisFit(headerHeight, l.height-xaxisHeight-gap, l.data)

	for i, hline := range hlines {
		fmt.Fprintf(
			w,
			"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
			yaxisWidth,
			l.width-rightMargin,
			convy(hline),
			convy(hline),
			lightAxisColor,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' font-size='8pt'>%s</text>",
			float64(gap),
			convy(hline),
			labels[i],
		)
	}

	// vertical lines
	dw := float64(l.width-yaxisWidth-gap*2-rightMargin) / float64(len(l.xaxis)-1)
	for i := 0; i < len(l.xaxis); i++ {

		fmt.Fprintf(
			w,
			"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
			float64(yaxisWidth+gap)+dw*float64(i),
			float64(yaxisWidth+gap)+dw*float64(i),
			headerHeight,
			l.height-xaxisHeight,
			lightAxisColor,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' font-size='8pt' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			float64(yaxisWidth+gap)+dw*float64(i),
			float64(l.height-xaxisHeight/2),
			l.xaxis[i],
		)
	}

	// xaxis
	fmt.Fprintf(
		w,
		"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
		yaxisWidth,
		l.width,
		float64(l.height-xaxisHeight-gap),
		float64(l.height-xaxisHeight-gap),
		darkerAxisColor,
	)

	// yaxis
	fmt.Fprintf(
		w,
		"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
		float64(yaxisWidth+gap),
		float64(yaxisWidth+gap),
		headerHeight,
		l.height-xaxisHeight,
		darkerAxisColor,
	)

	// series

	for s, serie := range l.data {
		points := ""
		for i := 0; i < len(serie); i++ {
			points += fmt.Sprintf(
				"%f,%f ",
				float64(yaxisWidth+gap)+dw*float64(i),
				convy(serie[i]),
			)
		}
		fmt.Fprintf(
			w,
			"<polyline points='%s' fill='none' stroke='%s' stroke-width='2' marker-start='url(#dot%d)' marker-mid='url(#dot%d)'  marker-end='url(#dot%d)'/>",
			points,
			l.colors.ColorPalette[s%len(l.colors.ColorPalette)],
			s%markerModulo, s%markerModulo, s%markerModulo,
		)
	}

	endSVG(w)

	return nil
}
