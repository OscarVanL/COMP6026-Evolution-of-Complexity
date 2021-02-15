from utils import EvolutionUtils

TARGET_STRING = 'methinks it is like a weasel'
POP_SIZE = 500


class CrossoverEvolution:
    def __init__(self, target, population):
        L = len(target)
        MUTATION_RATE = 1 / L

        # Generate random individuals
        pop = [EvolutionUtils.gen_str_individual(L) for _ in range(population)]
        # Evaluate similarity to target string
        pop = {i:EvolutionUtils.eval_str_individual(i, TARGET_STRING) for i in pop}

        best_fitness = 0
        while best_fitness != L:
            new_pop = {}
            # Mutate the population
            for individual, fit in pop.items():
                mutated = EvolutionUtils.mutate_str_individual(individual, MUTATION_RATE)
                new_fit = EvolutionUtils.eval_str_individual(mutated, TARGET_STRING)
                if new_fit > fit:
                    new_pop[mutated] = new_fit
                else:
                    new_pop[individual] = fit

            # Get fitness
            new_best_fitness = max(new_pop.values())
            if new_best_fitness != best_fitness:
                best_fitness = new_best_fitness
                print(best_fitness)

            pop = new_pop

        print("Match found!")
        print(sorted(pop, key=pop.get, reverse=True))



if __name__ == "__main__":
    CrossoverEvolution(TARGET_STRING, POP_SIZE)
