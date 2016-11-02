package norms

import (
	"math"
)

// the max N value supported for exponents
const MAX_N float64 = 100.0


func array_min(arr []float64) (float64, error) {
	smallest := arr[0]
	for _, v := range arr {
		if v < smallest {
			smallest = v
		}
	}
	return smallest, nil
}

func array_max(arr []float64) (float64, error) {
	largest := arr[0]
	for _, v := range arr {
		if v > largest {
			largest = v
		}
	}
	return largest, nil
}

func LN(arr []float64, N float64) (float64, error){
	// we put a max value on N and assume L-infinity above that
	if N >= MAX_N {
		return array_max(arr)
	} else {
		// square all elements of the array
		sqrs := []float64{}
		for _, v := range arr {
			sqrs = append(sqrs, math.Pow(v, N))
		}

		// sum it
		var sum_squared = 0.0
		for _, v := range sqrs {
			sum_squared += v
		}

		// take the square root
		norm := math.Pow(sum_squared, 1.0 / N)

		return norm, nil
	}
}

// l2 norm is the sqrt of sum of squares of a vector
func L2(arr []float64) (float64, error){
	return LN(arr, 2.0)
}

func L1(arr []float64) (float64, error){
	return LN(arr, 1.0)
}

func Linfinity(arr []float64) (float64, error){
	return LN(arr, MAX_N + 1)
}