package coersion

import (
	"strconv"
)

// csv to float vector converts a single comma seperated string into
// an array of 64 bit floats
func CSV2FloatArray (arr []string) ([]float64, error) {
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
