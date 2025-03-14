package utils

import "math/rand"

func Contains[T comparable](slice []T, item T) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RandomGenerator(min float64, max float64) float64 {
	randIndex := rand.Float64()
	randomGenerator := (min + randIndex*(max-min))
	if randomGenerator == max {
		return randomGenerator - 1
	}
	return randomGenerator
}