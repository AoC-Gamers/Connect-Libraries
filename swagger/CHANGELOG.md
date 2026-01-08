# Changelog - swagger

Todos los cambios notables en esta biblioteca serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.0.0] - 2026-01-08

### Added
- Sistema de detección automática de endpoints desde Chi Router
- Configuración flexible mediante `Config` y métodos fluent
- Detección automática de seguridad (JWT, API Keys) desde middlewares
- Sistema de tags configurable con reglas por patrón
- Soporte para rutas públicas y rutas a omitir
- Exportación JSON de rutas detectadas
- Handler HTTP para exponer rutas (`/swagger/routes`)
- Documentación completa en README.md con ejemplos por servicio
- Convenciones de nombres para tags por servicio

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias internas de Connect (solo Chi Router)
- Compatible con descarga directa desde GitHub
- Diseñado para reemplazar detectores duplicados en Auth/Core/Lobby/RT
