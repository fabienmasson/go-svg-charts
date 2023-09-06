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
	xaxisLegend     string
	yaxisLegend     string
	showMarkers     bool
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

func (l *LineChart) SetXaxisLegend(xaxisLegend string) *LineChart {
	l.xaxisLegend = xaxisLegend
	return l
}

func (l *LineChart) SetYaxisLegend(yaxisLegend string) *LineChart {
	l.yaxisLegend = yaxisLegend
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

func (l *LineChart) SetShowMarkers(showMarkers bool) *LineChart {
	l.showMarkers = showMarkers
	return l
}

func (l *LineChart) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const rightMargin = 20
	const textHeight = 15

	startSVG(w, l.width, l.height, l.colors)
	writeFontStyle(w)
	markerModulo := 7
	if l.showMarkers {
		markerModulo = writeDefsMarkers(w, 8.0, len(l.series), l.colors)
	}
	headerHeight := writeLineSeriesLegend(w, l.width, markerModulo, l.series, l.colors)

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
			"<text x='%f' y='%f'>%s</text>",
			float64(gap)+textHeight,
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
			"<text x='%f' y='%f' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			float64(yaxisWidth+gap)+dw*float64(i),
			float64(l.height-xaxisHeight+gap),
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
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' class='axislegend' dominant-baseline='middle' text-anchor='middle'>%s</text>",
		float64(yaxisWidth+(l.width-yaxisWidth-rightMargin)/2),
		float64(l.height-xaxisHeight+gap+textHeight),
		l.xaxisLegend,
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
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' transform='rotate(270, %f, %f)' class='axislegend' text-anchor='middle' alignment-baseline='middle'>%s</text>",
		float64(textHeight),
		float64(l.height)/2,
		float64(textHeight),
		float64(l.height)/2,
		l.yaxisLegend,
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
