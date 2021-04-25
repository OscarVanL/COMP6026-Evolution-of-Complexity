// Go bit flip operations taken from https://stackoverflow.com/a/23192263/6008271
// Two point crossover logic adapted from: https://stackoverflow.com/a/11705889/6008271

package common

import (
	"errors"
	"fmt"
	"math"
)

const (
	//twoPointMask1 = uint16((0x0fff >> 4) << 4) // 0000 1111 1111 0000 or 0x0FF0
	//twoPointMask2 = 0xffff ^ twoPointMask1     // 1111 0000 0000 1111 or 0xF00F
	twoPointMask1 = uint16(0x0FF0)
	twoPointMask2 = uint16(0xF00F)
)

// TwoPointCrossover creates two offspring from parent A and parent B's genes using two-point crossover
func TwoPointCrossover(parentA uint16, parentB uint16) (uint16, uint16) {
	// Two point crossover bitwise operations adapted from: https://stackoverflow.com/a/11705889/6008271
	output1 := (parentA & twoPointMask1) ^ (parentB & twoPointMask2)
	output2 := (parentA & twoPointMask2) ^ (parentB & twoPointMask1)
	return output1, output2
}

// TwoPointCrossoverGA for multi-gene crossover used in standard GA.
func TwoPointCrossoverGA(parentA []uint16, parentB []uint16) ([]uint16, []uint16, error) {
	if len(parentA) != len(parentB) {
		return nil, nil, errors.New("gene lengths must match")
	}
	if len(parentA) % 2 != 0 {
		return nil, nil, errors.New("must use even number of genes for Two-Point Crossover")
	}

	mask1 := make([]uint16, len(parentA))
	mask2 := make([]uint16, len(parentA))
	offspringA := make([]uint16, len(parentA))
	offspringB := make([]uint16, len(parentA))

	// Calculate parameters required for constructing bit masks
	widthLR := 0.25 * float64(len(parentA)) * 4   // Number of bytes to keep in left and right parts of 2-point-crossover
	fullMaskLR := int(math.Floor(widthLR / 4))    // Number of full-bytes to keep on left and right parts
	boundaryMaskRequired := int(widthLR) % 4 == 2 // Whether we have half-bytes at boundary
	widthC := 0.5 * float64(len(parentA)) * 4  // Number of bytes to keep in central part of 2-point crossover
	fullMaskC := int(math.Floor(widthC / 4))  // Number of full-bytes in central part

	if boundaryMaskRequired {
		fullMaskC -= 1
	}

	// Create masks for two-point-crossover upon parent A and parent B
	bytesAdded := 0
	// Left part of crossover
	for i:=bytesAdded; i<bytesAdded+fullMaskLR; i++ {
		mask1[i] = 0xFFFF
		mask2[i] = 0x0000
	}
	bytesAdded += fullMaskLR
	// Boundary between left-part and central-part of crossover
	if boundaryMaskRequired {
		mask1[bytesAdded] = 0xFF00
		mask2[bytesAdded] = 0x00FF
		bytesAdded += 1
	}
	// Central crossover point
	for i:=bytesAdded; i<bytesAdded+fullMaskC; i++ {
		mask1[i] = 0x0000
		mask2[i] = 0xFFFF
	}
	bytesAdded += fullMaskC
	// Boundary between central-part and right-part of crossover
	if boundaryMaskRequired {
		mask1[bytesAdded] = 0x00FF
		mask2[bytesAdded] = 0xFF00
		bytesAdded += 1
	}
	// Right part of crossover
	for i:=bytesAdded; i<bytesAdded+fullMaskLR; i++ {
		mask1[i] = 0xFFFF
		mask2[i] = 0x0000
	}
	bytesAdded += fullMaskLR

	// Create offspring
	for i:=0; i<len(parentA); i++ {
		offspringA[i] = (parentA[i] & mask1[i]) ^ (parentB[i] & mask2[i])
		offspringB[i] = (parentA[i] & mask2[i]) ^ (parentB[i] & mask1[i])
	}

	return offspringA, offspringB, nil
}

func printGenesHex(gene []uint16) {
	for i:=0; i<len(gene); i++ {
		fmt.Printf("%s, ", fmt.Sprintf("%04X", gene[i]))
	}
	fmt.Println()
}

// CalculateFMax finds the new FMax value for Scaling Window calculations
func CalculateFMax(worstFitnessHistory []float64, W int) float64 {
	var worstFitnessWindow []float64
	if len(worstFitnessHistory) < W {
		worstFitnessWindow = worstFitnessHistory
	} else {
		worstFitnessWindow = worstFitnessHistory[len(worstFitnessHistory)-W:]
	}
	fMax := worstFitnessWindow[0]
	for _, v := range worstFitnessWindow {
		if v > fMax {
			fMax = v
		}
	}
	return fMax
}

// Bit manipulation adapted from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271

// HasBit Checks if bit is set as position n
func HasBit(n uint16, pos uint) bool {
	val := n & (1 << pos)
	return val > 0
}

// SetBit Sets bit at index pos to 1
func SetBit(n uint16, pos uint) uint16 {
	n |= (1 << pos)
	return n
}

// ClearBit Sets bit at index pos to 0
func ClearBit(n uint16, pos uint) uint16 {
	mask := ^(1 << pos)
	nTemp := int(n)
	nTemp &= mask
	return uint16(nTemp)
}