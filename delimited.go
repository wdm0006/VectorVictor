package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// CSV2FloatArray converts a comma-separated string into an array of float64.
func CSV2FloatArray(stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, ",")
}

// TSV2FloatArray converts a tab-separated string into an array of float64.
func TSV2FloatArray(stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, "\t")
}

// PSV2FloatArray converts a pipe-separated string into an array of float64.
func PSV2FloatArray(stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, "|")
}

// Delimited2FloatArray converts a delimited string into an array of float64.
// It splits on the delimiter, trims surrounding whitespace from each token, and
// parses each non-empty token unchanged. A token containing unsupported
// characters returns the ParseFloat error instead of being silently rewritten.
// NaN and infinite tokens are rejected.
// Empty input and empty fields are skipped, yielding an empty slice.
func Delimited2FloatArray(stringvec string, delimiter string) ([]float64, error) {
	parts := strings.Split(stringvec, delimiter)
	result := make([]float64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		val, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return result, err
		}
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return result, fmt.Errorf("non-finite value %q: only finite numbers are accepted", part)
		}
		result = append(result, val)
	}

	return result, nil
}
