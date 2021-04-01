package ccga

import (
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
	"testing"
)

func TestMutate(t *testing.T) {
	N := f.RastriginN
	function := f.Rastrigin
	mutationP := f.RastriginMutationP

	species := InitSpecies(N, 100)
	species.InitCoevolutions()
	species.EvalFitness(function)

	for i:=0; i<2000; i++ {
		species.Mutate(mutationP)
	}
}