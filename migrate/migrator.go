package migrate

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/rs/zerolog/log"
)

// Config holds the configuration for the migrator
type Config struct {
	ServiceName    string   // e.g., "Connect-Auth", "Connect-Core", "Connect-RT"
	SchemaName     string   // e.g., "auth", "core", "rt"
	MigrationsDir  string   // e.g., "migrations_sql"
	DataDir        string   // e.g., "data_sql" (optional, for initial data)
	ApplyData      bool     // true = apply data files automatically after migrations
	CriticalTables []string // Tables to verify after migration
}

// Migrator handles database migrations
type Migrator struct {
	config    Config
	schemaSQL string
	db        *sql.DB
}

// New creates a new Migrator instance
func New(config Config) (*Migrator, error) {
	if config.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if config.SchemaName == "" {
		return nil, fmt.Errorf("schema name is required")
	}
	schemaSQL, err := quoteIdentifier(config.SchemaName)
	if err != nil {
		return nil, fmt.Errorf("invalid schema name: %w", err)
	}
	if config.MigrationsDir == "" {
		config.MigrationsDir = "migrations_sql"
	}

	return &Migrator{
		config:    config,
		schemaSQL: schemaSQL,
	}, nil
}

// Connect establishes database connection using environment variables
func (m *Migrator) Connect() error {
	dbHost := getEnvRequired("POSTGRES_HOST")
	dbPort := getEnvRequired("POSTGRES_PORT")
	dbUser := getEnvRequired("POSTGRES_USER")
	dbPassword := getEnvRequired("POSTGRES_PASSWORD")
	dbName := getEnvRequired("POSTGRES_NAME")
	sslMode := getEnvRequired("POSTGRES_SSLMODE")

	// Construir DSN para pgx
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode,
	)

	// Conectar directamente con database/sql usando pgx driver
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	m.db = db

	return nil
}

// Close closes the database connection
func (m *Migrator) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// Run executes the full migration process
func (m *Migrator) Run() error {
	// Ensure schema exists
	if err := m.ensureSchema(); err != nil {
		return err
	}

	if err := m.setSearchPath(); err != nil {
		return err
	}

	// Create migrations tracking table
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Apply migrations
	appliedCount, _, err := m.applyMigrations()
	if err != nil {
		return err
	}

	// Verify critical tables if any migrations were applied
	if appliedCount > 0 && len(m.config.CriticalTables) > 0 {
		if err := m.verifyCriticalTables(); err != nil {
			return err
		}
	}

	// Apply data files if configured
	if m.config.ApplyData && m.config.DataDir != "" {
		if err := m.ApplyFixtures(); err != nil {
			return fmt.Errorf("failed to apply data files: %w", err)
		}
	}

	return nil
}

// ensureSchema creates the schema if it doesn't exist
func (m *Migrator) ensureSchema() error {
	// #nosec G202 -- schema identifier is strictly validated by quoteIdentifier() in New()
	query := "CREATE SCHEMA IF NOT EXISTS " + m.schemaSQL
	if _, err := m.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (m *Migrator) setSearchPath() error {
	searchPath := m.config.SchemaName + ",public"
	if _, err := m.db.Exec(`SELECT set_config('search_path', $1, false)`, searchPath); err != nil {
		return fmt.Errorf("failed to set search_path: %w", err)
	}

	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		)
	`

	if _, err := m.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// applyMigrations reads and applies all pending migrations
func (m *Migrator) applyMigrations() (int, int, error) {
	files, migrationsRoot, err := m.getSortedMigrationFiles()
	if err != nil {
		return 0, 0, err
	}
	defer migrationsRoot.Close()

	if len(files) == 0 {
		log.Warn().Msg("⚠️  No migration files found")
		return 0, 0, nil
	}

	appliedCount := 0
	skippedCount := 0

	for _, file := range files {
		applied, err := m.applySingleMigration(file, migrationsRoot)
		if err != nil {
			return appliedCount, skippedCount, err
		}

		if applied {
			appliedCount++
		} else {
			skippedCount++
		}
	}

	if appliedCount > 0 {
		log.Info().Int("count", appliedCount).Msg("Migrations applied successfully")
	} else {
		log.Info().Msg("All migrations up to date")
	}

	return appliedCount, skippedCount, nil
}

func (m *Migrator) getSortedMigrationFiles() ([]string, *os.Root, error) {
	if _, err := os.Stat(m.config.MigrationsDir); os.IsNotExist(err) {
		return nil, nil, fmt.Errorf("migrations directory not found: %s", m.config.MigrationsDir)
	}

	migrationsRoot, err := os.OpenRoot(m.config.MigrationsDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open migrations root: %w", err)
	}

	entries, err := os.ReadDir(m.config.MigrationsDir)
	if err != nil {
		if closeErr := migrationsRoot.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("failed to close migrations root")
		}
		return nil, nil, fmt.Errorf("failed to read migration files: %w", err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		files = append(files, entry.Name())
	}

	sort.Strings(files)
	return files, migrationsRoot, nil
}

// applySingleMigration applies a single migration file if not already applied
func (m *Migrator) applySingleMigration(file string, migrationsRoot *os.Root) (bool, error) {
	basename := filepath.Base(file)

	// Check if already applied
	var exists bool
	checkSQL := "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)"
	if err := m.db.QueryRow(checkSQL, basename).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check migration status for %s: %w", basename, err)
	}

	if exists {
		log.Debug().Msgf("Skipped: %s (already applied)", basename)
		return false, nil
	}

	// Read file content
	log.Debug().Msgf("Applying: %s", basename)
	fileHandle, err := migrationsRoot.Open(file)
	if err != nil {
		return false, fmt.Errorf("failed to open file %s: %w", basename, err)
	}
	defer fileHandle.Close()

	content, err := io.ReadAll(fileHandle)
	if err != nil {
		return false, fmt.Errorf("failed to read file %s: %w", basename, err)
	}

	// Execute SQL
	// #nosec G701 -- migration SQL is loaded from trusted local migration files under migrations root
	if _, err := m.db.Exec(string(content)); err != nil {
		return false, fmt.Errorf("migration failed for %s: %w", basename, err)
	}

	// Mark as applied
	insertSQL := "INSERT INTO schema_migrations (version) VALUES ($1)"
	if _, err := m.db.Exec(insertSQL, basename); err != nil {
		return false, fmt.Errorf("failed to mark migration as applied for %s: %w", basename, err)
	}

	log.Info().Msgf("Applied: %s", basename)
	return true, nil
}

// verifyCriticalTables checks that all critical tables exist
func (m *Migrator) verifyCriticalTables() error {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.tables 
			WHERE table_schema = $1 
			  AND table_name = $2
		)
	`

	for _, table := range m.config.CriticalTables {
		var exists bool
		if err := m.db.QueryRow(query, m.config.SchemaName, table).Scan(&exists); err != nil {
			return fmt.Errorf("verification query failed for table %s: %w", table, err)
		}

		if !exists {
			return fmt.Errorf("table %s.%s not found", m.config.SchemaName, table)
		}

		log.Debug().Msgf("Verified: %s.%s", m.config.SchemaName, table)
	}

	log.Debug().Msg("All critical tables verified")
	return nil
}

// getEnvRequired gets an environment variable or fails if not set
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal().Msgf("❌ Required environment variable not set: %s", key)
	}
	return value
}
