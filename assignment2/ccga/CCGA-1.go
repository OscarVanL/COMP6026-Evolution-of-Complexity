// Elements of roulette selection were adapted from: https://stackoverflow.com/questions/177271/roulette-selection-in-genetic-algorithms



package ccga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// CrossoverP is probability of performing crossover
const CrossoverP = 0.6

// InitCoevolutions creates initial subpopulations by coevolving with random individuals from each other species.
func (pop Species) InitCoevolutions() {
	// Evaluate each species
	for s :=0; s<len(pop); s++ {
		species := pop[s]
		// Evaluate each individual in the species
		for i := 0; i<len(species); i++ {
			individual := species[i]

			// Combining individual with random individuals from each other species to calculate initial fitness
			tmpPop := make([]uint16, len(pop))
			for N:=0; N <len(pop); N++ {
				if s == N {
					// Keep the evaluated individual for its gene
					tmpPop[N] = individual.Gene
				} else {
					// Use a random individual from the other species genes
					tmpPop[N] = pop[N][rand.Intn(len(pop[N]))].Gene
				}
			}
			individual.Coevolution = tmpPop
			pop[s][i] = individual
		}
	}
}

func (pop Species) Mutate(MutationP float32) {
	type empty struct{}
	mutate := make(chan empty, len(pop))

	for s:=0; s<len(pop); s++ {
		go func(s int) {
			so := rand.NewSource(time.Now().UnixNano())
			r := rand.New(so)

			// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
			for i:=1; i<len(pop[s]); i++ {
				individual := pop[s][i]
				//mutatedCoevolution := pop[s][i].coevolution
				for g:=0; g<len(pop[s][i].Coevolution); g++ {
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
				pop[s][i] = individual
			}
			mutate <- empty{}
		} (s)

	}

	for i:=0; i<len(pop); i++ { <- mutate }
}

func (pop Species) CoevolveRoulette(crossoverP float32) {

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(pop))

	// Do Roulette to pick specific gene (N) for each coevolution
	for N :=0; N <len(pop); N++ {
		so := rand.NewSource(time.Now().UnixNano())
		r := rand.New(so)

		// Todo: Do all roulette setups simultaneously as goroutine
		pop[N].RouletteSetup()

		go func(N int) {
			// Update Coevolution for each species on this gene
			for sp := 0; sp < len(pop); sp++ {

				// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
				for i := 1; i < len(pop[sp]); i++ {

					// Two cases for updating Coevolutions:
					//	1. We're updating the subpop member's own gene:
					//		-> TwoPointCrossover with its existing gene & roulette-selected gene from the same subpopulation
					//  2. We're picking genes for the coevolution from other subpopulations:
					//      -> Select gene using roulette selection from these subpopulations. No crossover.

					if pop[sp][i].SpeciesId == N {
						// Coevolution Case 1

						// Whether to use crossover (otherwise, do nothing)
						if r.Float32() < crossoverP {
							offspringA, offspringB := common.TwoPointCrossover(pop[sp][i].Gene, pop[sp].RouletteSelection(r).Gene)

							// Randomly select one of the offspring to use
							if r.Intn(2) == 0 {
								pop[sp][i].Coevolution[N] = offspringA
							} else {
								pop[sp][i].Coevolution[N] = offspringB
							}
						}

					} else {
						// Coevolution Case 2
						// Use roulette selection: adapted from https://stackoverflow.com/a/177278/6008271
						pop[sp][i].Coevolution[N] = pop[sp].RouletteSelection(r).Gene
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
func (pop Species) EvalFitness(fitness f.Fitness, fMax float64) int {
	type empty struct{}
	eval := make(chan empty, len(pop))

	for s:=0; s<len(pop); s++ {
		species := pop[s]

		// Evaluate each individual's fitness
		go func(s int) {
			for i:=0; i<len(species); i++ {
				// Calculate fitness while applying fMax scaling window
				pop[s][i].Fitness = fitness(pop[s][i].Coevolution)
				pop[s][i].ScaledFitness = math.Abs(fMax - pop[s][i].Fitness)
			}

			eval <- empty{}
		} (s)
	}

	// Wait until all individuals have had fitness evaluated
	for i:=0; i<len(pop); i++ { <- eval }
	return len(pop) * len(pop[0])
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
func (pop Species) GetBestFitness() (float64, []uint16) {
	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	for s:=0; s<len(pop); s++ {
		if pop[s][0].Fitness < bestFitness {
			bestFitness = pop[s][0].Fitness
			bestCoevolution = pop[s][0].Coevolution
		}
	}
	return bestFitness, bestCoevolution
}

// GetWorstFitness finds the individual with the least fit score amongst the species
// Note: Run this after SortFitness so fitnesses are pre-sorted
func (pop Species) GetWorstFitness() (float64, []uint16) {
	popSize := len(pop[0])
	worstFitness := 0.0
	var worstCoevolution []uint16

	for s:=0; s<len(pop); s++ {
		if pop[s][popSize-1].Fitness > worstFitness {
			worstFitness = pop[s][popSize-1].Fitness
			worstCoevolution = pop[s][popSize-1].Coevolution
		}
	}
	return worstFitness, worstCoevolution
}