package main

import (
	"flag"
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/ccga"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/common"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/ga"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"github.com/cheggaaa/pb"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"time"
)

const iterations = 100000
const popSize = 100
const W = 5  // Scaling Window width

type Algorithm int

const (
	Rastrigin Algorithm = iota
	Schwefel
	Griewangk
	Ackley
	Rosenbrock
)

var cpuprofile = flag.String("cpuprofile", "assignment2.prof", "write cpu profile to file")


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

	var results []chart.EvolutionResults

	fmt.Println("Benchmarking Rastrigin Function")
	results = append(results, DoGeneticAlgorithms(Rastrigin))
	fmt.Println("Benchmarking Schwefel Function")
	results = append(results, DoGeneticAlgorithms(Schwefel))
	fmt.Println("Benchmarking Griewangk Function")
	results = append(results, DoGeneticAlgorithms(Griewangk))
	fmt.Println("Benchmarking Ackley Function")
	results = append(results, DoGeneticAlgorithms(Ackley))

	fmt.Println("Creating Charts")
	chart.PlotResults(results)

}

func DoGeneticAlgorithms(algo Algorithm) chart.EvolutionResults {
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



	return chart.EvolutionResults{label,
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
	var fMax float64  // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestCoevolution []uint16
	var xVal []float64
	var worstFitnessHistory []float64  // Track worst fitness for each generation
	var bestFitnessHistory []float64  // Track best overall fitness across all generations

	species := ccga.InitSpecies(N, popSize)
	species.InitCoevolutions()
	species.EvalFitness(function, 0)
	fMax, _ = species.GetWorstFitness(popSize)  // Set initial value of f'max

	for i:=0; i<iterations; i++ {
		// Todo: Track the number of function evaluations, not GA iterations
		xVal = append(xVal, float64(i+1))  // Evolution iteration for X-Axis
		// Coevolves individuals with the best (mutated) genes from each species
		species.CoevolveRoulette()
		// Mutates each individual's genes
		species.Mutate(mutationP)
		// Re-evaluates bestGenFitness
		species.EvalFitness(function, fMax)
		// Finds individual with best bestGenFitness & genes in this generation
		bestGenFitness, bestGenCoevolution := species.GetBestFitness()
		worstGenFitness, _ := species.GetWorstFitness(popSize)

		if bestGenFitness < bestFitness {
			bestFitness = bestGenFitness
			bestCoevolution = bestGenCoevolution
		}
		worstFitnessHistory = append(worstFitnessHistory, worstGenFitness)
		bestFitnessHistory = append(bestFitnessHistory, bestFitness)

		fMax = common.CalculateFMax(worstFitnessHistory, W)

		bar.Increment()
	}
	bar.Finish()

	fmt.Println("Best Coevolution fitness:", bestFitness, ". Parameters:")
	for i:=0; i<len(bestCoevolution); i++ {
		fmt.Print(bestCoevolution[i], ", ")
	}
	fmt.Println()

	return xVal, bestFitnessHistory, bestFitness, bestCoevolution
}

func GA(N int, function f.Fitness, mutationP float32) ([]float64, []float64, float64, []uint16) {
	bar := pb.New(iterations)
	bar.SetRefreshRate(time.Second)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	bestFitness := math.MaxFloat64
	var fMax float64  // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestGenes []uint16
	var xVal []float64
	var worstFitnessHistory []float64  // Track worst fitness for each generation
	var bestFitnessHistory []float64  // Track best overall fitness across all generations

	population := ga.InitPopulation(N, popSize, time.Now().Unix())
	population.EvalFitness(function, 0)
	fMax = population[len(population)-1].Fitness  // Set initial value of f'max
	fmt.Println("First fMax:", fMax)

	for i:=0; i<iterations; i++ {
		// Todo: Track the number of function evaluations, not GA iterations
		xVal = append(xVal, float64(i+1))
		//population.SortScaledFitness()
		population.Crossover()
		population.Mutate(mutationP)
		population.EvalFitness(function, fMax)
		bestGenFitness, bestGenGene := population[0].Fitness, population[0].Genes
		worstGenFitness := population[len(population)-1].Fitness
		if bestGenFitness < bestFitness {
			bestFitness = bestGenFitness
			bestGenes = bestGenGene
		}

		worstFitnessHistory = append(worstFitnessHistory, worstGenFitness)
		bestFitnessHistory = append(bestFitnessHistory, bestFitness)

		fMax = common.CalculateFMax(worstFitnessHistory, W)

		bar.Increment()
	}
	bar.Finish()

	fmt.Println("Best GA fitness:", bestFitness, ". Parameters:")
	for i:=0; i<len(bestGenes); i++ {
		fmt.Print(bestGenes[i], ", ")
	}
	fmt.Println()

	return xVal, bestFitnessHistory, bestFitness, bestGenes
}