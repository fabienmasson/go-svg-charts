package charts_test

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	charts "github.com/fabienmasson/go-svg-charts"
)

func TestGeoMap(t *testing.T) {

	data := make(map[string]float64)
	for i := 0; i < 96; i++ {
		data[fmt.Sprintf("%02d", i)] = math.Round(rand.Float64() * 100000)
	}
	data["2A"] = math.Round(rand.Float64() * 100000)
	data["2B"] = math.Round(rand.Float64() * 100000)

	lc := charts.NewGeoMap(
		"france.departments",
		data,
	).SetInteractive(true)

	file, err := os.Create("examples/geomap.svg")
	if err != nil {
		t.Errorf("os.Create error: %s", err)
	}
	defer file.Close()

	err = lc.RenderSVG(file)
	if err != nil {
		t.Errorf("Error rendering SVG: %s", err)
	}

}
