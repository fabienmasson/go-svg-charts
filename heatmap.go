package charts

import (
	"fmt"
	"io"
)

type HeatMap struct {
	Dimension
	xaxis         []string
	yaxis         []string
	data          [][]float64
	numberFormat  string
	colorScheme   *ColorScheme
	xaxisLegend   string
	yaxisLegend   string
	showValues    bool
	isInteractive bool
}

func NewHeatMap(
	width int,
	height int,
	xaxis []string,
	yaxis []string,
	data [][]float64,
) *HeatMap {
	return &HeatMap{
		Dimension: Dimension{
			width:  width,
			height: height,
		},
		colorScheme:   &DefaultColorScheme,
		xaxis:         xaxis,
		yaxis:         yaxis,
		data:          data,
		showValues:    false,
		isInteractive: false,
	}
}

func (hm *HeatMap) SetColorDcheme(colorScheme *ColorScheme) *HeatMap {
	hm.colorScheme = colorScheme
	return hm
}

func (hm *HeatMap) SetXaxisLegend(xaxisLegend string) *HeatMap {
	hm.xaxisLegend = xaxisLegend
	return hm
}

func (hm *HeatMap) SetYaxisLegend(yaxisLegend string) *HeatMap {
	hm.yaxisLegend = yaxisLegend
	return hm
}

func (hm *HeatMap) SetNumberFormat(numberFormat string) *HeatMap {
	hm.numberFormat = numberFormat
	fmt.Println(hm.colorScheme.Background)
	return hm
}

func (hm *HeatMap) SetInteractive(interactive bool) *HeatMap {
	hm.isInteractive = interactive
	return hm
}
func (hm *HeatMap) SetShowValue(showValues bool) *HeatMap {
	hm.showValues = showValues
	return hm
}

func (hm *HeatMap) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const textHeight = 15
	const barGap = 20

	startSVG(w, hm.width, hm.height, hm.colorScheme)
	writeFontStyle(w, hm.isInteractive)

	chartHeight := float64(hm.height) - float64(xaxisHeight) - 2.0*gap
	chartWidth := float64(hm.width) - float64(yaxisWidth) - gap*2.0
	dh := chartHeight / float64(len(hm.yaxis))
	dw := chartWidth / float64(len(hm.xaxis))
	convytop := func(i int) float64 {
		return float64(hm.height) - float64(xaxisHeight) - gap - dh*float64(i+1)
	}
	convybottom := func(i int) float64 {
		return float64(hm.height) - float64(xaxisHeight) - gap - dh*float64(i)
	}
	convymiddle := func(i int) float64 {
		return (convytop(i) + convybottom(i)) / 2.0
	}
	convxstart := func(i int) float64 {
		return float64(yaxisWidth) + gap + dw*float64(i)
	}
	convxend := func(i int) float64 {
		return float64(yaxisWidth) + gap + dw*float64(i+1)
	}
	convxmiddle := func(i int) float64 {
		return (convxstart(i) + convxend(i)) / 2.0
	}

	min, max := hm.data[0][0], hm.data[0][0]
	for _, row := range hm.data {
		for _, value := range row {
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
	}

	// horizontal labels
	for i, label := range hm.yaxis {
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f'>%s</text>",
			float64(gap)+textHeight,
			convymiddle(i),
			label,
		)
	}

	// vertical lines
	for i, label := range hm.xaxis {
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			convxmiddle(i),
			float64(hm.height-xaxisHeight+gap),
			label,
		)
	}

	// xaxis
	fmt.Fprintf(
		w,
		"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
		yaxisWidth,
		hm.width,
		float64(hm.height-xaxisHeight-gap),
		float64(hm.height-xaxisHeight-gap),
		hm.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' class='axislegend' dominant-baseline='middle' text-anchor='middle'>%s</text>",
		float64(yaxisWidth)+float64(hm.width-yaxisWidth)/2.0,
		float64(hm.height-xaxisHeight+gap+textHeight),
		hm.xaxisLegend,
	)

	// yaxis
	fmt.Fprintf(
		w,
		"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
		float64(yaxisWidth+gap),
		float64(yaxisWidth+gap),
		0,
		hm.height-xaxisHeight,
		hm.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' transform='rotate(270, %f, %f)' class='axislegend' text-anchor='middle' alignment-baseline='middle'>%s</text>",
		float64(textHeight),
		float64(hm.height)/2,
		float64(textHeight),
		float64(hm.height)/2,
		hm.yaxisLegend,
	)

	// series

	for i := 0; i < len(hm.data); i++ {
		for j := 0; j < len(hm.data[0]); j++ {
			fmt.Fprintf(
				w,
				"<rect x='%f' y='%f' width='%f' height='%f' fill='%s' opacity='%f' stroke='%s' stroke-width='1'/>",
				convxstart(i),
				convytop(j),
				dw,
				dh,
				hm.colorScheme.ColorPalette(0),
				(hm.data[i][j]-min)/(max-min),
				hm.colorScheme.Background,
			)
		}
	}

	for i := 0; i < len(hm.data); i++ {
		for j := 0; j < len(hm.data[0]); j++ {
			if hm.isInteractive {
				fmt.Fprintf(
					w,
					"<rect class='hovercircle' x='%f' y='%f' width='%f' height='%f' fill-opacity='0'/>",
					convxstart(i),
					convytop(j),
					dw,
					dh,
				)
			}
			if hm.showValues || hm.isInteractive {
				fmt.Fprintf(
					w,
					"<text style='paint-order:stroke fill' class='value' x='%f' y='%f' text-anchor='middle' alignment-baseline='middle' stroke='#fff' stroke-width='10' fill='#555'>%g</text>",
					convxmiddle(i),
					convymiddle(j),
					hm.data[i][j],
				)
			}

		}
	}

	endSVG(w)

	return nil
}
