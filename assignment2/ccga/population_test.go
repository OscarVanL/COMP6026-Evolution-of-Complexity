package ccga

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitPopulationPopSize(t *testing.T) {
	population := InitSpecies(1, 5, 0)

	expected := Species{
		Population{
			Individual{0, 37889, 0, 0, 0, nil},
			Individual{0, 16832, 0, 0, 0, nil},
			Individual{0, 51315, 0, 0, 0, nil},
			Individual{0, 42338, 0, 0, 0, nil},
			Individual{0, 11594, 0, 0, 0, nil},
		},
	}

	assert.Equal(t, 1, len(population), "InitPopulation did not create desired number of Species")
	assert.Equal(t, 5, len(population[0]), "InitPopulation did not create desired PopSize")
	assert.Equal(t, expected, population, "InitPopulation did not create population as expected")
}

func TestInitPopulationGenes(t *testing.T) {
	population := InitSpecies(3, 1, 0)

	expected := Species{
		Population{
			Individual{0, 37889, 0, 0, 0, nil},
		},
		Population{
			Individual{1, 16832, 0, 0, 0, nil},
		},
		Population{
			Individual{2, 51315, 0, 0, 0, nil},
		},
	}

	assert.Equal(t, 3, len(population), fmt.Sprintf("InitPopulation did not create desired number of species"))
	assert.Equal(t, 1, len(population[0]), fmt.Sprintf("InitPopulation did not create desired PopSize"))
	assert.Equal(t, expected, population, fmt.Sprintf("InitPopulation did not create population as expected"))
}