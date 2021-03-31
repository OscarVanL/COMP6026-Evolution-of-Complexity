package evolution

const (
	twoPointMask1 = uint16((0x0fff >> 4) << 4) // 0000 1111 1111 0000
	twoPointMask2 = 0xffff ^ twoPointMask1     // 1111 0000 0000 1111
)

// TwoPointCrossover creates two offspring from parent A and parent B's genes using two-point crossover
func TwoPointCrossover(parentA uint16, parentB uint16) (uint16, uint16) {
	// Two point crossover bitwise operations adapted from: https://stackoverflow.com/a/11705889/6008271
	output1 := (parentA & twoPointMask1) ^ (parentB & twoPointMask2)
	output2 := (parentA & twoPointMask2) ^ (parentB & twoPointMask1)
	return output1, output2
}