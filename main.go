package main

import (
    	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
	"strings"
	"time"
	"./norms"
)

// Wrap response is a helper function to wrap the content block of
// a response with info blocks (containing things like local time
// and the version of the code, and placing the content in a nested
// content block.
func wrap_response(content gin.H) (gin.H){
	response := gin.H{
		"info": gin.H{"now": time.Now()},
		"content": content,
	}
	return response
}

// csv to float vector converts a single comma seperated string into
// an array of 64 bit floats
func csv2floatarray (arr []string) ([]float64, error) {
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

// Index simply responds with an indication that the system is up.
func index (c *gin.Context){
	content := gin.H{"message": "Vector-Victor: go do some math."}
	response := wrap_response(content)
	c.JSON(200, response)
}

// Square takes in a single value, and returns it's square.
func square (c *gin.Context){
	x, err := strconv.ParseFloat(c.Param("x"), 64)
	if err != nil {
		fmt.Println("some error")
	}
	content := gin.H{"val": x, "square": x * x}

	response := wrap_response(content)
	c.JSON(200, response)
}

// norm takes the l2 norm of a vector
func norm (c *gin.Context){
	// parse the input into a vector of floats
	stringvec := c.Param("v")
	arr, err := csv2floatarray(strings.Split(stringvec, ","))
	if err != nil {
		panic(err)
	}

	norm, err := norms.L2(arr)
	if err != nil {
		panic(err)
	}

	content := gin.H{"vector": arr, "norm": norm}
	response := wrap_response(content)
	c.JSON(200, response)
}


func main(){
	app := gin.Default()

	// just a vanilla get with no parameters
	app.GET("/", index)

	// a get with some parameters
	app.GET("/square/:x", square)
	app.GET("/norm/:v", norm)

	app.Run(":8000")
}