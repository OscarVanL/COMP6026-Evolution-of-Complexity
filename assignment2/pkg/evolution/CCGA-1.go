package evolution

import "math/rand"

// Mutation probability
const CrossoverP = 0.6

// Fitness functions that can be used to evaluate several individuals
type Fitness func(x []Individual) float64

// EvalInitialFitness calculates the initial fitness each member of newly generated species.
func (pop Species) EvalInitialFitness(fitness Fitness) {
	// Evaluate each species
	for s :=0; s<len(pop); s++ {
		species := pop[s]
		// Evaluate each individual in the species
		for i := 0; i<len(species); i++ {
			individual := species[i]

			// Combining individual with random individuals from each other species to calculate initial fitness
			tmpPop := make([]Individual, len(pop))
			for N :=0; N <len(pop); N++ {
				if s == N {
					// Keep the evaluated individual for its gene
					tmpPop[N] = individual
				} else {
					// Use a random individual from the other species genes
					tmpPop[N] = pop[N][rand.Intn(len(pop[N]))]
				}
			}
			individual.fitness = fitness(tmpPop)
			pop[s][i] = individual
		}
	}
}