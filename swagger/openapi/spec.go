package openapi

import (
	"encoding/json"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
	"github.com/AoC-Gamers/connect-libraries/swagger/detector"
	"github.com/AoC-Gamers/connect-libraries/swagger/schema"
)

// GenerateSpec generates a complete OpenAPI/Swagger specification
func GenerateSpec(cfg *config.Config, routes []detector.Route, schemas *schema.Registry) ([]byte, error) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":       cfg.ServiceName,
			"version":     cfg.Version,
			"description": cfg.Description,
		},
		"paths":      generatePaths(routes, schemas),
		"components": generateComponents(),
	}

	// Add contact if configured
	if cfg.ContactName != "" || cfg.ContactEmail != "" {
		contact := map[string]interface{}{}
		if cfg.ContactName != "" {
			contact["name"] = cfg.ContactName
		}
		if cfg.ContactURL != "" {
			contact["url"] = cfg.ContactURL
		}
		if cfg.ContactEmail != "" {
			contact["email"] = cfg.ContactEmail
		}
		spec["info"].(map[string]interface{})["contact"] = contact
	}

	// Add license if configured
	if cfg.LicenseName != "" {
		license := map[string]interface{}{
			"name": cfg.LicenseName,
		}
		if cfg.LicenseURL != "" {
			license["url"] = cfg.LicenseURL
		}
		spec["info"].(map[string]interface{})["license"] = license
	}

	return json.MarshalIndent(spec, "", "  ")
}

// generateComponents generates the OpenAPI components section
func generateComponents() map[string]interface{} {
	return map[string]interface{}{
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
				"name":        "X-Internal-API-Key",
				"description": "API Key for service-to-service authentication",
			},
		},
	}
}
