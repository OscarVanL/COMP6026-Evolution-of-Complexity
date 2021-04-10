package ga

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

// Mutate performs bit-flip mutation on each of the individual's genes
func (pop Population) Mutate(MutationP float32) {
	// Todo: To implement eletist strategy, skip mutation for the best individual
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i:=0; i<len(pop); i++ {
		individual := pop[i]

		mutatedGenes := individual.Genes

		// Mutate each of the individual's genes
		for g:=0; g<len(mutatedGenes); g++ {
			// Mutate each of the 16 bits in the gene
			for b:=0; b<16; b++ {
				// P probability of mutation
				if r.Float32() < MutationP {
					// Perform bit-flip
					if common.HasBit(mutatedGenes[g], uint(b)) {
						mutatedGenes[g] = common.ClearBit(mutatedGenes[g], uint(b))
					} else {
						mutatedGenes[g] = common.SetBit(mutatedGenes[g], uint(b))
					}
				}
			}
		}
		// Replace individual's old genes with mutated ones
		pop[i].Genes = mutatedGenes
	}
}

func (pop Population) Crossover() {
	// Todo: To implement eletist strategy, skip mutation for the best individual
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// Get individual with best Fitness
	// Todo: Replace this with another approach, eg: roulette wheel
	bestIndividual := pop[0]

	for i:=0; i<len(pop); i++ {
		// Do crossover for each gene
		for g:=0; g<len(bestIndividual.Genes); g++ {
			if r.Float32() < CrossoverP {
				// Perform two-point crossover
				offspringA, offspringB := common.TwoPointCrossover(pop[i].Genes[g], bestIndividual.Genes[g])

				// Randomly select one of the offspring to use
				if r.Intn(2) == 0 {
					pop[i].Genes[g] = offspringA
				} else {
					pop[i].Genes[g] = offspringB
				}
			}
		}
	}
}

func (pop Population) EvalFitness(fitness f.Fitness, fMax float64) {
	for i:=0; i<len(pop); i++ {
		// Calculate individual's Fitness
		pop[i].Fitness = fitness(pop[i].Genes)
		pop[i].ScaledFitness = math.Abs(fMax - pop[i].Fitness)
	}

	pop.SortFitness()

	//fmt.Println("fMax:", fMax, "BestIndividual fitness:", pop[0].Fitness, ", scaledFitness", pop[0].ScaledFitness)
}

func (pop Population) SortFitness() {
	// Sort the populations individuals by fittest (smallest) to least fit (largest)
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fitness < pop[j].Fitness
	})
}

func (pop Population) SortScaledFitness() {
	// In the case of the scaled fitness, larger fitness values are considered more fit.
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].ScaledFitness > pop[j].ScaledFitness
	})
}

