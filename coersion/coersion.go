package coersion

import (
	"strconv"
	"strings"
)

// csv to float vector converts a single comma separated string into
// an array of 64 bit floats
func CSV2FloatArray (stringvec string) ([]float64, error) {
	whitelist := strings.Split("0123456789.,", "")
	var clean_stringvec string = ""
	for char, _ := range clean_stringvec {
		for wl, _ := range whitelist {
			if wl == char {
				clean_stringvec += char
				break
			}
		}
	}
	var arr = strings.Split(clean_stringvec, ",")
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
