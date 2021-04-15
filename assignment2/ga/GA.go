package ga

import (
	"fmt"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/cheggaaa/pb"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

// Crossover probability
const CrossoverP = 0.6

func Run(iterations int, popSize int, N int, function f.Fitness, mutationP float32, W int) ([]chart.BestFitness, float64, []uint16) {
	bar := pb.New(iterations)
	bar.SetRefreshRate(time.Second)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.Start()

	bestFitness := math.MaxFloat64
	var evals int
	var fMax float64  // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestGenes []uint16
	var worstFitnessHistory []float64  // Track worst fitness for each generation
	var bestFitnessHistory []chart.BestFitness
	//var bestFitnessHistory []float64  // Track best overall fitness across all generations

	population := InitPopulation(N, popSize, time.Now().Unix())
	evals += population.EvalFitness(function, 0)
	population.SortFitness()
	bestFitnessHistory = append(bestFitnessHistory, chart.BestFitness{X: evals, Fitness: population[0].Fitness})
	fMax = population[len(population)-1].Fitness  // Set initial value of f'max
	fmt.Println("First fMax:", fMax)

	for evals<iterations {
		// Todo: Track the number of function evaluations, not GA iterations?
		//population.SortScaledFitness()
		population.Crossover(CrossoverP)
		population.Mutate(mutationP)
		evals += population.EvalFitness(function, fMax)
		population.SortFitness()
		bestGenFitness, bestGenGene := population[0].Fitness, population[0].Genes
		worstGenFitness := population[len(population)-1].Fitness
		if bestGenFitness < bestFitness {
			bestFitness = bestGenFitness
			bestGenes = bestGenGene
			bestFitnessHistory = append(bestFitnessHistory, chart.BestFitness{X: evals, Fitness: population[0].Fitness})
		}
		worstFitnessHistory = append(worstFitnessHistory, worstGenFitness)
		fMax = common.CalculateFMax(worstFitnessHistory, W)

		bar.Set(evals)
	}
	bar.Finish()

	return bestFitnessHistory, bestFitness, bestGenes
}

// Mutate performs bit-flip mutation on each of the individual's genes
func (pop Population) Mutate(MutationP float32) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i:=1; i<len(pop); i++ {
		individual := pop[i]

		mutatedGenes := individual.Genes

		// Mutate each of the individual's genes
		for g:=0; g<len(mutatedGenes); g++ {
			// Mutate each of the 16 bits in the gene
			for b:=0; b<16; b++ {
				// P probability of mutation
				if r.Float32() < MutationP {
					// Perform bit-flip
					if common.HasBit(mutatedGenes[g], uint(b)) {
						mutatedGenes[g] = common.ClearBit(mutatedGenes[g], uint(b))
					} else {
						mutatedGenes[g] = common.SetBit(mutatedGenes[g], uint(b))
					}
				}
			}
		}
		// Replace individual's old genes with mutated ones
		pop[i].Genes = mutatedGenes
	}
}

func (pop Population) Crossover(crossoverP float32) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	pop.RouletteSetup()

	for i:=1; i<len(pop); i++ {
		rouletteGenes := pop.RouletteSelection(r).Genes
		if r.Float32() < crossoverP {
			// Perform two-point crossover
			offspringA, offspringB, err := common.TwoPointCrossoverGA(pop[i].Genes, rouletteGenes)
			if err != nil{
				_, _ = fmt.Fprintf(os.Stderr, "Error during two-point crossover: %v\n", err)
				os.Exit(1)
			}

			// Todo: Try picking the best offspring, rather than randomly select.
			if r.Intn(2) == 0 {
				pop[i].Genes = offspringA
			} else {
				pop[i].Genes = offspringB
			}
		}
	}

	for i:=1; i<len(pop); i++ {
		rouletteGenes1 := pop.RouletteSelection(r).Genes
		//rouletteGenes2 := pop.RouletteSelection(r).Genes

		// Check whether to use crossover for this individual
		if r.Float32() < crossoverP {
			// Do crossover for each gene
			for g:=0; g<len(pop[i].Genes); g++ {
				// Perform two-point crossover
				offspringA, offspringB := common.TwoPointCrossover(pop[i].Genes[g], rouletteGenes1[g])

				// Randomly select one of the offspring to use
				// Todo: Try picking the best offspring, rather than randomly select.
				if r.Intn(2) == 0 {
					pop[i].Genes[g] = offspringA
				} else {
					pop[i].Genes[g] = offspringB
				}
			}
		}
	}
}

// RouletteSetup calculates population selection probabilities from ScaledFitness scores, required before using RouletteSelection.
// Adapted from: https://stackoverflow.com/a/177278/6008271
func (pop Population) RouletteSetup() {
	// Calculate total fitness for subpopulation
	var fitnessSum float64
	for i:=0; i<len(pop); i++ {
		fitnessSum += pop[i].ScaledFitness
	}

	// Calculate fitness proportionate probability
	var accumulatedProbability float64
	for i:=0; i<len(pop); i++ {
		pop[i].SelectProbability = pop[i].ScaledFitness / fitnessSum
		pop[i].SelectProbability += accumulatedProbability
		accumulatedProbability = pop[i].SelectProbability
	}
}

// RouletteSelection uses a roulette approach to apply higher selective pressure for individuals with better fitness
// Adapted from: https://stackoverflow.com/a/177278/6008271
func (pop Population) RouletteSelection(r *rand.Rand) Individual {
	// Todo: Use binary search here, instead of linear search.
	number := r.Float64()
	for p:=0; p<len(pop); p++ {
		if p == 0 {
			// First entry on roulette wheel, range 0.0 - Select Probability
			if number < pop[p].SelectProbability {
				return pop[p]
			}
		} else {
			// SelectProbability greater than last individual, but within this individual's probability range
			if number > pop[p-1].SelectProbability && number < pop[p].SelectProbability {
				return pop[p]
			}
		}
	}
	return pop[0]
}

// EvalFitness checks the fitness of an individual's genes and updates its Fitness & ScaledFitness scores.
// Return number of fitness evaluations
func (pop Population) EvalFitness(fitness f.Fitness, fMax float64) int {
	for i:=0; i<len(pop); i++ {
		// Calculate individual's Fitness
		pop[i].Fitness = fitness(pop[i].Genes)
		pop[i].ScaledFitness = math.Abs(fMax - pop[i].Fitness)
	}
	return len(pop)
}

func (pop Population) SortFitness() {
	// Sort the populations individuals by fittest (smallest) to least fit (largest)
	sort.Slice(pop, func(i, j int) bool {
		return pop[i].Fitness < pop[j].Fitness
	})
}

