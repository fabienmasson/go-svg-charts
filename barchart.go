package charts

import (
	"fmt"
	"io"
)

type BarChart struct {
	Dimension
	xaxis           []string
	series          []string
	data            [][]float64
	horizontalLines int
	numberFormat    string
	colorScheme     *ColorScheme
	xaxisLegend     string
	yaxisLegend     string
	showZero        bool
	showValues      bool
	isInteractive   bool
}

func NewBarChart(
	width int,
	height int,
	xaxis []string,
	series []string,
	data [][]float64,
) *BarChart {
	return &BarChart{
		Dimension: Dimension{
			width:  width,
			height: height,
		},
		colorScheme:     &DefaultColorScheme,
		horizontalLines: 8,
		xaxis:           xaxis,
		series:          series,
		data:            data,
		showZero:        true,
		showValues:      false,
		isInteractive:   false,
	}
}

func (bc *BarChart) SetColorDcheme(colorScheme *ColorScheme) *BarChart {
	bc.colorScheme = colorScheme
	return bc
}

func (bc *BarChart) SetXaxisLegend(xaxisLegend string) *BarChart {
	bc.xaxisLegend = xaxisLegend
	return bc
}

func (bc *BarChart) SetYaxisLegend(yaxisLegend string) *BarChart {
	bc.yaxisLegend = yaxisLegend
	return bc
}

func (bc *BarChart) SetNumberFormat(numberFormat string) *BarChart {
	bc.numberFormat = numberFormat
	fmt.Println(bc.colorScheme.Background)
	return bc
}

func (bc *BarChart) SetHorizontalLines(horizontalLines int) *BarChart {
	bc.horizontalLines = horizontalLines
	return bc
}

func (bc *BarChart) SetShowZero(showZero bool) *BarChart {
	bc.showZero = showZero
	return bc
}

func (bc *BarChart) SetInteractive(interactive bool) *BarChart {
	bc.isInteractive = interactive
	return bc
}
func (bc *BarChart) SetShowValue(showValues bool) *BarChart {
	bc.showValues = showValues
	return bc
}

func (bc *BarChart) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const rightMargin = 20
	const textHeight = 15
	const barGap = 20

	startSVG(w, bc.width, bc.height, bc.colorScheme)
	writeFontStyle(w, bc.isInteractive)
	headerHeight := writeBarSeriesLegend(w, bc.width, bc.series, bc.colorScheme)

	// horizontal lines and labels
	labels, hlines, convy := yAxisFit(headerHeight, bc.height-xaxisHeight-gap, bc.data, bc.showZero)

	for i, hline := range hlines {
		fmt.Fprintf(
			w,
			"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
			yaxisWidth,
			bc.width-rightMargin,
			convy(hline),
			convy(hline),
			bc.colorScheme.LightAxisColor,
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
	dw := float64(bc.width-yaxisWidth-gap*2-rightMargin) / float64(len(bc.xaxis))
	for i := 0; i < len(bc.xaxis); i++ {
		fmt.Fprintf(
			w,
			"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
			float64(yaxisWidth+gap)+dw/2.0+dw*float64(i),
			float64(yaxisWidth+gap)+dw/2.0+dw*float64(i),
			headerHeight,
			bc.height-xaxisHeight,
			bc.colorScheme.LightAxisColor,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			float64(yaxisWidth+gap)+dw/2.0+dw*float64(i),
			float64(bc.height-xaxisHeight+gap),
			bc.xaxis[i],
		)
	}

	// xaxis
	fmt.Fprintf(
		w,
		"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
		yaxisWidth,
		bc.width,
		float64(bc.height-xaxisHeight-gap),
		float64(bc.height-xaxisHeight-gap),
		bc.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' class='axislegend' dominant-baseline='middle' text-anchor='middle'>%s</text>",
		float64(yaxisWidth+(bc.width-yaxisWidth-rightMargin)/2),
		float64(bc.height-xaxisHeight+gap+textHeight),
		bc.xaxisLegend,
	)

	// yaxis
	fmt.Fprintf(
		w,
		"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
		float64(yaxisWidth+gap),
		float64(yaxisWidth+gap),
		headerHeight,
		bc.height-xaxisHeight,
		bc.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' transform='rotate(270, %f, %f)' class='axislegend' text-anchor='middle' alignment-baseline='middle'>%s</text>",
		float64(textHeight),
		float64(bc.height)/2,
		float64(textHeight),
		float64(bc.height)/2,
		bc.yaxisLegend,
	)

	// series
	bw := (dw - barGap) / float64(len(bc.series))
	relativeStart := (dw - barGap) / 2
	for s, serie := range bc.data {
		for i := 0; i < len(serie); i++ {
			fmt.Fprintf(
				w,
				"<rect x='%f' y='%f' fill='%s' width='%f' height='%f'/>",
				float64(yaxisWidth+gap)+dw/2.0+dw*float64(i)-relativeStart+bw*float64(s),
				convy(serie[i]),
				bc.colorScheme.ColorPalette(s),
				bw,
				(float64(bc.height)-xaxisHeight-gap)-convy(serie[i]),
			)
		}
	}

	for s, serie := range bc.data {
		for i := 0; i < len(serie); i++ {
			if bc.isInteractive {
				fmt.Fprintf(
					w,
					"<rect class='hovercircle' x='%f' y='%f' width='%f' height='%f' fill-opacity='0' />",
					float64(yaxisWidth+gap)+dw/2.0+dw*float64(i)-relativeStart+bw*float64(s),
					convy(serie[i]),
					bw,
					(float64(bc.height)-xaxisHeight-gap)-convy(serie[i]),
				)
			}
			if bc.showValues || bc.isInteractive {
				fmt.Fprintf(
					w,
					"<text style='paint-order:stroke fill' class='value' x='%f' y='%f' text-anchor='middle' alignment-baseline='middle' stroke='#fff' stroke-width='10' fill='#555'>%g</text>",
					float64(yaxisWidth+gap)+dw/2.0+dw*float64(i)-relativeStart+bw*float64(s)+10,
					convy(serie[i])-10.0,
					serie[i],
				)
			}
		}
	}

	endSVG(w)

	return nil
}
