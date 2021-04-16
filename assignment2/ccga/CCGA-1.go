// Elements of roulette selection were adapted from: https://stackoverflow.com/questions/177271/roulette-selection-in-genetic-algorithms



package ccga

import (
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/cheggaaa/pb"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const(
	CrossoverP = 0.6  // Probability of performing crossover
	W = 5  // Scaling Window width
)

func Run(evaluations int, generations int, popSize int, N int, function f.Fitness, mutationP float32) ([]chart.BestFitness, float64, []uint16) {
	bestFitness := math.MaxFloat64
	var evals int
	var fMax float64  // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestCoevolution []uint16
	var worstFitnessHistory []float64  // Track worst fitness for each generation
	var bestFitnessHistory []chart.BestFitness

	// Initialise CCGA-1's population
	species := InitSpecies(N, popSize, time.Now().Unix())
	species.InitCoevolutions()
	evals += species.EvalFitness(function, 0)
	species.SortFitness()
	fitness, _ := species.GetBestFitness()
	bestFitnessHistory = append(bestFitnessHistory, chart.BestFitness{X: evals, Fitness: fitness})
	fMax, _ = species.GetWorstFitness()  // Set initial value of f'max

	if evaluations != 0 {
		// Run CCGA for N function evaluations
		fmt.Println("Running CCGA-1 for", evaluations, "function evaluations.")
		bar := pb.New(evaluations)
		bar.SetRefreshRate(time.Second)
		bar.ShowTimeLeft = true
		bar.ShowSpeed = true
		bar.Start()

		for evals<evaluations {
			species.doGeneration(function, mutationP, 0, &evals, &fMax, &bestFitness, &bestCoevolution, &bestFitnessHistory, &worstFitnessHistory)
			bar.Set(evals)
		}
		bar.Finish()
	} else if generations != 0 {
		// Run CCGA for N generations
		fmt.Println("Running CCGA-1 for", generations, "generations.")
		bar := pb.New(generations)
		bar.SetRefreshRate(time.Second)
		bar.ShowTimeLeft = true
		bar.ShowSpeed = true
		bar.Start()

		for gen:=0; gen<generations; gen++ {
			species.doGeneration(function, mutationP, gen, &evals, &fMax, &bestFitness, &bestCoevolution, &bestFitnessHistory, &worstFitnessHistory)
			bar.Increment()
		}
		bar.Finish()
	}

	return bestFitnessHistory, bestFitness, bestCoevolution
}

func (spec Species) doGeneration(function f.Fitness, mutationP float32, gen int, evals *int, fMax *float64, bestFitness *float64, bestCoevolution *[]uint16, bestFitnessHistory *[]chart.BestFitness, worstFitnessHistory *[]float64) {
	// Coevolves individuals with the best (mutated) genes from each species
	spec.CoevolveRoulette(CrossoverP)
	// Mutates each individual's genes
	spec.Mutate(mutationP)
	// Re-evaluates individual fitnesses
	*evals += spec.EvalFitness(function, *fMax)
	// Sort the population's individuals by fittest (smallest) to least fit (largest)
	spec.SortFitness()
	// Finds individual with best fitness & genes in this generation
	bestGenFitness, bestGenCoevolution := spec.GetBestFitness()
	worstGenFitness, _ := spec.GetWorstFitness()

	if bestGenFitness < *bestFitness {
		*bestFitness = bestGenFitness
		*bestCoevolution = bestGenCoevolution
		if gen != 0 {
			*bestFitnessHistory = append(*bestFitnessHistory, chart.BestFitness{X: gen, Fitness: bestGenFitness})
		} else {
			*bestFitnessHistory = append(*bestFitnessHistory, chart.BestFitness{X: *evals, Fitness: bestGenFitness})
		}
	}
	*worstFitnessHistory = append(*worstFitnessHistory, worstGenFitness)
	*fMax = common.CalculateFMax(*worstFitnessHistory, W)
}

// InitCoevolutions creates initial subpopulations by coevolving with random individuals from each other species.
func (spec Species) InitCoevolutions() {
	// Evaluate each species
	for s :=0; s<len(spec); s++ {
		species := spec[s]
		// Evaluate each individual in the species
		for i := 0; i<len(species); i++ {
			individual := species[i]

			// Combining individual with random individuals from each other species to calculate initial fitness
			tmpPop := make([]uint16, len(spec))
			for N:=0; N <len(spec); N++ {
				if s == N {
					// Keep the evaluated individual for its gene
					tmpPop[N] = individual.Gene
				} else {
					// Use a random individual from the other species genes
					tmpPop[N] = spec[N][rand.Intn(len(spec[N]))].Gene
				}
			}
			individual.Coevolution = tmpPop
			spec[s][i] = individual
		}
	}
}

func (spec Species) Mutate(MutationP float32) {
	type empty struct{}
	mutate := make(chan empty, len(spec))

	for s:=0; s<len(spec); s++ {
		go func(s int) {
			so := rand.NewSource(time.Now().UnixNano())
			r := rand.New(so)

			// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
			for i:=1; i<len(spec[s]); i++ {
				individual := spec[s][i]
				//mutatedCoevolution := pop[s][i].coevolution
				for g:=0; g<len(spec[s][i].Coevolution); g++ {
					// Mutate each of the 16 bits in the individual's uint16 gene
					for b:=0; b<16; b++ {
						// P probability of mutation
						if r.Float32() < MutationP {
							// Perform bit-flip
							if common.HasBit(individual.Coevolution[g], uint(b)) {
								individual.Coevolution[g] = common.ClearBit(individual.Coevolution[g], uint(b))
							} else {
								individual.Coevolution[g] = common.SetBit(individual.Coevolution[g], uint(b))
							}
						}
					}
				}

				//Update individual's own mutated gene too
				individual.Gene = individual.Coevolution[individual.SpeciesId]
				spec[s][i] = individual
			}
			mutate <- empty{}
		} (s)

	}

	for i:=0; i<len(spec); i++ { <- mutate }
}

