package norms

import (
	"math"
)

// l2 norm is the sqrt of sum of squares of a vector TODO: any exception
// handling whatsoever.
func L2(arr []float64) (float64, error){
	// square all elements of the array
	sqrs := []float64{}
	for _, v := range arr {
		sqrs = append(sqrs, v * v)
	}

	// sum it
	var sum_squared = 0.0
	for _, v := range sqrs {
		sum_squared += v
	}

	// take the square root
	l2_norm := math.Sqrt(sum_squared)

	return l2_norm, nil
}