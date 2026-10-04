package main

import (
	"fmt"
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
	result, err := arrayExp(result, 2.0)
	if err != nil {
		return nil, err
	}
	for i, v := range result {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return nil, fmt.Errorf("square of element %d (%g) is not representable as a finite number", i, arr[i])
		}
	}
	return result, nil
}
