package ga

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitPopulationPopSize(t *testing.T) {
	population := InitPopulation(1, 5, 0)

	assert.Equal(t, 5, len(population), fmt.Sprintf("InitPopulation did not create desired PopSize"))

	expected := Population{
		Individual{[]uint16{37889}, 0, 0, 0},
		Individual{[]uint16{16832}, 0, 0, 0},
		Individual{[]uint16{51315}, 0, 0, 0},
		Individual{[]uint16{42338}, 0, 0, 0},
		Individual{[]uint16{11594}, 0, 0, 0},
	}

	assert.Equal(t, expected, population, fmt.Sprintf("InitPopulation did not create population as expected"))
}

func TestInitPopulationGenes(t *testing.T) {
	population := InitPopulation(3, 1, 0)

	assert.Equal(t, 3, len(population[0].Genes), fmt.Sprintf("InitPopulation did not create desired number of genes per individual"))

	expected := Population{
		Individual{[]uint16{37889, 16832, 51315}, 0, 0, 0},
	}

	assert.Equal(t, expected, population, fmt.Sprintf("InitPopulation did not create population as expected"))
}
