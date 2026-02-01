package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

// Version is the application version.
const Version = "1.0.0"

// maxN is the maximum N value supported for exponents before falling back to L-infinity.
const maxN float64 = 100.0

// wrapResponse wraps the content block of a response with metadata including
// timestamp, version info, and any errors.
func wrapResponse(content gin.H, err error) gin.H {
	var errStr interface{}
	if err != nil {
		errStr = err.Error()
	}
	return gin.H{
		"info":    gin.H{"now": time.Now()},
		"content": content,
		"errors":  errStr,
		"version": Version,
	}
}

// index returns a JSON response indicating the system is up.
func index(c *gin.Context) {
	content := gin.H{"message": "Vector-Victor: go do some math."}
	response := wrapResponse(content, nil)
	c.JSON(http.StatusOK, response)
}

// square takes in a vector and returns its element-wise square.
func square(c *gin.Context) {
	stringvec := c.Query("v")

	arr, err := CSV2FloatArray(stringvec)
	if err != nil {
		content := gin.H{"val": arr}
		c.JSON(http.StatusBadRequest, wrapResponse(content, err))
		return
	}

	result, err := Square(arr)
	if err != nil {
		content := gin.H{"val": arr}
		c.JSON(http.StatusInternalServerError, wrapResponse(content, err))
		return
	}

	content := gin.H{"val": result}
	c.JSON(http.StatusOK, wrapResponse(content, nil))
}

// norm calculates the specified norm of a vector.
func norm(c *gin.Context) {
	stringvec := c.Query("v")
	kind := strings.ToLower(c.DefaultQuery("kind", "l2"))

	arr, err := CSV2FloatArray(stringvec)
	if err != nil {
		content := gin.H{"vector": arr, "norm": nil, "kind": kind}
		c.JSON(http.StatusBadRequest, wrapResponse(content, err))
		return
	}

	var normVal float64
	switch kind {
	case "l0":
		normVal, err = L0(arr)
	case "l1":
		normVal, err = L1(arr)
	case "l2":
		normVal, err = L2(arr)
	case "linfinity":
		normVal, err = Linfinity(arr)
	case "lhalf", "l0.5":
		normVal, err = Lhalf(arr)
	case "lp":
		// Get p parameter, default to 2
		pStr := c.DefaultQuery("p", "2")
		p, parseErr := strconv.ParseFloat(pStr, 64)
		if parseErr != nil {
			err = fmt.Errorf("invalid p value: %s", pStr)
		} else {
			normVal, err = Lp(arr, p)
		}
	case "weighted", "weightedl2":
		// Get weights parameter
		weightsStr := c.Query("weights")
		var weights []float64
		if weightsStr != "" {
			weights, err = CSV2FloatArray(weightsStr)
			if err != nil {
				content := gin.H{"vector": arr, "norm": nil, "kind": kind}
				c.JSON(http.StatusBadRequest, wrapResponse(content, fmt.Errorf("invalid weights: %v", err)))
				return
			}
		}
		normVal, err = WeightedL2(arr, weights)
	case "mahalanobis":
		// Get variances parameter
		variancesStr := c.Query("variances")
		var variances []float64
		if variancesStr != "" {
			variances, err = CSV2FloatArray(variancesStr)
			if err != nil {
				content := gin.H{"vector": arr, "norm": nil, "kind": kind}
				c.JSON(http.StatusBadRequest, wrapResponse(content, fmt.Errorf("invalid variances: %v", err)))
				return
			}
		}
		normVal, err = Mahalanobis(arr, variances)
	default:
		err = fmt.Errorf("unknown norm kind: %s (valid: l0, l1, l2, linfinity, lhalf, lp, weighted, mahalanobis)", kind)
	}

	if err != nil {
		content := gin.H{"vector": arr, "norm": normVal, "kind": kind}
		c.JSON(http.StatusBadRequest, wrapResponse(content, err))
		return
	}

	content := gin.H{"vector": arr, "norm": normVal, "kind": kind}
	c.IndentedJSON(http.StatusOK, wrapResponse(content, nil))
}

// loadTemplates loads HTML templates from the embedded filesystem.
func loadTemplates(list ...string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	for _, name := range list {
		templateBytes, err := templatesFS.ReadFile("templates/" + name)
		if err != nil {
			log.Fatalf("failed to read template %s: %v", name, err)
		}

		tmpl, err := template.New(name).Parse(string(templateBytes))
		if err != nil {
			log.Fatalf("failed to parse template %s: %v", name, err)
		}

		r.Add(name, tmpl)
	}
	return r
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.New()
	app.Use(gin.Logger())
	app.Use(gin.Recovery())

	app.HTMLRender = loadTemplates("index.tmpl", "square.tmpl", "norms.tmpl")

	// Index page
	app.POST("/", index)
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "VectorVictor",
		})
	})

	// Square endpoint
	app.POST("/square", square)
	app.GET("/square", func(c *gin.Context) {
		c.HTML(http.StatusOK, "square.tmpl", gin.H{
			"title": "VectorVictor: Square",
		})
	})

	// Norm endpoint
	app.POST("/norm", norm)
	app.GET("/norm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "norms.tmpl", gin.H{
			"title": "VectorVictor: Norms",
		})
	})

	// Health check endpoint
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	srv := &http.Server{
		Addr:           ":8080",
		Handler:        app,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
