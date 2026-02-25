package apikey

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	middlewareHeaderAPIKey = "X-Internal-API-Key"
	middlewareKeyValid     = "valid-key"
	middlewareCoreKey      = "core-key"
	middlewareEnvCoreKey   = "CORE_API_KEY"
)

type testResponder struct {
	unauthorizedCalled bool
	forbiddenCalled    bool
}

func (r *testResponder) Unauthorized(w http.ResponseWriter, detail string) {
	r.unauthorizedCalled = true
	w.WriteHeader(http.StatusUnauthorized)
}

func (r *testResponder) InsufficientPermissions(w http.ResponseWriter, action string) {
	r.forbiddenCalled = true
	w.WriteHeader(http.StatusForbidden)
}

func TestRequireAPIKeyWithResponderUnauthorized(t *testing.T) {
	validator := NewValidator(map[string]string{middlewareKeyValid: testServiceCore})
	responder := &testResponder{}
	mw := RequireAPIKeyWithResponder(validator, responder)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	if !responder.unauthorizedCalled {
		t.Fatalf("expected unauthorized responder to be called")
	}
}

func TestRequireAPIKeyWithResponderSetsContext(t *testing.T) {
	validator := NewValidator(map[string]string{middlewareKeyValid: testServiceAuth})
	mw := RequireAPIKeyWithResponder(validator, nil)

	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsServiceAuthenticated(r) {
			t.Fatalf("expected service authenticated")
		}
		if got := GetServiceNameFromContext(r); got != testServiceAuth {
			t.Fatalf("unexpected service name: %s", got)
		}
		if !IsAuthService(r) {
			t.Fatalf("expected auth service")
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(middlewareHeaderAPIKey, middlewareKeyValid)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestRequireConnectServiceWithResponder(t *testing.T) {
	t.Setenv(middlewareEnvCoreKey, middlewareCoreKey)
	responder := &testResponder{}

	allowedMW := RequireConnectServiceWithResponder(responder, "connect-core")
	allowedHandler := allowedMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	allowedReq := httptest.NewRequest(http.MethodGet, "/", nil)
	allowedReq.Header.Set(middlewareHeaderAPIKey, middlewareCoreKey)
	allowedRR := httptest.NewRecorder()
	allowedHandler.ServeHTTP(allowedRR, allowedReq)

	if allowedRR.Code != http.StatusOK {
		t.Fatalf("expected 200 for allowed service, got %d", allowedRR.Code)
	}

	notAllowedMW := RequireConnectServiceWithResponder(responder, "connect-auth")
	notAllowedHandler := notAllowedMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	notAllowedReq := httptest.NewRequest(http.MethodGet, "/", nil)
	notAllowedReq.Header.Set(middlewareHeaderAPIKey, middlewareCoreKey)
	notAllowedRR := httptest.NewRecorder()
	notAllowedHandler.ServeHTTP(notAllowedRR, notAllowedReq)

	if notAllowedRR.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for non-allowed service, got %d", notAllowedRR.Code)
	}
	if !responder.forbiddenCalled {
		t.Fatalf("expected insufficient permissions responder to be called")
	}
}
