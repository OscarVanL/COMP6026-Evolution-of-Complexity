package evolution

import (
	"math"
	"math/rand"
	"sort"
)

// Crossover probability
const CrossoverP = 0.6

// Fitness functions that can be used to evaluate several individuals
type Fitness func(x []uint16) float64

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
			individual.coevolution = tmpPop
			pop[s][i] = individual
		}
	}
}

// EvolveSpecies does crossover for each individual with the best other individuals and mutates the coevolved offspring.
func (pop Species) EvolveSpecies(MutationP float64) {
	// Evolve each species
	for s:=0; s<len(pop); s++ {
		species := pop[s]

		// Evolve each individual in the species
		for i:=0; i<len(species); i++ {
			individual := species[i]

			// Combine individual with best individuals from other species, if CrossoverP is met
			tmpPop := make([]uint16, len(pop))
			for N:=0; N<len(pop); N++ {
				if s == N {
					tmpPop[N] = individual.Gene
				} else {
					//Todo: Two-point crossover should be used here
					if rand.Float64() < CrossoverP {
						// Get the species' best individual's gene (already sorted)
						tmpPop[N] = pop[N][0].Gene
					} else {
						// Keep the existing gene, no crossover
						tmpPop[N] = individual.coevolution[N]
					}

				}
			}
			individual.coevolution = tmpPop
			//Todo: Do mutation
			individual.MutateIndividual(MutationP)

			pop[s][i] = individual
		}
	}
}

// EvalFitness calculates the fitness score for each coevolved individual
func (pop Species) EvalFitness(fitness Fitness) {
	type empty struct{}
	spec := make(chan empty, len(pop))

	for s:=0; s<len(pop); s++ {
		species := pop[s]

		// Evaluate each individual's fitness
		go func(s int) {
			for i:=0; i<len(species); i++ {
				individual := species[i]
				individual.fitness = fitness(individual.coevolution)
				pop[s][i] = individual
			}

			// Sort the population's individuals by fittest (smallest) to least fit (largest)
			sort.Slice(pop[s], func(i, j int) bool {
				return pop[s][i].fitness < pop[s][j].fitness
			})
			spec <- empty{}
		} (s)
	}

	// Wait until all individuals have had fitness evaluated
	for i:=0; i<len(pop); i++ { <- spec }
}


// GetBestFitness finds the individual with the fittest (smallest) fitness score amongst the species
func (pop Species) GetBestFitness() (float64, []uint16) {
	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	for s:=0; s<len(pop); s++ {
		for i:=0; i<len(pop[s]); i++ {
			if pop[s][i].fitness < bestFitness {
				bestFitness = pop[s][i].fitness
				bestCoevolution = pop[s][i].coevolution
			}
		}
	}
	return bestFitness, bestCoevolution
}