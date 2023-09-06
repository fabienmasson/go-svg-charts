package charts_test

import (
	"math/rand"
	"os"
	"testing"
	"trankiloubilou/charts"
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
		600,
		[]string{"Janvier", "Fevrier", "Mars", "Avril", "Mai", "Juin", "Juillet", "Aout", "Septembre", "Octobre", "Novembre", "Decembre"},
		[]string{"Equipe 1", "Equipe 2" /*"Equipe 3", "Equipe 4", "Equipe 5"*/},
		[][]float64{caeq1, caeq2 /*caeq3, caeq4, caeq5*/},
	).
		SetXaxisLegend("Mois").
		SetYaxisLegend("CA")

	file, err := os.Create("barchart.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	lc.RenderSVG(file)

}