func (spec Species) CoevolveRoulette(crossoverP float32) {

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(spec))

	// Do Roulette to pick specific gene (N) for each coevolution
	for N :=0; N <len(spec); N++ {
		so := rand.NewSource(time.Now().UnixNano())
		r := rand.New(so)

		// Todo: Do all roulette setups simultaneously as goroutine
		spec[N].RouletteSetup()

		go func(N int) {
			// Update Coevolution for each species on this gene
			for sp := 0; sp < len(spec); sp++ {

				// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
				for i := 1; i < len(spec[sp]); i++ {

					// Two cases for updating Coevolutions:
					//	1. We're updating the subpop member's own gene:
					//		-> TwoPointCrossover with its existing gene & roulette-selected gene from the same subpopulation
					//  2. We're picking genes for the coevolution from other subpopulations:
					//      -> Select gene using roulette selection from these subpopulations. No crossover.

					if spec[sp][i].SpeciesId == N {
						// Coevolution Case 1

						// Whether to use crossover (otherwise, do nothing)
						if r.Float32() < crossoverP {
							offspringA, offspringB := common.TwoPointCrossover(spec[sp][i].Gene, spec[sp].RouletteSelection(r).Gene)

							// Randomly select one of the offspring to use
							if r.Intn(2) == 0 {
								spec[sp][i].Coevolution[N] = offspringA
							} else {
								spec[sp][i].Coevolution[N] = offspringB
							}
						}

					} else {
						// Coevolution Case 2
						// Use roulette selection: adapted from https://stackoverflow.com/a/177278/6008271
						spec[sp][i].Coevolution[N] = spec[sp].RouletteSelection(r).Gene
					}

				}
			}
			waitGroup.Done()
		}(N)
	}

	waitGroup.Wait()
}

// RouletteSetup calculates population selection probabilities from ScaledFitness scores, required before using RouletteSelection.
// Adapted from: https://stackoverflow.com/a/177278/6008271
func (subpop Population) RouletteSetup() {
	// Calculate total fitness for subpopulation
	var fitnessSum float64
	for i:=0; i<len(subpop); i++ {
		fitnessSum += subpop[i].ScaledFitness
	}

	// Calculate fitness proportionate probability
	var accumulatedProbability float64
	for i:=0; i<len(subpop); i++ {
		subpop[i].SelectProbability = subpop[i].ScaledFitness / fitnessSum
		subpop[i].SelectProbability += accumulatedProbability
		accumulatedProbability = subpop[i].SelectProbability
	}
}

// RouletteSelection uses a roulette approach to apply higher selective pressure for individuals with better fitness
// Adapted from: https://stackoverflow.com/a/177278/6008271
func (subpop Population) RouletteSelection(r *rand.Rand) Individual {
	// Todo: Use binary search here, instead of linear search.
	number := r.Float64()
	for p:=0; p<len(subpop); p++ {
		if p == 0 {
			// First entry on roulette wheel, range 0.0 - Select Probability
			if number < subpop[p].SelectProbability {
				return subpop[p]
			}
		} else {
			// SelectProbability greater than last individual, but within this individual's probability range
			if number > subpop[p-1].SelectProbability && number < subpop[p].SelectProbability {
				return subpop[p]
			}
		}
	}
	return subpop[0]
}

// EvalFitness checks the fitness of each coevolved individual's genes and updates its Fitness & ScaledFitness scores.
// Return number of fitness evaluations
func (spec Species) EvalFitness(fitness f.Fitness, fMax float64) int {
	type empty struct{}
	eval := make(chan empty, len(spec))

	for s:=0; s<len(spec); s++ {
		species := spec[s]

		// Evaluate each individual's fitness
		go func(s int) {
			for i:=0; i<len(species); i++ {
				// Calculate fitness while applying fMax scaling window
				spec[s][i].Fitness = fitness(spec[s][i].Coevolution)
				spec[s][i].ScaledFitness = math.Abs(fMax - spec[s][i].Fitness)
			}

			eval <- empty{}
		} (s)
	}

	// Wait until all individuals have had fitness evaluated
	for i:=0; i<len(spec); i++ { <- eval }
	return len(spec) * len(spec[0])
}

func (spec Species) SortFitness() {
	// Sort the population's individuals by fittest (smallest) to least fit (largest)
	for s:=0; s<len(spec); s++ {
		sort.Slice(spec[s], func(i, j int) bool {
			return spec[s][i].Fitness <  spec[s][j].Fitness
		})
	}

}

// GetBestFitness finds the individual with the fittest (smallest) fitness score amongst the species
// Note: Run this after SortFitness so fitnesses are pre-sorted
func (spec Species) GetBestFitness() (float64, []uint16) {
	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	for s:=0; s<len(spec); s++ {
		if spec[s][0].Fitness < bestFitness {
			bestFitness = spec[s][0].Fitness
			bestCoevolution = spec[s][0].Coevolution
		}
	}
	return bestFitness, bestCoevolution
}

// GetWorstFitness finds the individual with the least fit score amongst the species
// Note: Run this after SortFitness so fitnesses are pre-sorted
func (spec Species) GetWorstFitness() (float64, []uint16) {
	popSize := len(spec[0])
	worstFitness := 0.0
	var worstCoevolution []uint16

	for s:=0; s<len(spec); s++ {
		if spec[s][popSize-1].Fitness > worstFitness {
			worstFitness = spec[s][popSize-1].Fitness
			worstCoevolution = spec[s][popSize-1].Coevolution
		}
	}
	return worstFitness, worstCoevolution
}