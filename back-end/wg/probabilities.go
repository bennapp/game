package wg

import (
	"fmt"
	"math/rand"
)

func weightsToProbMap(mapping map[string]float64) map[string]float64 {
	sum := 0.0

	for _, v := range mapping {
		sum += v
	}

	probmap := make(map[string]float64)
	for k, v := range mapping {
		probmap[k] = float64(v) / float64(sum)
	}

	return probmap
}

func reWeight(weights map[string]float64, probs map[string]float64) map[string]float64 {
	weightedProbMap := weightsToProbMap(weights)

	probmap := make(map[string]float64)
	for k := range probs {
		probmap[k] = (weightedProbMap[k] + probs[k]) / 2
	}

	return probmap
}

func scaleProbMap(probMap map[string]float64, scale float64) map[string]float64 {
	scaledProbMap := make(map[string]float64)
	for k, v := range probMap {
		scaledProbMap[k] = v * scale
	}

	return scaledProbMap
}

func addProbMap(firstProbMap map[string]float64, secondProbMap map[string]float64) map[string]float64 {
	summedProbMap := make(map[string]float64)
	for k := range firstProbMap {
		summedProbMap[k] = firstProbMap[k] + secondProbMap[k]
	}

	return summedProbMap
}

func randomValueFromProbMap(probMap map[string]float64) string {
	randomNumber := rand.Float64()

	var value string

	for k, v := range probMap {
		randomNumber -= v
		value = k

		if randomNumber < 0 {
			break
		}
	}

	return value
}

func Weights() {
	weightmap := make(map[string]float64)
	weightmap["red"] = 2
	weightmap["blue"] = 3
	weightmap["green"] = 5

	probmap := weightsToProbMap(weightmap)
	fmt.Println(probmap)

	weights := make(map[string]float64)
	weights["red"] = 4
	weights["blue"] = 2
	weights["green"] = 1

	weightedProbMap := reWeight(weights, probmap)

	scale := 0.5
	scaledMap := scaleProbMap(weightedProbMap, scale)

	fmt.Println(scaledMap)

	summedMap := addProbMap(scaledMap, scaledMap)
	fmt.Println(summedMap)

	value := randomValueFromProbMap(summedMap)
	fmt.Println(value)
}
