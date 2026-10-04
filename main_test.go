package main

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
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
		{"spaced input", "?v=1, 2, 3", http.StatusOK, []float64{1, 4, 9}},
		{"malformed token", "?v=1oops2,3", http.StatusBadRequest, nil},
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
		{"signed and scientific L1", "?v=-3,4e0&kind=l1", http.StatusOK, 7.0},
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
		{"malformed v token", "?v=1oops2,3&kind=l2"},
		{"invalid weights token", "?v=3,4&kind=weighted&weights=1x,2"},
		{"invalid variances token", "?v=3,4&kind=mahalanobis&variances=1x,2"},
		{"zero variance", "?v=5,5&kind=mahalanobis&variances=0,1"},
		{"negative weight beyond vector", "?v=1&kind=weighted&weights=1,-5"},
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

	// Negative weights/variances used to produce NaN, which is not
	// JSON-serializable, causing the encoder to emit an empty-body 200.
	// Assert a proper 400 with a populated errors field and non-empty body.
	negativeCases := []struct {
		name  string
		query string
	}{
		{"negative weights", "?v=3,4&kind=weighted&weights=-1,-1"},
		{"negative variances", "?v=3,4&kind=mahalanobis&variances=-1,-1"},
	}

	for _, tc := range negativeCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/norm"+tc.query, nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("POST /norm%s returned status %d, want %d", tc.query, w.Code, http.StatusBadRequest)
			}

			if w.Body.Len() == 0 {
				t.Fatalf("POST /norm%s returned an empty body, want a populated error response", tc.query)
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

func TestCalculatorPagesAvoidHTMLSinks(t *testing.T) {
	router := setupRouter()

	for _, path := range []string{"/norm", "/square"} {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", path, nil)
		router.ServeHTTP(w, req)

		body := w.Body.String()
		for _, sink := range []string{"innerHTML", "insertAdjacentHTML", "outerHTML", "document.write"} {
			if strings.Contains(body, sink) {
				t.Errorf("GET %s page contains HTML-parsing sink %q", path, sink)
			}
		}
	}
}

func TestNonFiniteInputsReturnJSONError(t *testing.T) {
	router := setupRouter()
	queries := []string{
		"/norm?v=NaN&kind=l2",
		"/norm?v=1e308,1e308&kind=l1",
		"/norm?v=1&kind=weighted&weights=NaN",
		"/norm?v=1&kind=mahalanobis&variances=NaN",
		"/norm?v=1e200&kind=weighted&weights=1e300",
		"/square?v=1e308",
		"/square?v=Inf",
	}
	for _, q := range queries {
		t.Run(q, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", q, nil)
			router.ServeHTTP(w, req)

			if w.Code < 400 {
				t.Errorf("status = %d, want error status", w.Code)
			}
			if w.Body.Len() == 0 {
				t.Fatal("empty body")
			}
			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("invalid JSON: %v", err)
			}
			if s, _ := response["errors"].(string); s == "" {
				t.Errorf("errors = %v, want non-empty string", response["errors"])
			}
		})
	}
}

func TestNormL2LargeFiniteVector(t *testing.T) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/norm?v=1e308,1e308&kind=l2", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200: %s", w.Code, w.Body.String())
	}
	var response struct {
		Content struct {
			Norm float64 `json:"norm"`
		} `json:"content"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	want := math.Hypot(1e308, 1e308)
	if math.Abs(response.Content.Norm-want) > want*1e-12 {
		t.Errorf("norm = %v, want %v", response.Content.Norm, want)
	}
}

func TestNormLpExtremeMagnitudes(t *testing.T) {
	cases := []struct {
		query string
		want  float64
	}{
		{"?v=1e200&kind=lp&p=3", 1e200},
		{"?v=1e-200&kind=lp&p=3", 1e-200},
	}
	for _, tc := range cases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/norm"+tc.query, nil)
		setupRouter().ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("%s: status = %d: %s", tc.query, w.Code, w.Body.String())
		}
		var response struct {
			Content struct {
				Norm float64 `json:"norm"`
			} `json:"content"`
			Errors any `json:"errors"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("%s: invalid JSON: %v", tc.query, err)
		}
		if response.Errors != nil {
			t.Errorf("%s: errors = %v, want null", tc.query, response.Errors)
		}
		if math.Abs(response.Content.Norm-tc.want) > tc.want*1e-12 {
			t.Errorf("%s: norm = %v, want %v", tc.query, response.Content.Norm, tc.want)
		}
	}
}
