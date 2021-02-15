import string
from random import random, choice


# Generate random population and return [(individual, fitness)]
from typing import List, Tuple


def gen_rand_pop(population: int, length: int, target_string: str) -> List[Tuple[str, int]]:
    # Generate random individuals
    pop = [gen_str_individual(length) for _ in range(population)]
    # Evaluate similarity to target string
    return [(i, eval_str_individual(i, target_string)) for i in pop]


# Randomly generate ascii character
def randchar() -> str:
    return choice(string.ascii_uppercase + string.ascii_lowercase + string.digits + string.punctuation + ' ')


# Generate a string individual
def gen_str_individual(len: int) -> str:
    return ''.join(randchar() for _ in range(len))


# Calculate individual based on string character matches
def eval_str_individual(individual: str, target: str) -> int:
    return sum([1 for (i,t) in zip(individual, target) if i == t])


# Mutate an individual, with a probability for each character
def mutate_str_individual(individual: str, probability: float) -> str:
    return ''.join(c if random()>probability else randchar() for c in individual)


# Create uniform crossover offspring from two parents
def gen_crossover(p1: str, p2: str) -> str:
    assert len(p1) == len(p2)
    return ''.join(p1[i] if random()<0.5 else p2[i] for i in range(len(p1)))