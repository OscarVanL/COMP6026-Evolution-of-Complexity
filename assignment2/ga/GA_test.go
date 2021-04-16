package ga

import (
	"fmt"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

// TestPopulation_Mutate_ZeroProbability ensures no mutations occur when MutationP is 0.0
func TestPopulation_Mutate_ZeroProbability(t *testing.T) {
	input := Population{
		Individual{[]uint16{0x0000, 0xFFFF}, 0, 0, 0.0},
		Individual{[]uint16{0xFFFF, 0x0000}, 0, 0, 0.0},
	}

	input.Mutate(0.0)
	assert.Equal(t, uint16(0xFFFF), input[1].Genes[0], "Mutation with mutationP=0 should not change any bits")
	assert.Equal(t, uint16(0x0000), input[1].Genes[1], "Mutation with mutationP=0 should not change any bits")
}

// TestPopulation_Mutate_OneProbability ensures all bits are mutated when MutationP is 1.0
func TestPopulation_Mutate_OneProbability(t *testing.T) {
	// To make test deterministic use 100% probability, meaning every bit should be flipped
	input := Population{
		Individual{[]uint16{0x0000, 0xFFFF}, 0, 0, 0.0},
		Individual{[]uint16{0xFFFF, 0x0000}, 0, 0, 0.0},
	}

	input.Mutate(1.0)
	assert.Equal(t, uint16(0x0000), input[1].Genes[0], "Mutation with mutationP=1 should flip every bit")
	assert.Equal(t, uint16(0xFFFF), input[1].Genes[1], "Mutation with mutationP=1 should flip every bit")
}

// TestPopulation_Mutate_Elitist ensures elitist strategy is applied by skipping mutation on 0-index individual
// When pre-sorted, the 0-index individual is the one with highest fitness, so is preserved without mutation.
func TestPopulation_Mutate_Elitist(t *testing.T) {
	// To make test deterministic use 100% probability, meaning every bit should be flipped
	input := Population{
		Individual{[]uint16{0x0000, 0xFFFF}, 0, 0, 0.0},
		Individual{[]uint16{0xFFFF, 0x0000}, 0, 0, 0.0},
	}

	input.Mutate(1.0)
	assert.Equal(t, uint16(0x0000), input[0].Genes[0], "Mutation should follow elitist strategy")
}

// TestPopulation_Crossover checks that crossover operations are completed as expected
func TestPopulation_Crossover_OneProbability(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Population{
		Individual{[]uint16{0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF}, 0, 0, 1.0},
		Individual{[]uint16{0x0000, 0x0000, 0x0000, 0x0000}, 0, 0, 0.0},
	}

	// Crossover with 100% probability
	input.Crossover(1.0, f.Rosenbrock)

	assert.Equal(t, []uint16{0x0000, 0xFFFF, 0xFFFF, 0x0000}, input[1].Genes, "Genes were not crossed over as expected")
}

// TestPopulation_Crossover checks that genes are kept constant with 0 crossover probability
func TestPopulation_Crossover_ZeroProbability(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Population{
		Individual{[]uint16{0x0000, 0xFFFF}, 0, 0, 1.0},
		Individual{[]uint16{0xFFFF, 0x0000}, 0, 0, 0.0},
	}

	// Crossover with 100% probability
	input.Crossover(0.0, f.Rosenbrock)

	assert.Equal(t, uint16(0xFFFF), input[1].Genes[0], "Genes were modified when they shouldn't")
	assert.Equal(t, uint16(0x0000), input[1].Genes[1], "Genes were modified when they shouldn't")
}

// TestPopulation_Crossover_Elitist ensures elitist strategy is applied by skipping crossover on the 0-index individual
// When pre-sorted, the 0-index individual is the one with highest fitness, and is preserved without genetic functions.
func TestPopulation_Crossover_Elitist(t *testing.T) {
	// To make test deterministic, create individuals with certain roulette selection probability
	input := Population{
		Individual{[]uint16{0x0000, 0xFFFF}, 0, 0, 0.0},
		Individual{[]uint16{0xFFFF, 0x0000}, 0, 0, 1.0},
	}

	// Crossover with 100% probability
	input.Crossover(1.0, f.Rosenbrock)
	
	assert.Equal(t, uint16(0x0000), input[0].Genes[0], "Genes for 0-index individual should remain unchanged")
}

// TestPopulation_RouletteSetup tests if the roulette setup assigns SelectProbabilities correctly.
func TestPopulation_RouletteSetup(t *testing.T) {
	// A population with pre-defined ScaledFitness values.
	input := Population{
		Individual{[]uint16{0}, 0, 100, 0},
		Individual{[]uint16{1}, 0, 700, 0},
		Individual{[]uint16{2}, 0, 100, 0},
		Individual{[]uint16{3}, 0, 100, 0},
	}
	expectedSelectProbability := []float64{0.1, 0.8, 0.9, 1.0}

	input.RouletteSetup()

	for i:=0; i<len(input); i++ {
		fmt.Println(input[i].SelectProbability)
		assert.InDelta(t, expectedSelectProbability[i], input[i].SelectProbability, 0.01, "Roulette Selection Probabilities were not assigned as expected")
	}
}

// TestPopulation_RouletteSelection tests if the roulette selection selects individuals in the correct proportion, relative to their SelectProbability.
func TestPopulation_RouletteSelection(t *testing.T) {
	// A population with pre-defined selection probabilities
	input := Population{
		Individual{[]uint16{0}, 0, 0, 0.1},
		Individual{[]uint16{1}, 0, 0, 0.8},
		Individual{[]uint16{2}, 0, 0, 0.9},
		Individual{[]uint16{3}, 0, 0, 1.0},
	}
	expectedSelectionRatio := []float64{0.1, 0.7, 0.1, 0.1}

	s := rand.NewSource(0)
	r := rand.New(s)
	// Run roulette selection a large number of times
	results := make(map[uint16]float64)
	for i:=0; i<100000; i++ {
		results[input.RouletteSelection(r).Genes[0]] += 1.0
	}

	// Calculate ratios of each selection by roulette wheel
	for k := range results {
		results[k] /= 100000
	}

	// Check probabilities of selection are correct, within a 5% margin of error to accommodate randomness
	for i:=0; i<len(results); i++ {
		assert.InDelta(t, expectedSelectionRatio[i], results[uint16(i)], 0.05, "Roulette Selection did not sample proportionately")
	}
}


func TestPopulation_EvalFitness(t *testing.T) {
	input := Population{
		Individual{[]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0, 0, 0},
	}
	input.EvalFitness(f.Schwefel, 3000)
	assert.InDelta(t, 2392.9928386744673, input[0].Fitness, 0.01, "Fitness was not calculated properly")
	assert.InDelta(t, 607.0071613255327, input[0].ScaledFitness, 0.01, "Fitness was not calculated properly")
}

func TestPopulation_SortFitness(t *testing.T) {
	input := Population{
		Individual{[]uint16{37889}, 2345, 0, 0},
		Individual{[]uint16{16832}, 1234, 0, 0},
		Individual{[]uint16{51315}, 23, 0, 0},
		Individual{[]uint16{42338}, 64, 0, 0},
		Individual{[]uint16{11594}, -1, 0, 0},
	}

	expected := Population{
		Individual{[]uint16{11594}, -1, 0, 0},
		Individual{[]uint16{51315}, 23, 0, 0},
		Individual{[]uint16{42338}, 64, 0, 0},
		Individual{[]uint16{16832}, 1234, 0, 0},
		Individual{[]uint16{37889}, 2345, 0, 0},
	}

	input.SortFitness()

	assert.Equal(t, expected, input, "Population was not correctly sorted")

}