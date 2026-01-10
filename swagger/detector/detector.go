package detector

import (
	"net/http"
	"reflect"
	"runtime"
	"strings"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
	"github.com/AoC-Gamers/connect-libraries/swagger/schema"
	"github.com/go-chi/chi/v5"
)

// Detector automatically detects and documents endpoints from Chi router
type Detector struct {
	config  *config.Config
	routes  []Route
	schemas *schema.Registry
}

// Route represents an automatically detected endpoint
type Route struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	Security    []string `json:"security"`
	Tags        []string `json:"tags"`
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
}

// New creates a new detector with the given configuration
func New(cfg *config.Config) *Detector {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	return &Detector{
		config:  cfg,
		routes:  make([]Route, 0),
		schemas: schema.NewRegistry(),
	}
}

// GetSchemaRegistry returns the schema registry for manual registration
func (d *Detector) GetSchemaRegistry() *schema.Registry {
	return d.schemas
}

// GetConfig returns the detector's configuration
func (d *Detector) GetConfig() *config.Config {
	return d.config
}

// ScanRouter walks the Chi router and automatically detects endpoints
func (d *Detector) ScanRouter(r chi.Router) error {
	d.routes = make([]Route, 0)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		// Clean up route path
		route = strings.Replace(route, "/*/", "/", -1)

		// Skip internal routes
		if d.shouldSkipRoute(route) {
			return nil
		}

		detected := Route{
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

// GetRoutes returns all detected routes
func (d *Detector) GetRoutes() []Route {
	return d.routes
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
