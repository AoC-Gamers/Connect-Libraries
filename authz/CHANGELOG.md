# Changelog - authz

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [2.0.2] - 2026-02-25

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

## [2.0.1] - 2026-02-07

### Fixed
- Module path actualizado a `/v2` para Semantic Import Versioning (Go modules)

## [2.0.0] - 2026-02-06

### Changed
- BREAKING: Permisos renombrados a CamelCase en authz/permissions (sin aliases), manteniendo las claves ALL_CAPS
- Normalización de MissionList/GamemodeList en TEAM y COMMUNITY
- Helpers comunes para roles/permisos y orden estable en GetAll*PermissionNames

## [1.0.1] - 2026-01-27

### Added
- Módulo de permisos (authz/permissions) con bitmasks, grupos, roles y helpers
- Constantes WEB/COMMUNITY/TEAM/LOBBY y utilidades de autorización

## [1.0.0] - 2026-01-07

### Added
- Sistema de autorización basado en roles (RBAC)
- Definiciones de roles predefinidos (Admin, Moderator, User, etc.)
- Sistema de permisos granular
- Validación de permisos para operaciones específicas
- Mapeo de roles a permisos

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives  
- Compatible con descarga directa desde GitHub
