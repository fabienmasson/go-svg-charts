package charts_test

import (
	"math"
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
)

func TestLineChart(t *testing.T) {

	caeq1 := make([]float64, 0)
	caeq2 := make([]float64, 0)
	caeq3 := make([]float64, 0)
	caeq4 := make([]float64, 0)
	caeq5 := make([]float64, 0)

	for i := 0; i < 12; i++ {
		caeq1 = append(caeq1, math.Round(rand.Float64()*10000)/10000)
		caeq2 = append(caeq2, math.Round(rand.Float64()*10000)/10000)
		caeq3 = append(caeq3, math.Round(rand.Float64()*10000)/10000)
		caeq4 = append(caeq4, math.Round(rand.Float64()*10000)/10000)
		caeq5 = append(caeq5, math.Round(rand.Float64()*10000)/10000)
	}

	lc := charts.NewLineChart(
		800,
		500,
		[]string{"Jan", "Feb", "Mar", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		[]string{"Team 1", "Team 2", "Team 3", "Team 4", "Team 5"},
		[][]float64{caeq1, caeq2, caeq3, caeq4, caeq5},
	).
		SetXaxisLegend("Month").
		SetYaxisLegend("Net growth").
		SetShowMarkers(true).
		SetInteractive(true)

	file, err := os.Create("linechart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
