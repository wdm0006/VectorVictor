package main

import (
    	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"./norms"
	"./coersion"
	"fmt"
	"errors"
	"net/http"
	"strings"
)

// Wrap response is a helper function to wrap the content block of
// a response with info blocks (containing things like local time
// and the version of the code, and placing the content in a nested
// content block.
func wrap_response(content gin.H, errors error) (gin.H){
	response := gin.H{
		"info": gin.H{"now": time.Now()},
		"content": content,
		"errors": errors,
	}
	return response
}

// a middleware for adding in the start time of the request
func TimerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		t := time.Now()

		// process request
		c.Next()

		// after request
		latency := time.Since(t)
		// TODO: append this latency to the response itself
		print(latency.String())
	}
}



// Index simply responds with an indication that the system is up.
func index (c *gin.Context){
	content := gin.H{"message": "Vector-Victor: go do some math."}
	response := wrap_response(content, nil)
	c.JSON(200, response)
}

// Square takes in a single value, and returns it's square.
func square (c *gin.Context){
	x, err := strconv.ParseFloat(c.Param("x"), 64)

	var response gin.H
	var code int

	if err != nil {
		content := gin.H{"val": x}
		response = wrap_response(content, err)
		code = 500
	} else {
		content := gin.H{"val": x, "square": x * x}
		response = wrap_response(content, nil)
		code = 200
	}

	c.JSON(code, response)
}

// norm takes the l2 norm of a vector
func norm (c *gin.Context) {
	// parse the input into a vector of floats
	stringvec := c.Query("v")
	kind := c.DefaultQuery("kind", "l2")
	kind = strings.ToLower(kind)

	arr, err := coersion.CSV2FloatArray(stringvec)

	var response gin.H
	var code int
	if err != nil {
		content := gin.H{"vector": arr, "norm": nil, "kind": kind}
		response = wrap_response(content, err)
		code = 500
	} else {
		var norm float64
		var err error
		if kind == "l2" {
			norm, err = norms.L2(arr)
		} else if kind == "l1"{
			norm, err = norms.L1(arr)
		} else if kind == "linfinity" {
			norm, err = norms.Linfinity(arr)
		} else {
			err = errors.New(fmt.Sprintf("unknown norm kind: %s", kind))
		}


		if err != nil {
			content := gin.H{"vector": arr, "norm": norm, "kind": kind}
			response = wrap_response(content, err)
			code = 500
		} else {
			content := gin.H{"vector": arr, "norm": norm, "kind": kind}
			response = wrap_response(content, nil)
			code = 200
		}
	}

	c.IndentedJSON(code, response)

}


func main(){
	app := gin.New()

	// middlwares
	app.Use(gin.Logger())
	app.Use(gin.Recovery())
	app.Use(TimerMiddleWare())

	app.LoadHTMLGlob("templates/**/*")

	// index page, post returns json status, get will return an html status page
	app.POST("/", index)
        app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "public/index.tmpl", gin.H{
			"title": "hello world",
		})
	})

	// square: a post will return json response and get will be an html page
	app.POST("/square/:x", square)
	app.GET("/square/:x", func(c *gin.Context) {
		c.HTML(http.StatusOK, "public/square.tmpl", gin.H{
			"title": "hello world",
		})
	})

	// norm: a post will return json response and get will be an html page
	app.POST("/norm", norm)
	app.GET("/norm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "public/norms.tmpl", gin.H{
			"title": "hello world",
		})
	})

	// run the server
	s := &http.Server{
		Addr:           ":8000",
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}