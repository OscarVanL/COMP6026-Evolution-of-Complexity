package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMasksValid(t *testing.T) {
	// Check Two Point Crossover Mask 1 is 0000 1111 1111 0000
	assert.Equal(t, uint16(0x0FF0), twoPointMask1, fmt.Sprintf("Invalid two-point crossover mask"))
	// Check Two Point Crossover Mask 2 is 1111 0000 0000 1111
	assert.Equal(t, uint16(0xF00F), twoPointMask2, fmt.Sprintf("Invalid two-point crossover mask"))
}

func TestTwoPointCrossover(t *testing.T) {
	output1, output2 := TwoPointCrossover(uint16(0xFFFF), uint16(0x0000))

	assert.Equal(t, output1, uint16(0x0FF0), fmt.Sprintf("Invalid two-point crossover result"))
	assert.Equal(t, output2, uint16(0xF00F), fmt.Sprintf("Invalid two-point crossover result"))
}

func TestTwoPointCrossoverArbitrary(t *testing.T) {
	output1, output2 := TwoPointCrossover(uint16(61345), uint16(1235))

	assert.Equal(t, output1, uint16(4003), fmt.Sprintf("Invalid two-point crossover result"))
	assert.Equal(t, output2, uint16(58577), fmt.Sprintf("Invalid two-point crossover result"))
}

// TestTwoPointCrossoverGA_HalfMask tests the crossover where the crossover point crosses within a uint16
func TestTwoPointCrossoverGA_HalfMask(t *testing.T) {
	inputA := []uint16{0xFFFF, 0x0000, 0x0F0F, 0x0FF0, 0x0000, 0xFFFF}
	inputB := []uint16{0xAAAA, 0xBBBB, 0xCCCC, 0xDDDD, 0xEEEE, 0xFFFF}
	offspringA, offspringB, err := TwoPointCrossoverGA(inputA, inputB)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	expectedA := []uint16{0xFFFF, 0x00BB, 0xCCCC, 0xDDDD, 0xEE00, 0xFFFF}
	expectedB := []uint16{0xAAAA, 0xBB00, 0x0F0F, 0x0FF0, 0x00EE, 0xFFFF}
	assert.Equal(t, offspringA, expectedA, "Offspring should have two-point crossover")
	assert.Equal(t, offspringB, expectedB, "Offspring should have two-point crossover")
}

// TestTwoPointCrossoverGA_FullMask tests the crossover where it can be done neatly without mask tricks (size of genes are a multiple of 4)
func TestTwoPointCrossoverGA_FullMask(t *testing.T) {
	// Test 1: Ensure correct bits chosen from each parent
	inputA := []uint16{0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF}
	inputB := []uint16{0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000, 0x0000}
	offspringA, offspringB, err := TwoPointCrossoverGA(inputA, inputB)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	expectedA := []uint16{0xFFFF, 0xFFFF, 0x0000, 0x0000, 0x0000, 0x0000, 0xFFFF, 0xFFFF}
	expectedB := []uint16{0x0000, 0x0000, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF, 0x0000, 0x0000}
	assert.Equal(t, offspringA, expectedA, "Offspring should have two-point crossover")
	assert.Equal(t, offspringB, expectedB, "Offspring should have two-point crossover")

	// Test 2: Ensure masks do not overwrite content
	inputA = []uint16{0x0000, 0x1111, 0x2222, 0x3333}
	inputB = []uint16{0xAAAA, 0xBBBB, 0xCCCC, 0xDDDD}
	offspringA, offspringB, err = TwoPointCrossoverGA(inputA, inputB)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	expectedA = []uint16{0x0000, 0xBBBB, 0xCCCC, 0x3333}
	expectedB = []uint16{0xAAAA, 0x1111, 0x2222, 0xDDDD}
	assert.Equal(t, offspringA, expectedA, "Offspring should have two-point crossover")
	assert.Equal(t, offspringB, expectedB, "Offspring should have two-point crossover")
}

func TestCalculateFMax(t *testing.T) {
	fitnesswindow := []float64{123.0, 321.0, 242.0}
	fMax := CalculateFMax(fitnesswindow, 5)
	assert.Equal(t, 321.0, fMax, "CalculateFMax did not get largest fitness score")

	fitnesswindow = []float64{123.0, 321.0, 242.0, 452.0, 12.0}
	fMax = CalculateFMax(fitnesswindow, 5)
	assert.Equal(t, 452.0, fMax, "CalculateFMax did not get largest fitness score")

	fitnesswindow = []float64{999.0, 123.0, 999.0, 242.0, 452.0, 12.0, 125.0, 124.0}
	fMax = CalculateFMax(fitnesswindow, 5)
	assert.Equal(t, 452.0, fMax, "CalculateFMax did not get largest fitness score within window")
}

func TestHasBit(t *testing.T) {
	// Test all bits
	for pos := 0; pos < 16; pos++ {
		assert.True(t, HasBit(uint16(0xFFFF), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
	// Test no bits
	for pos := 0; pos < 16; pos++ {
		assert.False(t, HasBit(uint16(0x0000), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
	// Test number '1234'.
	bits1234 := []bool{false, true, false, false, true, false, true, true, false, false, true, false, false, false, false, false}
	for pos := 0; pos < 16; pos++ {
		assert.Equal(t, bits1234[pos], HasBit(uint16(1234), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
}

func TestSetBit(t *testing.T) {
	// Test setting all bits {
	value := uint16(0x0000)
	for pos := 0; pos < 16; pos++ {
		value = SetBit(value, uint(pos))
	}
	assert.Equal(t, uint16(0xFFFF), value, fmt.Sprintf("SetBit did not set all bits"))

	// Test setting one bit from test number 1234
	value2 := uint16(1234)
	value2 = SetBit(value2, uint(0))
	assert.Equal(t, uint16(1235), value2, fmt.Sprintf("SetBit did not set correct bit"))
}

func TestClearBit(t *testing.T) {
	// Test clearing all bits
	value := uint16(0xFFFF)
	for pos := 0; pos < 16; pos++ {
		value = ClearBit(value, uint(pos))
	}
	assert.Equal(t, uint16(0x0000), value, fmt.Sprintf("ClearBit did not clear all bits"))

	// Test clearing one bit from test number 1234
	value2 := uint16(1234)
	value2 = ClearBit(value2, uint(1))
	assert.Equal(t, uint16(1232), value2, fmt.Sprintf("ClearBit did not clear correct bit"))
}
