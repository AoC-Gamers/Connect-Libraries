# Changelog - errors

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.0.3] - 2026-02-02

### Changed
- Actualización de dependencias (go.sum)

## [1.0.2] - 2026-01-27

### Changed
- Actualización de dependencias (go.sum)

## [1.0.1] - 2026-01-25

### Changed
- Actualización de dependencias (go.sum)

## [1.0.0] - 2026-01-07

### Added
- Sistema de manejo de errores estandarizado siguiendo RFC 7807
- Códigos de error consistentes entre servicios
- Helpers para creación y manejo de errores
- Errores tipados para casos comunes (validación, autenticación, autorización, etc.)
- Serialización JSON de errores
- Contexto y detalles adicionales en errores
- Sistema de errores internos separado

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
