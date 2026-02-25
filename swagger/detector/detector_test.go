package detector

import (
	"net/http"
	"strings"
	"testing"

	"github.com/AoC-Gamers/connect-libraries/swagger/config"
)

const (
	usersPath   = "/users"
	userIDPath  = "/users/{id}"
	privatePath = "/private"
)

func requireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func noopMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestDetectorBasicsAndHelpers(t *testing.T) {
	d := New(nil)
	if d == nil || d.GetConfig() == nil || d.GetSchemaRegistry() == nil {
		t.Fatal("expected detector initialized with default config and schema registry")
	}

	if !d.shouldSkipRoute("/v1/swagger/index.html") {
		t.Fatal("expected swagger route to be skipped")
	}
	if d.shouldSkipRoute("/v1/users") {
		t.Fatal("did not expect /v1/users to be skipped")
	}
}

func TestGenerateSummaryAndDescription(t *testing.T) {
	d := New(config.DefaultConfig())

	if summary := d.generateSummary("GET", usersPath); summary != "List Users" {
		t.Fatalf("unexpected summary: %s", summary)
	}
	if summary := d.generateSummary("GET", userIDPath); summary != "Get Users by ID" {
		t.Fatalf("unexpected summary for path param: %s", summary)
	}
	if summary := d.generateSummary("DELETE", userIDPath); summary != "Delete Users" {
		t.Fatalf("unexpected delete summary: %s", summary)
	}
	if summary := d.generateSummary("OPTIONS", usersPath); summary != "OPTIONS /users" {
		t.Fatalf("unexpected fallback summary: %s", summary)
	}

	publicDesc := d.generateDescription("GET", "/health", nil)
	if publicDesc == "" || publicDesc[len(publicDesc)-15:] != "Public endpoint" {
		t.Fatalf("unexpected public description: %s", publicDesc)
	}
	securedDesc := d.generateDescription("GET", usersPath, []string{"BearerAuth", "ApiKeyAuth"})
	if securedDesc == "" || !strings.Contains(securedDesc, "JWT authentication") || !strings.Contains(securedDesc, "API key authentication") {
		t.Fatalf("unexpected secured description: %s", securedDesc)
	}
}

func TestSecurityAndTagInference(t *testing.T) {
	cfg := config.DefaultConfig().
		WithTagRules([]config.TagRule{{PathPattern: usersPath, Tags: []string{"Users"}}}).
		WithDefaultTag("General")
	d := New(cfg)

	security := d.detectSecurityFromMiddlewares([]func(http.Handler) http.Handler{requireAuthMiddleware}, usersPath)
	if len(security) != 1 || security[0] != "BearerAuth" {
		t.Fatalf("expected BearerAuth from middleware, got %+v", security)
	}

	defaultSecurity := d.detectSecurityFromMiddlewares([]func(http.Handler) http.Handler{noopMiddleware}, "/private/resource")
	if len(defaultSecurity) != 1 || defaultSecurity[0] != cfg.DefaultSecurity {
		t.Fatalf("expected default security %s, got %+v", cfg.DefaultSecurity, defaultSecurity)
	}

	publicSecurity := d.detectSecurityFromMiddlewares(nil, "/health/live")
	if len(publicSecurity) != 0 {
		t.Fatalf("expected no security for public path, got %+v", publicSecurity)
	}

	tags := d.inferTagsFromPath(userIDPath)
	if len(tags) != 1 || tags[0] != "Users" {
		t.Fatalf("expected Users tag, got %+v", tags)
	}
	defaultTags := d.inferTagsFromPath("/unknown")
	if len(defaultTags) != 1 || defaultTags[0] != "General" {
		t.Fatalf("expected default tag General, got %+v", defaultTags)
	}

	if !d.isPublicEndpoint("/health") {
		t.Fatal("expected /health to be public")
	}
	if d.isPublicEndpoint(privatePath) {
		t.Fatal("did not expect /private to be public")
	}
}

func TestUtilityFunctions(t *testing.T) {
	if got := extractResource(userIDPath); got != "Users" {
		t.Fatalf("expected Users, got %s", got)
	}
	if got := extractResource("/"); got != "Resource" {
		t.Fatalf("expected fallback Resource, got %s", got)
	}

	if !contains([]string{"a", "b"}, "a") || contains([]string{"a", "b"}, "c") {
		t.Fatal("contains helper returned unexpected result")
	}
}
