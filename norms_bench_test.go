package main

import (
	"math"
	"testing"
)

// Benchmark data
var benchVector = make([]float64, 10000)

func init() {
	for i := range benchVector {
		benchVector[i] = float64(i%100) - 50 // Mix of positive and negative
	}
}

// Current implementations use math.Pow for everything
func BenchmarkL2Current(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = L2(benchVector)
	}
}

func BenchmarkL1Current(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = L1(benchVector)
	}
}

func BenchmarkLinfinityCurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Linfinity(benchVector)
	}
}

// Optimized versions for comparison
func l2Optimized(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += v * v // Direct multiplication instead of math.Pow
	}
	return math.Sqrt(sum)
}

func l1Optimized(arr []float64) float64 {
	var sum float64
	for _, v := range arr {
		sum += math.Abs(v) // Direct abs instead of math.Pow
	}
	return sum
}

func lInfinityOptimized(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	}
	maxVal := math.Abs(arr[0])
	for _, v := range arr[1:] {
		if absV := math.Abs(v); absV > maxVal {
			maxVal = absV
		}
	}
	return maxVal
}

func BenchmarkL2Optimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l2Optimized(benchVector)
	}
}

func BenchmarkL1Optimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l1Optimized(benchVector)
	}
}

func BenchmarkLInfinityOptimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lInfinityOptimized(benchVector)
	}
}

// Benchmark the other norms too
func BenchmarkL0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = L0(benchVector)
	}
}

func BenchmarkLhalf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Lhalf(benchVector)
	}
}

func BenchmarkWeightedL2(b *testing.B) {
	weights := make([]float64, len(benchVector))
	for i := range weights {
		weights[i] = 1.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = WeightedL2(benchVector, weights)
	}
}

func BenchmarkMahalanobis(b *testing.B) {
	variances := make([]float64, len(benchVector))
	for i := range variances {
		variances[i] = 1.0
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Mahalanobis(benchVector, variances)
	}
}

func BenchmarkLpGeneral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Lp(benchVector, 3.0)
	}
}
