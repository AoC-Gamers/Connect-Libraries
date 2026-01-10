package detector

import (
	"fmt"
	"strings"
)

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
