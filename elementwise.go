package main

import (
	"math"
)

// arrayExp raises each element of the array to the given exponent.
// Note: This modifies the input array in place and returns it.
func arrayExp(arr []float64, exp float64) ([]float64, error) {
	for idx, v := range arr {
		arr[idx] = math.Pow(v, exp)
	}
	return arr, nil
}

// Square returns a new array with each element squared.
func Square(arr []float64) ([]float64, error) {
	// Create a copy to avoid modifying the original
	result := make([]float64, len(arr))
	copy(result, arr)
	return arrayExp(result, 2.0)
}
