package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// FixturesManager handles initial data loading
type FixturesManager struct {
	db         *sql.DB
	schemaName string
	dataDir    string
}

// ApplyFixtures applies all pending fixtures from data_sql directory
func (m *Migrator) ApplyFixtures() error {
	if m.config.DataDir == "" {
		log.Info().Msg("Data directory not configured, skipping fixtures")
		return nil
	}

	// Check if data directory exists
	if _, err := os.Stat(m.config.DataDir); os.IsNotExist(err) {
		log.Info().
			Str("dataDir", m.config.DataDir).
			Msg("Data directory does not exist, skipping fixtures")
		return nil
	}

	fm := &FixturesManager{
		db:         m.db,
		schemaName: m.config.SchemaName,
		dataDir:    m.config.DataDir,
	}

	return fm.applyAllFixtures()
}

// ensureFixturesTable creates the tracking table for fixtures
func (fm *FixturesManager) ensureFixturesTable() error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s.schema_migrations_data (
			version VARCHAR(50) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`, fm.schemaName)

	_, err := fm.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations_data table: %w", err)
	}

	log.Debug().
		Str("table", fmt.Sprintf("%s.schema_migrations_data", fm.schemaName)).
		Msg("Fixtures tracking table ensured")

	return nil
}

// applyAllFixtures applies all pending data fixtures
func (fm *FixturesManager) applyAllFixtures() error {
	// 1. Ensure tracking table exists
	if err := fm.ensureFixturesTable(); err != nil {
		return err
	}

	// 2. Get list of fixture files
	fixtureFiles, err := fm.getFixtureFiles()
	if err != nil {
		return err
	}

	if len(fixtureFiles) == 0 {
		log.Debug().Msg("No fixture files found")
		return nil
	}

	// 3. Get already applied fixtures
	appliedFixtures, err := fm.getAppliedFixtures()
	if err != nil {
		return err
	}

	// 4. Apply each pending fixture
	appliedCount := 0
	for _, fixtureFile := range fixtureFiles {
		version := fm.extractVersion(fixtureFile)

		// Skip if already applied
		if _, exists := appliedFixtures[version]; exists {
			log.Debug().
				Str("fixture", filepath.Base(fixtureFile)).
				Msg("Fixture already applied, skipping")
			continue
		}

		// Apply fixture
		if err := fm.applyFixtureFile(fixtureFile); err != nil {
			return fmt.Errorf("error applying fixture %s: %w", fixtureFile, err)
		}
		appliedCount++
	}

	if appliedCount > 0 {
		log.Info().
			Int("count", appliedCount).
			Msg("Fixtures applied successfully")
	} else {
		log.Info().Msg("All fixtures up to date")
	}

	return nil
}

// getFixtureFiles returns sorted list of fixture files
func (fm *FixturesManager) getFixtureFiles() ([]string, error) {
	files, err := os.ReadDir(fm.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	var fixtureFiles []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Only .sql files
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Skip README and other non-numbered files
		if !strings.Contains(file.Name(), "_") {
			continue
		}

		fixtureFiles = append(fixtureFiles, filepath.Join(fm.dataDir, file.Name()))
	}

	// Sort by filename (ensures 005_xxx comes before 006_xxx)
	sort.Strings(fixtureFiles)

	return fixtureFiles, nil
}

// getAppliedFixtures returns map of already applied fixtures
func (fm *FixturesManager) getAppliedFixtures() (map[string]bool, error) {
	query := fmt.Sprintf(`
		SELECT version FROM %s.schema_migrations_data
	`, fm.schemaName)

	rows, err := fm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied fixtures: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// applyFixtureFile reads and executes a fixture file
func (fm *FixturesManager) applyFixtureFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read fixture file: %w", err)
	}

	fileName := filepath.Base(filePath)
	version := fm.extractVersion(filePath)

	log.Info().
		Str("fixture", fileName).
		Str("version", version).
		Msg("Applying fixture")

	// Execute fixture SQL
	if _, err := fm.db.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute fixture SQL: %w", err)
	}

	// Record fixture as applied
	query := fmt.Sprintf(`
		INSERT INTO %s.schema_migrations_data (version, name, applied_at)
		VALUES ($1, $2, NOW())
	`, fm.schemaName)

	if _, err := fm.db.Exec(query, version, fileName); err != nil {
		return fmt.Errorf("failed to record fixture: %w", err)
	}

	log.Info().
		Str("fixture", fileName).
		Msg("Fixture applied successfully")

	return nil
}

// extractVersion extracts version number from filename
// Example: "005_missions_data.sql" -> "005_missions_data"
func (fm *FixturesManager) extractVersion(filePath string) string {
	fileName := filepath.Base(filePath)
	// Remove .sql extension
	version := strings.TrimSuffix(fileName, ".sql")
	return version
}
