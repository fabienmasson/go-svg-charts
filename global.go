package charts

import (
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
// - a slice of float64 containing the horizontal lines height
// - A in Ax+B
// - B in Ax+B
func yAxisDimensions(height float64, data [][]float64) []float64 {
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

}
