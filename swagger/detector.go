package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Detector automatically detects and documents endpoints from Chi router
type Detector struct {
	config *Config
	routes []DetectedRoute
}

// DetectedRoute represents an automatically detected endpoint
type DetectedRoute struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Security    []string `json:"security"`
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
}

// NewDetector creates a new detector with the given configuration
func NewDetector(config *Config) *Detector {
	if config == nil {
		config = DefaultConfig()
	}
	return &Detector{
		config: config,
		routes: make([]DetectedRoute, 0),
	}
}

// ScanRouter walks the Chi router and automatically detects endpoints
func (d *Detector) ScanRouter(r chi.Router) error {
	d.routes = make([]DetectedRoute, 0)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Clean up route path
		route = strings.Replace(route, "/*/", "/", -1)

		// Skip internal routes
		if d.shouldSkipRoute(route) {
			return nil
		}

		detected := DetectedRoute{
			Method: method,
			Path:   route,
		}

		// Detect security from middlewares
		detected.Security = d.detectSecurityFromMiddlewares(middlewares, route)

		// Infer tags from path using configured rules
		detected.Tags = d.inferTagsFromPath(route)

		// Generate summary and description
		detected.Summary = d.generateSummary(method, route)
		detected.Description = d.generateDescription(method, route, detected.Security)

		d.routes = append(d.routes, detected)
		return nil
	}

	return chi.Walk(r, walkFunc)
}

// shouldSkipRoute determines if a route should be skipped
func (d *Detector) shouldSkipRoute(route string) bool {
	for _, skip := range d.config.SkipPaths {
		if strings.Contains(route, skip) {
			return true
		}
	}
	return false
}

// detectSecurityFromMiddlewares analyzes middleware chain to detect JWT vs API Key
func (d *Detector) detectSecurityFromMiddlewares(middlewares []func(http.Handler) http.Handler, path string) []string {
	security := []string{}

	for _, mw := range middlewares {
		if mw == nil {
			continue
		}

		// Get function name to identify middleware type
		funcName := getFunctionName(mw)

		// Check against configured security patterns
		for pattern, secType := range d.config.SecurityPatterns {
			if strings.Contains(funcName, pattern) {
				if !contains(security, secType) {
					security = append(security, secType)
				}
			}
		}
	}

	// If no specific security detected but it's not a public endpoint, use default
	if len(security) == 0 && !d.isPublicEndpoint(path) && d.config.DefaultSecurity != "" {
		security = append(security, d.config.DefaultSecurity)
	}

	return security
}

// isPublicEndpoint determines if an endpoint should be public
func (d *Detector) isPublicEndpoint(path string) bool {
	for _, publicPath := range d.config.PublicPaths {
		if strings.HasPrefix(path, publicPath) || strings.Contains(path, publicPath) {
			return true
		}
	}
	return false
}

// inferTagsFromPath categorizes endpoints based on configured tag rules
func (d *Detector) inferTagsFromPath(path string) []string {
	// Check tag rules in order (first match wins)
	for _, rule := range d.config.TagRules {
		if strings.Contains(path, rule.PathPattern) {
			return rule.Tags
		}
	}

	// Default fallback
	return []string{d.config.DefaultTag}
}

// generateSummary creates a summary for the endpoint
func (d *Detector) generateSummary(method, path string) string {
	// Clean path for summary
	cleanPath := strings.Replace(path, "{", "", -1)
	cleanPath = strings.Replace(cleanPath, "}", "", -1)

	resource := extractResource(path)

	switch method {
	case "GET":
		if strings.Contains(path, "/{") {
			return fmt.Sprintf("Get %s by ID", resource)
		}
		return fmt.Sprintf("List %s", resource)
	case "POST":
		return fmt.Sprintf("Create %s", resource)
	case "PUT":
		return fmt.Sprintf("Update %s", resource)
	case "PATCH":
		return fmt.Sprintf("Patch %s", resource)
	case "DELETE":
		return fmt.Sprintf("Delete %s", resource)
	default:
		return fmt.Sprintf("%s %s", method, cleanPath)
	}
}

