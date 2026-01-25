# Changelog - core-types

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.2.0] - 2026-01-25

### Added
- Campo `team_memberships` en `UserAllMembershipsResponse`

## [1.1.0] - 2026-01-25

### Added
- Endpoint `AuthUpdateMembership` para actualizar membresías
- Modelos `UpdateMembershipRequest` y `UpdateMembershipResponse` para actualización de membresías

## [1.0.0] - 2026-01-07

### Added
- Tipos compartidos para endpoints HTTP
- Modelos de dominio compartidos entre servicios
- Definiciones de errores comunes
- DTOs para comunicación entre servicios
- Tipos de datos para entidades core

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
