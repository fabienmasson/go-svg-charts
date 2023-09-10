package charts_test

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
)

func TestTreemapChart(t *testing.T) {

	ca := make([]float64, 0)
	labels := make([]string, 0)

	for i := 0; i < 10; i++ {
		ca = append(ca, math.Round(math.Pow(rand.Float64()*1000, 2))/100)
		labels = append(labels, fmt.Sprintf("Team %d", i))
	}

	tm := charts.NewTreemapChart(
		800,
		500,
		labels,
		ca,
	).
		SetInteractive(false).
		SetShowValue(true)

	file, err := os.Create("examples/treemapchart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	tm.RenderSVG(file)

}
