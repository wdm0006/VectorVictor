package main

import (
	"strconv"
	"strings"
)

// csv to float vector converts a single comma separated string into
// an array of 64 bit floats
func CSV2FloatArray (stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, ",")
}

func TSV2FloatArray (stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, "\t")
}

func PSV2FloatArray (stringvec string) ([]float64, error) {
	return Delimited2FloatArray(stringvec, "|")
}

func whitelist_string(string_in string, whitelist []string) (string, error) {
	var clean_stringvec string = ""
	var stringvec []string = strings.Split(string_in, "")
 	for _, c_char := range stringvec {
		for _, wl_char := range whitelist {
			if wl_char == c_char {
				clean_stringvec += c_char
				break
			}
		}
	}
	return clean_stringvec, nil
}

func Delimited2FloatArray (stringvec string, delimiter string) ([]float64, error) {
	// first clean the string
	var whitelist []string = strings.Split("0123456789.,", "")
	clean_stringvec, err := whitelist_string(stringvec, whitelist)
	if err != nil {
		return []float64{0.0}, nil	// TODO: return an error here
	}

	// then parse it
	var arr = strings.Split(clean_stringvec, delimiter)
	var v = []float64{}
	for _, i := range arr {
		j, err := strconv.ParseFloat(i, 64)
		if err != nil {
			return v, err
		}
		v = append(v, j)
	}
	return v, nil
}
