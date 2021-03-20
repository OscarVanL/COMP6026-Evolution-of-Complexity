package evolution

import (
	"math/rand"
	"time"
)


// Individual holds the gene and the evaluated fitness
type Individual struct {
	Gene uint16
	fitness float64
}

// Population keeps the individuals (Individual) that make a subpopulation
type Population []Individual

// Species keeps the subpopulations (Population) that are combined to make the target function parameters
type Species []Population

// GenerateSpecies will generate SpeciesN number of species, each of PopSize population
func GenerateSpecies(SpeciesN int, PopSize int) Species {
	rand.Seed(time.Now().Unix())

	species := make(Species, SpeciesN)
	for i:=0; i<SpeciesN; i++ {
		species[i] = GeneratePopulation(PopSize)
	}
	return species
}

// GeneratePopulation will generate a random population of PopSize individuals
func GeneratePopulation(PopSize int) Population {
	pop := make(Population, PopSize)
	for i:=0; i<PopSize; i++ {
		pop[i] = Individual{uint16(rand.Int()), 0.0}
	}
	return pop
}
