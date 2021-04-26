package result

import (
	"encoding/json"
	"github.com/OscarVanL/COMP6026-Evolution-of-Complexity/assignment2/chart"
	"io/ioutil"
	"log"
	"math"
)

type FunctionResults struct {
	Function string
	GA       Result
	CCGA1    Result
	CCGAHC   Result
}

type Result struct {
	Fitnesses []float64
	Mean      float64
	StdDev    float64
}

// WriteResults is used to write the experiment results to a JSON file, including mean and standard deviation calculations.
func WriteResults(output string, res [][]chart.EvolutionResults) {
	var allResults []FunctionResults

	// Create results JSON for each function's runs
	for i := 0; i < len(res); i++ {
		currResult := res[i]

		GAFitnesses, CCGAFitnesses, CCGAHCFitnesses := getFinalFitnesses(currResult)
		GAMean, CCGAMean, CCGAHCMean := getMean(GAFitnesses, CCGAFitnesses, CCGAHCFitnesses)
		GAStdDev, CCGAStdDev, CCGAHCStdDev := getStdDev(GAFitnesses, GAMean, CCGAFitnesses, CCGAMean, CCGAHCFitnesses, CCGAHCMean)

		result := FunctionResults{
			Function: currResult[0].Title,
			GA: Result{
				Fitnesses: GAFitnesses,
				Mean:      GAMean,
				StdDev:    GAStdDev,
			},
			CCGA1: Result{
				Fitnesses: CCGAFitnesses,
				Mean:      CCGAMean,
				StdDev:    CCGAStdDev,
			},
			CCGAHC: Result{
				Fitnesses: CCGAHCFitnesses,
				Mean:      CCGAHCMean,
				StdDev:    CCGAHCStdDev,
			}}

		allResults = append(allResults, result)
	}

	b, err := json.MarshalIndent(allResults, "", "    ")
	if err != nil {
		log.Fatal("Unable to encode results as JSON:", err)
	}
	err = ioutil.WriteFile(output+".json", b, 0644)
	if err != nil {
		log.Fatal("Unable to write JSON to file:", err)
	}
}

func getFinalFitnesses(result []chart.EvolutionResults) ([]float64, []float64, []float64) {
	var GAFitnesses, CCGAFitnesses, CCGAHCFitnesses []float64
	for i := 0; i < len(result); i++ {
		hist := result[i].GAFitnessHistory
		GAFitnesses = append(GAFitnesses, hist[len(hist)-1].Fitness)
		hist = result[i].CCGAFitnessHistory
		CCGAFitnesses = append(CCGAFitnesses, hist[len(hist)-1].Fitness)
		hist = result[i].CCGAHCFitnessHistory
		CCGAHCFitnesses = append(CCGAHCFitnesses, hist[len(hist)-1].Fitness)
	}
	return GAFitnesses, CCGAFitnesses, CCGAHCFitnesses
}

func getMean(GAFitnesses []float64, CCGAFitnesses []float64, CCGAHCFitnesses []float64) (float64, float64, float64) {
	var GASum, CCGASum, CCGAHCSum float64
	for i := 0; i < len(GAFitnesses); i++ {
		GASum += GAFitnesses[i]
		CCGASum += CCGAFitnesses[i]
		CCGAHCSum += CCGAHCFitnesses[i]
	}

	return GASum / float64(len(GAFitnesses)), CCGASum / float64(len(CCGAFitnesses)), CCGAHCSum / float64(len(CCGAHCFitnesses))
}

func getStdDev(GAFitnesses []float64, GAMean float64, CCGAFitnesses []float64, CCGAMean float64, CCGAHCFitnesses []float64, CCGAHCMean float64) (float64, float64, float64) {
	var GAStdDev, CCGAStdDev, CCGAHCStdDev float64

	for i := 0; i < len(GAFitnesses); i++ {
		GAStdDev += math.Pow(GAFitnesses[i]-GAMean, 2)
		CCGAStdDev += math.Pow(CCGAFitnesses[i]-CCGAMean, 2)
		CCGAHCStdDev += math.Pow(CCGAHCFitnesses[i]-CCGAHCMean, 2)
	}

	GAStdDev /= float64(len(GAFitnesses))
	CCGAStdDev /= float64(len(CCGAFitnesses))
	CCGAHCStdDev /= float64(len(CCGAHCFitnesses))

	return GAStdDev, CCGAStdDev, CCGAHCStdDev
}
