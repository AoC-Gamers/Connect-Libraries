package openapi

import (
	"encoding/json"
	"testing"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
	"github.com/AoC-Gamers/connect-libraries/swagger/detector"
	"github.com/AoC-Gamers/connect-libraries/swagger/schema"
)

func TestGenerateSpecBasicStructure(t *testing.T) {
	cfg := config.DefaultConfig().
		WithServiceInfo("Connect API", "1.2.3").
		WithDescription("Service docs").
		WithContact("Team", "https://example.com", "team@example.com").
		WithLicense("MIT", "https://opensource.org/licenses/MIT")

	routes := []detector.Route{
		{
			Method:      "GET",
			Path:        "/users/{id}",
			Security:    []string{"BearerAuth"},
			Tags:        []string{"Users"},
			Summary:     "Get Users by ID",
			Description: "Endpoint for GET /users/{id}",
		},
	}

	data, err := GenerateSpec(cfg, routes, schema.NewRegistry())
	if err != nil {
		t.Fatalf("expected spec generation without error, got %v", err)
	}

	var spec map[string]interface{}
	if err := json.Unmarshal(data, &spec); err != nil {
		t.Fatalf("could not parse generated spec: %v", err)
	}

	if spec["openapi"] != "3.0.0" {
		t.Fatalf("expected openapi 3.0.0, got %v", spec["openapi"])
	}
	info := spec["info"].(map[string]interface{})
	if info["title"] != "Connect API" || info["version"] != "1.2.3" {
		t.Fatal("expected info title/version to be set")
	}
	if info["contact"] == nil || info["license"] == nil {
		t.Fatal("expected contact and license in info")
	}

	paths := spec["paths"].(map[string]interface{})
	usersPath := paths["/users/{id}"].(map[string]interface{})
	getOp := usersPath["get"].(map[string]interface{})
	if getOp["summary"] != "Get Users by ID" {
		t.Fatalf("unexpected operation summary %v", getOp["summary"])
	}
	if getOp["parameters"] == nil {
		t.Fatal("expected auto-detected path parameter")
	}
	if getOp["responses"] == nil || getOp["security"] == nil {
		t.Fatal("expected responses and security in operation")
	}

	components := spec["components"].(map[string]interface{})
	securitySchemes := components["securitySchemes"].(map[string]interface{})
	if securitySchemes["BearerAuth"] == nil || securitySchemes["ApiKeyAuth"] == nil {
		t.Fatal("expected default security schemes")
	}
}
