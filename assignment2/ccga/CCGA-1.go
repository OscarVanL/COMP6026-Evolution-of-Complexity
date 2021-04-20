// Elements of roulette selection were adapted from: https://stackoverflow.com/questions/177271/roulette-selection-in-genetic-algorithms

package ccga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
)

const (
	CrossoverP = 0.6 // Probability of performing crossover
	W          = 5   // Scaling Window width
)

func Run(evaluations int, generations int, popSize int, N int, function f.Fitness, mutationP float32) ([]chart.BestFitness, float64, []uint16) {
	bestFitness := math.MaxFloat64
	var evals int
	var fMax float64 // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestCoevolution []uint16
	var worstFitnessHistory []float64 // Track worst fitness for each generation
	var bestFitnessHistory []chart.BestFitness

	// Initialise CCGA-1's population
	species := InitSpecies(N, popSize, time.Now().Unix())
	species.InitCoevolutions()
	species.EvalFitness(function, 0)
	species.SortFitness()
	fitness, _ := species.GetBestFitness()
	bestFitnessHistory = append(bestFitnessHistory, chart.BestFitness{X: evals, Fitness: fitness})
	fMax, _ = species.GetWorstFitness() // Set initial value of f'max

	if evaluations != 0 {
		// Run CCGA for N function evaluations
		for evals < evaluations {
			species.doGeneration(function, mutationP, 0, &evals, &fMax, &bestFitness, &bestCoevolution, &bestFitnessHistory, &worstFitnessHistory)
		}
	} else if generations != 0 {
		// Run CCGA for N generations
		for gen := 0; gen < generations; gen++ {
			species.doGeneration(function, mutationP, gen, &evals, &fMax, &bestFitness, &bestCoevolution, &bestFitnessHistory, &worstFitnessHistory)
		}
	}
	return bestFitnessHistory, bestFitness, bestCoevolution
}

