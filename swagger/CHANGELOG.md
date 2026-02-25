# Changelog - swagger

Todos los cambios notables en esta biblioteca serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/es-ES/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/lang/es/).

## [Unreleased]

### Changed
- Actualización del workflow de CI para ejecución y validación por librería.
- Estandarización de la secuencia de calidad `lint`, `test` y `gosec`.
- Generación de reportes locales por módulo en `reports/` con salida de `gosec` legible (`gosec.log` sin códigos ANSI).

### Security
- Ajuste de construcción de `SecurityPatterns` para evitar falso positivo de credenciales hardcodeadas (`G101`) sin cambiar el comportamiento de detección.
- Manejo explícito del error de escritura en `ServeSwaggerSpec` (`G104`) para robustecer la respuesta HTTP.
- Validación de `gosec` en módulo `swagger` con resultado final `Issues: 0`.

## [1.3.1] - 2026-01-27

### Changed
- Actualización de dependencias (go.sum)

## [1.3.0] - 2026-01-10

### Added
- **BuildSchemaFromStruct()** - Genera schemas OpenAPI completos automáticamente desde structs Go
- **ExtractParamsFromStruct()** - Extrae parámetros con detección automática de tipos y formatos
- **buildStructProperties()** - Helper para procesar campos de structs en propiedades OpenAPI
- **buildFieldSchema()** - Helper para construir schema de campos individuales
- Detección automática de tipos: string, integer (int32/int64), number (float/double), boolean, array, object
- Detección de formatos especializados: int32, int64, float, double, date-time para `time.Time`
- Soporte completo para arrays y slices con detección de tipos de elementos
- Soporte para objetos anidados con recursión completa
- Nuevos tags soportados: `description`, `example`, `default`
- Sistema híbrido: detección automática + configuración manual opcional
- Documentación completa en `HYBRID_DETECTION.md`

### Changed
- `reflection.go`: Expandido de 235 líneas a 370+ líneas con funciones de detección híbrida
- Función `goTypeToSwaggerType()` ahora incluida en detección híbrida (no eliminada)
- Función `goTypeToSwaggerFormat()` mejorada para detectar formatos especializados
- README.md actualizado con sección de Detección Híbrida y ejemplos completos

### Improved
- **Organización del código**: Simplificación y limpieza de estructura
- `swagger.go`: Reducido de 137 a 67 líneas (-70 líneas)
  - Eliminadas re-exportaciones innecesarias de tipos
  - Eliminadas constantes no utilizadas
  - Movidos métodos `ExportJSON()` y `ExportSpec()` desde detector para evitar import cycle
  - API pública simplificada: solo funciones esenciales
- `schema/generator.go` → `schema/registry.go`: Renombrado para mejor claridad
  - Eliminados métodos Register* no utilizados (RegisterQueryParams, RegisterPathParams, RegisterRequestBody, RegisterResponse)
  - Eliminado método MarshalToJSON() no utilizado
  - Mantenidos solo métodos Get* (GetRequestBody, GetResponse, GetQueryParams, GetPathParams)
  - Reducido de ~100 líneas a 55 líneas
- `detector/detector.go`: Limpieza de dependencias circulares
  - Eliminadas importaciones de encoding/json y openapi (movidas a swagger.go)
  - Eliminados métodos ExportJSON(), ExportSpec(), ServeHTTP() (movidos a swagger.go)
  - Foco exclusivo en lógica de detección de endpoints

### Fixed
- Import cycle: detector ↔ openapi resuelto moviendo métodos de exportación a swagger.go
- Claridad de nombres: registry es más descriptivo que generator
- Separación de responsabilidades: cada paquete con un propósito claro

### Technical
- Reflection avanzado usando paquete `reflect` de Go
- Análisis recursivo de tipos para estructuras anidadas
- Parsing de tags struct usando `reflect.StructTag`
- Detección de required fields basada en presencia de `omitempty`
- Zero breaking changes: 100% compatible con v1.2.0
- Arquitectura mejorada con responsabilidades bien definidas:
  - `swagger.go` → API pública y exportación
  - `detector/` → Detección de endpoints desde Chi router
  - `schema/` → Registro manual y detección híbrida de schemas
  - `openapi/` → Generación de especificación OpenAPI

### Migration Notes
```go
// v1.2.0 y anteriores siguen funcionando sin cambios
detector := swagger.New(config)
router.Get("/swagger/routes", detector.ServeHTTP)

// v1.3.0 - Nueva funcionalidad híbrida (opcional)
import "github.com/AoC-Gamers/connect-libraries/swagger/schema"

// Generar schema automático
schema := schema.BuildSchemaFromStruct(MyStruct{})

// Extraer parámetros con tipos
params := schema.ExtractParamsFromStruct(MyParams{}, schema.InQuery)
```

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
