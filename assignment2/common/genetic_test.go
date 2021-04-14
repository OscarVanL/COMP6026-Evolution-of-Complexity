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
	for pos:=0; pos<16; pos++ {
		assert.True(t, HasBit(uint16(0xFFFF), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
	// Test no bits
	for pos:=0; pos<16; pos++ {
		assert.False(t, HasBit(uint16(0x0000), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
	// Test number '1234'.
	bits1234 := []bool{false, true, false, false, true, false, true, true, false, false, true, false, false, false, false, false}
	for pos:=0; pos<16; pos++ {
		assert.Equal(t, bits1234[pos], HasBit(uint16(1234), uint(pos)), fmt.Sprintf("Unexpected bit at index: %d", pos))
	}
}

func TestSetBit(t *testing.T) {
	// Test setting all bits {
	value := uint16(0x0000)
	for pos:=0; pos<16; pos++ {
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
	for pos:=0; pos<16; pos++ {
		value = ClearBit(value, uint(pos))
	}
	assert.Equal(t, uint16(0x0000), value, fmt.Sprintf("ClearBit did not clear all bits"))

	// Test clearing one bit from test number 1234
	value2 := uint16(1234)
	value2 = ClearBit(value2, uint(1))
	assert.Equal(t, uint16(1232), value2, fmt.Sprintf("ClearBit did not clear correct bit"))
}