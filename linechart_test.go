package charts_test

import (
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
)

func TestLineChart(t *testing.T) {

	caeq1 := make([]float64, 0)
	caeq2 := make([]float64, 0)

	for i := 0; i < 12; i++ {
		caeq1 = append(caeq1, rand.Float64()*10)
		caeq2 = append(caeq2, rand.Float64()*20)
	}

	lc := charts.NewLineChart(
		800,
		600,
		[]string{"Janvier", "Fevrier", "Mars", "Avril", "Mai", "Juin", "Juillet", "Aout", "Septembre", "Octobre", "Novembre", "Decembre"},
		[]string{"Equipe 1", "Equipe 2"},
		[][]float64{caeq1, caeq2},
	)

	file, err := os.Create("linechart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
