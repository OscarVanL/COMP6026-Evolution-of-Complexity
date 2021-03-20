package optimisation

import (
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/pkg/evolution"
	"math"
)

const RastriginN = 20
const RastriginMin = -5.12
const RastriginMax = 5.12
const RastriginMutationP = 1/RastriginN

func Rastrigin(x []evolution.Individual) float64 {
	xScaled := scaleInputs(x[:], RastriginMin, RastriginMax)
	sum := 0.0
	for i:=0; i< RastriginN; i++ {
		sum += math.Pow(xScaled[i], 2) - 3*math.Cos(2*math.Pi*xScaled[i])
	}
	return 3*RastriginN + sum
}

const SchwefelN = 10
const SchwefelMin = -500.0
const SchwefelMax = 500.0
const SchwefelMutationP = 1/SchwefelN

// Schwefel function differs to that in the paper, the paper has a mistake in a sign (+ve instead of -ve)
func Schwefel(x []evolution.Individual) float64 {
	xScaled := scaleInputs(x[:], SchwefelMin, SchwefelMax)
	sum := 0.0
	for i:=0; i< SchwefelN; i++ {
		sum += xScaled[i] * math.Sin(math.Sqrt(math.Abs(xScaled[i])))
	}
	return 418.9829*SchwefelN - sum
}

const GriewangkN = 10
const GriewangkMin = -600.0
const GriewangkMax = 600.0
const GriewangkMutationP = 1/GriewangkN

func Griewangk(x []evolution.Individual) float64 {
	xScaled := scaleInputs(x[:], GriewangkMin, GriewangkMax)
	sigma := 0.0
	product := 1.0
	for i:=0; i< GriewangkN; i++ {
		sigma += math.Pow(xScaled[i], 2) / 4000
		product *= math.Cos(xScaled[i]/math.Sqrt(float64(i+1)))
	}

	return 1.0 + sigma - product
}

const AckleyN = 30
const AckleyMin = -30.0
const AckleyMax = 30.0
const AckleyMutationP = 1/GriewangkN


func Ackley(x []evolution.Individual) float64 {
	xScaled := scaleInputs(x[:], AckleyMin, AckleyMax)
	sumA, sumB := 0.0, 0.0
	for i:=0; i< AckleyN; i++ {
		sumA += math.Pow(xScaled[i], 2)
		sumB += math.Cos(2*math.Pi*xScaled[i])
	}
	sumA *= 1/ AckleyN
	sumB *= 1/ AckleyN

	return 20 + math.E - 20*math.Exp(-0.2*math.Sqrt(sumA)) - math.Exp(sumB)
}

func Rosenbrock() float64 {
	return 0.0
}

// scaleInputs scales an slice of uint16s between two ranges
func scaleInputs(x []evolution.Individual, min float64, max float64) []float64 {
	var scaled []float64
	for i:=0; i<len(x); i++ {
		scaled = append(scaled, (float64(x[i].Gene) / 65535) * (max - min) + min)
	}
	return scaled
}