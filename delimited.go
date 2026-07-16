package main

import (
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

// whitelistString removes any characters not in the whitelist from the input string.
func whitelistString(input string, whitelist string) string {
	var result strings.Builder
	result.Grow(len(input))

	for _, c := range input {
		if strings.ContainsRune(whitelist, c) {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// Delimited2FloatArray converts a delimited string into an array of float64.
// It first sanitizes the input to allow signed decimal and scientific notation.
func Delimited2FloatArray(stringvec string, delimiter string) ([]float64, error) {
	// Sanitize the input string
	cleanString := whitelistString(stringvec, "0123456789.,-+eE"+delimiter)

	if cleanString == "" {
		return []float64{}, nil
	}

	// Parse the cleaned string
	parts := strings.Split(cleanString, delimiter)
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
		result = append(result, val)
	}

	return result, nil
}
