package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

type EvolutionResults struct {
	Label string  // Label to represent result
	XValsCCGA []float64  // Iteration values for X-Axis, for CCGA-1
	YValsCCGA []float64  // Fitness values for Y-Axis, for CCGA-1
	BestFitnessCCGA float64  // Best Fitness from CCGA-1
	BestAssignmentCCGA []uint16  // Best assignment of function variables from CCGA-1

	XValsGA []float64  // Iteration values for X-Axis, standard GA
	YValsGA []float64  // Fitness values for Y-Axis, standard GA
	BestFitnessGA float64  // Best Fitness from standard GA
	BestAssignmentGA []uint16  // Best assignment of function variables from standard GA
}

func PlotResults(res []EvolutionResults) {
	page := components.NewPage()

	for i:=0; i<len(res); i++ {
		result := res[i]
		line := charts.NewLine()

		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: result.Label,
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "best individual",
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "function evalusations",
			}),
		)

		line.SetXAxis(result.XValsGA).
			AddSeries("Standard GA", convertLineData(result.YValsGA)).
			//	charts.WithLabelOpts(opts.Label{Show: true, Position: "right"})).
			AddSeries("CCGA-1", convertLineData(result.YValsCCGA)).
			//	charts.WithLabelOpts(opts.Label{Show: true, Position: "right"})).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: true,
				}),
		)

		page.AddCharts(line)
	}

	f, err := os.Create("charts.html")
	if err != nil {
		panic(err)
	}
	_ = page.Render(io.MultiWriter(f))
}

func convertLineData(elems []float64) []opts.LineData {
	points := make([]opts.LineData, 0)
	for i := 0; i < len(elems); i++ {
		points = append(points, opts.LineData{Value: elems[i]})
	}
	return points
}


//func PlotResults(res []EvolutionResults) {
//	// Create a ContinuousSeries for each of the results
//	for i:=0; i<len(res); i++ {
//		result := res[i]
//
//		// Add results to a titled chart
//		graph := chart.Chart {
//			Title: result.Label,
//			TitleStyle: chart.Style {
//				FontSize: 16,
//			},
//			XAxis: chart.XAxis{
//				Name: "function evaluations",
//			},
//			YAxis: chart.YAxis{
//				NameStyle: chart.Style {
//					TextRotationDegrees: 270,
//				},
//				AxisType: chart.YAxisSecondary,
//				Name: "best individual",
//			},
//			Background: chart.Style{
//				Padding: chart.Box{
//					Top: 80,
//					Left: 30,
//				},
//			},
//			Series: []chart.Series{
//				chart.ContinuousSeries{
//					Name: "CCGA-1",
//					XValueFormatter: chart.IntValueFormatter,
//					XValues: result.XValsCCGA,
//					YValues: result.YValsCCGA,
//				},
//				chart.ContinuousSeries{
//					Name: "Standard GA",
//					XValueFormatter: chart.IntValueFormatter,
//					XValues: result.XValsGA,
//					YValues: result.YValsGA,
//				},
//			},
//		}
//
//		// Create legend for results
//		graph.Elements = []chart.Renderable{chart.LegendThin(&graph)}
//
//		// Save charts to PNG files
//		file, _ := os.Create(fmt.Sprintf("%s.png", result.Label))
//
//		graph.Width = 700
//		graph.Height = 500
//		err := graph.Render(chart.PNG, file)
//
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}