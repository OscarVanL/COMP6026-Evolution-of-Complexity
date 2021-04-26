package ccga

import (
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSpecies_InitCoevolutions(t *testing.T) {
	input := Species{
		Population{
			Individual{0, 1234, 0, 0, 0.0, []uint16{}},
			Individual{0, 5678, 0, 0, 0.0, []uint16{}},
		},
		Population{
			Individual{1, 2345, 0, 0, 0.0, []uint16{}},
			Individual{1, 6789, 0, 0, 0.0, []uint16{}},
		},
	}
	input.InitCoevolutions()

	// Coevolutions are initialised randomly, but should keep the individual's own gene.
	assert.Equal(t, uint16(1234), input[0][0].Coevolution[0], "Coevolution initialisation should keep individual's own gene")
	assert.Equal(t, uint16(5678), input[0][1].Coevolution[0], "Coevolution initialisation should keep individual's own gene")
	assert.Equal(t, uint16(2345), input[1][0].Coevolution[1], "Coevolution initialisation should keep individual's own gene")
	assert.Equal(t, uint16(6789), input[1][1].Coevolution[1], "Coevolution initialisation should keep individual's own gene")
}

func TestSpecies_Mutate_OneProbability(t *testing.T) {
	input := Population{
		Individual{0, 0x0000, 0, 0, 0.0, []uint16{0xFFFF, 0x0000}},
		Individual{0, 0xFFFF, 0, 0, 0.0, []uint16{0x0000, 0xFFFF}},
	}
	input.Mutate(0.0)

	assert.Equal(t, uint16(0x0000), input[1].Gene, "Mutation with mutationP=1 should change all bits")
	assert.Equal(t, uint16(0x0000), input[1].Coevolution[0], "Mutation with mutationP=1 should change all bits")
}

func TestSpecies_Mutate_ZeroProbability(t *testing.T) {
	input := Population{
		Individual{0, 0xFFFF, 0, 0, 0.0, []uint16{0xFFFF, 0x0000}},
		Individual{0, 0x0000, 0, 0, 0.0, []uint16{0x0000, 0x0000}},
	}
	input.Mutate(0.0)

	assert.Equal(t, uint16(0x0000), input[1].Gene, "Mutation with mutationP=0 should not change any bits")
	assert.Equal(t, uint16(0x0000), input[1].Coevolution[0], "Mutation with mutationP=0 should not change any bits")
}

func TestSpecies_Mutate_Elitist(t *testing.T) {
	input := Population{
		Individual{0, 0xFFFF, 0, 0, 0, []uint16{0x0000, 0xFFFF}},
	}

	input.Mutate(1.0)

	assert.Equal(t, uint16(0xFFFF), input[0].Gene, "Mutation should follow elitist strategy")
	assert.Equal(t, uint16(0x0000), input[0].Coevolution[0], "Mutation should follow elitist strategy")
}

func TestSpecies_CoevolveRoulette_OneProbability(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Species{
		Population{
			Individual{0, 0xFFFF, 0, 1, 1.0, []uint16{0x0000, 0xFFFF}},
			Individual{0, 0x0000, 0, 0, 0.0, []uint16{0xFFFF, 0x0000}},
		},
		Population{
			Individual{1, 0x0000, 0, 1, 1.0, []uint16{0xFFFF, 0x0000}},
			Individual{1, 0xFFFF, 0, 0, 0.0, []uint16{0x0000, 0xFFFF}},
		},
	}
	input[0].CoevolveRoulette(1.0, input, f.TestFunc)
	input[1].CoevolveRoulette(1.0, input, f.TestFunc)

	expectedGene1 := (input[0][1].Coevolution[0] == 0x0FF0) || (input[0][1].Coevolution[0] == 0xF00F)
	assert.True(t, expectedGene1, "Genes were not crossed over as expected")

	expectedGene2 := (input[1][1].Coevolution[1] == 0x0FF0) || (input[1][1].Coevolution[1] == 0xF00F)
	assert.True(t, expectedGene2, "Genes were not crossed over as expected")
}

func TestSpecies_CoevolveRoulette_ZeroProbability(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Species{
		Population{
			Individual{0, 0x0000, 0, 0, 0.0, []uint16{0x0000, 0xFFFF}},
			Individual{0, 0xFFFF, 0, 0, 1.0, []uint16{0xFFFF, 0xFFFF}},
		},
		Population{
			Individual{1, 0xFFFF, 0, 0, 0.0, []uint16{0xFFFF, 0x0000}},
			Individual{1, 0x0000, 0, 0, 1.0, []uint16{0x0000, 0xFFFF}},
		},
	}
	input[0].CoevolveRoulette(0.0, input, f.TestFunc)
	input[1].CoevolveRoulette(0.0, input, f.TestFunc)

	assert.Equal(t, uint16(0xFFFF), input[0][1].Coevolution[0], "Coevolved genes should not change when crossoverP is 0")
	assert.Equal(t, uint16(0xFFFF), input[1][1].Coevolution[1], "Coevolved genes should not change when crossoverP is 0")
}

// TestSpecies_CoevolveRoulette_Elitist ensures elitist strategy is applied by skipping crossover on the 0-index individual of each species
// When pre-sorted, the 0-index individual is the one with highest fitness, and is preserved without genetic functions.
func TestSpecies_CoevolveRoulette_Elitist(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Species{
		Population{
			Individual{0, 0x0000, 0, 0, 0.0, []uint16{0x0000, 0xFFFF}},
			Individual{0, 0xFFFF, 0, 0, 1.0, []uint16{0xFFFF, 0xFFFF}},
		},
		Population{
			Individual{1, 0xFFFF, 0, 0, 0.0, []uint16{0xFFFF, 0x0000}},
			Individual{1, 0x0000, 0, 0, 1.0, []uint16{0x0000, 0xFFFF}},
		},
	}

	// Crossover with 100% probability
	input[0].CoevolveRoulette(1.0, input, f.TestFunc)
	input[1].CoevolveRoulette(1.0, input, f.TestFunc)

	assert.Equal(t, uint16(0x0000), input[0][0].Coevolution[0], "Genes for 0-index individual should remain unchanged")
	assert.Equal(t, uint16(0xFFFF), input[0][0].Coevolution[1], "Genes for 0-index individual should remain unchanged")
	assert.Equal(t, uint16(0xFFFF), input[1][0].Coevolution[0], "Genes for 0-index individual should remain unchanged")
	assert.Equal(t, uint16(0x0000), input[1][0].Coevolution[1], "Genes for 0-index individual should remain unchanged")
}

// TestPopulation_RouletteSetup tests if the roulette setup assigns SelectProbabilities correctly.
func TestPopulation_RouletteSetup(t *testing.T) {
	// A population with pre-defined ScaledFitness values.
	input := Population{
		Individual{0, 0, 0, 100, 0, nil},
		Individual{0, 1, 0, 700, 0, nil},
		Individual{0, 2, 0, 100, 0, nil},
		Individual{0, 3, 0, 100, 0, nil},
	}
	expectedSelectProbability := []float64{0.1, 0.8, 0.9, 1.0}
	input.RouletteSetup()

	for i := 0; i < len(input); i++ {
		assert.InDelta(t, expectedSelectProbability[i], input[i].SelectProbability, 0.01, "Roulette Selection Probabilities were not assigned as expected")
	}
}

// TestPopulation_RouletteSelection tests if the roulette selection selects individuals in the correct proportion, relative to their SelectProbability.
func TestPopulation_RouletteSelection(t *testing.T) {
	// A population with pre-defined selection probabilities
	input := Population{
		Individual{0, 0, 0, 0, 0.1, nil},
		Individual{0, 1, 0, 0, 0.8, nil},
		Individual{0, 2, 0, 0, 0.9, nil},
		Individual{0, 3, 0, 0, 1.0, nil},
	}
	expectedSelectionRatio := []float64{0.1, 0.7, 0.1, 0.1}

	s := rand.NewSource(0)
	r := rand.New(s)
	// Run roulette selection a large number of times
	results := make(map[uint16]float64)
	for i := 0; i < 100000; i++ {
		results[input.RouletteSelection(r).Gene] += 1.0
	}

	// Calculate ratios of each selection by roulette wheel
	for k := range results {
		results[k] /= 100000
	}

	// Check probabilities of selection are correct, within a 5% margin of error to accommodate randomness
	for i := 0; i < len(results); i++ {
		assert.InDelta(t, expectedSelectionRatio[i], results[uint16(i)], 0.05, "Roulette Selection did not sample proportionately")
	}
}

func TestSpecies_EvalFitness(t *testing.T) {
	input := Species{
		Population{
			Individual{0, 1, 0, 0, 0, []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		},
		Population{
			Individual{1, 2, 0, 0, 0, []uint16{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},
		},
	}

	input.EvalFitness(f.Schwefel, 3000)

	assert.InDelta(t, 2392.9928386744673, input[0][0].Fitness, 0.01, "Fitness was not calculated properly")
	assert.InDelta(t, 607.0071613255327, input[0][0].ScaledFitness, 0.01, "Fitness was not calculated properly")

	assert.InDelta(t, 2409.4667586600553, input[1][0].Fitness, 0.01, "Fitness was not calculated properly")
	assert.InDelta(t, 590.5332413399447, input[1][0].ScaledFitness, 0.01, "Fitness was not calculated properly")
}

func TestSpecies_SortFitness(t *testing.T) {
	input := Species{
		Population{
			Individual{0, 1234, 1000, 0, 0, []uint16{1111, 1111, 1111}},
			Individual{0, 5678, 500, 0, 0, []uint16{2222, 2222, 2222}},
			Individual{0, 9012, 250, 0, 0, []uint16{3333, 3333, 3333}},
			Individual{0, 3456, 2000, 0, 0, []uint16{4444, 4444, 4444}},
		},
		Population{
			Individual{1, 1234, 1, 0, 0, []uint16{1111, 1111, 1111}},
			Individual{1, 5678, 500, 0, 0, []uint16{2222, 2222, 2222}},
			Individual{1, 9012, 250, 0, 0, []uint16{3333, 3333, 333}},
			Individual{1, 3456, 1, 0, 0, []uint16{4444, 4444, 4444}},
		},
	}
	input.SortFitness()

	assert.Equal(t, Individual{0, 9012, 250, 0, 0, []uint16{3333, 3333, 3333}}, input[0][0], "The 0-index does not contain the most fit individual (smallest fitness)")
	assert.Equal(t, Individual{0, 3456, 2000, 0, 0, []uint16{4444, 4444, 4444}}, input[0][3], "The last index does not contain the last fit individual (largest fitness)")
	assert.True(t, (input[1][0].Gene == uint16(1234) || input[1][1].Gene == uint16(1234)) && (input[1][0].Gene == uint16(3456) || input[1][1].Gene == uint16(3456)), "Sorting may not work properly for individuals with identical fitness")
	assert.Equal(t, Individual{1, 5678, 500, 0, 0, []uint16{2222, 2222, 2222}}, input[1][3], "The last index does not contain the last fit individual (largest fitness)")
}

func TestSpecies_GetBestFitness(t *testing.T) {
	input := Species{
		Population{
			Individual{0, 1234, 500, 0, 0, []uint16{1111, 1111, 1111}},
			Individual{0, 5678, 1000, 0, 0, []uint16{2222, 2222, 2222}},
		},
		Population{
			Individual{1, 2345, 1000, 0, 0, []uint16{3333, 3333, 3333}},
			Individual{1, 6789, 2000, 0, 0, []uint16{4444, 4444, 4444}},
		},
		Population{
			Individual{2, 4567, 600, 0, 0, []uint16{5555, 5555, 5555}},
			Individual{2, 8901, 750, 0, 0, []uint16{6666, 6666, 6666}},
		},
	}

	fitness, coevolution := input.GetBestFitness()

	assert.Equal(t, input[0][0].Fitness, fitness, "Did not get best fitness")
	assert.Equal(t, input[0][0].Coevolution, coevolution, "Did not get coevolution associated with best fitness")
}

func TestSpecies_GetWorstFitness(t *testing.T) {
	input := Species{
		Population{
			Individual{0, 1234, 500, 0, 0, []uint16{1111, 1111, 1111}},
			Individual{0, 5678, 1000, 0, 0, []uint16{2222, 2222, 2222}},
		},
		Population{
			Individual{1, 2345, 1000, 0, 0, []uint16{3333, 3333, 3333}},
			Individual{1, 6789, 2000, 0, 0, []uint16{4444, 4444, 4444}},
		},
		Population{
			Individual{2, 4567, 600, 0, 0, []uint16{5555, 5555, 5555}},
			Individual{2, 8901, 750, 0, 0, []uint16{6666, 6666, 6666}},
		},
	}

	fitness, coevolution := input.GetWorstFitness()

	assert.Equal(t, input[1][1].Fitness, fitness, "Did not get worst fitness")
	assert.Equal(t, input[1][1].Coevolution, coevolution, "Did not get coevolution associated with worst fitness")
}
