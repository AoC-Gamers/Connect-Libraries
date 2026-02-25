# Changelog - migrate

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.0.2] - 2026-02-25

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

### Security
- Endurecimiento de lectura de archivos SQL usando `os.OpenRoot` para acotar el acceso a `migrations_sql` y `data_sql`.
- Validación estricta de identificadores SQL (schema) antes de construir DDL.
- Refactor de consultas de tracking para evitar SQL dinámico innecesario (uso de `search_path` y queries constantes).
- Resolución de hallazgos `gosec` en `migrator.go` y `fixtures.go` (`G201`, `G202`, `G304`, `G701`) manteniendo el comportamiento del migrador.

## [1.0.1] - 2026-01-27

### Changed
- Actualización de dependencias (go.sum)

## [1.0.0] - 2026-01-07

### Added
- Sistema de migraciones para PostgreSQL
- Tracking de migraciones aplicadas
- Soporte para archivos SQL de migración
- Sistema de fixtures para datos de prueba
- Rollback de migraciones
- Logger integrado para seguimiento de operaciones

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
