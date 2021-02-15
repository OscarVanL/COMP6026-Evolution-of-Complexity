import string
import random


# Randomly generate ascii character
def randchar() -> str:
    return random.choice(string.ascii_uppercase + string.ascii_lowercase + string.digits + string.punctuation + ' ')


# Generate a string individual
def gen_str_individual(len) -> str:
    return ''.join(randchar() for _ in range(len))


# Calculate individual based on string character matches
def eval_str_individual(individual, target) -> int:
    return sum([1 for (i,t) in zip(individual, target) if i == t])


# Mutate an individual, with a probability for each character
def mutate_str_individual(individual, probability) -> str:
    return ''.join(c if random.random()>probability else randchar() for c in individual)