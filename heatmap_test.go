package charts_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
)

func TestHeatMap(t *testing.T) {

	activity := make([][]float64, 12)
	months := []string{"Jan", "Feb", "Mar", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	days := make([]string, 0)
	for i := 0; i < 31; i++ {
		days = append(days, fmt.Sprintf("%d", i))
	}
	for i := 0; i < 12; i++ {
		activity[i] = make([]float64, 31)
		for j := 0; j < 31; j++ {
			if j < 31 || (i != 1 && i != 3 && i != 5 && i != 8 && i != 10) {
				activity[i][j] = rand.Float64() * 10
			}
		}
	}

	lc := charts.NewHeatMap(
		800,
		800,
		months, days, activity,
	).
		SetXaxisLegend("Month").
		SetYaxisLegend("Day").
		SetInteractive(true)

	file, err := os.Create("examples/heatmap.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
