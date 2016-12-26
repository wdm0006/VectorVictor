package main

import (
	"math"
)

func aexp(arr []float64, exp float64) ([]float64, error) {
	for idx, v := range arr {
		arr[idx] = math.Pow(v, exp)
	}
	return arr, nil
}

func Square(arr []float64) ([]float64, error) {
	return aexp(arr, 2.0)
}