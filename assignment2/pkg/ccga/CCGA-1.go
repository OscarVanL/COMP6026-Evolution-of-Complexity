package ccga

import (
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"math"
	"math/rand"
	"sort"
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

func (pop Species) Mutate(MutationP float64) {
	// N.B. Bit manipulation inner-functions are taken from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271

	//Checks if bit is set as position n
	hasBit := func(n uint16, pos uint) bool {
		val := n & (1 << pos)
		return val > 0
	}

	// Sets bit at index pos to 1
	setBit := func(n uint16, pos uint) uint16 {
		n |= (1 << pos)
		return n
	}

	// Sets bit at index pos to 0
	clearBit := func(n uint16, pos uint) uint16 {
		mask := ^(1 << pos)
		nTemp := int(n)
		nTemp &= mask
		return uint16(nTemp)
	}

	for s:=0; s<len(pop); s++ {
		species := pop[s]

		for i:=0; i<len(species); i++ {
			individual := species[i]

			mutatedGene := individual.Gene
			// Mutate each of the 16 bits in the individual's uint16 gene
			for b:=0; b<16; b++ {
				// P probability of mutation
				if rand.Float64() < MutationP {
					// Perform bit-flip
					if hasBit(mutatedGene, uint(b)) {
						mutatedGene = clearBit(mutatedGene, uint(b))
					} else {
						mutatedGene = setBit(mutatedGene, uint(b))
					}
				}
			}
			// Replace individual's old gene with mutated one
			pop[s][i].Gene = mutatedGene
		}
	}

}

// Coevolve does crossover for each individual with the best other individuals and mutates the coevolved offspring.
func (pop Species) Coevolve(MutationP float64) {
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
			pop[s][i] = individual
		}
	}
}

// EvalFitness calculates the fitness score for each coevolved individual
func (pop Species) EvalFitness(fitness f.Fitness) {
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