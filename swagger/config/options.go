package config

// Builder methods for Config

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
