package testhelpers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// NewMockDB creates a sqlmock-backed *sql.DB and returns the mock.
// This is useful for code that uses database/sql directly (not sqlx).
// The cleanup is automatic via t.Cleanup().
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
		_ = db.Close()
	})

	return db, mock
}

// NewSQLMock creates a sqlmock-backed *sqlx.DB and returns the sqlmock and a cleanup func.
// This is a generic helper that can be used by any module for testing SQL operations.
func NewSQLMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	cleanup := func() {
		// ensure expectations are met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sqlmock expectations: %v", err)
		}
		_ = sqlxDB.Close()
	}
	return sqlxDB, mock, cleanup
}

// NewMiniredisClient starts a miniredis server and returns a standard go-redis *redis.Client
// along with a cleanup function. This uses public types and can be used across all modules.
// If you need a module-specific wrapper (like Connect-Core's internal/redis.Client),
// wrap this in your module's internal/testhelpers.
func NewMiniredisClient(t *testing.T) (*redis.Client, *miniredis.Miniredis, func()) {
	mrs, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}

	rclient := redis.NewClient(&redis.Options{
		Addr: mrs.Addr(),
	})

	cleanup := func() {
		_ = rclient.Close()
		mrs.Close()
	}
	return rclient, mrs, cleanup
}

// MakeAuthServer returns an httptest.Server that responds to Connect-Auth permission checks
// with hasPermission equal to the provided value.
func MakeAuthServer(has bool) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/auth/internal/permissions/check" {
			resp := map[string]bool{"hasPermission": has}
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		http.NotFound(w, r)
	}))
}
