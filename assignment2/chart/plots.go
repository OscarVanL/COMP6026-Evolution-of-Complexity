package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

type EvolutionResults struct {
	Title string  // Title to represent result
	XLabel string  // Label to give X Axis
	Iterations int  // Number of function evaluations represented in charts
	CCGAFitnessHistory []BestFitness  // Best fitness over time for CCGA
	BestFitnessCCGA float64  // Best Fitness from CCGA-1

	GAFitnessHistory []BestFitness  // Best fitness over time for GA
	BestFitnessGA float64  // Best Fitness from standard GA
}



type BestFitness struct {
	X int
	Fitness float64
}

func PlotResults(output string, res []EvolutionResults) {
	page := components.NewPage()
	xVals := initXValsSlice(res[0].Iterations)

	for i:=0; i<len(res); i++ {
		result := res[i]
		// Take best fitness history and fill in the gaps so data can be plotted
		yValsCCGA := fillMissingPoints(result.Iterations, result.CCGAFitnessHistory)
		yValsGA := fillMissingPoints(result.Iterations, result.GAFitnessHistory)

		line := charts.NewLine()

		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle: "Comparison of standard GA and CCGA-1 performance",
				Width: "625px",
				Height: "450px",
			}),
			charts.WithTitleOpts(opts.Title{
				Title: result.Title,
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "best individual",
				Max: int(yValsCCGA[0]),
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: result.XLabel,
			}),
			charts.WithLegendOpts(opts.Legend{
				Show: true,
			}),
		)

		line.SetXAxis(xVals).
			AddSeries("Standard GA", convertLineData(yValsGA)).
			AddSeries("CCGA-1", convertLineData(yValsCCGA)).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: true,
				}),
		)

		page.AddCharts(line)
	}

	f, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	_ = page.Render(io.MultiWriter(f))
}

func initXValsSlice(iterations int) []int {
	xVals := make([]int, iterations)
	for i:=0; i<iterations; i++ {
		xVals[i] = i
	}
	return xVals
}


// fillMissingPoints fills in gaps in results data so that scores are properly spaced on plots
func fillMissingPoints(iterations int, BestFitnessHistory []BestFitness) []float64 {
	yVals := make([]float64, iterations)

	lastStart := 0
	for score:=0; score<len(BestFitnessHistory); score++ {
		bestScore := BestFitnessHistory[score].Fitness

		if score == len(BestFitnessHistory)-1 {
			for i:=lastStart; i<iterations; i++ {
				yVals[i] = bestScore
			}
		} else {
			for i:=lastStart; i<BestFitnessHistory[score+1].X; i++ {
				yVals[i] = bestScore
			}
			lastStart = BestFitnessHistory[score+1].X
		}
	}

	return yVals
}

func convertLineData(elems []float64) []opts.LineData {
	points := make([]opts.LineData, 0)
	for i := 0; i < len(elems); i++ {
		points = append(points, opts.LineData{Value: elems[i]})
	}
	return points
}