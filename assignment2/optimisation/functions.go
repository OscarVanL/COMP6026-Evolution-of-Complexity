package optimisation

import (
	"errors"
	"math"
)

type Fitness func(x []uint16) float64

type Params struct {
	Function  Fitness
	Label     string
	N         int
	MutationP float32
	ScaleMin  float64
	ScaleMax  float64
}

func GetParams(algo string) (Params, error) {
	switch algo {
	case "rastrigin":
		return Params{
			Function:  Rastrigin,
			Label:     RastriginLabel,
			N:         RastriginN,
			MutationP: RastriginMutationP,
			ScaleMin:  RastriginMin,
			ScaleMax:  RastriginMax,
		}, nil
	case "schwefel":
		return Params{
			Function:  Schwefel,
			Label:     SchwefelLabel,
			N:         SchwefelN,
			MutationP: SchwefelMutationP,
			ScaleMin:  SchwefelMin,
			ScaleMax:  SchwefelMax,
		}, nil
	case "griewangk":
		return Params{
			Function:  Griewangk,
			Label:     GriewangkLabel,
			N:         GriewangkN,
			MutationP: GriewangkMutationP,
			ScaleMin:  GriewangkMin,
			ScaleMax:  GriewangkMax,
		}, nil
	case "ackley":
		return Params{
			Function:  Ackley,
			Label:     AckleyLabel,
			N:         AckleyN,
			MutationP: AckleyMutationP,
			ScaleMin:  AckleyMin,
			ScaleMax:  AckleyMax,
		}, nil
	case "rosenbrock":
		return Params{
			Function:  Rosenbrock,
			Label:     RosenbrockLabel,
			N:         RosenbrockN,
			MutationP: RosenbrockMutationP,
			ScaleMin:  RosenbrockMin,
			ScaleMax:  RosenbrockMax,
		}, nil
	}

	return Params{}, errors.New("invalid function passed to GetParams")
}

const (
	RastriginLabel     = "Rastrigin Function"
	RastriginN         = 20
	RastriginMin       = -5.12
	RastriginMax       = 5.12
	RastriginMutationP = float32(1) / (float32(16) * RastriginN)
)

func Rastrigin(x []uint16) float64 {
	xScaled := ScaleInputs(x[:], RastriginMin, RastriginMax)
	sum := 0.0
	for i := 0; i < RastriginN; i++ {
		sum += math.Pow(xScaled[i], 2.0) - 3.0*math.Cos(2.0*math.Pi*xScaled[i])
	}
	return 3.0*float64(RastriginN) + sum
}

const (
	SchwefelLabel     = "Schwefel Function"
	SchwefelN         = 10
	SchwefelMin       = -500.0
	SchwefelMax       = 500.0
	SchwefelMutationP = float32(1) / (float32(16) * SchwefelN)
)

// Schwefel Function differs to that in the paper, the paper has a mistake in a sign (+ve instead of -ve)
func Schwefel(x []uint16) float64 {
	xScaled := ScaleInputs(x[:], SchwefelMin, SchwefelMax)
	sum := 0.0
	for i := 0; i < SchwefelN; i++ {
		sum += xScaled[i] * math.Sin(math.Sqrt(math.Abs(xScaled[i])))
	}
	return 418.9829*float64(SchwefelN) - sum
}

const (
	GriewangkLabel     = "Griewangk Function"
	GriewangkN         = 10
	GriewangkMin       = -600.0
	GriewangkMax       = 600.0
	GriewangkMutationP = float32(1) / (float32(16) * GriewangkN)
)

func Griewangk(x []uint16) float64 {
	xScaled := ScaleInputs(x[:], GriewangkMin, GriewangkMax)
	sigma := 0.0
	product := 1.0
	for i := 0; i < GriewangkN; i++ {
		sigma += math.Pow(xScaled[i], 2) / 4000
		product *= math.Cos(xScaled[i] / math.Sqrt(float64(i+1)))
	}

	return 1.0 + sigma - product
}

const (
	AckleyLabel     = "Ackley Function"
	AckleyN         = 30
	AckleyMin       = -30.0
	AckleyMax       = 30.0
	AckleyMutationP = float32(1) / (float32(16) * AckleyN)
)

func Ackley(x []uint16) float64 {
	xScaled := ScaleInputs(x[:], AckleyMin, AckleyMax)
	sumA, sumB := 0.0, 0.0
	for i := 0; i < AckleyN; i++ {
		sumA += math.Pow(xScaled[i], 2)
		sumB += math.Cos(2.0 * math.Pi * xScaled[i])
	}

	sumA *= 1 / float64(AckleyN)
	sumB *= 1 / float64(AckleyN)

	return 20.0 + math.E - 20.0*math.Exp(-0.2*math.Sqrt(sumA)) - math.Exp(sumB)
}

const (
	RosenbrockLabel     = "Rosenbrock Function"
	RosenbrockN         = 40 // TODO: I made up this N dimensionality, it may need to be tuned if problem too easy or hard
	RosenbrockMin       = -2.048
	RosenbrockMax       = 2.048
	RosenbrockMutationP = float32(1) / (float32(16) * RosenbrockN)
)

func Rosenbrock(x []uint16) float64 {
	xScaled := ScaleInputs(x[:], RosenbrockMin, RosenbrockMax)
	sum := 0.0
	for i := 0; i < RosenbrockN/2; i++ {
		sum += math.Pow(100.0*(xScaled[2*i]-xScaled[2*i+1]), 2) + math.Pow(xScaled[2*i]-1, 2)
	}

	return sum
}

func TestFunc(x []uint16) float64 {
	return 0.0
}

// ScaleInputs scales an slice of uint16s between two ranges
func ScaleInputs(x []uint16, min float64, max float64) []float64 {
	var scaled []float64
	for i := 0; i < len(x); i++ {
		scaled = append(scaled, (float64(x[i])/65535)*(max-min)+min)
	}
	return scaled
}
