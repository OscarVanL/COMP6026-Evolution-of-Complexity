package ga

import (
	"math/rand"
	"time"
)


//// GAGenerations > GAPopulation > GAIndividual

// Generation tracks the species across iterations of the GA
type Generations []Population

// GAPopulation keeps the individuals (GAIndividual) that make a the GA's population
type Population []Individual

// GAIndividual holds the genes and fitness for an individual
type Individual struct {
	Genes []uint16
	fitness float64
}

// InitPopulation will generate the initial population for the standard GA, each with N genes
func InitPopulation(N int, PopSize int) Population {
	rand.Seed(time.Now().Unix())

	pop := make(Population, PopSize)
	for i:=0; i<PopSize; i++ {
		genes := make([]uint16, N)
		for n:=0; n<N; n++ {
			genes[n] = uint16(rand.Int())
		}
		pop[i] = Individual{genes, 0.0}

	}
	return pop
}