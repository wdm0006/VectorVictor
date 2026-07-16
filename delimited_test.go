package main

import (
	"math"
	"testing"
)

func TestCSV2FloatArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []float64
		wantErr  bool
	}{
		{"simple integers", "1,2,3", []float64{1.0, 2.0, 3.0}, false},
		{"with spaces", "1, 2, 3", []float64{1.0, 2.0, 3.0}, false},
		{"decimals", "1.5,2.5,3.5", []float64{1.5, 2.5, 3.5}, false},
		{"single value", "42", []float64{42.0}, false},
		{"empty string", "", []float64{}, false},
		{"negative numbers", "-1,-2,-3", []float64{-1.0, -2.0, -3.0}, false},
		{"mixed signs", "-1,0,1", []float64{-1.0, 0.0, 1.0}, false},
		{"scientific notation", "1e3,2", []float64{1000.0, 2.0}, false},
		{"scientific notation with negative exponent", "1.5e-3", []float64{0.0015}, false},
		{"scientific notation with uppercase exponent", "2E2", []float64{200.0}, false},
		{"explicit positive sign", "+2,3", []float64{2.0, 3.0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CSV2FloatArray(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CSV2FloatArray(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("CSV2FloatArray(%q) unexpected error: %v", tt.input, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("CSV2FloatArray(%q) = %v (len %d), want %v (len %d)",
					tt.input, result, len(result), tt.expected, len(tt.expected))
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("CSV2FloatArray(%q)[%d] = %v, want %v", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestTSV2FloatArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []float64
	}{
		{"tab separated", "1\t2\t3", []float64{1.0, 2.0, 3.0}},
		{"single value", "42", []float64{42.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TSV2FloatArray(tt.input)
			if err != nil {
				t.Errorf("TSV2FloatArray(%q) unexpected error: %v", tt.input, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("TSV2FloatArray(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("TSV2FloatArray(%q)[%d] = %v, want %v", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestPSV2FloatArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []float64
	}{
		{"pipe separated", "1|2|3", []float64{1.0, 2.0, 3.0}},
		{"single value", "42", []float64{42.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PSV2FloatArray(tt.input)
			if err != nil {
				t.Errorf("PSV2FloatArray(%q) unexpected error: %v", tt.input, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("PSV2FloatArray(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("PSV2FloatArray(%q)[%d] = %v, want %v", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestWhitelistString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		whitelist string
		expected  string
	}{
		{"remove letters", "a1b2c3", "0123456789", "123"},
		{"keep all", "123", "0123456789", "123"},
		{"remove all", "abc", "0123456789", ""},
		{"mixed", "1.5,2.5", "0123456789.,", "1.5,2.5"},
		{"empty input", "", "0123456789", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := whitelistString(tt.input, tt.whitelist)
			if result != tt.expected {
				t.Errorf("whitelistString(%q, %q) = %q, want %q", tt.input, tt.whitelist, result, tt.expected)
			}
		})
	}
}

func TestDelimited2FloatArray(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []float64
		wantErr   bool
	}{
		{"comma delimiter", "1,2,3", ",", []float64{1.0, 2.0, 3.0}, false},
		{"semicolon delimiter", "1;2;3", ";", []float64{1.0, 2.0, 3.0}, false},
		{"with extra characters filtered", "1a,2b,3c", ",", []float64{1.0, 2.0, 3.0}, false},
		{"empty string", "", ",", []float64{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Delimited2FloatArray(tt.input, tt.delimiter)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Delimited2FloatArray(%q, %q) expected error", tt.input, tt.delimiter)
				}
				return
			}

			if err != nil {
				t.Errorf("Delimited2FloatArray(%q, %q) unexpected error: %v", tt.input, tt.delimiter, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Delimited2FloatArray(%q, %q) = %v, want %v", tt.input, tt.delimiter, result, tt.expected)
				return
			}

			for i := range result {
				if math.Abs(result[i]-tt.expected[i]) > 1e-10 {
					t.Errorf("Delimited2FloatArray(%q, %q)[%d] = %v, want %v", tt.input, tt.delimiter, i, result[i], tt.expected[i])
				}
			}
		})
	}
}
