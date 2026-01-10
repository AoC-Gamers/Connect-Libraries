package config

// Config holds the configuration for Swagger detector
type Config struct {
	// Service information
	ServiceName string
	Version     string
	Description string

	// Contact information
	ContactName  string
	ContactURL   string
	ContactEmail string

	// License information
	LicenseName string
	LicenseURL  string

	// Tag rules define how to categorize endpoints
	// Rules are evaluated in order - first match wins
	TagRules []TagRule

	// Security patterns map middleware function names to security types
	SecurityPatterns map[string]string

	// Public paths that don't require authentication
	PublicPaths []string

	// Skip paths that should not appear in swagger
	SkipPaths []string

	// Default tag for endpoints that don't match any rule
	DefaultTag string

	// Default security for protected endpoints
	DefaultSecurity string
}

// TagRule defines a pattern-based tag assignment
type TagRule struct {
	PathPattern string   // URL path pattern to match (uses strings.Contains)
	Tags        []string // Tags to assign if pattern matches
}

// DefaultConfig returns a basic configuration
func DefaultConfig() *Config {
	return &Config{
		ServiceName: "API Service",
		Version:     "1.0.0",
		TagRules:    []TagRule{},
		SecurityPatterns: map[string]string{
			"RequireAuth":             "BearerAuth",
			"JWTAuth":                 "BearerAuth",
			"RequireAPIKey":           "ApiKeyAuth",
			"APIKey":                  "ApiKeyAuth",
			"RequireInternalServices": "ApiKeyAuth",
			"RequireWebPermission":    "BearerAuth",
			"Permission":              "BearerAuth",
			"RequireVIP":              "BearerAuth",
			"RequireAdmin":            "BearerAuth",
		},
		PublicPaths: []string{
			"/health",
			"/status",
		},
		SkipPaths: []string{
			"/swagger",
			"/debug",
		},
		DefaultTag:      "General",
		DefaultSecurity: "BearerAuth",
	}
}
