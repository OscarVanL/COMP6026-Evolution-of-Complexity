package ga

import (
	"fmt"
	"math/rand"
)

func (pop Population) Mutate(MutationP float64) {
	// N.B. Bit manipulation inner-functions are taken from Stack Overflow. Source: https://stackoverflow.com/a/23192263/6008271

	//Checks if bit is set as position n
	hasBit := func(n uint16, pos uint) bool {
		val := n & (1 << pos)
		return val > 0
	}

	// Sets bit at index pos to 1
	setBit := func(n uint16, pos uint) uint16 {
		n |= (1 << pos)
		return n
	}

	// Sets bit at index pos to 0
	clearBit := func(n uint16, pos uint) uint16 {
		mask := ^(1 << pos)
		nTemp := int(n)
		nTemp &= mask
		return uint16(nTemp)
	}


	for i:=0; i<len(pop); i++ {
		individual := pop[i]

		mutatedGenes := individual.Genes
		fmt.Println("Before mutation:", mutatedGenes)

		// Mutate each of the individual's genes
		for g:=0; g<len(mutatedGenes); g++ {
			// Mutate each of the 16 bits in the gene
			for b:=0; b<16; b++ {
				// P probability of mutation
				if rand.Float64() < MutationP {
					fmt.Println("Mutating")
					// Perform bit-flip
					if hasBit(mutatedGenes[g], uint(b)) {
						mutatedGenes[g] = clearBit(mutatedGenes[g], uint(b))
					} else {
						mutatedGenes[g] = setBit(mutatedGenes[g], uint(b))
					}
				} else {
					fmt.Println("Not mutating")
				}
			}
		}
		fmt.Println("After mutation:", mutatedGenes)
		// Replace individual's old genes with mutated ones
		pop[i].Genes = mutatedGenes
	}


}
