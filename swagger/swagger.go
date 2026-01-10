// Package swagger provides automatic OpenAPI/Swagger documentation generation
// for Chi routers with support for automatic endpoint detection and schema registration.
//
// This package automatically detects routes from a Chi router and generates
// OpenAPI 3.0 specifications with support for:
//   - Automatic endpoint detection
//   - Security scheme detection (JWT, API Key)
//   - Tag-based organization
//   - Parameter and schema registration
//   - Request/Response body schemas
//
// Example usage:
//
//	cfg := swagger.DefaultConfig().
//		WithServiceInfo("My API", "1.0.0").
//		WithDescription("My awesome API").
//		AddTagRule("/users", "Users").
//		AddTagRule("/posts", "Posts")
//
//	detector := swagger.NewDetector(cfg)
//	detector.ScanRouter(router)
//
//	spec, _ := detector.ExportSwaggerSpec()
package swagger

import (
	"encoding/json"
	"net/http"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
	"github.com/AoC-Gamers/connect-libraries/swagger/detector"
	"github.com/AoC-Gamers/connect-libraries/swagger/openapi"
	"github.com/AoC-Gamers/connect-libraries/swagger/schema"
	"github.com/go-chi/chi/v5"
)

// Re-export main types for backward compatibility
type (
	// Config holds the configuration for Swagger detector
	Config = config.Config

	// TagRule defines a pattern-based tag assignment
	TagRule = config.TagRule

	// Detector automatically detects and documents endpoints from Chi router
	Detector = detector.Detector

	// DetectedRoute represents an automatically detected endpoint
	DetectedRoute = detector.Route

	// SchemaRegistry maintains schemas registered manually
	SchemaRegistry = schema.Registry

	// ParamSchema defines the structure of a parameter
	ParamSchema = schema.ParamSchema

	// ParamLocation specifies where the parameter is located
	ParamLocation = schema.ParamLocation

	// RequestBodySchema defines the schema of a request body
	RequestBodySchema = schema.RequestBodySchema

	// ResponseSchema defines the schema of a response
	ResponseSchema = schema.ResponseSchema
)

// Re-export constants
const (
	ParamInPath   = schema.ParamInPath
	ParamInQuery  = schema.ParamInQuery
	ParamInHeader = schema.ParamInHeader
)

// DefaultConfig returns a basic configuration
func DefaultConfig() *Config {
	return config.DefaultConfig()
}

// NewDetector creates a new detector with the given configuration
func NewDetector(cfg *Config) *Detector {
	return detector.New(cfg)
}

// ExportJSON exports detected routes as JSON
func ExportJSON(d *Detector) ([]byte, error) {
	routes := d.GetRoutes()
	return json.MarshalIndent(routes, "", "  ")
}

// ExportSwaggerSpec exports a complete OpenAPI/Swagger specification
func ExportSwaggerSpec(d *Detector) ([]byte, error) {
	cfg := d.GetConfig()
	routes := d.GetRoutes()
	schemas := d.GetSchemaRegistry()

	return openapi.GenerateSpec(cfg, routes, schemas)
}

// ServeSwaggerRoutes serves the detected routes as JSON (for /swagger/routes endpoint)
func ServeSwaggerRoutes(d *Detector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := ExportJSON(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

// Re-export helper functions from schema
var (
	// BuildSchemaFromStruct constructs a JSON schema from a struct using reflection
	BuildSchemaFromStruct = schema.BuildSchemaFromStruct

	// ExtractPathParamsFromRoute automatically detects path params from a route
	ExtractPathParamsFromRoute = schema.ExtractPathParamsFromRoute
)

// Compatibility aliases for old method names

// Route is an alias for DetectedRoute (backward compatibility)
type Route = DetectedRoute

// NewSchemaRegistry creates a new schema registry
func NewSchemaRegistry() *SchemaRegistry {
	return schema.NewRegistry()
}

// ScanRouter is a convenience function that creates a detector and scans a router
func ScanRouter(r chi.Router, cfg *Config) (*Detector, error) {
	d := NewDetector(cfg)
	err := d.ScanRouter(r)
	return d, err
}