func (spec Species) doGeneration(function f.Fitness, mutationP float32, gen int, evals *int, fMax *float64, bestFitness *float64, bestCoevolution *[]uint16, bestFitnessHistory *[]chart.BestFitness, worstFitnessHistory *[]float64) {
	// Coevolve in a round-robin fashion for each subpopulation
	for s := 0; s < len(spec); s++ {
		subpop := spec[s]
		// Coevolve a subpopulation with the best genes from each species
		subpop.CoevolveRoulette(CrossoverP, spec, function)
		// Mutate each individual's genes
		subpop.Mutate(mutationP)
		// Re-evaluate subpopulation member's fitnesses
		*evals += subpop.EvalFitness(function, *fMax)
		// Sort sub-population by fittest (smallest) to least fit (largest)
		subpop.SortFitness()
		// Finds individual with best fitness & genes in this subpopulation
		bestGenFitness, bestGenCoevolution := subpop.GetBestFitness()
		// Todo: Perhaps this worst fitness needs to be evaluated per generation, not per round-robin coevolution stage
		worstGenFitness, _ := subpop.GetWorstFitness()

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
}

// InitCoevolutions creates initial subpopulations by coevolving with random individuals from each other species.
func (spec Species) InitCoevolutions() {
	// Evaluate each species
	for s := 0; s < len(spec); s++ {
		species := spec[s]
		// Evaluate each individual in the species
		for i := 0; i < len(species); i++ {
			individual := species[i]

			// Combining individual with random individuals from each other species to calculate initial fitness
			tmpPop := make([]uint16, len(spec))
			for N := 0; N < len(spec); N++ {
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

func (subpop Population) Mutate(MutationP float32) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(subpop)-1)



	// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
	for i := 1; i < len(subpop); i++ {
		so := rand.NewSource(time.Now().UnixNano())
		r := rand.New(so)

		go func(i int, r *rand.Rand) {

			//mutatedCoevolution := pop[s][i].coevolution
			for g := 0; g < len(subpop[i].Coevolution); g++ {
				// Mutate each of the 16 bits in the individual's uint16 gene
				for b := 0; b < 16; b++ {
					// P probability of mutation
					if r.Float32() < MutationP {
						// Perform bit-flip
						if common.HasBit(subpop[i].Coevolution[g], uint(b)) {
							subpop[i].Coevolution[g] = common.ClearBit(subpop[i].Coevolution[g], uint(b))
						} else {
							subpop[i].Coevolution[g] = common.SetBit(subpop[i].Coevolution[g], uint(b))
						}
					}
				}
			}

			//Update individual's own mutated gene too
			subpop[i].Gene = subpop[i].Coevolution[subpop[i].SpeciesId]

			waitGroup.Done()
		} (i, r)
	}

	waitGroup.Wait()
}

func (subpop Population) CoevolveRoulette(crossoverP float32, spec Species, fitness f.Fitness) {
	var waitGroup sync.WaitGroup
	NGenes := len(subpop[0].Coevolution)
	waitGroup.Add(NGenes)

	subpop.RouletteSetup()

	// Do Roulette to pick each gene (N) contained in the coevolution
	for N := 0; N < NGenes; N++ {
		so := rand.NewSource(time.Now().UnixNano())
		r := rand.New(so)

		go func(N int, r *rand.Rand) {

			// Index starts at 1 to skip the most fit individual. Elitist strategy preserving fittest individual from each subspecies.
			for i := 1; i < len(subpop); i++ {
				// Two cases for updating Coevolutions:
				//	1. We're updating the subpop member's own gene:
				//		-> TwoPointCrossover with its existing gene & roulette-selected gene from the same subpopulation
				//  2. We're picking genes for the coevolution from other subpopulations:
				//      -> Select current best subcomponents

				if subpop[i].SpeciesId == N {
					// Coevolution Case 1

					// Whether to use crossover (otherwise, do nothing)
					if r.Float32() < crossoverP {
						offspringA, offspringB := common.TwoPointCrossover(subpop[i].Gene, subpop.RouletteSelection(r).Gene)

						// Pick best offspring
						subpop[i].Coevolution[N] = offspringA
						fitnessA := fitness(subpop[i].Coevolution)
						subpop[i].Coevolution[N] = offspringB
						fitnessB := fitness(subpop[i].Coevolution)
						if fitnessA > fitnessB {
							subpop[i].Coevolution[N] = offspringB
						} else {
							subpop[i].Coevolution[N] = offspringA
						}
					}

				} else {
					// Coevolution Case 2
					subpop[i].Coevolution[N] = spec[N][0].Gene
				}
			}
			waitGroup.Done()
		}(N, r)

	}

	waitGroup.Wait()
}



// RouletteSetup calculates population selection probabilities from ScaledFitness scores, required before using RouletteSelection.
// Adapted from: https://stackoverflow.com/a/177278/6008271
func (subpop Population) RouletteSetup() {
	// Calculate total fitness for subpopulation
	var fitnessSum float64
	for i := 0; i < len(subpop); i++ {
		fitnessSum += subpop[i].ScaledFitness
	}

	// Calculate fitness proportionate probability
	var accumulatedProbability float64
	for i := 0; i < len(subpop); i++ {
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
	for p := 0; p < len(subpop); p++ {
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
func (spec Species) EvalFitness(fitness f.Fitness, fMax float64) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(spec))

	for s:=0; s<len(spec); s++ {
		species := spec[s]

		// Evaluate each individual's fitness
		go func(s int) {
			for i:=0; i<len(species); i++ {
				// Calculate fitness while applying fMax scaling window
				spec[s][i].Fitness = fitness(spec[s][i].Coevolution)
				spec[s][i].ScaledFitness = math.Abs(fMax - spec[s][i].Fitness)
			}

			waitGroup.Done()
		} (s)
	}

	// Wait until all individuals have had fitness evaluated
	waitGroup.Wait()
}

// EvalFitness checks the fitness of each coevolved individual's genes and updates its Fitness & ScaledFitness scores.
// Return number of fitness evaluations
func (subpop Population) EvalFitness(fitness f.Fitness, fMax float64) int {
	for i := 0; i < len(subpop); i++ {
		// Calculate fitness while applying fMax scaling window
		subpop[i].Fitness = fitness(subpop[i].Coevolution)
		subpop[i].ScaledFitness = math.Abs(fMax - subpop[i].Fitness)
	}

	return len(subpop)
}

func (spec Species) SortFitness() {
	// Sort the population's individuals by fittest (smallest) to least fit (largest)
	for s:=0; s<len(spec); s++ {
		sort.Slice(spec[s], func(i, j int) bool {
			return spec[s][i].Fitness <  spec[s][j].Fitness
		})
	}
}

func (subpop Population) SortFitness() {
	// Sort the subpopulation's individuals by fittest (smallest) to least fit (largest)
	sort.Slice(subpop, func(i, j int) bool {
		return subpop[i].Fitness < subpop[j].Fitness
	})
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

// GetBestFitness finds the individual with the fittest (smallest) fitness score amongst the subpopulation
// Note: Run this after SortFitness so fitnesses are pre-sorted
func (subpop Population) GetBestFitness() (float64, []uint16) {
	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	if subpop[0].Fitness < bestFitness {
		bestFitness = subpop[0].Fitness
		bestCoevolution = subpop[0].Coevolution
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

// GetWorstFitness finds the individual with the least fit score amongst the subpopulation
// Note: Run this after SortFitness so fitnesses are pre-sorted
func (subpop Population) GetWorstFitness() (float64, []uint16) {
	popSize := len(subpop)
	worstFitness := 0.0
	var worstCoevolution []uint16

	for s := 0; s < len(subpop); s++ {
		if subpop[popSize-1].Fitness > worstFitness {
			worstFitness = subpop[popSize-1].Fitness
			worstCoevolution = subpop[popSize-1].Coevolution
		}
	}
	return worstFitness, worstCoevolution
}