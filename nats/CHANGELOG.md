# Changelog - nats

Todos los cambios notables en esta librería serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.0.2] - 2026-02-25

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

### Security
- Endurecimiento de lectura del certificado CA TLS en `config.go` usando `os.OpenRoot` sobre el directorio declarado.
- Incorporación de helper seguro para lectura de archivos (`readFileFromDeclaredDir`) con sanitización de ruta.
- Resolución de hallazgos `gosec` en NATS (`G304` y `G703`) manteniendo compatibilidad del flujo de conexión TLS.

## [1.0.1] - 2026-01-27

### Changed
- Actualización de dependencias (go.sum)

## [1.0.0] - 2026-01-07

### Added
- Cliente NATS con soporte para JetStream
- Gestión automática de conexiones
- Publisher para publicación de eventos
- Configuración mediante variables de entorno
- Helpers para streams y consumers
- Manejo de reconexiones automáticas

### Technical
- Módulo Go independiente con versionado propio
- Sin dependencias de replace directives
- Compatible con descarga directa desde GitHub
