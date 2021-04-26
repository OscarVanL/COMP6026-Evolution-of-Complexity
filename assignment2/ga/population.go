package ga

import (
	"math/rand"
)

//// GAPopulation > GAIndividual

// Population keeps the individuals (GAIndividual) that make a the GA's population
type Population []Individual

// Individual holds the genes and Fitness for an individual
type Individual struct {
	Genes             []uint16
	Fitness           float64
	ScaledFitness     float64
	SelectProbability float64 // Probability of selection in Roulette Wheel
}

// InitPopulation will generate the initial population for the standard GA, each with N genes
func InitPopulation(N int, PopSize int, seed int64) Population {
	rand.Seed(seed)

	pop := make(Population, PopSize)
	for i := 0; i < PopSize; i++ {
		genes := make([]uint16, N)
		for n := 0; n < N; n++ {
			genes[n] = uint16(rand.Int())
		}
		pop[i] = Individual{genes, 0.0, 0.0, 0.0}

	}
	return pop
}
