package ccga

import (
	"math/rand"
	"time"
)

//// CCGAGenerations > CCGASpecies > CCGAPopulation > CCGAIndividual

// CCGAGenerations tracks the species across iterations of CCGA
type Generations []Species

// CCGASpecies keeps the subpopulations (CCGAPopulation) that are combined to make the target function parameters
type Species []Population

// CCGAPopulation keeps the individuals (CCGAIndividual) that make a subpopulation
type Population []Individual

// CCGAIndividual holds the gene, evaluated fitness, and combination of individuals that achieved this fitness
type Individual struct {
	Gene uint16
	fitness float64
	coevolution []uint16  // Combination of individuals that produced this fitness outcome
}

// InitSpecies will generate SpeciesN number of species, each of PopSize population
func InitSpecies(SpeciesN int, PopSize int) Species {
	rand.Seed(time.Now().Unix())

	species := make(Species, SpeciesN)
	for i:=0; i<SpeciesN; i++ {
		species[i] = InitPopulation(PopSize)
	}
	return species
}

// InitPopulation will generate a random population of PopSize individuals
func InitPopulation(PopSize int) Population {
	pop := make(Population, PopSize)
	for i:=0; i<PopSize; i++ {
		pop[i] = Individual{uint16(rand.Int()), 0.0, nil}
	}
	return pop
}