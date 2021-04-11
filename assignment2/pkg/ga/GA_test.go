package ga

import (
	"fmt"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestPopulation_Mutate(t *testing.T) {

}

func TestPopulation_Crossover(t *testing.T) {

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
	for k, _ := range results {
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