package main

import (
	"math"
	"testing"
)

func TestArrayMin(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"single element", []float64{5.0}, 5.0},
		{"positive numbers", []float64{3.0, 1.0, 4.0, 1.0, 5.0}, 1.0},
		{"negative numbers", []float64{-3.0, -1.0, -4.0}, -4.0},
		{"mixed numbers", []float64{-2.0, 0.0, 2.0}, -2.0},
		{"empty array", []float64{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := arrayMin(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("arrayMin(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestArrayMax(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"single element", []float64{5.0}, 5.0},
		{"positive numbers", []float64{3.0, 1.0, 4.0, 1.0, 5.0}, 5.0},
		{"negative numbers", []float64{-3.0, -1.0, -4.0}, -1.0},
		{"mixed numbers", []float64{-2.0, 0.0, 2.0}, 2.0},
		{"empty array", []float64{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := arrayMax(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("arrayMax(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestL2(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"3-4-5 triangle", []float64{3.0, 4.0}, 5.0},
		{"unit vector", []float64{1.0, 0.0, 0.0}, 1.0},
		{"equal components", []float64{1.0, 1.0, 1.0, 1.0}, 2.0},
		{"single element", []float64{5.0}, 5.0},
		{"empty array", []float64{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := L2(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("L2(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestL1(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"positive numbers", []float64{1.0, 2.0, 3.0}, 6.0},
		{"single element", []float64{5.0}, 5.0},
		{"zeros", []float64{0.0, 0.0, 0.0}, 0.0},
		{"empty array", []float64{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := L1(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("L1(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLinfinity(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"positive numbers", []float64{1.0, 5.0, 3.0}, 5.0},
		{"single element", []float64{7.0}, 7.0},
		{"all same", []float64{2.0, 2.0, 2.0}, 2.0},
		{"empty array", []float64{}, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Linfinity(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("Linfinity(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLN(t *testing.T) {
	tests := []struct {
		name     string
		arr      []float64
		n        float64
		expected float64
	}{
		{"L1 via LN", []float64{1.0, 2.0, 3.0}, 1.0, 6.0},
		{"L2 via LN", []float64{3.0, 4.0}, 2.0, 5.0},
		{"high N falls back to max", []float64{1.0, 5.0, 3.0}, 150.0, 5.0},
		{"positive infinity falls back to max", []float64{1.0, 5.0, 3.0}, math.Inf(1), 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LN(tt.arr, tt.n)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("LN(%v, %v) = %v, want %v", tt.arr, tt.n, result, tt.expected)
			}
		})
	}
}

func TestLNInvalidExponent(t *testing.T) {
	tests := []struct {
		name string
		arr  []float64
		n    float64
	}{
		{"zero N", []float64{1.0, 2.0}, 0.0},
		{"negative N", []float64{1.0, 2.0}, -1.0},
		{"negative infinity N", []float64{1.0, 2.0}, math.Inf(-1)},
		{"NaN N", []float64{1.0, 2.0}, math.NaN()},
		{"empty array with invalid N", []float64{}, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LN(tt.arr, tt.n)
			if err == nil {
				t.Errorf("LN(%v, %v) = %v, want an error", tt.arr, tt.n, result)
			}
			if result != 0 {
				t.Errorf("LN(%v, %v) = %v, want 0 alongside the error", tt.arr, tt.n, result)
			}
		})
	}
}

func TestL0(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"all non-zero", []float64{1.0, 2.0, 3.0}, 3.0},
		{"with zeros", []float64{1.0, 0.0, 3.0, 0.0, 5.0}, 3.0},
		{"all zeros", []float64{0.0, 0.0, 0.0}, 0.0},
		{"single non-zero", []float64{5.0}, 1.0},
		{"single zero", []float64{0.0}, 0.0},
		{"empty array", []float64{}, 0.0},
		{"negative values", []float64{-1.0, 0.0, -3.0}, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := L0(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("L0(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLp(t *testing.T) {
	tests := []struct {
		name     string
		arr      []float64
		p        float64
		expected float64
	}{
		{"Lp with p=1 equals L1", []float64{1.0, 2.0, 3.0}, 1.0, 6.0},
		{"Lp with p=2 equals L2", []float64{3.0, 4.0}, 2.0, 5.0},
		{"Lp with p=3", []float64{1.0, 2.0, 2.0}, 3.0, math.Pow(17, 1.0/3.0)},
		{"Lp with p=0 counts non-zero", []float64{1.0, 0.0, 3.0}, 0.0, 2.0},
		{"empty array", []float64{}, 2.0, 0.0},
		{"negative values with p=2", []float64{-3.0, -4.0}, 2.0, 5.0},
		{"Lp with p=0.5 equals Lhalf", []float64{1.0, 4.0}, 0.5, 9.0},
		{"high p falls back to max", []float64{1.0, 5.0, 3.0}, 150.0, 5.0},
		{"positive infinity falls back to max", []float64{1.0, 5.0, 3.0}, math.Inf(1), 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Lp(tt.arr, tt.p)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("Lp(%v, %v) = %v, want %v", tt.arr, tt.p, result, tt.expected)
			}
		})
	}
}

func TestLpInvalidExponent(t *testing.T) {
	tests := []struct {
		name string
		arr  []float64
		p    float64
	}{
		{"negative p", []float64{1.0, 2.0}, -1.0},
		{"small negative p", []float64{1.0, 2.0}, -0.5},
		{"negative infinity p", []float64{1.0, 2.0}, math.Inf(-1)},
		{"NaN p", []float64{1.0, 2.0}, math.NaN()},
		{"empty array with invalid p", []float64{}, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Lp(tt.arr, tt.p)
			if err == nil {
				t.Errorf("Lp(%v, %v) = %v, want an error", tt.arr, tt.p, result)
			}
			if result != 0 {
				t.Errorf("Lp(%v, %v) = %v, want 0 alongside the error", tt.arr, tt.p, result)
			}
		})
	}
}

func TestLhalf(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{"simple case", []float64{1.0, 4.0}, 9.0}, // (sqrt(1) + sqrt(4))^2 = (1+2)^2 = 9
		{"single element", []float64{4.0}, 4.0},   // (sqrt(4))^2 = 4
		{"zeros", []float64{0.0, 0.0}, 0.0},
		{"empty array", []float64{}, 0.0},
		{"unit values", []float64{1.0, 1.0, 1.0}, 9.0}, // (1+1+1)^2 = 9
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Lhalf(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("Lhalf(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWeightedL2(t *testing.T) {
	tests := []struct {
		name     string
		arr      []float64
		weights  []float64
		expected float64
	}{
		{"uniform weights", []float64{3.0, 4.0}, []float64{1.0, 1.0}, 5.0},
		{"double first weight", []float64{3.0, 4.0}, []float64{4.0, 1.0}, math.Sqrt(36 + 16)}, // sqrt(4*9 + 1*16)
		{"no weights defaults to 1", []float64{3.0, 4.0}, []float64{}, 5.0},
		{"partial weights", []float64{3.0, 4.0}, []float64{1.0}, 5.0}, // second element uses weight 1
		{"empty array", []float64{}, []float64{}, 0.0},
		{"zero weights", []float64{3.0, 4.0}, []float64{0.0, 1.0}, 4.0}, // only second element counts
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WeightedL2(tt.arr, tt.weights)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("WeightedL2(%v, %v) = %v, want %v", tt.arr, tt.weights, result, tt.expected)
			}
		})
	}
}

func TestWeightedL2NegativeWeight(t *testing.T) {
	tests := []struct {
		name    string
		arr     []float64
		weights []float64
	}{
		{"single negative weight", []float64{3.0, 4.0}, []float64{-1.0, -1.0}},
		{"negative weight second position", []float64{3.0, 4.0}, []float64{1.0, -2.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := WeightedL2(tt.arr, tt.weights)
			if err == nil {
				t.Errorf("WeightedL2(%v, %v) = %v, expected non-nil error for negative weight", tt.arr, tt.weights, result)
			}
			if math.IsNaN(result) {
				t.Errorf("WeightedL2(%v, %v) returned NaN, expected 0 with error", tt.arr, tt.weights)
			}
		})
	}
}

func TestMahalanobis(t *testing.T) {
	tests := []struct {
		name      string
		arr       []float64
		variances []float64
		expected  float64
	}{
		{"uniform variances", []float64{3.0, 4.0}, []float64{1.0, 1.0}, 5.0},
		{"scaled variances", []float64{6.0, 8.0}, []float64{4.0, 4.0}, 5.0}, // sqrt(36/4 + 64/4) = sqrt(9+16) = 5
		{"no variances defaults to 1", []float64{3.0, 4.0}, []float64{}, 5.0},
		{"partial variances", []float64{3.0, 4.0}, []float64{1.0}, 5.0},
		{"empty array", []float64{}, []float64{}, 0.0},
		{"high variance reduces contribution", []float64{10.0, 10.0}, []float64{100.0, 100.0}, math.Sqrt(2)}, // sqrt(100/100 + 100/100)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Mahalanobis(tt.arr, tt.variances)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if math.Abs(result-tt.expected) > 1e-10 {
				t.Errorf("Mahalanobis(%v, %v) = %v, want %v", tt.arr, tt.variances, result, tt.expected)
			}
		})
	}
}

func TestMahalanobisNegativeVariance(t *testing.T) {
	tests := []struct {
		name      string
		arr       []float64
		variances []float64
	}{
		{"single negative variance", []float64{3.0, 4.0}, []float64{-1.0, -1.0}},
		{"negative variance second position", []float64{3.0, 4.0}, []float64{1.0, -2.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Mahalanobis(tt.arr, tt.variances)
			if err == nil {
				t.Errorf("Mahalanobis(%v, %v) = %v, expected non-nil error for negative variance", tt.arr, tt.variances, result)
			}
			if math.IsNaN(result) {
				t.Errorf("Mahalanobis(%v, %v) returned NaN, expected 0 with error", tt.arr, tt.variances)
			}
		})
	}
}