// generateDescription creates a description including security info
func (d *Detector) generateDescription(method, path string, security []string) string {
	desc := fmt.Sprintf("Endpoint for %s %s", method, path)

	if len(security) > 0 {
		securityTypes := make([]string, 0)
		for _, sec := range security {
			switch sec {
			case "BearerAuth":
				securityTypes = append(securityTypes, "JWT authentication")
			case "ApiKeyAuth":
				securityTypes = append(securityTypes, "API key authentication")
			}
		}

		if len(securityTypes) > 0 {
			desc += ". Requires: " + strings.Join(securityTypes, " and ")
		}
	} else {
		desc += ". Public endpoint"
	}

	return desc
}

// GetRoutes returns all detected routes
func (d *Detector) GetRoutes() []DetectedRoute {
	return d.routes
}

// ExportJSON exports detected routes as JSON
func (d *Detector) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(d.routes, "", "  ")
}

// ExportSwaggerSpec exports a complete OpenAPI/Swagger specification
func (d *Detector) ExportSwaggerSpec() ([]byte, error) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       d.config.ServiceName,
			"version":     d.config.Version,
			"description": d.config.Description,
		},
		"paths": d.generatePaths(),
		"components": map[string]interface{}{
			"securitySchemes": map[string]interface{}{
				"BearerAuth": map[string]interface{}{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
					"description":  "JWT authentication token",
				},
				"ApiKeyAuth": map[string]interface{}{
					"type":        "apiKey",
					"in":          "header",
					"name":        "X-API-Key",
					"description": "API Key for service-to-service authentication",
				},
			},
		},
	}

	// Add contact if configured
	if d.config.ContactName != "" || d.config.ContactEmail != "" {
		contact := map[string]interface{}{}
		if d.config.ContactName != "" {
			contact["name"] = d.config.ContactName
		}
		if d.config.ContactURL != "" {
			contact["url"] = d.config.ContactURL
		}
		if d.config.ContactEmail != "" {
			contact["email"] = d.config.ContactEmail
		}
		spec["info"].(map[string]interface{})["contact"] = contact
	}

	// Add license if configured
	if d.config.LicenseName != "" {
		license := map[string]interface{}{
			"name": d.config.LicenseName,
		}
		if d.config.LicenseURL != "" {
			license["url"] = d.config.LicenseURL
		}
		spec["info"].(map[string]interface{})["license"] = license
	}

	return json.MarshalIndent(spec, "", "  ")
}

// generatePaths converts detected routes to OpenAPI paths format
func (d *Detector) generatePaths() map[string]interface{} {
	paths := make(map[string]interface{})

	for _, route := range d.routes {
		if _, exists := paths[route.Path]; !exists {
			paths[route.Path] = make(map[string]interface{})
		}

		operation := map[string]interface{}{
			"summary":     route.Summary,
			"description": route.Description,
			"tags":        route.Tags,
			"responses": map[string]interface{}{
				"200": map[string]interface{}{
					"description": "Successful response",
				},
				"400": map[string]interface{}{
					"description": "Bad request",
				},
				"401": map[string]interface{}{
					"description": "Unauthorized",
				},
				"500": map[string]interface{}{
					"description": "Internal server error",
				},
			},
		}

		// Add security if required
		if len(route.Security) > 0 {
			security := make([]map[string][]string, 0)
			for _, sec := range route.Security {
				security = append(security, map[string][]string{sec: {}})
			}
			operation["security"] = security
		}

		// Add to paths
		method := strings.ToLower(route.Method)
		paths[route.Path].(map[string]interface{})[method] = operation
	}

	return paths
}

// ServeHTTP serves the detected routes as JSON (for /swagger/routes endpoint)
func (d *Detector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := d.ExportJSON()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// Helper functions

func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractResource(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 {
		// Get the last meaningful part (before parameters)
		for i := len(parts) - 1; i >= 0; i-- {
			if !strings.Contains(parts[i], "{") && parts[i] != "" {
				return strings.Title(parts[i])
			}
		}
	}
	return "Resource"
}
