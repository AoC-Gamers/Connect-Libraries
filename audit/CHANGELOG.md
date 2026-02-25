# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.4] - 2026-02-25

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

## [1.0.0] - 2026-01-11

### Added
- Initial release of audit library
- **Core package**: Shared functionality for all audit types
  - `Filters` struct with validation and default values
  - SQL query building helpers (`ApplyFilters`, `ApplyPagination`)
  - Common validators (`ValidateEntryTime`, `ValidateAction`)
  - SQL constants to avoid duplication
  - Common error definitions
- **Community entity**: Audit support for communities
  - 12 action constants (CREATED, UPDATED, DELETED, etc.)
  - SQL query builders (SELECT, INSERT, COUNT, DELETE)
  - Payload formatting helpers
  - Action validation
- **Team entity**: Audit support for teams
  - 7 action constants (CREATED, UPDATED, DELETED, etc.)
  - SQL query builders
  - Payload formatting helpers
  - Action validation
- **Web entity**: Audit support for web/system operations
  - 13 action constants (LOGIN, LOGOUT, CONFIG_UPDATED, etc.)
  - SQL query builders with optional scope_id support
  - Payload formatting helpers
  - Action validation
  - `ApplyWebFilters` for web-specific filtering

### Features
- ✅ Eliminates ~600 lines of duplicated code across backends
- ✅ Centralized action constants (type-safe)
- ✅ Consistent validation and filtering logic
- ✅ Reusable query builders
- ✅ Extensible architecture for future entities (lobby, rt, etc.)
- ✅ Comprehensive README with usage examples

### Architecture
- Modular structure: `core/` + `entities/{community,team,web}/`
- Each entity has: `actions.go`, `queries.go`, `helpers.go`
- Easy to extend with new audit types
- Backwards compatible with existing code through aliases

[1.0.0]: https://github.com/AoC-Gamers/Connect-Libraries/releases/tag/audit/v1.0.0
