# Changelog - testhelpers

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

## [1.0.1] - 2026-01-27

### Changed
- Actualización de dependencias (go.sum)

## [1.0.0] - 2026-01-07

### Added
- Utilidades para testing y mocks
- Helpers para setup y teardown de tests
- Mocks de servicios externos
- Generadores de datos de prueba
- Assertions personalizadas
- Helpers para testing de HTTP handlers
- Limpieza automática de recursos en tests

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
