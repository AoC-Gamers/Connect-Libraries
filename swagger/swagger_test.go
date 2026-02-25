package swagger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestPublicAPIsAndServeSwaggerSpec(t *testing.T) {
	cfg := DefaultConfig().WithServiceInfo("Swagger Test", "0.0.1")
	d := New(cfg)
	if d == nil {
		t.Fatal("expected detector instance")
	}

	router := chi.NewRouter()
	router.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	if err := d.ScanRouter(router); err != nil {
		t.Fatalf("unexpected scan error: %v", err)
	}

	routesJSON, err := ExportJSON(d)
	if err != nil {
		t.Fatalf("unexpected ExportJSON error: %v", err)
	}
	if len(routesJSON) == 0 {
		t.Fatal("expected non-empty routes JSON")
	}

	specJSON, err := ExportSpec(d)
	if err != nil {
		t.Fatalf("unexpected ExportSpec error: %v", err)
	}
	var spec map[string]interface{}
	if err := json.Unmarshal(specJSON, &spec); err != nil {
		t.Fatalf("could not parse exported spec: %v", err)
	}

	handler := ServeSwaggerSpec(d)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/swagger.json", nil)
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected content type: %s", rr.Header().Get("Content-Type"))
	}
	if rr.Body.Len() == 0 {
		t.Fatal("expected non-empty response body")
	}
}
