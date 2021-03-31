package ga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"math/rand"
	"sort"
)

// Crossover probability
const CrossoverP = 0.6

func (pop Population) Mutate(MutationP float64) {
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


	for i:=0; i<len(pop); i++ {
		individual := pop[i]

		mutatedGenes := individual.Genes

		// Mutate each of the individual's genes
		for g:=0; g<len(mutatedGenes); g++ {
			// Mutate each of the 16 bits in the gene
			for b:=0; b<16; b++ {
				// P probability of mutation
				if rand.Float64() < MutationP {
					// Perform bit-flip
					if hasBit(mutatedGenes[g], uint(b)) {
						mutatedGenes[g] = clearBit(mutatedGenes[g], uint(b))
					} else {
						mutatedGenes[g] = setBit(mutatedGenes[g], uint(b))
					}
				}
			}
		}
		// Replace individual's old genes with mutated ones
		pop[i].Genes = mutatedGenes
	}
}

func (pop Population) Evolve() {
	// Get individual with best Fitness
	bestIndividual := pop[0]
	for i:=0; i<len(pop); i++ {
		// Do crossover for each gene
		for g:=0; g<len(bestIndividual.Genes); g++ {
			if rand.Float64() < CrossoverP {
				// Perform two-point crossover
				offspringA, offspringB := evolution.TwoPointCrossover(pop[i].Genes[g], bestIndividual.Genes[g])

				// Randomly select one of the offspring to use
				if rand.Intn(2) == 0 {
					pop[i].Genes[g] = offspringA
				} else {
					pop[i].Genes[g] = offspringB
				}
			}
		}
	}
}

func (pop Population) EvalFitness(fitness f.Fitness) {
	for i:=0; i<len(pop); i++ {
		// Calculate individual's Fitness
		pop[i].Fitness = fitness(pop[i].Genes)
	}

	// Sort the populations individuals by fittest (smallest) to least fit (largest)
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fitness < pop[j].Fitness
	})
}

