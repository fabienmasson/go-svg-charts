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
	return l
}

func (l *LineChart) RenderSVG(w io.Writer) error {

	headerHeight := 70
	xaxisHeight := 50
	yaxisWidth := 50
	gap := 10

	startTag(w, "svg",
		map[string]string{
			"version": "1.1",
			"xmlns":   "http://www.w3.org/2000/svg",
			"width":   fmt.Sprintf("%d", l.width),
			"height":  fmt.Sprintf("%d", l.height),
		},
	)

	// background
	fmt.Fprintf(w, "<rect width=\"100%%\" height=\"100%%\" fill=\"%s\" />", l.colors.Background)

	var color string

	//dh := float64(l.height-headerHeight-xaxisHeight-gap*2) / float64(l.horizontalLines-1)
	// horizontal lines
	/*
		for i := 0; i < l.horizontalLines-1; i++ {
			color = "#eee"
			fmt.Fprintf(
				w,
				"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
				yaxisWidth,
				l.width,
				float64(headerHeight+gap)+dh*float64(i),
				float64(headerHeight+gap)+dh*float64(i),
				color,
			)
		}
	*/

	// vertical lines
	dw := float64(l.width-yaxisWidth-gap*2) / float64(len(l.xaxis)-1)
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

	// min-max
	min, max := l.data[0][0], l.data[0][0]
	for i := range l.data {
		for _, v := range l.data[i] {
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}
	}

	endTag(w, "svg")
	return nil
}
