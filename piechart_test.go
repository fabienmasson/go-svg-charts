package charts_test

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
)

func TestPieChart(t *testing.T) {

	ca := make([]float64, 0)
	labels := make([]string, 0)

	for i := 0; i < 4; i++ {
		ca = append(ca, math.Round(rand.Float64()*1000)/100)
		labels = append(labels, fmt.Sprintf("Team %d", i))
	}

	lc := charts.NewPieChart(
		500,
		500,
		labels,
		ca,
	).
		SetInteractive(false).
		SetShowValue(true)

	file, err := os.Create("piechart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
