package main

import (
	"flag"
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/ccga"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/ga"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"log"
	"os"
	"runtime/pprof"
)

const iterations = 1000000
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
		prof, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(prof)
		defer pprof.StopCPUProfile()
	}

	var results []chart.EvolutionResults

	fmt.Println("Benchmarking Rastrigin Function")
	results = append(results, RunGAs(Rastrigin))
	fmt.Println("Benchmarking Schwefel Function")
	results = append(results, RunGAs(Schwefel))
	fmt.Println("Benchmarking Griewangk Function")
	results = append(results, RunGAs(Griewangk))
	fmt.Println("Benchmarking Ackley Function")
	results = append(results, RunGAs(Ackley))

	fmt.Println("Creating Charts")
	chart.PlotResults(results)

}

func RunGAs(algo Algorithm) chart.EvolutionResults {
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
	YValsGA, BestFitnessGA, BestAssignmentGA := ga.Run(iterations, popSize, N, function, mutationP, W)
	fmt.Println("Starting CCGA-1 with params: N:", N, "MutationP:", mutationP)
	YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA := ccga.Run(iterations, popSize, N, function, mutationP, W)

	fmt.Println("Best GA fitness:", BestFitnessGA, ". Parameters:")
	for i:=0; i<len(BestAssignmentGA); i++ {
		fmt.Print(BestAssignmentGA[i], ", ")
	}
	fmt.Println()

	fmt.Println("Best CCGA fitness:", BestFitnessCCGA, ". Parameters:")
	for i:=0; i<len(BestAssignmentCCGA); i++ {
		fmt.Print(BestAssignmentCCGA[i], ", ")
	}
	fmt.Println()

	return chart.EvolutionResults{
		Label: label,
		Iterations: iterations,
		CCGAFitnessHistory: YValsCCGA,
		BestFitnessCCGA: BestFitnessCCGA,
		GAFitnessHistory: YValsGA,
		BestFitnessGA: BestFitnessGA,
	}
}



