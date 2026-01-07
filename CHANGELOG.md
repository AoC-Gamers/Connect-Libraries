# Changelog

Todos los cambios notables en este proyecto serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.0.0] - 2026-01-07

### Added
- **apikey**: Sistema de autenticación y validación de API Keys
- **auth-lib**: Sistema de autenticación JWT, parsing de claims y gestión de permisos
- **authz**: Sistema de autorización basado en roles (RBAC) con definiciones de roles y permisos
- **core-types**: Tipos compartidos incluyendo endpoints, modelos de dominio y definiciones de errores
- **errors**: Sistema de manejo de errores estandarizado siguiendo RFC 7807
- **middleware**: Middlewares HTTP para framework Chi (autenticación, autorización, logging)
- **migrate**: Sistema de migraciones para PostgreSQL con tracking y soporte de fixtures
- **nats**: Cliente NATS/JetStream con gestión de conexiones y publicación de eventos
- **service-clients**: Clientes HTTP tipados para comunicación entre microservicios (Auth, Core, RT)
- **testhelpers**: Utilidades para testing incluyendo helpers y mocks con limpieza automática

### Changed
- Fragmentación de biblioteca `internal/` en tres módulos independientes:
  - `core-types`: Tipos compartidos (endpoints, models, errors)
  - `service-clients`: Clientes HTTP para inter-service communication
  - `authz`: Sistema de autorización y roles

### Removed
- Biblioteca `internal/` (fragmentada en módulos especializados)

---

## Formato de Versionado

Este repositorio usa **versionado semántico unificado**:

- **MAJOR** (x.0.0): Cambios incompatibles en la API
- **MINOR** (1.x.0): Nueva funcionalidad compatible con versiones anteriores
- **PATCH** (1.0.x): Correcciones de bugs compatibles con versiones anteriores

### Ejemplo de futuras versiones:

```markdown
## [1.1.0] - YYYY-MM-DD

### Added
- **auth-lib**: Soporte para refresh tokens
- **nats**: Retry automático con backoff exponencial

### Changed
- **errors**: Mejoras en mensajes de error

### Fixed
- **middleware**: Corrección en validación de permisos

## [1.0.1] - YYYY-MM-DD

### Fixed
- **apikey**: Corrección en validación de API keys vacías
- **migrate**: Fix en tracking de migraciones ejecutadas
```

[Unreleased]: https://github.com/AoC-Gamers/connect-libraries/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/AoC-Gamers/connect-libraries/releases/tag/v1.0.0
