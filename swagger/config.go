package swagger

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

// WithServiceInfo sets service name and version
func (c *Config) WithServiceInfo(name, version string) *Config {
	c.ServiceName = name
	c.Version = version
	return c
}

// WithDescription sets the API description
func (c *Config) WithDescription(description string) *Config {
	c.Description = description
	return c
}

// WithContact sets the contact information
func (c *Config) WithContact(name, url, email string) *Config {
	c.ContactName = name
	c.ContactURL = url
	c.ContactEmail = email
	return c
}

// WithLicense sets the license information
func (c *Config) WithLicense(name, url string) *Config {
	c.LicenseName = name
	c.LicenseURL = url
	return c
}

// WithTagRules sets custom tag rules (replaces existing)
func (c *Config) WithTagRules(rules []TagRule) *Config {
	c.TagRules = rules
	return c
}

// AddTagRule adds a single tag rule
func (c *Config) AddTagRule(pathPattern string, tags ...string) *Config {
	c.TagRules = append(c.TagRules, TagRule{
		PathPattern: pathPattern,
		Tags:        tags,
	})
	return c
}

// AddPublicPath adds a public path
func (c *Config) AddPublicPath(path string) *Config {
	c.PublicPaths = append(c.PublicPaths, path)
	return c
}

// AddSkipPath adds a path to skip
func (c *Config) AddSkipPath(path string) *Config {
	c.SkipPaths = append(c.SkipPaths, path)
	return c
}

// WithDefaultSecurity sets the default security scheme
func (c *Config) WithDefaultSecurity(security string) *Config {
	c.DefaultSecurity = security
	return c
}

// WithDefaultTag sets the default tag
func (c *Config) WithDefaultTag(tag string) *Config {
	c.DefaultTag = tag
	return c
}
