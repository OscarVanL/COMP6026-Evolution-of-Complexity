package ga

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/common"
	f "github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/optimisation"
	"github.com/jinzhu/copier"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	CrossoverP = 0.6  // Probability of performing crossover
	W = 5  // Scaling Window width
)

func Run(evaluations int, generations int, popSize int, N int, function f.Fitness, mutationP float32) ([]chart.BestFitness, float64, []uint16) {
	bestFitness := math.MaxFloat64
	var evals int
	var fMax float64  // Scaling Window f'max as per https://ieeexplore.ieee.org/document/4075583
	var bestGenes []uint16
	var worstFitnessHistory []float64  // Track worst fitness for each generation
	var bestFitnessHistory []chart.BestFitness

	// Initialise GA's population
	population := InitPopulation(N, popSize, time.Now().Unix())
	population.EvalFitness(function, 0)
	population.SortFitness()
	bestFitnessHistory = append(bestFitnessHistory, chart.BestFitness{X: evals, Fitness: population[0].Fitness})
	fMax = population[len(population)-1].Fitness  // Set initial value of f'max

	if evaluations != 0 {
		// Run GA for N function evaluations
		for evals<evaluations {
			population.doGeneration(function, mutationP, 0, &evals, &fMax, &bestFitness, &bestGenes, &bestFitnessHistory, &worstFitnessHistory)
		}
	} else if generations != 0 {
		// Run GA for N generations
		for gen:=0; gen<generations; gen++ {
			population.doGeneration(function, mutationP, gen, &evals, &fMax, &bestFitness, &bestGenes, &bestFitnessHistory, &worstFitnessHistory)
		}
	}

	return bestFitnessHistory, bestFitness, bestGenes
}


func (pop Population) doGeneration(function f.Fitness, mutationP float32, gen int, evals *int, fMax *float64, bestFitness *float64, bestGenes *[]uint16, bestFitnessHistory *[]chart.BestFitness, worstFitnessHistory *[]float64) {
	// Select members to be in new population via tournament selection
	//pop.SelectNewPopulation()
	// Perform two-point crossover for each individual
	pop.Crossover(CrossoverP, function)
	// Mutate each individual's genes
	pop.Mutate(mutationP)
	// Re-evaluates individual fitness
	*evals += pop.EvalFitness(function, *fMax)
	// Sort the population's individuals by fittest (smallest) to least fit (largest)
	pop.SortFitness()
	// Finds individual with best fitness & genes in this generation
	bestGenFitness, bestGenGene := pop[0].Fitness, pop[0].Genes
	worstGenFitness := pop[len(pop)-1].Fitness

	if bestGenFitness < *bestFitness {
		*bestFitness = bestGenFitness
		*bestGenes = bestGenGene
		if gen != 0 {
			*bestFitnessHistory = append(*bestFitnessHistory, chart.BestFitness{X: gen, Fitness: bestGenFitness})
		} else {
			*bestFitnessHistory = append(*bestFitnessHistory, chart.BestFitness{X: *evals, Fitness: bestGenFitness})
		}

	}
	*worstFitnessHistory = append(*worstFitnessHistory, worstGenFitness)
	*fMax = common.CalculateFMax(*worstFitnessHistory, W)
}

// SelectNewPopulation updates the individuals in the population using tournament selection
func (pop Population) SelectNewPopulation() {
	// Make deep copy of last generation's population
	lastGeneration := Population{}
	err := copier.CopyWithOption(&lastGeneration, &pop, copier.Option{DeepCopy: true})
	if err != nil {
		log.Fatal("Unable to make clone of last population via deep copy:", err)
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// Perform tournament selection (Experiment 9)
	for i:=1; i<len(pop); i++ {
		individualA, individualB := lastGeneration[r.Intn(len(pop))], lastGeneration[r.Intn(len(pop))]
		if individualA.Fitness > individualB.Fitness {
			pop[i] = individualB
		} else {
			pop[i] = individualA
		}
	}
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

func (pop Population) Crossover(crossoverP float32, fitness f.Fitness) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	pop.RouletteSetup()

	for i:=1; i<len(pop); i++ {
		if r.Float32() < crossoverP {
			// Select individual for crossover from last generation
			rouletteGenes := pop.RouletteSelection(r).Genes
			// Perform two-point crossover
			offspringA, offspringB := make([]uint16, len(pop[i].Genes)), make([]uint16, len(pop[i].Genes))
			for g:=0; g<len(pop[i].Genes); g++ {
				offspringA[g], offspringB[g] = common.TwoPointCrossover(pop[i].Genes[g], rouletteGenes[g])
			}

			// Pick best offspring
			fitnessA, fitnessB := fitness(offspringA), fitness(offspringB)
			if fitnessA > fitnessB {
				pop[i].Genes = offspringB
			} else {
				pop[i].Genes = offspringA
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

