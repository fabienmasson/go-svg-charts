package charts

import (
	"fmt"
	"io"
	"math"
)

type Dimension struct {
	width, height int
}

func startTag(w io.Writer, tag string, properties map[string]string) {
	w.Write([]byte("<"))
	w.Write([]byte(tag))
	for k, v := range properties {
		w.Write([]byte(" "))
		w.Write([]byte(k))
		w.Write([]byte("=\""))
		w.Write([]byte(v))
		w.Write([]byte("\""))
	}
	w.Write([]byte(">"))
}

func endTag(w io.Writer, tag string) {
	w.Write([]byte("</"))
	w.Write([]byte(tag))
	w.Write([]byte(">"))
}

func tag(w io.Writer, tag string, properties map[string]string) {
	w.Write([]byte("<"))
	w.Write([]byte(tag))
	for k, v := range properties {
		w.Write([]byte(" "))
		w.Write([]byte(k))
		w.Write([]byte("=\""))
		w.Write([]byte(v))
		w.Write([]byte("\""))
	}
	w.Write([]byte("/>"))
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
func yAxisDimensions(height float64, data [][]float64) ([]string, []float64, [][]float64) {
	min, max := data[0][0], data[0][0]
	for i := 1; i < len(data); i++ {
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

	top := 0.05 * height    // where min value goes
	bottom := 0.95 * height // where max value goes

	conv := func(val float64) float64 {
		return bottom - (bottom-top)*(val-min)/(max-min)
	}

	k := 0
	labels := make([]string, 0)
	lines := make([]float64, 0)
	convdata := make([][]float64, len(data))
	for {
		val := float64(int(min/interval)+1+k) * interval
		if val < max {
			labels = append(labels, fmt.Sprintf("%f", val))
			lines = append(lines, conv(val))
		} else {
			break
		}
	}
	for m, _ := range data {
		convdata[m] = make([]float64, len(data[m]))
		for n, _ := range data[m] {
			convdata[m][n] = conv(data[m][n])
		}
	}

	return labels, lines, convdata
}
