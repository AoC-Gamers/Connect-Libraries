package errors

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	errParseInternalResponse = "could not parse internal response: %v"
	errParseResponseBody     = "could not parse response body: %v"
	errExpectedStatus        = "expected status %d, got %d"
	errExpectedCode          = "expected code %s, got %s"
	errDatabaseFailed        = "Database operation failed"
	errConnectionLost        = "connection lost"
)

func TestDetectServiceName(t *testing.T) {
	t.Run("uses service header when present", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/any/path", nil)
		req.Header.Set("X-Service-Name", "custom-service")

		service := detectServiceName(req)
		if service != "custom-service" {
			t.Fatalf("expected custom-service, got %s", service)
		}
	})

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{name: "auth path", path: "/auth/session", expected: "connect-auth"},
		{name: "core path", path: "/core/users", expected: "connect-core"},
		{name: "lobby path", path: "/lobby/chat", expected: "connect-lobby"},
		{name: "rt path", path: "/rt/events", expected: "connect-rt"},
		{name: "unknown path", path: "/unknown/path", expected: "connect-service"},
		{name: "empty path", path: "/", expected: "connect-service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			service := detectServiceName(req)
			if service != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, service)
			}
		})
	}
}

func TestGetInternalStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		code     InternalErrorCode
		expected int
	}{
		{name: "unauthorized", code: InternalUnauthorized, expected: http.StatusUnauthorized},
		{name: "forbidden", code: InternalForbidden, expected: http.StatusForbidden},
		{name: "not found", code: InternalNotFound, expected: http.StatusNotFound},
		{name: "conflict", code: InternalConflict, expected: http.StatusConflict},
		{name: "validation", code: InternalValidation, expected: http.StatusBadRequest},
		{name: "bad request", code: InternalBadRequest, expected: http.StatusBadRequest},
		{name: "rate limit", code: InternalRateLimit, expected: http.StatusTooManyRequests},
		{name: "service down", code: InternalServiceDown, expected: http.StatusServiceUnavailable},
		{name: "database", code: InternalDatabase, expected: http.StatusInternalServerError},
		{name: "timeout", code: InternalTimeout, expected: http.StatusInternalServerError},
		{name: "server error", code: InternalServerError, expected: http.StatusInternalServerError},
		{name: "default", code: InternalAuthzCheck, expected: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := getInternalStatusCode(tt.code)
			if status != tt.expected {
				t.Fatalf("code %s: expected %d, got %d", tt.code, tt.expected, status)
			}
		})
	}
}

func TestRespondInternalServiceError(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/core/teams", nil)

	details := map[string]interface{}{"reason": "policy_mismatch"}
	RespondInternalServiceError(rr, req, InternalForbidden, "service not authorized", details)

	if rr.Code != http.StatusForbidden {
		t.Fatalf(errExpectedStatus, http.StatusForbidden, rr.Code)
	}

	var response InternalErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf(errParseInternalResponse, err)
	}

	if response.Code != InternalForbidden {
		t.Fatalf(errExpectedCode, InternalForbidden, response.Code)
	}
	if response.Service != "connect-core" {
		t.Fatalf("expected service connect-core, got %s", response.Service)
	}
	if response.Status != http.StatusForbidden {
		t.Fatalf(errExpectedStatus, http.StatusForbidden, response.Status)
	}
}

func TestRespondJSON(t *testing.T) {
	t.Run("writes status without body when data is nil", func(t *testing.T) {
		rr := httptest.NewRecorder()

		RespondJSON(rr, http.StatusNoContent, nil)

		if rr.Code != http.StatusNoContent {
			t.Fatalf(errExpectedStatus, http.StatusNoContent, rr.Code)
		}
		if rr.Body.Len() != 0 {
			t.Fatalf("expected empty body, got %q", rr.Body.String())
		}
	})

	t.Run("encodes provided payload", func(t *testing.T) {
		rr := httptest.NewRecorder()

		payload := map[string]string{"result": "ok"}
		RespondJSON(rr, http.StatusOK, payload)

		if rr.Code != http.StatusOK {
			t.Fatalf(errExpectedStatus, http.StatusOK, rr.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
			t.Fatalf(errParseResponseBody, err)
		}
		if response["result"] != "ok" {
			t.Fatalf("expected result=ok, got %q", response["result"])
		}
	})
}

func TestErrorHelpersTokenInvalid(t *testing.T) {
	rr := httptest.NewRecorder()
	RespondTokenInvalid(rr, "signature mismatch")

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf(errParseResponseBody, err)
	}

	if response.Code != CodeTokenInvalid {
		t.Fatalf(errExpectedCode, CodeTokenInvalid, response.Code)
	}
	if shouldReauth, ok := response.Meta["should_reauth"].(bool); !ok || !shouldReauth {
		t.Fatalf("expected should_reauth=true, got %v", response.Meta["should_reauth"])
	}
}

func TestErrorHelpersValidationErrors(t *testing.T) {
	rr := httptest.NewRecorder()
	validationErrors := map[string]string{"email": "required"}

	RespondValidationErrors(rr, validationErrors)

	var response ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf(errParseResponseBody, err)
	}

	if response.Code != CodeValidationError {
		t.Fatalf(errExpectedCode, CodeValidationError, response.Code)
	}
	if response.Meta == nil || response.Meta["errors"] == nil {
		t.Fatal("expected errors metadata to be present")
	}
}

func TestErrorHelpersDatabaseError(t *testing.T) {
	rrNil := httptest.NewRecorder()
	RespondDatabaseError(rrNil, nil)

	var nilResponse ErrorResponse
	if err := json.Unmarshal(rrNil.Body.Bytes(), &nilResponse); err != nil {
		t.Fatalf(errParseResponseBody, err)
	}
	if nilResponse.Detail != errDatabaseFailed {
		t.Fatalf("unexpected nil error detail: %s", nilResponse.Detail)
	}

	rrErr := httptest.NewRecorder()
	RespondDatabaseError(rrErr, errors.New(errConnectionLost))

	var errResponse ErrorResponse
	if err := json.Unmarshal(rrErr.Body.Bytes(), &errResponse); err != nil {
		t.Fatalf(errParseResponseBody, err)
	}
	if errResponse.Detail != errDatabaseFailed+": "+errConnectionLost {
		t.Fatalf("unexpected error detail: %s", errResponse.Detail)
	}
}
