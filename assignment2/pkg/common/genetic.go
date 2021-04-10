// Go bit flip operations taken from https://stackoverflow.com/a/23192263/6008271
// Two point crossover logic adapted from: https://stackoverflow.com/a/11705889/6008271

package common

const (
	twoPointMask1 = uint16((0x0fff >> 4) << 4) // 0000 1111 1111 0000 or 0x0FF0
	twoPointMask2 = 0xffff ^ twoPointMask1     // 1111 0000 0000 1111 or 0xF00F
)

// TwoPointCrossover creates two offspring from parent A and parent B's genes using two-point crossover
func TwoPointCrossover(parentA uint16, parentB uint16) (uint16, uint16) {
	// Two point crossover bitwise operations adapted from: https://stackoverflow.com/a/11705889/6008271
	output1 := (parentA & twoPointMask1) ^ (parentB & twoPointMask2)
	output2 := (parentA & twoPointMask2) ^ (parentB & twoPointMask1)
	return output1, output2
}

// CalculateFMax finds the new FMax value for Scaling Window calculations
func CalculateFMax(worstFitnessHistory []float64, W int) float64 {
	var worstFitnessWindow []float64
	if len(worstFitnessHistory) < 5 {
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


// N.B. Bit manipulation inner-functions are taken from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271

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