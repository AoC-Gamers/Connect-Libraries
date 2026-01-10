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
//	detector := swagger.New(cfg)
//	detector.ScanRouter(router)
//
//	spec, _ := detector.ExportSpec()
package swagger

import (
	"encoding/json"
	"net/http"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
	"github.com/AoC-Gamers/connect-libraries/swagger/detector"
	"github.com/AoC-Gamers/connect-libraries/swagger/openapi"
)

// DefaultConfig returns a basic configuration
func DefaultConfig() *config.Config {
	return config.DefaultConfig()
}

// New creates a new detector with the given configuration
func New(cfg *config.Config) *detector.Detector {
	return detector.New(cfg)
}

// ExportJSON exports detected routes as JSON
func ExportJSON(d *detector.Detector) ([]byte, error) {
	routes := d.GetRoutes()
	return json.MarshalIndent(routes, "", "  ")
}

// ExportSpec exports a complete OpenAPI/Swagger specification
func ExportSpec(d *detector.Detector) ([]byte, error) {
	return openapi.GenerateSpec(d.GetConfig(), d.GetRoutes(), d.GetSchemaRegistry())
}

// ServeSwaggerSpec returns an http.Handler that serves the OpenAPI specification
func ServeSwaggerSpec(d *detector.Detector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data, err := ExportSpec(d)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}
