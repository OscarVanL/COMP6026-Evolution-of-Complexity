package optimisation

import (
	"math"
)

const Rastrigin_n = 20
const Rastrigin_min = -5.12
const Rastrigin_max = 5.12

func Rastrigin(x [Rastrigin_n]uint16) float64 {
	x_scaled := scaleInputs(x[:], Rastrigin_min, Rastrigin_max)
	sum := 0.0
	for i:=0; i< Rastrigin_n; i++ {
		x_i := float64(x_scaled[i])
		sum += math.Pow(x_i, 2) - 3*math.Cos(2*math.Pi*x_i)
	}
	return 3*Rastrigin_n + sum
}

const Schwefel_n = 10
const Schwefel_min = -500.0
const Schwefel_max = 500.0

// Schwefel function differs to that in the paper, the paper has a mistake in a sign (+ve instead of -ve)
func Schwefel(x [Schwefel_n]uint16) float64 {
	x_scaled := scaleInputs(x[:], Schwefel_min, Schwefel_max)
	sum := 0.0
	for i:=0; i<Schwefel_n; i++ {
		x_i := float64(x_scaled[i])
		sum += x_i*math.Sin(math.Sqrt(math.Abs(x_i)))
	}
	return 418.9829*Schwefel_n - sum
}

const Griewangk_n = 10
const Griewangk_min = -600.0
const Griewangk_max = 600.0

func Griewangk(x [Griewangk_n]uint16) float64 {
	x_scaled := scaleInputs(x[:], Griewangk_min, Griewangk_max)
	sigma := 0.0
	product := 1.0
	for i:=0; i<Griewangk_n; i++ {
		x_i := float64(x_scaled[i])
		sigma += math.Pow(x_i, 2) / 4000
		product *= math.Cos(x_i/math.Sqrt(float64(i+1)))
	}

	return 1.0 + sigma - product
}

const Ackley_n = 30
const Ackley_min = -30.0
const Ackley_max = 30.0


func Ackley(x [Ackley_n]uint16) float64 {
	x_scaled := scaleInputs(x[:], Ackley_min, Ackley_max)
	sumA, sumB := 0.0, 0.0
	for i:=0; i<Ackley_n; i++ {
		x_i := float64(x_scaled[i])
		sumA += math.Pow(x_i, 2)
		sumB += math.Cos(2*math.Pi*x_i)
	}
	sumA *= 1/Ackley_n
	sumB *= 1/Ackley_n

	return 20 + math.E - 20*math.Exp(-0.2*math.Sqrt(sumA)) - math.Exp(sumB)
}

func Rosenbrock() float64 {
	return 0.0
}

// scaleInputs scales an slice of uint16s between two ranges
func scaleInputs(x []uint16, min float64, max float64) []float64 {
	var scaled []float64
	for i:=0; i<len(x); i++ {
		scaled = append(scaled, (float64(x[i]) / 65535) * (max - min) + min)
	}
	return scaled
}