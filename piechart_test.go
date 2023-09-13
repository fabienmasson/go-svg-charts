package charts_test

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	charts "github.com/fabienmasson/go-svg-charts"
)

func TestPieChart(t *testing.T) {

	ca := make([]float64, 0)
	labels := make([]string, 0)

	for i := 0; i < 15; i++ {
		ca = append(ca, math.Round(rand.Float64()*1000)/100)
		labels = append(labels, fmt.Sprintf("Team %d", i))
	}

	lc := charts.NewPieChart(
		800,
		400,
		labels,
		ca,
	).
		SetInteractive(false).
		SetShowValue(true)

	file, err := os.Create("examples/piechart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
