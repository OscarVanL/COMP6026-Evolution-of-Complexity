package main

import (
	"fmt"
	pop "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
)

func main() {
	species := pop.GenerateSpecies(f.RastriginN, 100)
	species.EvalInitialFitness(f.Rastrigin)

	fmt.Println(species)
	
}
