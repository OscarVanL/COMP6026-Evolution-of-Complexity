package main

import (
	"fmt"
	pop "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/optimisation"
)

const iterations = 100000

func main() {
	var gen pop.Generations

	species := pop.GenerateSpecies(f.RastriginN, 100)
	species.InitCoevolutions()
	species.EvalFitness(f.Rastrigin)

	gen = append(gen, species)


	var bestCoevolution []uint16
	var fitnessHistory []float64
	for i:=0; i<iterations; i++ {
		if i % 1000 == 0 {
			fmt.Println(i)
		}

		//start := time.Now()
		species.EvolveSpecies(f.RastriginMutationP)
		//elapsed := time.Since(start)
		//log.Printf("EvolveSpecies took %s", elapsed)

		//start2 := time.Now()
		species.EvalFitness(f.Rastrigin)
		//elapsed2 := time.Since(start2)
		//log.Printf("EvalFitness took %s", elapsed2)

		//Todo: Do mutation & evolution cycle
		//start3 := time.Now()
		fitness, coevolution := species.GetBestFitness()
		//log.Println("GetBestFitness took:", time.Since(start3))

		fitnessHistory = append(fitnessHistory, fitness)
		bestCoevolution = coevolution
	}

	//fmt.Println("Fitness history:", fitnessHistory)
	fmt.Println("Best Coevolution parameters:")
	for i:=0; i<len(bestCoevolution); i++ {
		fmt.Print(bestCoevolution[i], ", ")
	}


}
