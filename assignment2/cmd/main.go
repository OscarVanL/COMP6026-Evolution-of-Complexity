package main

import (
	"fmt"
	pop "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"github.com/wcharczuk/go-chart"
	"math"
	"os"
)

const iterations = 1000

type Algorithm int

const (
	Rastrigin Algorithm = iota
	Schwefel
	Griewangk
	Ackley
	Rosenbrock
)

func main() {

	var results []EvolutionResults

	fmt.Println("Benchmarking Rastrigin Function")
	results = append(results, DoGeneticAlgorithms(Rastrigin))
	fmt.Println("Benchmarking Schwefel Function")
	results = append(results, DoGeneticAlgorithms(Schwefel))
	fmt.Println("Benchmarking Griewangk Function")
	results = append(results, DoGeneticAlgorithms(Griewangk))
	fmt.Println("Benchmarking Ackley Function")
	results = append(results, DoGeneticAlgorithms(Ackley))

	PlotResults(results)

}

func PlotResults(res []EvolutionResults) {
	// Create a ContinuousSeries for each of the results
	for i:=0; i<len(res); i++ {
		result := res[i]

		// Add results to a titled chart
		graph := chart.Chart {
			Title: result.label,
			Series: []chart.Series{
				chart.ContinuousSeries{
					Name: "CCGA-1",
					XValues: result.XValsCCGA,
					YValues: result.YValsCCGA,
				},
				chart.ContinuousSeries{
					Name: "Standard GA",
					XValues: result.XValsGA,
					YValues: result.YValsGA,
				},
			},
		}

		// Create legend for results
		graph.Elements = []chart.Renderable{chart.Legend(&graph)}

		// Save charts to PNG files
		file, _ := os.Create(fmt.Sprintf("%s.png", result.label))
		err := graph.Render(chart.PNG, file)

		if err != nil {
			fmt.Println(err)
		}
	}


}

type EvolutionResults struct {
	label string  // Label to represent result
	XValsCCGA []float64  // Iteration values for X-Axis, for CCGA-1
	YValsCCGA []float64  // Fitness values for Y-Axis, for CCGA-1
	BestFitnessCCGA float64  // Best Fitness from CCGA-1
	BestAssignmentCCGA []uint16  // Best assignment of function variables from CCGA-1

	XValsGA []float64  // Iteration values for X-Axis, standard GA
	YValsGA []float64  // Fitness values for Y-Axis, standard GA
	BestFitnessGA float64  // Best Fitness from standard GA
	BestAssignmentGA []uint16  // Best assignment of function variables from standard GA
}


func DoGeneticAlgorithms(algo Algorithm) EvolutionResults {
	var label string
	var N int
	var function pop.Fitness
	var mutationP float64

	switch algo {
	case Rastrigin:
		label = "Rastrigin Function"
		N = f.RastriginN
		function = f.Rastrigin
		mutationP = f.RastriginMutationP
	case Schwefel:
		label = "Schwefel Function"
		N = f.SchwefelN
		function = f.Schwefel
		mutationP = f.SchwefelMutationP
	case Griewangk:
		label = "Griewangk Function"
		N = f.GriewangkN
		function = f.Griewangk
		mutationP = f.GriewangkMutationP
	case Ackley:
		label = "Ackley Function"
		N = f.AckleyN
		function = f.Ackley
		mutationP = f.AckleyMutationP
	case Rosenbrock:
		// Todo: Write Rosenbrock function
		label = "Rosenbrock Function"
		N = 10
		function = f.Rosenbrock
		mutationP = 1/10
	}

	XValsCCGA, YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA := CCGA1(N, function, mutationP)
	XValsGA, YValsGA, BestFitnessGA, BestAssignmentGA := GA(N, function, mutationP)

	// Todo: Add results from standard GA
	return EvolutionResults{label,
		XValsCCGA, YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA,
		XValsGA, YValsGA, BestFitnessGA, BestAssignmentGA}
}

func CCGA1(N int, function pop.Fitness, mutationP float64) ([]float64, []float64, float64, []uint16) {
	species := pop.GenerateSpecies(N, 100)
	species.InitCoevolutions()
	species.EvalFitness(function)

	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	var xVal []float64
	var fitnessHistory []float64

	for i:=0; i<iterations; i++ {
		xVal = append(xVal, float64(i+1))

		species.EvolveSpecies(mutationP)
		species.EvalFitness(function)
		//Todo: Do mutation & evolution cycle
		fitness, coevolution := species.GetBestFitness()

		fitnessHistory = append(fitnessHistory, fitness)
		if fitness < bestFitness {
			fmt.Println("New best fitness:", fitness)
			bestFitness = fitness
			bestCoevolution = coevolution
		}

		if i % 1000 == 0 {
			fmt.Println(i, ":", fitness)
		}
	}

	//fmt.Println("Fitness history:", fitnessHistory)
	fmt.Println("Best Coevolution fitness:", bestFitness, ". Parameters:")
	for i:=0; i<len(bestCoevolution); i++ {
		fmt.Print(bestCoevolution[i], ", ")
	}
	fmt.Println()

	return xVal, fitnessHistory, bestFitness, bestCoevolution
}

func GA(N int, function pop.Fitness, mutationP float64) ([]float64, []float64, float64, []uint16) {
	// TODO: Do normal GA
	return []float64{}, []float64{}, 0.0, []uint16{}
}