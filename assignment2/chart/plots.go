package chart

import (
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

type EvolutionResults struct {
	Title              string        // Title to represent result
	XLabel             string        // Label to give X Axis
	Iterations         int           // Number of function evaluations represented in charts
	CCGAFitnessHistory []BestFitness // Best fitness over time for CCGA
	BestFitnessCCGA    float64       // Best Fitness from CCGA-1
	BestAssignmentCCGA []uint16      // Best assignment of genes

	GAFitnessHistory []BestFitness // Best fitness over time for GA
	BestFitnessGA    float64       // Best Fitness from standard GA
	BestAssignmentGA []uint16      // Best assignment of genes
}

type BestFitness struct {
	X       int
	Fitness float64
}

func PlotResults(output string, res [][]EvolutionResults) {
	page := components.NewPage()

	xVals := initXValsSlice(res[0][0].Iterations)



	for i := 0; i < len(res); i++ {
		result := res[i][0]
		var YMax int
		// Set Y-axis limits to same as those in paper
		switch result.Title {
		case optimisation.RastriginLabel:
			YMax = 40
		case optimisation.SchwefelLabel:
			YMax = 400
		case optimisation.GriewangkLabel:
			YMax = 8
		case optimisation.AckleyLabel:
			YMax = 16
		}
		// Calculate average result, filling in any gaps in the data
		yValsGA, yValsCCGA := averageResults(result.Iterations, res[i])
		fmt.Println(result.Title, ": Best Average Fitness GA for:", yValsGA[len(yValsCCGA)-1])
		fmt.Println(result.Title, ": Best Average Fitness CCGA:", yValsCCGA[len(yValsCCGA)-1])

		line := charts.NewLine()

		line.SetGlobalOptions(
			charts.WithInitializationOpts(opts.Initialization{
				PageTitle: "Comparison of standard GA and CCGA-1 performance",
				Width:     "625px",
				Height:    "450px",
			}),
			charts.WithTitleOpts(opts.Title{
				Title: result.Title,
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "best individual",
				//Max: int(math.Min(yValsCCGA[0], yValsGA[0])),
				Max:  YMax,
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
	for i := 0; i < iterations; i++ {
		xVals[i] = i
	}
	return xVals
}

func averageResults(iterations int, results []EvolutionResults) ([]float64, []float64) {
	allYValsGA, allYValsCCGA := make([][]float64, iterations), make([][]float64, iterations)
	yValsGAAveraged, yValsCCGAAveraged := make([]float64, iterations), make([]float64, iterations)

	// Fill in missing points for each result
	for res := 0; res < len(results); res++ {
		allYValsGA[res] = fillMissingPoints(iterations, results[res].GAFitnessHistory)
		allYValsCCGA[res] = fillMissingPoints(iterations, results[res].CCGAFitnessHistory)
	}

	// Calculate averages for each point
	for i := 0; i < iterations; i++ {
		pointSumGA, pointSumCCGA := 0.0, 0.0
		for res := 0; res < len(results); res++ {
			pointSumGA += allYValsGA[res][i]
			pointSumCCGA += allYValsCCGA[res][i]
		}
		yValsGAAveraged[i] = pointSumGA / float64(len(results))
		yValsCCGAAveraged[i] = pointSumCCGA / float64(len(results))
	}

	return yValsGAAveraged, yValsCCGAAveraged
}

// fillMissingPoints fills in gaps in results data so that scores are properly spaced on plots
func fillMissingPoints(iterations int, BestFitnessHistory []BestFitness) []float64 {
	yVals := make([]float64, iterations)

	lastStart := 0
	for score := 0; score < len(BestFitnessHistory); score++ {
		bestScore := BestFitnessHistory[score].Fitness

		if score == len(BestFitnessHistory)-1 {
			for i := lastStart; i < iterations; i++ {
				yVals[i] = bestScore
			}
		} else {
			for i := lastStart; i < BestFitnessHistory[score+1].X; i++ {
				if i < iterations {
					yVals[i] = bestScore
				}
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
