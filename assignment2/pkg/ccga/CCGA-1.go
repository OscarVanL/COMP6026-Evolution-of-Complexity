package ccga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
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
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

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

		// Todo: Run in parallel for each species using goroutine
		for i:=0; i<len(species); i++ {
			individual := species[i]

			mutatedCoevolution := individual.coevolution
			for g:=0; g<len(mutatedCoevolution); g++ {
				// Mutate each of the 16 bits in the individual's uint16 gene
				for b:=0; b<16; b++ {
					// P probability of mutation
					if r.Float32() < MutationP {
						// Perform bit-flip
						if hasBit(mutatedCoevolution[g], uint(b)) {
							mutatedCoevolution[g] = clearBit(mutatedCoevolution[g], uint(b))
						} else {
							mutatedCoevolution[g] = setBit(mutatedCoevolution[g], uint(b))
						}
					}
				}
			}

			// Replace individuals's coevolved gene with mutated one
			pop[s][i].coevolution = mutatedCoevolution

			//Update individual's own mutated gene too
			pop[s][i].Gene = mutatedCoevolution[pop[s][i].SpeciesId]
		}
	}

}

// Mutate performs bit-flip mutation on each of the individual's genes
//func (pop Species) Mutate(MutationP float32) {
//	// N.B. Bit manipulation inner-functions are taken from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271
//
//	//Checks if bit is set as position n
//	hasBit := func(n uint16, pos uint) bool {
//		val := n & (1 << pos)
//		return val > 0
//	}
//
//	// Sets bit at index pos to 1
//	setBit := func(n uint16, pos uint) uint16 {
//		n |= (1 << pos)
//		return n
//	}
//
//	// Sets bit at index pos to 0
//	clearBit := func(n uint16, pos uint) uint16 {
//		mask := ^(1 << pos)
//		nTemp := int(n)
//		nTemp &= mask
//		return uint16(nTemp)
//	}
//
//
//	for s:=0; s<len(pop); s++ {
//		species := pop[s]
//
//
//
//		m := make(chan Individual)
//		// Todo: Run in parallel for each species using goroutine
//		go func(pop Population, m chan Individual) {
//			for i := 0; i < len(pop); i++ {
//				ind := pop[i]
//
//				for g := 0; g < len(ind.coevolution); g++ {
//					// Mutate each of the 16 bits in the individual's uint16 gene
//					for b := 0; b < 16; b++ {
//						// P probability of mutation
//						if rand.Float32() < MutationP {
//							// Perform bit-flip
//							if hasBit(ind.coevolution[g], uint(b)) {
//								ind.coevolution[g] = clearBit(ind.coevolution[g], uint(b))
//							} else {
//								ind.coevolution[g] = setBit(ind.coevolution[g], uint(b))
//							}
//						}
//					}
//				}
//
//				ind.Gene = ind.coevolution[ind.SpeciesId]
//
//				m <- ind
//
//			}
//			close(m)
//		} (species, m)
//
//		tmp := make([]Individual, len(species))
//		for i:=0; i<len(species); i++ {
//			tmp[i] = <- m
//		}
//
//		pop[s] = tmp
//	}
//
//}

// Coevolve does crossover for each individual with the best other individuals and mutates the coevolved offspring.
func (pop Species) Coevolve() {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// Evolve each species
	for s:=0; s<len(pop); s++ {
		species := pop[s]

		// Evolve each individual in the species
		for i:=0; i<len(species); i++ {
			individual := species[i]

			// Combine individual with best individuals from other species, if CrossoverP is met
			tmpPop := make([]uint16, len(pop))
			for N:=0; N<len(pop); N++ {
				// Whether to perform crossover
				if r.Float32() < CrossoverP {
					// Perform two-point crossover with best gene and individual's existing gene
					offspringA, offspringB := evolution.TwoPointCrossover(pop[N][0].Gene, individual.coevolution[N])

					// Randomly select one of the offspring to use
					if r.Intn(2) == 0 {
						tmpPop[N] = offspringA
					} else {
						tmpPop[N] = offspringB
					}
				} else {
					// Keep the existing gene, no crossover
					tmpPop[N] = individual.coevolution[N]
				}
			}
			individual.coevolution = tmpPop
			pop[s][i] = individual
		}
	}
}

// EvalFitness calculates the fitness score for each coevolved individual. Sorts populations from fittest to least fit.
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