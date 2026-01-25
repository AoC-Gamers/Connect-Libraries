# Changelog - service-clients

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.1.1] - 2026-01-25

### Changed
- Actualizada dependencia a `core-types` v1.1.0

## [1.1.0] - 2026-01-25

### Added
- Método `UpdateMembership` en el cliente de Auth para actualizar rol y/o permisos

## [1.0.0] - 2026-01-07

### Added
- Clientes HTTP tipados para comunicación entre microservicios
- Cliente para Connect-Auth (autenticación y usuarios)
- Cliente para Connect-Core (comunidades, gamemodes, etc.)
- Cliente para Connect-RT (tiempo real y websockets)
- Manejo automático de autenticación entre servicios
- Retry automático con backoff exponencial
- Timeouts configurables

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
