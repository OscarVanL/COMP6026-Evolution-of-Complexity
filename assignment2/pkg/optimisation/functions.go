package optimisation

import (
	"math"
)

type Fitness func(x []uint16) float64

const(
	RastriginN = 20
	RastriginMin = -5.12
	RastriginMax = 5.12
	RastriginMutationP = float64(1)/RastriginN
)

func Rastrigin(x []uint16) float64 {
	xScaled := scaleInputs(x[:], RastriginMin, RastriginMax)
	sum := 0.0
	for i:=0; i< RastriginN; i++ {
		sum += math.Pow(xScaled[i], 2) - 3*math.Cos(2*math.Pi*xScaled[i])
	}
	return 3*RastriginN + sum
}

const(
	SchwefelN = 10
	SchwefelMin = -500.0
	SchwefelMax = 500.0
	SchwefelMutationP = float64(1)/SchwefelN
)

// Schwefel function differs to that in the paper, the paper has a mistake in a sign (+ve instead of -ve)
func Schwefel(x []uint16) float64 {
	xScaled := scaleInputs(x[:], SchwefelMin, SchwefelMax)
	sum := 0.0
	for i:=0; i< SchwefelN; i++ {
		sum += xScaled[i] * math.Sin(math.Sqrt(math.Abs(xScaled[i])))
	}
	return 418.9829*SchwefelN - sum
}

const(
	GriewangkN = 10
	GriewangkMin = -600.0
	GriewangkMax = 600.0
	GriewangkMutationP = float64(1)/GriewangkN
)

func Griewangk(x []uint16) float64 {
	xScaled := scaleInputs(x[:], GriewangkMin, GriewangkMax)
	sigma := 0.0
	product := 1.0
	for i:=0; i< GriewangkN; i++ {
		sigma += math.Pow(xScaled[i], 2) / 4000
		product *= math.Cos(xScaled[i]/math.Sqrt(float64(i+1)))
	}

	return 1.0 + sigma - product
}

const (
	AckleyN = 30
	AckleyMin = -30.0
	AckleyMax = 30.0
	AckleyMutationP = float64(1)/AckleyN
)

func Ackley(x []uint16) float64 {
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

func Rosenbrock(x []uint16) float64 {
	// Todo: Write Rosenbrock function
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