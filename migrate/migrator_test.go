package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
)

const testServiceName = "Connect-Test"

func TestNewMigratorValidationAndDefaults(t *testing.T) {
	if _, err := New(Config{}); err == nil {
		t.Fatal("expected error when service name is missing")
	}

	if _, err := New(Config{ServiceName: testServiceName}); err == nil {
		t.Fatal("expected error when schema name is missing")
	}

	migrator, err := New(Config{ServiceName: testServiceName, SchemaName: "test"})
	if err != nil {
		t.Fatalf("expected valid migrator, got error: %v", err)
	}
	if migrator.config.MigrationsDir != "migrations_sql" {
		t.Fatalf("expected default migrations dir migrations_sql, got %s", migrator.config.MigrationsDir)
	}
}

func TestApplyMigrationsDirectoryMissing(t *testing.T) {
	migrator, err := New(Config{
		ServiceName:   testServiceName,
		SchemaName:    "test",
		MigrationsDir: filepath.Join(t.TempDir(), "does-not-exist"),
	})
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	_, _, runErr := migrator.applyMigrations()
	if runErr == nil {
		t.Fatal("expected missing directory error")
	}
}

func TestApplyFixturesNoDataDirOrMissingDir(t *testing.T) {
	migrator, err := New(Config{ServiceName: testServiceName, SchemaName: "test"})
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	if applyErr := migrator.ApplyFixtures(); applyErr != nil {
		t.Fatalf("expected nil when data dir not configured, got %v", applyErr)
	}

	migrator.config.DataDir = filepath.Join(t.TempDir(), "missing-data-dir")
	if applyErr := migrator.ApplyFixtures(); applyErr != nil {
		t.Fatalf("expected nil when data dir does not exist, got %v", applyErr)
	}
}

func TestGetEnvFallback(t *testing.T) {
	t.Setenv("MIGRATE_TEST_ENV", "configured")
	if value := getEnv("MIGRATE_TEST_ENV", "fallback"); value != "configured" {
		t.Fatalf("expected configured value, got %s", value)
	}

	t.Setenv("MIGRATE_TEST_ENV", "")
	if value := getEnv("MIGRATE_TEST_ENV", "fallback"); value != "fallback" {
		t.Fatalf("expected fallback value, got %s", value)
	}
}

func TestSetupLoggerLevels(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("LOG_FORMAT", "json")
	SetupLogger()
	if zerolog.GlobalLevel() != zerolog.DebugLevel {
		t.Fatalf("expected debug level, got %s", zerolog.GlobalLevel())
	}

	t.Setenv("LOG_LEVEL", "warn")
	SetupLogger()
	if zerolog.GlobalLevel() != zerolog.WarnLevel {
		t.Fatalf("expected warn level, got %s", zerolog.GlobalLevel())
	}

	t.Setenv("LOG_LEVEL", "invalid")
	t.Setenv("LOG_FORMAT", "console")
	SetupLogger()
	if zerolog.GlobalLevel() != zerolog.InfoLevel {
		t.Fatalf("expected fallback info level, got %s", zerolog.GlobalLevel())
	}
}

func TestFixturesHelpers(t *testing.T) {
	dataDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dataDir, "005_users_data.sql"), []byte("SELECT 1;"), 0600); err != nil {
		t.Fatalf("failed to create fixture file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "README.md"), []byte("docs"), 0600); err != nil {
		t.Fatalf("failed to create readme: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "nonnumbered.sql"), []byte("SELECT 1;"), 0600); err != nil {
		t.Fatalf("failed to create non-numbered fixture file: %v", err)
	}

	fm := &FixturesManager{dataDir: dataDir, schemaName: "test"}
	files, err := fm.getFixtureFiles()
	if err != nil {
		t.Fatalf("unexpected getFixtureFiles error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 fixture file, got %d", len(files))
	}
	if filepath.Base(files[0]) != "005_users_data.sql" {
		t.Fatalf("unexpected fixture file %s", filepath.Base(files[0]))
	}

	version := fm.extractVersion(filepath.Join(dataDir, "010_seed_roles.sql"))
	if version != "010_seed_roles" {
		t.Fatalf("expected version 010_seed_roles, got %s", version)
	}
}
