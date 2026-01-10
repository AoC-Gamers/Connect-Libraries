# Detección Híbrida de Parámetros

La biblioteca swagger ahora soporta detección híbrida de parámetros, combinando:
1. **Detección automática de endpoints** desde el router Chi
2. **Detección automática de parámetros** usando reflection en structs de Go

## Funciones disponibles

### 1. BuildSchemaFromStruct
Construye un schema JSON desde un struct usando reflection.

```go
type UserResponse struct {
    ID        string    `json:"id" description:"User ID"`
    Name      string    `json:"name" description:"User's full name"`
    Email     string    `json:"email" description:"Email address"`
    CreatedAt time.Time `json:"created_at" description:"Account creation timestamp"`
    Premium   bool      `json:"premium,omitempty" description:"Premium status"`
}

schema := schema.BuildSchemaFromStruct(UserResponse{})
// Genera:
// {
//   "type": "object",
//   "properties": {
//     "id": {"type": "string", "description": "User ID"},
//     "name": {"type": "string", "description": "User's full name"},
//     "email": {"type": "string", "description": "Email address"},
//     "created_at": {"type": "string", "format": "date-time", "description": "Account creation timestamp"},
//     "premium": {"type": "boolean", "description": "Premium status"}
//   },
//   "required": ["id", "name", "email", "created_at"]
// }
```

### 2. ExtractParamsFromStruct
Extrae parámetros desde un struct para query params, headers, etc.

```go
type UserListQuery struct {
    Page     int    `json:"page" description:"Page number" default:"1"`
    PageSize int    `json:"page_size" description:"Items per page" default:"20"`
    Search   string `json:"search,omitempty" description:"Search term"`
    SortBy   string `json:"sort_by,omitempty" description:"Sort field" example:"created_at"`
}

params := schema.ExtractParamsFromStruct(UserListQuery{}, schema.ParamInQuery)
// Genera parámetros de query con tipos, descripciones, defaults y ejemplos
```

### 3. ExtractPathParamsFromRoute
Detecta automáticamente path params desde la ruta.

```go
params := schema.ExtractPathParamsFromRoute("/users/{userId}/posts/{postId}")
// Genera:
// [
//   {Name: "userId", In: "path", Type: "string", Required: true, Description: "User identifier"},
//   {Name: "postId", In: "path", Type: "string", Required: true, Description: "Post identifier"}
// ]
```

## Ejemplo completo de uso híbrido

```go
package main

import (
    "github.com/AoC-Gamers/connect-libraries/swagger"
    "github.com/AoC-Gamers/connect-libraries/swagger/schema"
    "github.com/go-chi/chi/v5"
)

// Definir structs para documentación automática
type CreateUserRequest struct {
    Name     string `json:"name" description:"User's full name" example:"John Doe"`
    Email    string `json:"email" description:"Email address" example:"john@example.com"`
    Password string `json:"password" description:"User password (min 8 chars)"`
}

type UserResponse struct {
    ID        string    `json:"id" description:"User ID"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type UserListQuery struct {
    Page   int    `json:"page" default:"1"`
    Limit  int    `json:"limit" default:"20"`
    Search string `json:"search,omitempty"`
}

func main() {
    r := chi.NewRouter()
    
    // Definir rutas normalmente
    r.Post("/users", createUserHandler)
    r.Get("/users", listUsersHandler)
    r.Get("/users/{userId}", getUserHandler)
    
    // Configurar swagger con detección automática
    cfg := swagger.DefaultConfig().
        WithServiceInfo("Users API", "1.0.0").
        WithDescription("API for user management").
        AddTagRule("/users", "Users")
    
    detector := swagger.New(cfg)
    detector.ScanRouter(r) // Detecta endpoints automáticamente
    
    // Registrar schemas automáticamente desde structs
    registry := detector.GetSchemaRegistry()
    
    // Opción 1: Registrar manualmente para endpoints específicos
    // (esto sobrescribe la detección automática)
    
    // Query params para GET /users
    queryParams := schema.ExtractParamsFromStruct(UserListQuery{}, schema.ParamInQuery)
    // Registrarlos manualmente si se desea...
    
    // Request body para POST /users
    requestSchema := schema.BuildSchemaFromStruct(CreateUserRequest{})
    // Registrarlo manualmente si se desea...
    
    // Response body
    responseSchema := schema.BuildSchemaFromStruct(UserResponse{})
    // Registrarlo manualmente si se desea...
    
    // Generar especificación OpenAPI
    spec, _ := swagger.ExportSpec(detector)
    
    // Servir documentación Swagger
    r.Get("/swagger.json", swagger.ServeSwaggerSpec(detector))
}
```

## Tags soportados en structs

- `json:"name"` - Nombre del campo en JSON (requerido)
- `json:"name,omitempty"` - Campo opcional
- `description:"..."` - Descripción del campo
- `example:"..."` - Valor de ejemplo
- `default:"..."` - Valor por defecto
- `binding:"required"` - Campo requerido (alternativa a omitir omitempty)

## Detección automática de tipos

La biblioteca detecta automáticamente:

- **Tipos básicos**: string, int, float, bool
- **Formatos**: int32, int64, float, double, byte
- **time.Time**: Convertido a string con formato "date-time"
- **Arrays/Slices**: Detecta el tipo de elementos
- **Structs anidados**: Genera schemas de objeto

## Inferencia de tipos por nombre

Para path params, la biblioteca infiere tipos basándose en el nombre:

- `*Id`, `*ID` → string (para soportar UUIDs)
- `page`, `limit`, `count`, `size` → integer
- Otros → string (por defecto)
