# Changelog - authz

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

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
