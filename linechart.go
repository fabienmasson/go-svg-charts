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
	colorScheme     *ColorScheme
	xaxisLegend     string
	yaxisLegend     string
	showMarkers     bool
	showValues      bool
	isInteractive   bool
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
		colorScheme:     &DefaultColorScheme,
		horizontalLines: 8,
		xaxis:           xaxis,
		series:          series,
		data:            data,
	}
}

func (l *LineChart) SetColorDcheme(colorScheme *ColorScheme) *LineChart {
	l.colorScheme = colorScheme
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
	fmt.Println(l.colorScheme.Background)
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

func (l *LineChart) SetInteractive(interactive bool) *LineChart {
	l.isInteractive = interactive
	return l
}
func (l *LineChart) SetShowValue(showValues bool) *LineChart {
	l.showValues = showValues
	return l
}

func (l *LineChart) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const rightMargin = 20
	const textHeight = 15

	startSVG(w, l.width, l.height, l.colorScheme)
	writeFontStyle(w, l.isInteractive)
	markerModulo := 7
	if l.showMarkers {
		markerModulo = writeDefsMarkers(w, 8.0, len(l.series), l.colorScheme)
	}
	headerHeight := writeLineSeriesLegend(w, l.width, markerModulo, l.series, l.colorScheme)

	// horizontal lines and labels
	labels, hlines, convy := yAxisFit(headerHeight, l.height-xaxisHeight-gap, l.data, false)

	for i, hline := range hlines {
		fmt.Fprintf(
			w,
			"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
			yaxisWidth,
			l.width-rightMargin,
			convy(hline),
			convy(hline),
			l.colorScheme.LightAxisColor,
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
			l.colorScheme.LightAxisColor,
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
		l.colorScheme.DarkerAxisColor,
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
		l.colorScheme.DarkerAxisColor,
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
			l.colorScheme.ColorPalette(s),
			s%markerModulo, s%markerModulo, s%markerModulo,
		)

	}

	for _, serie := range l.data {

		for i := 0; i < len(serie); i++ {
			if l.isInteractive {

				fmt.Fprintf(
					w,
					"<circle class='hovercircle' cx='%f' cy='%f' r='15' fill='#fff' fill-opacity='0' />",
					float64(yaxisWidth+gap)+dw*float64(i),
					convy(serie[i]),
				)
			}
			if l.isInteractive || l.showValues {
				fmt.Fprintf(
					w,
					"<text style='paint-order:stroke fill' class='value' x='%f' y='%f' text-anchor='middle' alignment-baseline='middle' stroke='#fff' stroke-width='10' fill='#555'>%g</text>",
					float64(yaxisWidth+gap)+dw*float64(i),
					convy(serie[i])-10.0,
					serie[i],
				)
			}
		}
	}

	endSVG(w)

	return nil
}
