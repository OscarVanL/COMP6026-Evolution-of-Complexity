package ccga

import (
	"math/rand"
)

//// CCGASpecies > CCGAPopulation > CCGAIndividual

// Species keeps the subpopulations (CCGAPopulation) that are combined to make the target function parameters
type Species []Population

// Population keeps the individuals (CCGAIndividual) that make a subpopulation
type Population []Individual

// Individual holds the gene, evaluated fitness, and combination of individuals that achieved this fitness
type Individual struct {
	SpeciesId         int
	Gene              uint16
	Fitness           float64
	ScaledFitness     float64
	SelectProbability float64  // Probability of selection in Roulette Wheel
	Coevolution       []uint16 // Combination of individuals that produced this fitness outcome
}

// InitSpecies will generate SpeciesN number of species, each of PopSize population
func InitSpecies(SpeciesN int, PopSize int, seed int64) Species {
	s := rand.NewSource(seed)
	r := rand.New(s)

	species := make(Species, SpeciesN)
	// Repeat process for N "genes" (species)
	for s := 0; s < SpeciesN; s++ {
		// Randomly generate a species
		pop := make(Population, PopSize)
		for i := 0; i < PopSize; i++ {
			pop[i] = Individual{s, uint16(r.Int()), 0.0, 0.0, 0.0, nil}
		}
		species[s] = pop
	}
	return species
}
