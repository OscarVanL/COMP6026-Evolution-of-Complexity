package ccga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Crossover probability
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
			individual.coevolution = tmpPop
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

			for i:=0; i<len(pop[s]); i++ {
				individual := pop[s][i]
				//mutatedCoevolution := pop[s][i].coevolution
				for g:=0; g<len(pop[s][i].coevolution); g++ {
					// Mutate each of the 16 bits in the individual's uint16 gene
					for b:=0; b<16; b++ {
						// P probability of mutation
						if r.Float32() < MutationP {
							// Perform bit-flip
							if common.HasBit(individual.coevolution[g], uint(b)) {
								individual.coevolution[g] = common.ClearBit(individual.coevolution[g], uint(b))
							} else {
								individual.coevolution[g] = common.SetBit(individual.coevolution[g], uint(b))
							}
						}
					}
				}

				//Update individual's own mutated gene too
				individual.Gene = individual.coevolution[individual.SpeciesId]
				pop[s][i] = individual
			}
			mutate <- empty{}
		} (s)

	}

	for i:=0; i<len(pop); i++ { <- mutate }
}


// Coevolve does crossover for each individual with the best other individuals and mutates the coevolved offspring.
func (pop Species) Coevolve() {
	type empty struct{}
	coevolve := make(chan empty, len(pop))

	// Evolve each species
	for s:=0; s<len(pop); s++ {
		species := pop[s]

		go func(s int) {
			so := rand.NewSource(time.Now().UnixNano())
			r := rand.New(so)

			// Evolve each individual in the species
			for i:=0; i<len(species); i++ {

				individual := species[i]

				// Combine individual with best individuals from other species, if CrossoverP is met
				for N:=0; N<len(pop); N++ {
					// Whether to perform crossover
					if r.Float32() < CrossoverP {
						// Perform two-point crossover with best gene and individual's existing gene
						offspringA, offspringB := common.TwoPointCrossover(pop[N][0].Gene, individual.coevolution[N])

						// Randomly select one of the offspring to use
						if r.Intn(2) == 0 {
							individual.coevolution[N] = offspringA
						} else {
							individual.coevolution[N] = offspringB
						}
					}
				}
				pop[s][i] = individual
			}
			coevolve <- empty{}
		} (s)

	}

	for i:=0; i<len(pop); i++ { <- coevolve }
}

// EvalFitness calculates the fitness score for each coevolved individual. Sorts populations from fittest to least fit.
func (pop Species) EvalFitness(fitness f.Fitness) {
	type empty struct{}
	eval := make(chan empty, len(pop))

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
			eval <- empty{}
		} (s)
	}

	// Wait until all individuals have had fitness evaluated
	for i:=0; i<len(pop); i++ { <- eval }
}

// GetBestFitness finds the individual with the fittest (smallest) fitness score amongst the species
// Note: Run this after EvalFitness so that fitnesses are already sorted
func (pop Species) GetBestFitness() (float64, []uint16) {
	bestFitness := math.MaxFloat64
	var bestCoevolution []uint16
	for s:=0; s<len(pop); s++ {
		if pop[s][0].fitness < bestFitness {
			bestFitness = pop[s][0].fitness
			bestCoevolution = pop[s][0].coevolution
		}
	}
	return bestFitness, bestCoevolution
}