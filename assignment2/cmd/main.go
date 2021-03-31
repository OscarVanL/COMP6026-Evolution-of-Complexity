package main

import (
	"flag"
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/ccga"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/ga"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"github.com/cheggaaa/pb"
	"github.com/wcharczuk/go-chart"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"
)

const iterations = 100000

type Algorithm int

const (
	Rastrigin Algorithm = iota
	Schwefel
	Griewangk
	Ackley
	Rosenbrock
)

var cpuprofile = flag.String("cpuprofile", "rastrigin.prof", "write cpu profile to file")


func main() {

	if *cpuprofile != "" {
		fmt.Println(*cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var results []EvolutionResults

	fmt.Println("Benchmarking Rastrigin Function")
	results = append(results, DoGeneticAlgorithms(Rastrigin))
	fmt.Println("Benchmarking Schwefel Function")
	results = append(results, DoGeneticAlgorithms(Schwefel))
	fmt.Println("Benchmarking Griewangk Function")
	results = append(results, DoGeneticAlgorithms(Griewangk))
	fmt.Println("Benchmarking Ackley Function")
	results = append(results, DoGeneticAlgorithms(Ackley))

	fmt.Println("Creating Charts")
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
	var function f.Fitness
	var mutationP float32

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


	fmt.Println("Starting standard GA with params: N:", N, "MutationP:", mutationP)
	XValsGA, YValsGA, BestFitnessGA, BestAssignmentGA := GA(N, function, mutationP)
	fmt.Println("Starting CCGA-1 with params: N:", N, "MutationP:", mutationP)
	XValsCCGA, YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA := CCGA1(N, function, mutationP)

	return EvolutionResults{label,
		XValsCCGA, YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA,
		XValsGA, YValsGA, BestFitnessGA, BestAssignmentGA}
}

func CCGA1(N int, function f.Fitness, mutationP float32) ([]float64, []float64, float64, []uint16) {
	bar := pb.New(iterations)
	bar.SetRefreshRate(time.Second)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	var xVal []float64
	var fitnessHistory []float64

	species := ccga.InitSpecies(N, 100)
	species.InitCoevolutions()
	species.EvalFitness(function)

	for i:=0; i<iterations; i++ {
		xVal = append(xVal, float64(i+1))  // Evolution iteration for X-Axis

		// Coevolves individuals with the best (mutated) genes from each species
		species.Coevolve()
		// Mutates each individual's genes
		species.Mutate(mutationP)
		// Re-evaluates fitness
		species.EvalFitness(function)
		// Finds individual with best fitness & genes in this generation
		fitness, coevolution := species.GetBestFitness()

		if fitness < bestFitness {
			bestFitness = fitness
			bestCoevolution = coevolution
		}
		fitnessHistory = append(fitnessHistory, bestFitness)
		bar.Increment()
	}
	bar.Finish()

	fmt.Println("Best Coevolution fitness:", bestFitness, ". Parameters:")
	for i:=0; i<len(bestCoevolution); i++ {
		fmt.Print(bestCoevolution[i], ", ")
	}
	fmt.Println()

	return xVal, fitnessHistory, bestFitness, bestCoevolution
}

func GA(N int, function f.Fitness, mutationP float32) ([]float64, []float64, float64, []uint16) {
	bar := pb.New(iterations)
	bar.SetRefreshRate(time.Second)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	bestFitness := math.MaxFloat64
	var bestGenes []uint16
	var xVal []float64
	var fitnessHistory []float64

	population := ga.InitPopulation(N, 100)
	population.EvalFitness(function)

	for i:=0; i<iterations; i++ {
		xVal = append(xVal, float64(i+1))
		population.Evolve()
		population.Mutate(mutationP)
		population.EvalFitness(function)
		fitness, gene := population[0].Fitness, population[0].Genes
		if fitness < bestFitness {
			bestFitness = fitness
			bestGenes = gene
		}

		fitnessHistory = append(fitnessHistory, bestFitness)
		bar.Increment()
	}
	bar.Finish()

	fmt.Println("Best GA fitness:", bestFitness, ". Parameters:")
	for i:=0; i<len(bestGenes); i++ {
		fmt.Print(bestGenes[i], ", ")
	}
	fmt.Println()

	return xVal, fitnessHistory, bestFitness, bestGenes
}