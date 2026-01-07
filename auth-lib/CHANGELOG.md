# Changelog - auth-lib

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

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
