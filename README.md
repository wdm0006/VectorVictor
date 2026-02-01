# VectorVictor

[![CI](https://github.com/wdm0006/VectorVictor/actions/workflows/ci.yml/badge.svg)](https://github.com/wdm0006/VectorVictor/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/wdm0006/VectorVictor)](https://goreportcard.com/report/github.com/wdm0006/VectorVictor)
[![Go Version](https://img.shields.io/github/go-mod/go-version/wdm0006/VectorVictor)](https://go.dev/)

> **Note:** This is a toy project created primarily as a learning exercise for Go web development. It's not intended for production use, but rather as a playground for experimenting with Go's HTTP handling, templates, testing, and project structure. Feel free to use it as a reference or starting point for your own learning!

A lightweight Go web service for vector mathematics operations. VectorVictor provides both a clean web interface and a REST API for calculating vector norms and performing element-wise operations.

## Why This Project?

This project was built to learn and experiment with:

- **Go web frameworks** - Using Gin for HTTP routing and middleware
- **Go templates** - Server-side HTML rendering with `html/template`
- **Go embed** - Bundling static assets into the binary (Go 1.16+)
- **Testing in Go** - Table-driven tests, HTTP handler testing with `httptest`
- **Project structure** - Organizing a Go web application
- **CI/CD** - GitHub Actions for automated testing and linting

## Background

### What are Vector Norms?

A **norm** is a function that assigns a non-negative length or size to vectors. Norms are fundamental in linear algebra, machine learning, and data science for measuring distances, regularization, and optimization.

VectorVictor supports **8 different norms**:

| Norm | Formula | Description | Use Cases |
|------|---------|-------------|-----------|
| **L0** (Zero) | count(xᵢ ≠ 0) | Number of non-zero elements | Sparsity, compressed sensing |
| **L1** (Manhattan) | Σ\|xᵢ\| | Sum of absolute values | Lasso regression, sparse solutions |
| **L2** (Euclidean) | √(Σxᵢ²) | Straight-line distance | Most common, Ridge regression |
| **L∞** (Maximum) | max(\|xᵢ\|) | Largest absolute component | Worst-case analysis |
| **L0.5** (Sub-unitary) | (Σ√\|xᵢ\|)² | Promotes extreme sparsity | Feature selection |
| **Lp** (General) | (Σ\|xᵢ\|^p)^(1/p) | Generalized p-norm | Flexible regularization |
| **Weighted L2** | √(Σwᵢ·xᵢ²) | L2 with per-element weights | Prioritizing dimensions |
| **Mahalanobis** | √(Σxᵢ²/σᵢ²) | Variance-normalized distance | Outlier detection |

**Example:** For vector `[3, 4]`:
- L0 norm = 2 (both non-zero)
- L1 norm = |3| + |4| = **7**
- L2 norm = √(3² + 4²) = √25 = **5**
- L∞ norm = max(|3|, |4|) = **4**

## Features

- **8 Vector Norms** - L0, L1, L2, L∞, L0.5, Lp, Weighted L2, and Mahalanobis
- **Element-wise Operations** - Square each element of a vector
- **Web Interface** - Clean Bootstrap 5 UI with real-time calculations
- **REST API** - JSON endpoints for programmatic access
- **Well-tested** - Comprehensive test suite with 65%+ coverage

## Quick Start

### Prerequisites

- Go 1.21 or later

### Installation

```bash
# Clone the repository
git clone https://github.com/wdm0006/VectorVictor.git
cd VectorVictor

# Build the binary
go build -o vectorvictor

# Run the server
./vectorvictor
```

The server starts on `http://localhost:8080`

### Using Go Run

```bash
go run .
```

## API Reference

### Endpoints Overview

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Home page (HTML) |
| `/` | POST | Status check (JSON) |
| `/norm` | GET | Norm calculator UI |
| `/norm` | POST | Calculate vector norm |
| `/square` | GET | Square calculator UI |
| `/square` | POST | Square vector elements |
| `/health` | GET | Health check |

### Calculate Vector Norm

**POST** `/norm`

Calculate various norms of a vector.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `v` | string | Yes | Comma-separated vector values |
| `kind` | string | No | Norm type (default: `l2`). Options: `l0`, `l1`, `l2`, `linfinity`, `lhalf`, `lp`, `weighted`, `mahalanobis` |
| `p` | float | No | p value for `lp` norm (default: 2) |
| `weights` | string | No | Comma-separated weights for `weighted` norm |
| `variances` | string | No | Comma-separated variances for `mahalanobis` norm |

**Example - L2 Norm (Euclidean):**

```bash
curl -X POST "http://localhost:8080/norm?v=3,4&kind=l2"
```

```json
{
  "content": {
    "kind": "l2",
    "norm": 5,
    "vector": [3, 4]
  },
  "errors": null,
  "info": {
    "now": "2025-02-01T12:00:00.000000Z"
  },
  "version": "1.0.0"
}
```

**Example - L1 Norm (Manhattan):**

```bash
curl -X POST "http://localhost:8080/norm?v=1,2,3,4&kind=l1"
```

```json
{
  "content": {
    "kind": "l1",
    "norm": 10,
    "vector": [1, 2, 3, 4]
  },
  "errors": null,
  "info": {
    "now": "2025-02-01T12:00:00.000000Z"
  },
  "version": "1.0.0"
}
```

**Example - L-infinity Norm (Maximum):**

```bash
curl -X POST "http://localhost:8080/norm?v=-5,3,4&kind=linfinity"
```

```json
{
  "content": {
    "kind": "linfinity",
    "norm": 4,
    "vector": [-5, 3, 4]
  },
  "errors": null,
  "info": {
    "now": "2025-02-01T12:00:00.000000Z"
  },
  "version": "1.0.0"
}
```

**Example - L0 Norm (Sparsity):**

```bash
curl -X POST "http://localhost:8080/norm?v=1,0,3,0,5&kind=l0"
# Returns: norm = 3 (three non-zero elements)
```

**Example - Lp Norm with custom p:**

```bash
curl -X POST "http://localhost:8080/norm?v=1,2,3&kind=lp&p=3"
# Returns: (1³ + 2³ + 3³)^(1/3) = 36^(1/3) ≈ 3.30
```

**Example - Weighted L2 Norm:**

```bash
curl -X POST "http://localhost:8080/norm?v=3,4&kind=weighted&weights=4,1"
# Returns: √(4×9 + 1×16) = √52 ≈ 7.21
```

**Example - Mahalanobis Distance:**

```bash
curl -X POST "http://localhost:8080/norm?v=6,8&kind=mahalanobis&variances=4,4"
# Returns: √(36/4 + 64/4) = √25 = 5
```

### Square Vector Elements

**POST** `/square`

Square each element of a vector.

**Parameters:**
| Name | Type | Required | Description |
|------|------|----------|-------------|
| `v` | string | Yes | Comma-separated vector values |

**Example:**

```bash
curl -X POST "http://localhost:8080/square?v=1,2,3,4,5"
```

```json
{
  "content": {
    "val": [1, 4, 9, 16, 25]
  },
  "errors": null,
  "info": {
    "now": "2025-02-01T12:00:00.000000Z"
  },
  "version": "1.0.0"
}
```

### Health Check

**GET** `/health`

Returns the service health status.

```bash
curl http://localhost:8080/health
```

```json
{
  "status": "healthy"
}
```

### Error Handling

When an error occurs, the API returns an appropriate HTTP status code with error details:

```bash
curl -X POST "http://localhost:8080/norm?v=1,2,3&kind=invalid"
```

```json
{
  "content": {
    "kind": "invalid",
    "norm": 0,
    "vector": [1, 2, 3]
  },
  "errors": "unknown norm kind: invalid",
  "info": {
    "now": "2025-02-01T12:00:00.000000Z"
  },
  "version": "1.0.0"
}
```

## Development

### Run Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run with coverage report
go test ./... -cover

# Run with race detection
go test ./... -race
```

Current test coverage: **67.9%** (54 tests)

### Lint

```bash
# Install golangci-lint (if needed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Build

```bash
# Standard build
go build -o vectorvictor

# Build with version info
go build -ldflags "-X main.Version=1.0.0" -o vectorvictor

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o vectorvictor-linux
```

## Project Structure

```
VectorVictor/
├── main.go              # HTTP server, routing, request handlers
├── main_test.go         # HTTP endpoint tests
├── norms.go             # Norm calculations (L1, L2, L-infinity, LN)
├── norms_test.go        # Norm function tests
├── elementwise.go       # Element-wise operations (Square, arrayExp)
├── elementwise_test.go  # Element-wise operation tests
├── delimited.go         # String parsing (CSV, TSV, PSV to float arrays)
├── delimited_test.go    # Parser tests
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── templates/           # HTML templates (embedded via go:embed)
│   ├── index.tmpl       # Home page
│   ├── norms.tmpl       # Norm calculator UI
│   └── square.tmpl      # Square calculator UI
├── .github/
│   └── workflows/
│       └── ci.yml       # GitHub Actions CI pipeline
└── README.md
```

## Technical Details

- **Framework**: [Gin](https://github.com/gin-gonic/gin) v1.10 - High-performance HTTP web framework
- **Templates**: Go 1.16+ `embed` directive for bundling HTML templates
- **Frontend**: Bootstrap 5.3, vanilla JavaScript (no jQuery)
- **Shutdown**: Graceful shutdown with SIGINT/SIGTERM signal handling
- **CI/CD**: GitHub Actions with automated testing, linting, and builds

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Run linter (`golangci-lint run`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Future Ideas

Here are additional operations that could be added:

### Additional Operations

| Operation | Description | Use Case |
|-----------|-------------|----------|
| **Normalize** | Scale vector to unit length | Preprocessing for ML |
| **Dot product** | Inner product of two vectors | Similarity measurement |
| **Cosine similarity** | Angle between vectors | Text similarity, recommendations |
| **Cross product** | 3D vector perpendicular to inputs | Physics, graphics |
| **Projection** | Project one vector onto another | Dimensionality reduction |

### Matrix Operations (Future Scope)

| Operation | Description |
|-----------|-------------|
| **Frobenius norm** | L2 norm extended to matrices |
| **Spectral norm** | Largest singular value |
| **Nuclear norm** | Sum of singular values |
| **Matrix multiplication** | Standard matrix product |
| **Determinant** | Scalar from square matrix |
| **Eigenvalues** | Characteristic values |

Contributions welcome! See the Contributing section above.

## Acknowledgments

- Inspired by the classic ["What's our vector, Victor?"](https://www.youtube.com/watch?v=fVq4_HhBK8Y) line from *Airplane!* (1980)

## References

- [Norm (mathematics) - Wikipedia](https://en.wikipedia.org/wiki/Norm_(mathematics))
- [Vector Norms - GeeksforGeeks](https://www.geeksforgeeks.org/maths/vector-norms/)
- [Understanding Vector Norms - Medium](https://medium.com/@manoj.pillai/understanding-vector-norms-af40e2a7f7ea)
- [Matrix Norms - Wikipedia](https://en.wikipedia.org/wiki/Matrix_norm)
