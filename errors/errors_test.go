package errors_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	errors "github.com/AoC-Gamers/connect-libraries/errors"
)

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		handler        func(w http.ResponseWriter)
		expectedStatus int
		expectedCode   errors.ErrorCode
		expectedError  string
	}{
		{
			name: "Permission Denied",
			handler: func(w http.ResponseWriter) {
				errors.RespondPermissionDenied(w, 1, "WEB", "WEB__MISSION_VIEW", false)
			},
			expectedStatus: http.StatusForbidden,
			expectedCode:   errors.CodePermissionDenied,
			expectedError:  "insufficient permissions",
		},
		{
			name: "Not Found",
			handler: func(w http.ResponseWriter) {
				errors.RespondNotFound(w, "mission", "test-mission")
			},
			expectedStatus: http.StatusNotFound,
			expectedCode:   errors.CodeNotFound,
			expectedError:  "not found",
		},
		{
			name: "Token Expired",
			handler: func(w http.ResponseWriter) {
				errors.RespondTokenExpired(w)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCode:   errors.CodeTokenExpired,
			expectedError:  "unauthorized",
		},
		{
			name: "Validation Error",
			handler: func(w http.ResponseWriter) {
				errors.RespondValidationError(w, "name", "exceeds 128 characters")
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   errors.CodeValidationError,
			expectedError:  "validation failed",
		},
		{
			name: "Internal Error",
			handler: func(w http.ResponseWriter) {
				errors.RespondInternalError(w, "database connection failed")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   errors.CodeInternalError,
			expectedError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			tt.handler(rr)

			// Check status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// Check content type
			if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("handler returned wrong content type: got %v want application/json", ct)
			}

			// Parse response
			var response errors.ErrorResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("could not parse response: %v", err)
			}

			// Check code
			if response.Code != tt.expectedCode {
				t.Errorf("handler returned wrong code: got %v want %v", response.Code, tt.expectedCode)
			}

			// Check error message
			if response.Error != tt.expectedError {
				t.Errorf("handler returned wrong error message: got %v want %v", response.Error, tt.expectedError)
			}

			// Check status in response
			if response.Status != tt.expectedStatus {
				t.Errorf("handler returned wrong status in response: got %v want %v", response.Status, tt.expectedStatus)
			}
		})
	}
}

func TestErrorCodeHTTPStatus(t *testing.T) {
	tests := []struct {
		code           errors.ErrorCode
		expectedStatus int
	}{
		{errors.CodeUnauthorized, 401},
		{errors.CodeTokenExpired, 401},
		{errors.CodePermissionDenied, 403},
		{errors.CodeNotFound, 404},
		{errors.CodeAlreadyExists, 409},
		{errors.CodeValidationError, 400},
		{errors.CodeRateLimitExceeded, 429},
		{errors.CodeServiceUnavailable, 503},
		{errors.CodeTimeout, 504},
		{errors.CodeInternalError, 500},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			if status := tt.code.HTTPStatus(); status != tt.expectedStatus {
				t.Errorf("code %v returned wrong HTTP status: got %v want %v", tt.code, status, tt.expectedStatus)
			}
		})
	}
}

func TestMetadata(t *testing.T) {
	rr := httptest.NewRecorder()

	errors.RespondPermissionDenied(rr, 123, "TEAM", "TEAM__MANAGE", true)

	var response errors.ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	// Check metadata exists
	if response.Meta == nil {
		t.Fatal("expected metadata to be present")
	}

	// Check specific metadata fields
	if scopeID, ok := response.Meta["scope_id"].(float64); !ok || int64(scopeID) != 123 {
		t.Errorf("expected scope_id to be 123, got %v", response.Meta["scope_id"])
	}

	if scopeType, ok := response.Meta["scope_type"].(string); !ok || scopeType != "TEAM" {
		t.Errorf("expected scope_type to be TEAM, got %v", response.Meta["scope_type"])
	}

	if permission, ok := response.Meta["required_permission"].(string); !ok || permission != "TEAM__MANAGE" {
		t.Errorf("expected required_permission to be TEAM__MANAGE, got %v", response.Meta["required_permission"])
	}

	if hasMembership, ok := response.Meta["has_membership"].(bool); !ok || !hasMembership {
		t.Errorf("expected has_membership to be true, got %v", response.Meta["has_membership"])
	}
}

func TestDetailedError(t *testing.T) {
	rr := httptest.NewRecorder()

	errors.RespondOutOfRange(rr, "age", 0, 120, 150)

	var response errors.ErrorResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	// Check detail is present
	if response.Detail == "" {
		t.Error("expected detail to be present")
	}

	// Check metadata
	if response.Meta == nil {
		t.Fatal("expected metadata to be present")
	}

	if field, ok := response.Meta["field"].(string); !ok || field != "age" {
		t.Errorf("expected field to be 'age', got %v", response.Meta["field"])
	}
}

func TestLegacyError(t *testing.T) {
	rr := httptest.NewRecorder()

	errors.RespondLegacyError(rr, http.StatusBadRequest, "invalid request")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	if response["error"] != "invalid request" {
		t.Errorf("expected error to be 'invalid request', got %v", response["error"])
	}

	// Legacy format should not have other fields
	if len(response) != 1 {
		t.Errorf("expected only 'error' field, got %d fields", len(response))
	}
}
