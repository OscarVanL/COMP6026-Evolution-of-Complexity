package evolution

import (
	"math/rand"
	"time"
)

//// Generation > Species > Population > Individual

// Generation tracks the species across iterations of the algorithm
type Generations []Species

// Species keeps the subpopulations (Population) that are combined to make the target function parameters
type Species []Population

// Population keeps the individuals (Individual) that make a subpopulation
type Population []Individual

// Individual holds the gene, evaluated fitness, and combination of individuals that achieved this fitness
type Individual struct {
	Gene uint16
	fitness float64
	coevolution []uint16  // Combination of individuals that produced this fitness outcome
}




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
		pop[i] = Individual{uint16(rand.Int()), 0.0, nil}
	}
	return pop
}

func (i Individual) MutateIndividual(MutationP float64) Individual {

	// N.B. Bit manipulation inner-functions are taken from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271

	//Checks if bit is set as position n
	//hasBit := func(n uint16, pos uint) bool {
	//	val := n & (1 << pos)
	//	return val > 0
	//}
	//
	//// Sets bit at index pos to 1
	//setBit := func(n uint16, pos uint) uint16 {
	//	n |= (1 << pos)
	//	return n
	//}
	//
	//// Sets bit at index pos to 0
	//clearBit := func(n uint16, pos uint) uint16 {
	//	mask := ^(1 << pos)
	//	nTemp := int(n)
	//	nTemp &= mask
	//	return uint16(nTemp)
	//}





	return i
}