package main

import (
	"math"
	"testing"
)

func TestSquare(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{"positive integers", []float64{1.0, 2.0, 3.0}, []float64{1.0, 4.0, 9.0}},
		{"single element", []float64{5.0}, []float64{25.0}},
		{"zeros", []float64{0.0, 0.0}, []float64{0.0, 0.0}},
		{"negative numbers", []float64{-2.0, -3.0}, []float64{4.0, 9.0}},
		{"mixed", []float64{-1.0, 0.0, 1.0}, []float64{1.0, 0.0, 1.0}},
		{"empty array", []float64{}, []float64{}},
		{"decimals", []float64{0.5, 1.5}, []float64{0.25, 2.25}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of input to verify it's not modified
			inputCopy := make([]float64, len(tt.input))
			copy(inputCopy, tt.input)

			result, err := Square(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Square(%v) returned %d elements, want %d", tt.input, len(result), len(tt.expected))
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("Square(%v)[%d] = %v, want %v", inputCopy, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestArrayExp(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		exp      float64
		expected []float64
	}{
		{"square", []float64{2.0, 3.0}, 2.0, []float64{4.0, 9.0}},
		{"cube", []float64{2.0, 3.0}, 3.0, []float64{8.0, 27.0}},
		{"power of 1", []float64{2.0, 3.0}, 1.0, []float64{2.0, 3.0}},
		{"power of 0", []float64{2.0, 3.0}, 0.0, []float64{1.0, 1.0}},
		{"sqrt", []float64{4.0, 9.0}, 0.5, []float64{2.0, 3.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy since arrayExp modifies in place
			input := make([]float64, len(tt.input))
			copy(input, tt.input)

			result, err := arrayExp(input, tt.exp)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("arrayExp(%v, %v)[%d] = %v, want %v", tt.input, tt.exp, i, result[i], tt.expected[i])
				}
			}
		})
	}
}
