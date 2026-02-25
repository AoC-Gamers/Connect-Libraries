package config

import "testing"

func TestDefaultConfigAndBuilderMethods(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ServiceName == "" || cfg.Version == "" {
		t.Fatal("expected default service info")
	}
	if cfg.DefaultTag == "" || cfg.DefaultSecurity == "" {
		t.Fatal("expected default tag and security")
	}
	if len(cfg.SecurityPatterns) == 0 {
		t.Fatal("expected default security patterns")
	}

	rules := []TagRule{{PathPattern: "/users", Tags: []string{"Users"}}}
	returned := cfg.
		WithServiceInfo("Connect API", "2.0.0").
		WithDescription("desc").
		WithContact("Team", "https://example.com", "team@example.com").
		WithLicense("MIT", "https://opensource.org/licenses/MIT").
		WithTagRules(rules).
		AddTagRule("/admin", "Admin").
		AddPublicPath("/public").
		AddSkipPath("/internal").
		WithDefaultSecurity("ApiKeyAuth").
		WithDefaultTag("Default")

	if returned != cfg {
		t.Fatal("expected builder methods to return same config pointer")
	}
	if cfg.ServiceName != "Connect API" || cfg.Version != "2.0.0" || cfg.Description != "desc" {
		t.Fatal("expected service info and description to be set")
	}
	if cfg.ContactName != "Team" || cfg.ContactEmail != "team@example.com" {
		t.Fatal("expected contact info to be set")
	}
	if cfg.LicenseName != "MIT" {
		t.Fatal("expected license to be set")
	}
	if len(cfg.TagRules) != 2 {
		t.Fatalf("expected 2 tag rules, got %d", len(cfg.TagRules))
	}
	if cfg.DefaultSecurity != "ApiKeyAuth" || cfg.DefaultTag != "Default" {
		t.Fatal("expected defaults to be overridden")
	}
}
