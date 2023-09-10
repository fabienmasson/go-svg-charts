package charts

import (
	"fmt"
	"io"
)

type AeraChart struct {
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
	isBezier        bool
	datasum         [][]float64
}

func NewAreaChart(
	width int,
	height int,
	xaxis []string,
	series []string,
	data [][]float64,
) *AeraChart {
	ac := &AeraChart{
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
	ac.datasum = make([][]float64, 0)
	for i := 0; i < len(series); i++ {
		ac.datasum = append(ac.datasum, make([]float64, len(data[0])))
		for j := 0; j < len(data[0]); j++ {
			if i == 0 {
				ac.datasum[i][j] = ac.data[i][j]
			} else {
				ac.datasum[i][j] = ac.datasum[i-1][j] + ac.data[i][j]
			}
		}
	}
	return ac
}

func (ac *AeraChart) SetColorDcheme(colorScheme *ColorScheme) *AeraChart {
	ac.colorScheme = colorScheme
	return ac
}

func (ac *AeraChart) SetXaxisLegend(xaxisLegend string) *AeraChart {
	ac.xaxisLegend = xaxisLegend
	return ac
}

func (ac *AeraChart) SetYaxisLegend(yaxisLegend string) *AeraChart {
	ac.yaxisLegend = yaxisLegend
	return ac
}

func (ac *AeraChart) SetNumberFormat(numberFormat string) *AeraChart {
	ac.numberFormat = numberFormat
	fmt.Println(ac.colorScheme.Background)
	return ac
}

func (ac *AeraChart) SetHorizontalLines(horizontalLines int) *AeraChart {
	ac.horizontalLines = horizontalLines
	return ac
}

func (ac *AeraChart) SetShowMarkers(showMarkers bool) *AeraChart {
	ac.showMarkers = showMarkers
	return ac
}

func (ac *AeraChart) SetInteractive(interactive bool) *AeraChart {
	ac.isInteractive = interactive
	return ac
}
func (ac *AeraChart) SetShowValue(showValues bool) *AeraChart {
	ac.showValues = showValues
	return ac
}
func (ac *AeraChart) SetBezier(isBezier bool) *AeraChart {
	ac.isBezier = isBezier
	return ac
}

func (ac *AeraChart) RenderSVG(w io.Writer) error {

	const xaxisHeight = 50
	const yaxisWidth = 50
	const gap = 10
	const rightMargin = 20
	const textHeight = 15

	startSVG(w, ac.width, ac.height, ac.colorScheme)
	writeDefsTxtBg(w, ac.colorScheme)
	writeFontStyle(w, ac.isInteractive)
	markerModulo := 7
	if ac.showMarkers {
		markerModulo = writeDefsMarkers(w, 8.0, len(ac.series), ac.colorScheme)
	}
	headerHeight := writeLineSeriesLegend(w, ac.width, markerModulo, ac.series, ac.colorScheme)

	// horizontal lines and labels
	labels, hlines, convy := yAxisFit(headerHeight, ac.height-xaxisHeight-gap, ac.datasum, false)

	for i, hline := range hlines {
		fmt.Fprintf(
			w,
			"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
			yaxisWidth,
			ac.width-rightMargin,
			convy(hline),
			convy(hline),
			ac.colorScheme.LightAxisColor,
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
	dw := float64(ac.width-yaxisWidth-gap*2-rightMargin) / float64(len(ac.xaxis)-1)
	convx := func(x float64) float64 {
		return float64(yaxisWidth+gap) + dw*x
	}
	for i := 0; i < len(ac.xaxis); i++ {

		fmt.Fprintf(
			w,
			"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
			convx(float64(i)),
			convx(float64(i)),
			headerHeight,
			ac.height-xaxisHeight,
			ac.colorScheme.LightAxisColor,
		)
		fmt.Fprintf(
			w,
			"<text x='%f' y='%f' dominant-baseline='middle' text-anchor='middle'>%s</text>",
			convx(float64(i)),
			float64(ac.height-xaxisHeight+gap),
			ac.xaxis[i],
		)
	}

	// xaxis
	fmt.Fprintf(
		w,
		"<line x1='%d' x2='%d' y1='%f' y2='%f' stroke='%s' stroke-width='1'/>",
		yaxisWidth,
		ac.width,
		float64(ac.height-xaxisHeight-gap),
		float64(ac.height-xaxisHeight-gap),
		ac.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' class='axislegend' dominant-baseline='middle' text-anchor='middle'>%s</text>",
		float64(yaxisWidth+(ac.width-yaxisWidth-rightMargin)/2),
		float64(ac.height-xaxisHeight+gap+textHeight),
		ac.xaxisLegend,
	)

	// yaxis
	fmt.Fprintf(
		w,
		"<line x1='%f' x2='%f' y1='%d' y2='%d' stroke='%s' stroke-width='1'/>",
		float64(yaxisWidth+gap),
		float64(yaxisWidth+gap),
		headerHeight,
		ac.height-xaxisHeight,
		ac.colorScheme.DarkerAxisColor,
	)
	fmt.Fprintf(
		w,
		"<text x='%f' y='%f' transform='rotate(270, %f, %f)' class='axislegend' text-anchor='middle' alignment-baseline='middle'>%s</text>",
		float64(textHeight),
		float64(ac.height)/2,
		float64(textHeight),
		float64(ac.height)/2,
		ac.yaxisLegend,
	)

	// series
	if ac.isBezier {

		allBesierPoints := make([][]*BezierPoint, 0)

		for _, serie := range ac.datasum {
			bezierPoints := make([]*BezierPoint, 0)
			for i := 0; i < len(serie); i++ {
				bezierPoint := BezierPoint{
					x:          convx(float64(i)),
					y:          convy(serie[i]),
					beforeCtlx: convx(float64(i) - 0.25),
					afterCtlx:  convx(float64(i) + 0.25),
				}
				bezierPoints = append(bezierPoints, &bezierPoint)
			}
			for i := 1; i < len(serie)-1; i++ {
				bezierPoints[i].beforeCtly = convy(serie[i] - (serie[i+1]-serie[i-1])/8.0)
				bezierPoints[i].afterCtly = convy(serie[i] + (serie[i+1]-serie[i-1])/8.0)
			}
			bezierPoints[0].afterCtly = bezierPoints[0].y
			bezierPoints[len(serie)-1].beforeCtly = bezierPoints[len(serie)-1].y
			bezierPoints[len(serie)-1].afterCtly = bezierPoints[len(serie)-1].y
			allBesierPoints = append(allBesierPoints, bezierPoints)
		}

		for s, serie := range ac.datasum {

			// area
			points := ""
			points += fmt.Sprintf(
				"M%f %f C %f %f,",
				allBesierPoints[s][0].x,
				allBesierPoints[s][0].y,
				allBesierPoints[s][0].afterCtlx,
				allBesierPoints[s][0].afterCtly,
			)
			for i := 1; i < len(allBesierPoints[s]); i++ {
				// start point
				points += fmt.Sprintf(
					" %f %f, %f %f ",
					allBesierPoints[s][i].beforeCtlx,
					allBesierPoints[s][i].beforeCtly,
					allBesierPoints[s][i].x,
					allBesierPoints[s][i].y,
				)
				// start control point
				if i < len(allBesierPoints[s])-1 {
					points += fmt.Sprintf("S")
				}
			}

			if s == 0 {
				points += fmt.Sprintf(
					"C %f %f, %f %f, %f %f",
					convx(float64(len(serie)-1)),
					float64(ac.height-xaxisHeight-gap),
					convx(float64(len(serie)-1)),
					allBesierPoints[s][len(serie)-1].y,
					convx(float64(len(serie)-1)),
					float64(ac.height-xaxisHeight-gap),
				)
				points += fmt.Sprintf(
					"C %f %f, %f %f, %f %f",
					convx(float64(0)),
					float64(ac.height-xaxisHeight-gap),
					convx(float64(len(serie)-1)),
					float64(ac.height-xaxisHeight-gap),
					convx(float64(0)),
					float64(ac.height-xaxisHeight-gap),
				)
			} else {
				points += fmt.Sprintf(
					"C %f %f, %f %f, %f %f ",
					convx(float64(len(serie)-1)),
					allBesierPoints[s-1][len(serie)-1].y,
					convx(float64(len(serie)-1)),
					allBesierPoints[s][len(serie)-1].y,
					convx(float64(len(serie)-1)),
					allBesierPoints[s-1][len(serie)-1].y,
				)
				points += fmt.Sprintf(
					"C %f %f,",
					allBesierPoints[s-1][len(serie)-1].beforeCtlx,
					allBesierPoints[s-1][len(serie)-1].beforeCtly,
				)
				for i := len(allBesierPoints[s]) - 1; i >= 0; i-- {
					// start point
					points += fmt.Sprintf(
						" %f %f, %f %f ",
						allBesierPoints[s-1][i].afterCtlx,
						allBesierPoints[s-1][i].afterCtly,
						allBesierPoints[s-1][i].x,
						allBesierPoints[s-1][i].y,
					)
					// start control point
					if i > 0 {
						points += fmt.Sprintf("S")
					}
				}

			}

			fmt.Fprintf(
				w,
				"<path d='%s' fill='%s' fill-opacity='0.5' stroke='none' stroke-width='2' />",
				points,
				ac.colorScheme.ColorPalette(s),
			)

			// plot
			points = ""
			points += fmt.Sprintf(
				"M%f %f C %f %f,",
				allBesierPoints[s][0].x,
				allBesierPoints[s][0].y,
				allBesierPoints[s][0].afterCtlx,
				allBesierPoints[s][0].afterCtly,
			)
			for i := 1; i < len(allBesierPoints[s]); i++ {
				// start point
				points += fmt.Sprintf(
					" %f %f, %f %f ",
					allBesierPoints[s][i].beforeCtlx,
					allBesierPoints[s][i].beforeCtly,
					allBesierPoints[s][i].x,
					allBesierPoints[s][i].y,
				)
				// start control point
				if i < len(allBesierPoints[s])-1 {
					points += fmt.Sprintf("S")
				}
			}

			fmt.Fprintf(
				w,
				"<path d='%s' fill='none' stroke='%s' stroke-width='2' marker-start='url(#dot%d)' marker-mid='url(#dot%d)'  marker-end='url(#dot%d)'/>",
				points,
				ac.colorScheme.ColorPalette(s),
				s%markerModulo, s%markerModulo, s%markerModulo,
			)

		}
	} else {
		for s, serie := range ac.datasum {

			// areas
			points := ""
			for i := 0; i < len(serie); i++ {
				points += fmt.Sprintf(
					"%f,%f ",
					convx(float64(i)),
					convy(serie[i]),
				)
			}
			if s == 0 {
				points += fmt.Sprintf(
					"%f,%f ",
					convx(float64(len(serie)-1)),
					float64(ac.height-xaxisHeight-gap),
				)
				points += fmt.Sprintf(
					"%f,%f ",
					convx(float64(0)),
					float64(ac.height-xaxisHeight-gap),
				)
			} else {
				for i := len(serie) - 1; i >= 0; i-- {
					points += fmt.Sprintf(
						"%f,%f ",
						convx(float64(i)),
						convy(ac.datasum[s-1][i]),
					)
				}
			}
			fmt.Fprintf(
				w,
				"<polyline points='%s' fill='%s' fill-opacity='0.5' stroke='none' stroke-width='2'/>",
				points,
				ac.colorScheme.ColorPalette(s),
			)

			// plot
			points = ""
			for i := 0; i < len(serie); i++ {
				points += fmt.Sprintf(
					"%f,%f ",
					convx(float64(i)),
					convy(serie[i]),
				)
			}

			fmt.Fprintf(
				w,
				"<polyline points='%s' fill='none' stroke='%s' stroke-width='2' marker-start='url(#dot%d)' marker-mid='url(#dot%d)'  marker-end='url(#dot%d)'/>",
				points,
				ac.colorScheme.ColorPalette(s),
				s%markerModulo, s%markerModulo, s%markerModulo,
			)

		}
	}

	for _, serie := range ac.datasum {

		for i := 0; i < len(serie); i++ {
			if ac.isInteractive {

				fmt.Fprintf(
					w,
					"<circle class='hovercircle' cx='%f' cy='%f' r='15' fill='#fff' fill-opacity='0' />",
					convx(float64(i)),
					convy(serie[i]),
				)
			}
			if ac.isInteractive || ac.showValues {
				fmt.Fprintf(
					w,
					"<text style='paint-order:stroke fill' class='value' x='%f' y='%f' text-anchor='middle' alignment-baseline='middle' filter='url(#textbg)'>%g</text>",
					convx(float64(i)),
					convy(serie[i])-10.0,
					serie[i],
				)
			}
		}
	}

	endSVG(w)

	return nil
}
