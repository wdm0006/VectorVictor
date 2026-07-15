package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	app := gin.New()
	app.HTMLRender = loadTemplates("index.tmpl", "square.tmpl", "norms.tmpl")

	app.POST("/", index)
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "VectorVictor"})
	})
	app.POST("/square", square)
	app.GET("/square", func(c *gin.Context) {
		c.HTML(http.StatusOK, "square.tmpl", gin.H{"title": "VectorVictor: Square"})
	})
	app.POST("/norm", norm)
	app.GET("/norm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "norms.tmpl", gin.H{"title": "VectorVictor: Norms"})
	})
	app.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return app
}

func TestIndexEndpoint(t *testing.T) {
	router := setupRouter()

	t.Run("POST returns JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("POST / returned status %d, want %d", w.Code, http.StatusOK)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse JSON response: %v", err)
		}

		if response["version"] != Version {
			t.Errorf("Version = %v, want %v", response["version"], Version)
		}
	})

	t.Run("GET returns HTML", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GET / returned status %d, want %d", w.Code, http.StatusOK)
		}

		contentType := w.Header().Get("Content-Type")
		if contentType != "text/html; charset=utf-8" {
			t.Errorf("Content-Type = %v, want text/html; charset=utf-8", contentType)
		}
	})
}

func TestSquareEndpoint(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedVal    []float64
	}{
		{"valid input", "?v=1,2,3", http.StatusOK, []float64{1, 4, 9}},
		{"single value", "?v=5", http.StatusOK, []float64{25}},
		{"empty input", "?v=", http.StatusOK, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/square"+tt.query, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("POST /square%s returned status %d, want %d", tt.query, w.Code, tt.expectedStatus)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
				return
			}

			if tt.expectedVal != nil {
				content := response["content"].(map[string]interface{})
				val := content["val"].([]interface{})
				if len(val) != len(tt.expectedVal) {
					t.Errorf("Result length = %d, want %d", len(val), len(tt.expectedVal))
				}
			}
		})
	}
}

func TestNormEndpoint(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedNorm   float64
	}{
		{"L2 norm 3-4-5", "?v=3,4&kind=l2", http.StatusOK, 5.0},
		{"L1 norm", "?v=1,2,3&kind=l1", http.StatusOK, 6.0},
		{"L-infinity", "?v=1,5,3&kind=linfinity", http.StatusOK, 5.0},
		{"default L2", "?v=3,4", http.StatusOK, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/norm"+tt.query, nil)
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("POST /norm%s returned status %d, want %d", tt.query, w.Code, tt.expectedStatus)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Failed to parse JSON response: %v", err)
				return
			}

			content := response["content"].(map[string]interface{})
			norm := content["norm"].(float64)
			if norm != tt.expectedNorm {
				t.Errorf("Norm = %v, want %v", norm, tt.expectedNorm)
			}
		})
	}
}

func TestNormEndpointErrors(t *testing.T) {
	router := setupRouter()

	t.Run("unknown norm kind", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/norm?v=1,2,3&kind=unknown", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("POST /norm with unknown kind returned status %d, want %d", w.Code, http.StatusBadRequest)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse JSON response: %v", err)
			return
		}

		if response["errors"] == nil {
			t.Error("Expected error in response, got nil")
		}
	})

	invalidP := []struct {
		name  string
		query string
	}{
		{"NaN p", "?v=1,2&kind=lp&p=NaN"},
		{"negative p", "?v=1,2&kind=lp&p=-1"},
	}

	for _, tt := range invalidP {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/norm"+tt.query, nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("POST /norm%s returned status %d, want %d", tt.query, w.Code, http.StatusBadRequest)
			}

			if w.Body.Len() == 0 {
				t.Fatalf("POST /norm%s returned an empty body", tt.query)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse JSON response: %v", err)
			}

			if response["errors"] == nil {
				t.Error("Expected error in response, got nil")
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GET /health returned status %d, want %d", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
		return
	}

	if response["status"] != "healthy" {
		t.Errorf("Health status = %v, want healthy", response["status"])
	}
}

func TestWrapResponse(t *testing.T) {
	t.Run("with no error", func(t *testing.T) {
		content := gin.H{"test": "value"}
		response := wrapResponse(content, nil)

		if response["errors"] != nil {
			t.Errorf("Expected nil errors, got %v", response["errors"])
		}
		if response["version"] != Version {
			t.Errorf("Version = %v, want %v", response["version"], Version)
		}
	})

	t.Run("with error", func(t *testing.T) {
		content := gin.H{"test": "value"}
		response := wrapResponse(content, http.ErrBodyNotAllowed)

		if response["errors"] == nil {
			t.Error("Expected error in response, got nil")
		}
	})
}
