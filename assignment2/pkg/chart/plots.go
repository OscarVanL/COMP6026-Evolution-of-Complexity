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
				//Max: int(result.YValsCCGA[0]),
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "function evalusations",
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: true,
			}),
		)

		line.SetXAxis(result.XValsGA).
			AddSeries("Standard GA", convertLineData(result.YValsGA)).
			AddSeries("CCGA-1", convertLineData(result.YValsCCGA)).
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