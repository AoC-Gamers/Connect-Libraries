package apikey

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testServiceCore      = "connect-core"
	testServiceAuth      = "connect-auth"
	testValidKey         = "valid-key"
	testAuthBearer       = "Bearer bearer-key"
	testQueryWithAPIKey  = "/?api_key=query-key"
	testAuthorizationHdr = "Authorization"
)

func TestDefaultConfigValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatalf("expected config, got nil")
	}
	if cfg.HeaderName != "X-Internal-API-Key" {
		t.Fatalf("unexpected header name: %s", cfg.HeaderName)
	}
	if !cfg.AllowBearer || !cfg.AllowQuery {
		t.Fatalf("expected bearer and query enabled by default")
	}
	if cfg.QueryParam != "api_key" {
		t.Fatalf("unexpected query param: %s", cfg.QueryParam)
	}
}

func TestValidateKey(t *testing.T) {
	validator := NewValidator(map[string]string{testValidKey: testServiceCore})

	service, ok := validator.ValidateKey(testValidKey)
	if !ok {
		t.Fatalf("expected key to be valid")
	}
	if service != testServiceCore {
		t.Fatalf("unexpected service: %s", service)
	}

	if _, ok = validator.ValidateKey(""); ok {
		t.Fatalf("expected empty key to be invalid")
	}
	if _, ok = validator.ValidateKey("invalid"); ok {
		t.Fatalf("expected unknown key to be invalid")
	}
}

func TestExtractAPIKeyPrecedence(t *testing.T) {
	validator := NewValidator(map[string]string{})
	cfg := DefaultConfig()

	req := httptest.NewRequest(http.MethodGet, testQueryWithAPIKey, nil)
	req.Header.Set(testAuthorizationHdr, testAuthBearer)
	req.Header.Set(cfg.HeaderName, "header-key")

	if got := validator.ExtractAPIKey(req, cfg); got != "header-key" {
		t.Fatalf("expected header key, got %q", got)
	}
}

func TestExtractAPIKeyBearerAndQuery(t *testing.T) {
	validator := NewValidator(map[string]string{})

	t.Run("bearer", func(t *testing.T) {
		cfg := DefaultConfig()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(testAuthorizationHdr, testAuthBearer)

		if got := validator.ExtractAPIKey(req, cfg); got != "bearer-key" {
			t.Fatalf("expected bearer key, got %q", got)
		}
	})

	t.Run("query", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.AllowBearer = false
		req := httptest.NewRequest(http.MethodGet, testQueryWithAPIKey, nil)

		if got := validator.ExtractAPIKey(req, cfg); got != "query-key" {
			t.Fatalf("expected query key, got %q", got)
		}
	})

	t.Run("none", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.AllowBearer = false
		cfg.AllowQuery = false
		req := httptest.NewRequest(http.MethodGet, testQueryWithAPIKey, nil)
		req.Header.Set(testAuthorizationHdr, testAuthBearer)

		if got := validator.ExtractAPIKey(req, cfg); got != "" {
			t.Fatalf("expected empty key, got %q", got)
		}
	})
}

func TestHasServiceAndListServices(t *testing.T) {
	validator := NewValidator(map[string]string{
		"k1": testServiceAuth,
		"k2": testServiceCore,
		"k3": testServiceCore,
	})

	if !validator.HasService(testServiceCore) {
		t.Fatalf("expected service %s to exist", testServiceCore)
	}
	if validator.HasService("connect-rt") {
		t.Fatalf("did not expect service connect-rt")
	}

	services := validator.ListServices()
	if len(services) != 2 {
		t.Fatalf("expected 2 unique services, got %d", len(services))
	}
}
