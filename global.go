package charts

import (
	"fmt"
	"io"
	"math"
)

type Dimension struct {
	width, height int
}

// yAxisDimensions calculates the y-axis dimensions for a given height and data.
//
// Parameters:
// - height: the height of the y-axis.
// - data: a 2D slice of float64 values representing the data.
//
// Returns:
// - slice of line labels
// - slice of lines y
// - slice of converted data y
func yAxisFit(start int, end int, data [][]float64) ([]string, []float64, func(float64) float64) {
	min, max := data[0][0], data[0][0]
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] < min {
				min = data[i][j]
			}
			if data[i][j] > max {
				max = data[i][j]
			}
		}
	}

	diff := max - min
	log10 := math.Log10(diff)
	i, f := math.Modf(log10)
	var interval float64
	if f < 0.3 {
		interval = math.Pow10(int(i)) / 5
	} else if f < 0.7 {
		interval = math.Pow10(int(i)) / 2
	} else {
		interval = math.Pow10(int(i))
	}

	height := float64(end - start)
	top := 0.0       // where max value goes
	bottom := height // where min value goes

	fmt.Printf("top:%f; bottom:%f; start:%d; end: %d; min: %f, max: %f, interval: %f\n", top, bottom, start, end, min, max, interval)

	conv := func(val float64) float64 {
		return float64(start) + bottom - (bottom-top)*(val-min)/(max-min)
	}

	k := 0
	labels := make([]string, 0)
	lines := make([]float64, 0)
	for {
		val := float64(int(min/interval)+1+k) * interval
		if val < max {
			if val > 1 || val <= -1 {
				labels = append(labels, fmt.Sprintf("%d", int(val)))
			} else {
				labels = append(labels, fmt.Sprintf("%f", val))
			}
			lines = append(lines, val)
		} else {
			break
		}
		k++
	}

	return labels, lines, conv
}

func seriesLegend(w io.Writer, height int, width int, series []string, colors *ColorScheme) {
	rectDimension := Dimension{30, 10}

	for i, serie := range series {

		fmt.Fprintf(
			w,
			"<rect width='%d' height='%d' fill='%s' />",
			rectDimension.width,
			rectDimension.height,
			colors.ColorPalette[i],
		)

		fmt.Fprintf(
			w,
			"<text>%s</text>",
			serie,
		)
	}

}
