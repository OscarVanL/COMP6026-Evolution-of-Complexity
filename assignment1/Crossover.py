from utils import EvolutionUtils as Evol
import random
import operator

TARGET_STRING = 'methinks it is like a weasel'
POP_SIZE = 500


class CrossoverEvolution:
    def __init__(self, target, population):
        self.L = len(target)
        self.MUTATION_RATE = 1 / self.L

        # Generate random individuals
        pop = [Evol.gen_str_individual(self.L) for _ in range(population)]
        # Evaluate similarity to target string
        self.pop = [(i, Evol.eval_str_individual(i, TARGET_STRING)) for i in pop]
        print("pop size:", len(pop))

        # soln, evals = self.mutation_hill_climber()
        soln, evals = self.GA_without_crossover()

        print("Match found in {} evaluations!".format(evals))
        print("Solution:", soln)


    def mutation_hill_climber(self):
        evals = 0
        best_fitness = max(self.pop, key=operator.itemgetter(1))
        while best_fitness[1] != self.L:
            # Mutate the population
            for i, p in enumerate(self.pop):
                individual, fit = p
                mutated = Evol.mutate_str_individual(individual, self.MUTATION_RATE)
                new_fit = Evol.eval_str_individual(mutated, TARGET_STRING)
                evals += 1
                # Replace original with mutation if it's an improvement
                if new_fit > fit:
                    self.pop[i] = (mutated, new_fit)

                    # Update best fitness
                    if new_fit > best_fitness[1]:
                        best_fitness = mutated, new_fit
                        print(best_fitness)

        return best_fitness, evals


    def GA_without_crossover(self):
        evals = 0
        best_fitness = max(self.pop, key=operator.itemgetter(1))
        while best_fitness[1] != self.L:
            # Pick 2 random parents
            A, A_fit = random.choice(self.pop)
            B, B_fit = random.choice(self.pop)

            # Pick the fitter parent
            if A_fit > B_fit:
                parent1 = A
            else:
                parent1 = B

            # Create mutation from fitter parent
            child = Evol.mutate_str_individual(parent1, self.MUTATION_RATE)
            C_fit = Evol.eval_str_individual(child, TARGET_STRING)
            evals += 1

            # Pick random parents to be replaced
            A_i = random.choice(range(len(self.pop)))
            B_i = random.choice(range(len(self.pop)))
            A, A_fit = self.pop[A_i]
            B, B_fit = self.pop[B_i]

            # Replace the less fit parent with the child (but only if that child does not already exist)
            if A_fit > B_fit:
                self.pop[B_i] = (child, C_fit)
            else:
                self.pop[A_i] = (child, C_fit)

            # Update fitness score
            if C_fit > best_fitness[1]:
                best_fitness = child, C_fit
                print(best_fitness)

        # Return solution
        return best_fitness, evals




if __name__ == "__main__":
    CrossoverEvolution(TARGET_STRING, POP_SIZE)
