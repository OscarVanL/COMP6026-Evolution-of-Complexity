package cmd

import (
	"errors"
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/ccga"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/ga"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/result"
	"github.com/cheggaaa/pb"
	"github.com/spf13/cobra"
	"github.com/stretchr/stew/slice"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "coevolve",
	Short: "Cooperative Coevolution implementation for COMP6026",
	Long:  `An implementation of DOI 10.1007/3-540-58484-6_269 'A Cooperative Coevolutionary Approach to Function Optimization'.'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("evaluations") && cmd.Flags().Changed("generations") {
			return errors.New("cannot set evaluations and generations, pick one")
		} else if !(cmd.Flags().Changed("evaluations") || cmd.Flags().Changed("generations")) {
			return errors.New("must set either evaluations or generations number")
		}
		if cmd.Flags().Changed("cpuprofile") && filepath.Ext(cpuprofile) != ".prof" {
			return errors.New("cpuprofile file argument must end with .prof extension")
		}
		if !(slice.Contains(algorithms, "ga") || slice.Contains(algorithms, "ccga") || slice.Contains(algorithms, "ccgahc")) {
			return errors.New("at least one algorithm must be configured with -a ga,ccga,ccgahc")
		}

		fmt.Println("Starting with algorithms:", algorithms)
		Start()
		return nil
	},
}

var algorithms []string
var evaluations int
var generations int
var popSize int
var repetitions int
var cpuprofile string
var output string

func init() {
	rootCmd.Flags().StringSliceVarP(&algorithms, "algorithms", "a", []string{"ga", "ccga", "ccgahc"}, "Which algorithms to compare (ga,ccga,ccgahc)")
	rootCmd.Flags().IntVarP(&evaluations, "evaluations", "e", 0, "Function evaluation limit")
	rootCmd.Flags().IntVarP(&generations, "generations", "g", 0, "Generations limit")
	rootCmd.Flags().IntVarP(&popSize, "population", "p", 100, "Population size")
	rootCmd.Flags().IntVarP(&repetitions, "repetitions", "r", 50, "Number of times to repeat experiment")
	rootCmd.Flags().StringVar(&cpuprofile, "cpuprofile", "", "Profile CPU usage to file (eg: assignment2.prof)")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Fitness Figure output file (eg: charts.html)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Start() {
	if cpuprofile != "" {
		fmt.Println("Profiling CPU usage to file:", cpuprofile)
		prof, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(prof)
		defer pprof.StopCPUProfile()
	}

	// Store results for each iteration of the function
	var results [][]chart.EvolutionResults

	fmt.Println("Benchmarking Rastrigin Function...")
	results = append(results, RunGAs("rastrigin"))
	fmt.Println("Benchmarking Schwefel Function...")
	results = append(results, RunGAs("schwefel"))
	fmt.Println("Benchmarking Griewangk Function...")
	results = append(results, RunGAs("griewangk"))
	fmt.Println("Benchmarking Ackley Function...")
	results = append(results, RunGAs("ackley"))
	//fmt.Println("Benchmarking Rosenbrock Function...")
	//results = append(results, RunGAs("rosenbrock"))


	if output != "" {
		fmt.Println("Creating Charts")
		chart.PlotResults(output, results)
		fmt.Println("Writing Results JSON")
		result.WriteResults(output, results)
	}
}

func RunGAs(function string) []chart.EvolutionResults {
	Params, err := f.GetParams(function)
	if err != nil {
		panic(err)
	}

	var results []chart.EvolutionResults

	bar := pb.New(repetitions)
	bar.SetRefreshRate(time.Second)
	bar.ShowTimeLeft = true
	bar.Start()
	var waitGroup sync.WaitGroup
	waitGroup.Add(repetitions)

	for i := 0; i < repetitions; i++ {
		// Run each separate GA repetition in its own goroutine
		go func() {

			var YValsGA, YValsCCGA, YValsCCGAHC []chart.BestFitness
			var BestFitnessGA, BestFitnessCCGA, BestFitnessCCGAHC float64
			var BestAssignmentGA, BestAssignmentCCGA, BestAssignmentCCGAHC []uint16

			// Start Standard Genetic Algorithm
			if slice.Contains(algorithms, "ga") {
				YValsGA, BestFitnessGA, BestAssignmentGA = ga.Run(evaluations, generations, popSize, Params.N, Params.Function, Params.MutationP)
			}
			// Start CCGA
			if slice.Contains(algorithms, "ccga") {
				YValsCCGA, BestFitnessCCGA, BestAssignmentCCGA = ccga.Run(false, evaluations, generations, popSize, Params.N, Params.Function, Params.MutationP)
			}
			// Start CCGAHC
			if slice.Contains(algorithms, "ccgahc") {
				YValsCCGAHC, BestFitnessCCGAHC, BestAssignmentCCGAHC = ccga.Run(true, evaluations, generations, popSize, Params.N, Params.Function, Params.MutationP)
			}

			if evaluations != 0 {
				results = append(results, chart.EvolutionResults{
					Title:              Params.Label,
					XLabel:             "function\nevals",
					Iterations:         evaluations,

					GAFitnessHistory:   YValsGA,
					BestFitnessGA:      BestFitnessGA,
					BestAssignmentGA:   BestAssignmentGA,

					CCGAFitnessHistory: YValsCCGA,
					BestFitnessCCGA:    BestFitnessCCGA,
					BestAssignmentCCGA: BestAssignmentCCGA,

					CCGAHCFitnessHistory: YValsCCGAHC,
					BestFitnessCCGAHC:    BestFitnessCCGAHC,
					BestAssignmentCCGAHC: BestAssignmentCCGAHC,
				})
			} else {
				results = append(results, chart.EvolutionResults{
					Title:              Params.Label,
					XLabel:             "gens",
					Iterations:         generations,

					CCGAFitnessHistory: YValsCCGA,
					BestFitnessCCGA:    BestFitnessCCGA,
					BestAssignmentCCGA: BestAssignmentCCGA,

					GAFitnessHistory:   YValsGA,
					BestFitnessGA:      BestFitnessGA,
					BestAssignmentGA:   BestAssignmentGA,

					CCGAHCFitnessHistory: YValsCCGAHC,
					BestFitnessCCGAHC:    BestFitnessCCGAHC,
					BestAssignmentCCGAHC: BestAssignmentCCGAHC,
				})
			}
			bar.Increment()
			waitGroup.Done()
		}()

		// Sleep 50ms between starting of each goroutine to avoid similar random numbers between runs
		// (as each goroutine's rand is seeded with the unix time)
		time.Sleep(50 * time.Millisecond)
	}
	waitGroup.Wait()
	bar.Finish()

	return results
}
