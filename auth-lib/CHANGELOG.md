# Changelog - auth-lib

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.2.1] - 2026-01-25

### Changed
- Actualización de dependencias (go.sum)

## [1.2.0] - 2026-01-24

### Added
- Permisos de Roles para WEB, COMMUNITY y TEAM
  - `WEB__ROLES_VIEW`, `WEB__ROLES_EDIT`
  - `COMMUNITY__ROLES_VIEW`, `COMMUNITY__ROLES_EDIT`
  - `TEAM__ROLES_VIEW`, `TEAM__ROLES_EDIT`
- Nuevas constantes de string: `PermWebRolesView`, `PermWebRolesEdit`, `PermCommunityRolesView`, `PermCommunityRolesEdit`, `PermTeamRolesView`, `PermTeamRolesEdit`
- Inclusión de permisos de roles en los grupos STAFF/OWNER correspondientes

## [1.1.0] - 2026-01-11

### Added
- **Permisos de Auditoría**: Nuevos permisos para visualizar logs de auditoría
  - `COMMUNITY__AUDIT_VIEW` (bit 13): Permite ver logs de auditoría de comunidades
  - `TEAM__AUDIT_VIEW` (bit 7): Permite ver logs de auditoría de equipos
- Permisos incluidos automáticamente en roles OWNER respectivos
  - `COMMUNITY__OWNER` ahora incluye `COMMUNITY__AUDIT_VIEW`
  - `TEAM__OWNER` ahora incluye `TEAM__AUDIT_VIEW`

### Technical
- Nuevas constantes de string: `PermCommunityAuditView`, `PermTeamAuditView`
- Compatibilidad mantenida con versiones anteriores (backward compatible)

## [1.0.0] - 2026-01-07

### Added
- Sistema de autenticación JWT completo
- Parsing y validación de claims JWT
- Gestión de permisos y contexto de autenticación
- Tipos de claims personalizados (Connect, Auth, Admin)
- Context helpers para extracción de información de autenticación
- Parser JWT con validación de firma y expiración
- Definiciones de políticas de autorización

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
