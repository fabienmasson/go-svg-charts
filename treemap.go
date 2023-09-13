/* Cf. https://marian-caikovski.medium.com/drawing-sectors-and-pie-charts-with-svg-paths-b99b5b6bf7bd */

package charts

import (
	"fmt"
	"io"
	"sort"
)

const margin = 10.0
const textMargin = 15.0

type TreemapChart struct {
	Dimension
	series        []string
	data          []float64
	numberFormat  string
	colorScheme   *ColorScheme
	showValues    bool
	isInteractive bool
}

func NewTreemapChart(
	width int,
	height int,
	series []string,
	data []float64,
) *TreemapChart {
	return &TreemapChart{
		Dimension: Dimension{
			width:  width,
			height: height,
		},
		colorScheme:   &DefaultColorScheme,
		series:        series,
		data:          data,
		showValues:    false,
		isInteractive: false,
	}
}

func (tm *TreemapChart) SetColorDcheme(colorScheme *ColorScheme) *TreemapChart {
	tm.colorScheme = colorScheme
	return tm
}

func (tm *TreemapChart) SetNumberFormat(numberFormat string) *TreemapChart {
	tm.numberFormat = numberFormat
	return tm
}

func (tm *TreemapChart) SetInteractive(interactive bool) *TreemapChart {
	tm.isInteractive = interactive
	return tm
}
func (tm *TreemapChart) SetShowValue(showValues bool) *TreemapChart {
	tm.showValues = showValues
	return tm
}

type tmSlice struct {
	value   float64
	label   string
	percent float64
	index   int
}

func (tm *TreemapChart) subRenderSVG(w io.Writer, x, y, width, height float64, tmSlices []tmSlice) error {

	subpercent := 0.0
	for _, tmSlice := range tmSlices {
		subpercent += tmSlice.percent
	}

	if width > height {
		ngroups := int(width/height + 1)
		// make ngroups groups
		n := 0
		groups := make([][]tmSlice, 0)
		groupPercents := []float64{0.0}
		groups = append(groups, make([]tmSlice, 0))
		for i := 0; i < len(tmSlices); i++ {
			groups[n] = append(groups[n], tmSlices[i])
			groupPercents[n] += tmSlices[i].percent / subpercent
			if groupPercents[n] >= 1.0/float64(ngroups) && i < len(tmSlices) {
				n++
				groups = append(groups, make([]tmSlice, 0))
				groupPercents = append(groupPercents, 0.0)
			}
		}
		// draw groups
		currentX := x
		currentY := y
		for n, group := range groups {
			currentWidth := width * groupPercents[n]
			if len(group) == 0 {
				continue
			} else if len(group) == 1 {
				fmt.Fprintf(
					w,
					"<rect x='%f' y='%f' width='%f' height='%f' fill='%s' stroke='%s' stroke-width='1'/>",
					currentX,
					currentY,
					currentWidth,
					height,
					tm.colorScheme.ColorPalette(groups[n][0].index),
					tm.colorScheme.Background,
				)
				fmt.Fprintf(
					w,
					"<text x='%f' y='%f' text-anchor='start' fill='#fff'><tspan x='%f' dy='1em'>%s</tspan> <tspan x='%f' dy='1em'>(%g)</tspan></text>",
					currentX+textMargin,
					currentY+textMargin,
					currentX+textMargin,
					groups[n][0].label,
					currentX+textMargin,
					groups[n][0].value,
				)
			} else {
				tm.subRenderSVG(w, currentX, currentY, currentWidth, height, groups[n])
			}
			currentX += currentWidth
		}

	} else {
		ngroups := int(height/width + 1)
		// make ngroups groups
		n := 0
		groups := make([][]tmSlice, 0)
		groupPercents := []float64{0.0}
		groups = append(groups, make([]tmSlice, 0))
		for i := 0; i < len(tmSlices); i++ {
			groups[n] = append(groups[n], tmSlices[i])
			groupPercents[n] += tmSlices[i].percent / subpercent
			if groupPercents[n] >= 1.0/float64(ngroups) && i < len(tmSlices) {
				n++
				groups = append(groups, make([]tmSlice, 0))
				groupPercents = append(groupPercents, 0.0)
			}
		}
		// draw groups
		currentX := x
		currentY := y
		for n, group := range groups {
			currentHeight := height * groupPercents[n]
			if len(group) == 0 {
				continue
			} else if len(group) == 1 {
				fmt.Fprintf(
					w,
					"<rect x='%f' y='%f' width='%f' height='%f' fill='%s' stroke='%s' stroke-width='1' />",
					currentX,
					currentY,
					width,
					currentHeight,
					tm.colorScheme.ColorPalette(groups[n][0].index),
					tm.colorScheme.Background,
				)
				fmt.Fprintf(
					w,
					"<text x='%f' y='%f' text-anchor='start' fill='#fff'><tspan x='%f' dy='1em'>%s</tspan> <tspan x='%f' dy='1em'>(%g)</tspan></text>",
					currentX+textMargin,
					currentY+textMargin,
					currentX+textMargin,
					groups[n][0].label,
					currentX+textMargin,
					groups[n][0].value,
				)
			} else {
				tm.subRenderSVG(w, currentX, currentY, width, currentHeight, groups[n])
			}
			currentY += currentHeight
		}
	}

	return nil
}

func (tm *TreemapChart) RenderSVG(w io.Writer) error {

	startSVG(w, tm.width, tm.height, tm.colorScheme)
	writeFontStyle(w, tm.isInteractive)
	writeBackground(w, tm.width, tm.height, tm.colorScheme)

	tmSlices := make([]tmSlice, len(tm.data))
	total := 0.0
	for i, v := range tm.data {
		tmSlices[i] = tmSlice{
			value: v,
			label: tm.series[i],
		}
		total += v
	}
	sort.SliceStable(tmSlices, func(i, j int) bool {
		return tmSlices[i].value > tmSlices[j].value
	})
	for i, _ := range tmSlices {
		tmSlices[i].percent = tmSlices[i].value / total
		tmSlices[i].index = i
	}

	err := tm.subRenderSVG(w, margin, margin, float64(tm.width)-2*margin, float64(tm.height)-2*margin, tmSlices)
	if err != nil {
		return err
	}

	endSVG(w)

	return nil
}
