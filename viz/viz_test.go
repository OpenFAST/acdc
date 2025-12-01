package viz_test

import (
	"acdc/viz"
	"fmt"
	"os"
	"testing"

	"github.com/wcharczuk/go-chart/v2" // This can be deleted later after rendering the graph in the frontend. Need v2 to show axes labels.
)

func TestBuildModeViz(t *testing.T) {

	data, err := viz.ParseModeData([]string{
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.001.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.002.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.003.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.004.vtp",
		"testdata/03_NREL_5MW-ED.Mode1.LinTime1.ED_TowerLn2Mesh_motion.005.vtp",
	})

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", data)

	PlotTipDeflection(data)
}

func PlotTipDeflection(data *viz.ModeData) {

	fmt.Println(data.Frames)

	// Get Component names - BD_BldMotion1, BD_BldMotion2, BD_BldMotion3.
	componentNames := make(map[string]struct{})
	for _, frame := range data.Frames {
		for k := range frame.Components {
			componentNames[k] = struct{}{}
		}
	}
	fmt.Println("Component Names:", componentNames)
	seriesList := make([]chart.Series, 0, len(componentNames)*2)

	// Loop over each blade
	for componentName := range componentNames {
		fmt.Println("Adding series of ", componentName)
		// For the Blade Tip, we only need the last point of each frame
		tipFlap := make([]float64, len(data.Frames))
		tipEdge := make([]float64, len(data.Frames))
		frames := make([]float64, len(data.Frames))
		for i, frame := range data.Frames {
			frames[i] = float64(i + 1) // Frame numbers start from 1
			if component, ok := frame.Components[componentName]; ok {
				if len(component.LocalLine) > 0 {
					tipFlap[i] = float64(component.LocalLine[len(component.LocalLine)-1].XYZ[0]) // X coordinate of the last point
					tipEdge[i] = float64(component.LocalLine[len(component.LocalLine)-1].XYZ[1]) // Y coordinate of the last point
				}
			}
		}

		fmt.Println("Frames:", frames)
		fmt.Println("Tip Flap:", tipFlap)
		fmt.Println("Tip Edge:", tipEdge)

		// Create a new series
		flapSeries := chart.ContinuousSeries{
			Name:    "Flap_" + componentName,
			XValues: frames,
			YValues: tipFlap,
		}

		edgeSeries := chart.ContinuousSeries{
			Name:    "Edge_" + componentName,
			XValues: frames,
			YValues: tipEdge,
		}

		seriesList = append(seriesList, flapSeries, edgeSeries)
	}

	// Create a new chart
	graph := chart.Chart{
		Title:  "",
		Series: seriesList,
		XAxis: chart.XAxis{
			Name: "Frames",
		},
		YAxis: chart.YAxis{
			Name: "Tip Deflection",
		},
	}

	// Add the legend to the chart
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	// Save the plot as a file
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

}
