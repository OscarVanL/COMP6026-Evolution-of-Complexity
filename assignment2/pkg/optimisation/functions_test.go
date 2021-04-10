package optimisation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRastrigin(t *testing.T) {
	input := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	answer := Rastrigin(input)
	assert.InDelta(t, 539.7935867419988, answer, 0.01, "Unexpected result from Rastrigin function")
}

func TestSchwefel(t *testing.T) {
	input := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	answer := Schwefel(input)
	assert.InDelta(t, 2392.9928386744673, answer, 0.01, "Unexpected result from Schwefel function")
}

func TestGriewangk(t *testing.T) {
	input := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	answer := Griewangk(input)
	assert.InDelta(t, 900.6978630212193, answer, 0.01, "Unexpected result from Griewangk function")
}

func TestAckley(t *testing.T) {
	input := []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30}
	answer := Ackley(input)
	assert.InDelta(t, 19.964400407804085, answer, 0.01, "Unexpected result from Ackley function")
}

func TestRosenbrock(t *testing.T) {
	// Todo: Finish writing Rosenbrock function & then write unit test
	var input []uint16
	answer := Rosenbrock(input)
	assert.InDelta(t, 0.0, answer, 0.01, "Unexpected result from Rosenbrock function")
}

func TestScaleInputs(t *testing.T) {
	// This test scales uint16s (0-65535) into range 0-6553.5, the expected behaviour is all values are 1/10th the original scale
	input := []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 65535}
	min := 0.0
	max := 6553.5
	scaled := scaleInputs(input, min, max)
	fmt.Println(scaled)
	for i:=0; i<len(input); i++ {
		assert.InDelta(t, float64(input[i])/10, scaled[i], 0.01, "ScaleInputs did not scale inputs correctly")
	}
}