package charts_test

import (
	"math/rand"
	"os"
	"testing"

	charts "github.com/fabienmasson/go-svg-charts"
)

func TestBarChart(t *testing.T) {

	caeq1 := make([]float64, 0)
	caeq2 := make([]float64, 0)
	//caeq3 := make([]float64, 0)
	//caeq4 := make([]float64, 0)
	//caeq5 := make([]float64, 0)

	for i := 0; i < 12; i++ {
		caeq1 = append(caeq1, rand.Float64()*10)
		caeq2 = append(caeq2, rand.Float64()*20)
		//caeq3 = append(caeq3, rand.Float64()*25)
		//caeq4 = append(caeq4, rand.Float64()*30)
		//caeq5 = append(caeq5, rand.Float64()*100)
	}

	lc := charts.NewBarChart(
		800,
		400,
		[]string{"Jan", "Feb", "Mar", "Apr", "Mai", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		[]string{"Team 1", "Team 2" /*"Team 3", "Team 4", "Team 5"*/},
		[][]float64{caeq1, caeq2 /*caeq3, caeq4, caeq5*/},
	).
		SetXaxisLegend("Month").
		SetYaxisLegend("Net growth").SetInteractive(true)

	file, err := os.Create("examples/barchart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
