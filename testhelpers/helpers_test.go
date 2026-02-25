package testhelpers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

const checkPermissionPath = "/auth/internal/permissions/check"
const selectOneQuery = "SELECT 1"

func TestNewMockDB(t *testing.T) {
	db, mock := NewMockDB(t)
	if db == nil || mock == nil {
		t.Fatal("expected non-nil sql.DB and sqlmock")
	}

	mock.ExpectExec(selectOneQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	if _, err := db.Exec(selectOneQuery); err != nil {
		t.Fatalf("expected exec success on mocked db, got %v", err)
	}
}

func TestNewSQLMock(t *testing.T) {
	sqlxDB, mock, cleanup := NewSQLMock(t)
	defer cleanup()

	if sqlxDB == nil || mock == nil {
		t.Fatal("expected non-nil sqlx db and sqlmock")
	}

	mock.ExpectExec(selectOneQuery).WillReturnResult(sqlmock.NewResult(0, 0))
	if _, err := sqlxDB.Exec(selectOneQuery); err != nil {
		t.Fatalf("expected exec success on sqlx mock, got %v", err)
	}
}

func TestNewMiniredisClient(t *testing.T) {
	client, miniRedis, cleanup := NewMiniredisClient(t)
	defer cleanup()

	if client == nil || miniRedis == nil {
		t.Fatal("expected non-nil redis client and miniredis")
	}

	ctx := context.Background()
	if err := client.Set(ctx, "sample-key", "sample-value", 0).Err(); err != nil {
		t.Fatalf("expected set success, got %v", err)
	}

	value, err := client.Get(ctx, "sample-key").Result()
	if err != nil {
		t.Fatalf("expected get success, got %v", err)
	}
	if value != "sample-value" {
		t.Fatalf("expected sample-value, got %s", value)
	}
}

func TestMakeAuthServer(t *testing.T) {
	server := MakeAuthServer(true)
	defer server.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+checkPermissionPath, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var payload map[string]bool
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to decode json response: %v", err)
	}
	if hasPermission := payload["hasPermission"]; !hasPermission {
		t.Fatalf("expected hasPermission=true, got %v", hasPermission)
	}

	notFoundResponse, notFoundErr := http.Get(server.URL + "/unknown")
	if notFoundErr != nil {
		t.Fatalf("failed to request unknown path: %v", notFoundErr)
	}
	defer notFoundResponse.Body.Close()
	if notFoundResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404 for unknown path, got %d", notFoundResponse.StatusCode)
	}
}
