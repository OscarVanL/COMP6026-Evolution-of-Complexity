package chart

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFillMissingPoints(t *testing.T) {
	iterations := 11
	fitnessHistory := []BestFitness{
		{X: 0, Fitness: 100},
		{X: 4, Fitness: 80},
		{X: 10, Fitness: 20},
	}

	yVals := fillMissingPoints(iterations, fitnessHistory)

	fmt.Println(yVals)
	assert.Equal(t, float64(100), yVals[0], "Fitness data not in correct place")
	assert.Equal(t, float64(80), yVals[4], "Fitness data not in correct place")
	assert.Equal(t, float64(20), yVals[10], "Fitness data not in correct place")

	assert.Equal(t, float64(100), yVals[1], "Fitness data should be filled in")
	assert.Equal(t, float64(100), yVals[3], "Fitness data should be filled in")
	assert.Equal(t, float64(80), yVals[9], "Fitness data should be filled in")
}

// TestFillMissingPoints_ExceedsIterations tests the rare case that the last BestFitness occurred beyond the iterations target.
// This can happen because in the final generation the final function evaluation count can exceed the iterations cap.
func TestFillMissingPoints_ExceedsIterations(t *testing.T) {
	iterations := 10
	fitnessHistory := []BestFitness {
		{X: 1, Fitness: 100},
		{X: 14, Fitness: 50},
	}
	yVals := fillMissingPoints(iterations, fitnessHistory)

	fmt.Print(yVals)

	assert.Equal(t, float64(50), yVals[9], "Fitness results exceeding iterations cap should be placed in last point of plot")
}