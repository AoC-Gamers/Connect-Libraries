package settingsruntime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewCoreHTTPClient_InvalidURL(t *testing.T) {
	if _, err := NewCoreHTTPClient("://bad-url", "", time.Second); err == nil {
		t.Fatal("expected invalid URL error")
	}
}

func TestCoreHTTPClient_Health(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewCoreHTTPClient(server.URL, "", time.Second)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("unexpected health error: %v", err)
	}
}

func TestCoreHTTPClient_GetSettingValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/core/internal/settings/CONFIG/rate_limit.enabled/value" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-Internal-API-Key") != "secret" {
			http.Error(w, "missing api key", http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"value":"true"}`))
	}))
	defer server.Close()

	client, err := NewCoreHTTPClient(server.URL, "secret", time.Second)
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}

	value, err := client.GetSettingValue(context.Background(), "config", "rate_limit.enabled")
	if err != nil {
		t.Fatalf("unexpected get setting error: %v", err)
	}
	if value != "true" {
		t.Fatalf("unexpected value: %q", value)
	}
}
