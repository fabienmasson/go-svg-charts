package charts

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//go:embed maps/*/*.svg
var folder embed.FS

func GetAvailableMaps() []string {
	return []string{
		"australia",
		"austria",
		"brazil",
		"cambodia",
		"cameroon",
		"canada",
		"canada.lambert-projection",
		"cape-verde",
		"china",
		"cli",
		"colombia",
		"denmark",
		"france.departments",
		"france.regions",
		"germany",
		"greece",
		"honduras",
		"hong-kong",
		"india",
		"indonesia",
		"israel",
		"italy",
		"japan",
		"kenya",
		"mexico",
		"moldova",
		"netherlands",
		"new-zealand",
		"nigeria",
		"pakistan.districts",
		"puerto-rico",
		"romania",
		"saudi-arabia",
		"south-korea",
		"spain",
		"sri-lanka",
		"sweden",
		"taiwan",
		"taiwan.main",
		"tanzania",
		"thailand",
		"tunisia",
		"uae",
		"ukraine",
		"usa",
		"usa.counties",
		"usa.florida",
		"usa.michigan",
		"usa.states-territories",
		"usa.utah",
		"uzbekistan",
		"world",
		"world.capitals",
		"zimbabwe",
	}
}

type GeoMap struct {
	Dimension
	mapName       string
	data          map[string]float64
	numberFormat  string
	colorScheme   *ColorScheme
	showValues    bool
	isInteractive bool
}

func NewGeoMap(
	mapName string,
	data map[string]float64,
) *GeoMap {
	return &GeoMap{
		colorScheme:   &DefaultColorScheme,
		mapName:       mapName,
		data:          data,
		showValues:    false,
		isInteractive: false,
	}
}

func (gm *GeoMap) SetColorDcheme(colorScheme *ColorScheme) *GeoMap {
	gm.colorScheme = colorScheme
	return gm
}

func (gm *GeoMap) SetNumberFormat(numberFormat string) *GeoMap {
	gm.numberFormat = numberFormat
	return gm
}

func (gm *GeoMap) SetInteractive(interactive bool) *GeoMap {
	gm.isInteractive = interactive
	return gm
}
func (gm *GeoMap) SetShowValue(showValues bool) *GeoMap {
	gm.showValues = showValues
	return gm
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func walk(nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.Nodes, f)
		}
	}
}

func (gm *GeoMap) RenderSVG(w io.Writer) error {

	templateMap, err := folder.ReadFile(fmt.Sprintf("maps/%s/%s.svg", gm.mapName, gm.mapName))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(templateMap)
	dec := xml.NewDecoder(buf)

	labelBuffer := new(bytes.Buffer)

	var n Node
	err = dec.Decode(&n)
	if err != nil {
		return err
	}

	walk([]Node{n}, func(n Node) bool {
		if n.XMLName.Local == "svg" {
			width, height := 1024, 1024
			for _, attr := range n.Attrs {
				if strings.ToUpper(attr.Name.Local) == "VIEWBOX" {
					s := strings.Split(attr.Value, " ")
					width, _ = strconv.Atoi(s[2])
					height, _ = strconv.Atoi(s[3])
					break
				}
			}
			startSVG(w, width, height, gm.colorScheme)
			writeDefsTxtBg(w, gm.colorScheme)
			writeFontStyle(w, gm.isInteractive)
			writeBackground(w, width, height, gm.colorScheme)

			gm.styleSVG(w)
			for _, node := range n.Nodes {
				fmt.Fprint(w, "<path ")
				var name, id string
				var labelx, labely float64
				for _, attr := range node.Attrs {
					fmt.Fprintf(w, "%s=\"%s\" ", attr.Name.Local, attr.Value)
					if attr.Name.Local == "name" {
						name = attr.Value
					}
					if attr.Name.Local == "id" {
						id = attr.Value
					}
					if attr.Name.Local == "d" {
						attr.Value = strings.ReplaceAll(attr.Value, "m", "")
						attr.Value = strings.Trim(attr.Value, " ")
						arr := strings.FieldsFunc(attr.Value, splitfn)
						labelx, _ = strconv.ParseFloat(arr[0], 64)
						labely, _ = strconv.ParseFloat(arr[1], 64)
					}
				}
				fmt.Fprint(w, "/>")
				if gm.showValues || gm.isInteractive {
					fmt.Fprint(labelBuffer, "<path ")
					if gm.isInteractive {
						fmt.Fprint(labelBuffer, "class='hovercircle' fill-opacity='0' ")
					}
					for _, attr := range node.Attrs {
						if strings.ToUpper(attr.Name.Local) != "ID" {
							fmt.Fprintf(labelBuffer, "%s=\"%s\" ", attr.Name.Local, attr.Value)
						}
					}
					fmt.Fprint(labelBuffer, "/>")
					fmt.Fprintf(
						labelBuffer,
						"<text style='paint-order:stroke fill' class='value' text-anchor='middle' alignment-baseline='middle' filter='url(#textbg)' x='%f' y='%f'>%s (%g)</text>",
						labelx,
						labely,
						name,
						gm.data[id],
					)
				}
			}
		}
		return true
	})
	if gm.showValues || gm.isInteractive {
		w.Write(labelBuffer.Bytes())
	}

	endSVG(w)

	return nil
}

func (hm *GeoMap) styleSVG(w io.Writer) error {

	var min, max float64
	var start = true
	for _, value := range hm.data {
		if start {
			min, max = value, value
			start = false
			continue
		}
		if value < min {
			min = value
		} else if value > max {
			max = value
		}
	}

	fmt.Fprint(w, "<style> ")
	/*
		fmt.Fprintf(w, " path { fill: %s; stroke: %s; stroke-width: 0.5; } \n",
			hm.colorScheme.Background,
			hm.colorScheme.LightAxisColor,
		)
	*/
	for k, v := range hm.data {

		fmt.Fprintf(w, " path[id='%s'] { fill: %s;  fill-opacity: %f; } \n",
			k,
			hm.colorScheme.ColorPalette(0),
			(v-min)/(max-min),
		)
	}
	fmt.Fprint(w, " </style>")
	return nil
}

func splitfn(r rune) bool {
	return r == ',' || r == ' '
}
