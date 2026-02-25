# Changelog - apikey

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

## [1.1.0] - 2026-01-27

### Changed
- Desacople de `errors` con `ErrorResponder` inyectable

## [1.0.2] - 2026-01-25

### Changed
- Header por defecto de API Key actualizado a X-Internal-API-Key.

## [1.0.1] - 2026-01-25

### Changed
- Actualización de dependencias (go.sum)

## [1.0.0] - 2026-01-07

### Added
- Sistema de autenticación y validación de API Keys
- Configuración mediante variables de entorno
- Validador de API keys con soporte para múltiples claves
- Helper para configuración centralizada
- Inicialización automática del sistema de validación

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
