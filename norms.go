package main

import (
	"math"
)

// arrayMin returns the minimum value in a slice of float64.
func arrayMin(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0, nil
	}
	smallest := arr[0]
	for _, v := range arr[1:] {
		if v < smallest {
			smallest = v
		}
	}
	return smallest, nil
}

// arrayMax returns the maximum value in a slice of float64.
func arrayMax(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0, nil
	}
	largest := arr[0]
	for _, v := range arr[1:] {
		if v > largest {
			largest = v
		}
	}
	return largest, nil
}

// L0 calculates the L0 "norm" of a vector.
// Counts the number of non-zero elements.
// Note: L0 is not a true norm (doesn't satisfy triangle inequality).
// Time complexity: O(n)
func L0(arr []float64) (float64, error) {
	var count float64
	for _, v := range arr {
		if v != 0 {
			count++
		}
	}
	return count, nil
}

// L1 calculates the L1 (Manhattan) norm of a vector.
// Formula: ||x||₁ = Σ|xᵢ|
// Time complexity: O(n)
func L1(arr []float64) (float64, error) {
	var sum float64
	for _, v := range arr {
		sum += math.Abs(v)
	}
	return sum, nil
}

// L2 calculates the L2 (Euclidean) norm of a vector.
// Formula: ||x||₂ = √(Σxᵢ²)
// Time complexity: O(n)
func L2(arr []float64) (float64, error) {
	var sum float64
	for _, v := range arr {
		sum += v * v
	}
	return math.Sqrt(sum), nil
}

// Linfinity calculates the L-infinity (maximum) norm of a vector.
// Formula: ||x||∞ = max(|xᵢ|)
// Time complexity: O(n)
func Linfinity(arr []float64) (float64, error) {
	if len(arr) == 0 {
		return 0, nil
	}
	maxVal := math.Abs(arr[0])
	for _, v := range arr[1:] {
		if absV := math.Abs(v); absV > maxVal {
			maxVal = absV
		}
	}
	return maxVal, nil
}

// Lhalf calculates the L0.5 (sub-unitary) quasi-norm of a vector.
// Formula: (Σ√|xᵢ|)²
// Promotes extreme sparsity; used in sparse signal recovery.
// Note: This is a quasi-norm, not a true norm.
// Time complexity: O(n)
func Lhalf(arr []float64) (float64, error) {
	var sum float64
	for _, v := range arr {
		sum += math.Sqrt(math.Abs(v))
	}
	return sum * sum, nil
}

// LN calculates the general L-N norm of a vector for any N > 0.
// Formula: ||x||ₙ = (Σ|xᵢ|^N)^(1/N)
// For N >= maxN, returns L-infinity norm.
// Time complexity: O(n), but slower than L1/L2 due to math.Pow
func LN(arr []float64, N float64) (float64, error) {
	if len(arr) == 0 {
		return 0, nil
	}

	// Fast paths for common cases
	switch N {
	case 1.0:
		return L1(arr)
	case 2.0:
		return L2(arr)
	}

	if N >= maxN {
		return Linfinity(arr)
	}

	var sumPowered float64
	for _, v := range arr {
		sumPowered += math.Pow(math.Abs(v), N)
	}
	return math.Pow(sumPowered, 1.0/N), nil
}

// Lp calculates the general Lp norm for any p >= 0.
// Formula: ||x||_p = (Σ|xᵢ|^p)^(1/p)
// Special cases: p=0 counts non-zero, p=1 is L1, p=2 is L2.
// Time complexity: O(n)
func Lp(arr []float64, p float64) (float64, error) {
	if len(arr) == 0 {
		return 0, nil
	}

	// Fast paths
	switch p {
	case 0:
		return L0(arr)
	case 0.5:
		return Lhalf(arr)
	case 1.0:
		return L1(arr)
	case 2.0:
		return L2(arr)
	}

	if math.IsInf(p, 1) || p >= maxN {
		return Linfinity(arr)
	}

	var sumPowered float64
	for _, v := range arr {
		sumPowered += math.Pow(math.Abs(v), p)
	}
	return math.Pow(sumPowered, 1.0/p), nil
}

// WeightedL2 calculates the weighted L2 norm of a vector.
// Formula: √(Σwᵢ·xᵢ²)
// If weights is shorter than arr, remaining elements use weight 1.0.
// Time complexity: O(n)
func WeightedL2(arr []float64, weights []float64) (float64, error) {
	var sum float64
	for i, v := range arr {
		w := 1.0
		if i < len(weights) {
			w = weights[i]
		}
		sum += w * v * v
	}
	return math.Sqrt(sum), nil
}

// Mahalanobis calculates a diagonal Mahalanobis distance.
// Formula: √(Σ(xᵢ²/σᵢ²))
// If variances is shorter than arr, remaining elements use variance 1.0.
// For full Mahalanobis, you need the complete covariance matrix.
// Time complexity: O(n)
func Mahalanobis(arr []float64, variances []float64) (float64, error) {
	var sum float64
	for i, v := range arr {
		variance := 1.0
		if i < len(variances) && variances[i] != 0 {
			variance = variances[i]
		}
		sum += (v * v) / variance
	}
	return math.Sqrt(sum), nil
}
