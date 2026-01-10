# Changelog - swagger

Todos los cambios notables en esta biblioteca serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

## [1.1.0] - 2026-01-10

### Added
- **Sistema híbrido de documentación de parámetros**: Auto-detección + registro manual
- Auto-detección de path parameters desde rutas Chi (`{id}`, `{steamid}`, etc.)
- Registro manual de query parameters mediante reflection de structs
- Registro manual de request body schemas con soporte completo de validaciones
- Registro manual de response schemas por código HTTP
- Nuevo módulo `schema.go` con `SchemaRegistry` para gestión de schemas
- Nuevo módulo `reflection.go` con helpers de conversión Go → OpenAPI
- Inferencia inteligente de tipos para path parameters (id → string, page → integer)
- Extracción automática de validaciones desde tags `binding` (min, max, email, etc.)
- Soporte para tags `json`, `binding`, y `description` en structs
- Generación de JSON Schema OpenAPI 3.0 completo desde structs Go
- Método `GetSchemaRegistry()` en `Detector` para acceso al sistema de schemas
- Documentación extendida en README.md con ejemplos de uso de schemas
- Nuevo archivo `EXAMPLE_OUTPUT.md` con ejemplos de JSON OpenAPI generado

### Changed
- `Detector` extendido con campo `schemas *SchemaRegistry`
- Método `generatePaths()` reescrito para incluir parámetros y schemas registrados
- Respuestas por defecto ahora incluyen schemas cuando están registrados
- README.md actualizado con sección completa sobre registro de schemas

### Technical
- Reflection avanzado para construir schemas desde tipos Go
- Conversión automática de tipos Go a tipos OpenAPI (string, integer, boolean, etc.)
- Soporte para tipos complejos (arrays, objetos anidados, time.Time)
- Generación automática de descripciones para path parameters
- 100% backward compatible: schemas son opcionales, funciona sin cambios si no se usan
- Implementado y probado en Connect-RT con 12+ endpoints documentados

### Migration Guide
```go
// v1.0.0 - Sin parámetros
detector := swaggerlib.NewDetector(cfg)

// v1.1.0 - Con parámetros (opcional)
detector := swaggerlib.NewDetector(cfg)
registry := detector.GetSchemaRegistry()

// Registrar request body
registry.RegisterRequestBody("/users", "POST",
    models.CreateUserRequest{}, "User data", true)

// Registrar response
registry.RegisterResponse("/users", "POST", 201,
    models.User{}, "User created")

// Path params se detectan automáticamente desde "/users/{id}"
```

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
