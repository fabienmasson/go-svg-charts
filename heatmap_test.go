package charts_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	charts "github.com/fabienmasson/go-svg-charts"
)

func TestHeatMap(t *testing.T) {

	activity := make([][]float64, 12)
	months := []string{"Jan", "Feb", "Mar", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	days := make([]string, 0)
	for i := 0; i < 31; i++ {
		days = append(days, fmt.Sprintf("%d", i+1))
	}
	for i := 0; i < 12; i++ {
		activity[i] = make([]float64, 31)
		for j := 0; j < 31; j++ {
			activity[i][j] = rand.Float64() * 10
		}
	}
	activity[1][30] = 0.0
	activity[1][29] = 0.0
	activity[1][28] = 0.0
	activity[3][30] = 0.0
	activity[5][30] = 0.0
	activity[8][30] = 0.0
	activity[10][30] = 0.0

	lc := charts.NewHeatMap(
		800,
		400,
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
