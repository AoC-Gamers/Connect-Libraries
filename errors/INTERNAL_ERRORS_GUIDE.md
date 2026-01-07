# Sistema de Errores Internos - Connect Backend

## Descripción

Este documento describe el nuevo sistema de errores internos implementado en la librería `connect-errors` para mejorar la comunicación entre servicios y facilitar el debugging.

## Características Principales

### 1. Códigos de Error Estructurados
```go
type InternalErrorCode string

const (
    // Authentication & Authorization
    InternalUnauthorized InternalErrorCode = "INTERNAL_UNAUTHORIZED"
    InternalForbidden    InternalErrorCode = "INTERNAL_FORBIDDEN"
    
    // Resource Management
    InternalNotFound     InternalErrorCode = "INTERNAL_NOT_FOUND"
    InternalConflict     InternalErrorCode = "INTERNAL_CONFLICT"
    
    // Validation & Input
    InternalValidation   InternalErrorCode = "INTERNAL_VALIDATION"
    InternalBadRequest   InternalErrorCode = "INTERNAL_BAD_REQUEST"
    
    // System & Infrastructure
    InternalDatabase     InternalErrorCode = "INTERNAL_DATABASE"
    InternalTimeout      InternalErrorCode = "INTERNAL_TIMEOUT"
    InternalServiceDown  InternalErrorCode = "INTERNAL_SERVICE_DOWN"
    InternalRateLimit    InternalErrorCode = "INTERNAL_RATE_LIMIT"
    
    // Generic
    InternalServerError  InternalErrorCode = "INTERNAL_SERVER_ERROR"
)
```

### 2. Respuesta Estructurada
```go
type InternalErrorResponse struct {
    Code    InternalErrorCode `json:"code"`
    Message string           `json:"message"`
    Service string           `json:"service"`
    Details interface{}      `json:"details,omitempty"`
    Status  int              `json:"status"`
}
```

### 3. Logging Automático
- Todas las respuestas de error generan logs estructurados con zerolog
- Información incluida: código interno, servicio, path, método, status, detalles
- Facilita el debugging y monitoreo

## Funciones Principales

### Función Base
```go
RespondInternalServiceError(c *gin.Context, code InternalErrorCode, message string, details interface{})
```

### Helpers Específicos

#### Errores de Autenticación
```go
// API key inválida o faltante (401)
RespondInternalUnauthorized(c *gin.Context, message string)

// Servicio no autorizado para endpoint (403)
RespondInternalForbidden(c *gin.Context, allowedServices []string, actualService string)
```

#### Errores de Recursos
```go
// Recurso no encontrado (404)
RespondInternalNotFound(c *gin.Context, resource string, id string)
```

#### Errores de Validación
```go
// Error de validación de campo (400)
RespondInternalValidation(c *gin.Context, field string, reason string)

// Request malformado (400)
RespondInternalBadRequest(c *gin.Context, reason string)
```

#### Errores de Sistema
```go
// Error de base de datos (500)
RespondInternalDatabase(c *gin.Context, operation string, err error)

// Timeout en operación (500)
RespondInternalTimeout(c *gin.Context, operation string, timeout string)

// Servicio dependiente no disponible (503)
RespondInternalServiceDown(c *gin.Context, service string)
```

## Ejemplos de Uso

### 1. En Handlers de API

```go
func GetUserHandler(c *gin.Context) {
    userID := c.Param("id")
    
    // Validar parámetros
    if userID == "" {
        errors.RespondInternalValidation(c, "user_id", "parameter is required")
        return
    }
    
    // Operación de base de datos
    user, err := getUserFromDB(userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            errors.RespondInternalNotFound(c, "user", userID)
            return
        }
        errors.RespondInternalDatabase(c, "get_user", err)
        return
    }
    
    c.JSON(200, user)
}
```

### 2. En Middleware de Autenticación

```go
func RequireAdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            errors.RespondInternalUnauthorized(c, "user not authenticated")
            c.Abort()
            return
        }
        
        isAdmin, err := checkUserIsAdmin(userID.(string))
        if err != nil {
            errors.RespondInternalDatabase(c, "check_admin_status", err)
            c.Abort()
            return
        }
        
        if !isAdmin {
            details := map[string]interface{}{
                "user_id": userID,
                "required_role": "admin",
            }
            errors.RespondInternalServiceError(c, errors.InternalForbidden, "admin access required", details)
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 3. En Llamadas a Servicios Externos

```go
func GetUserProfileHandler(c *gin.Context) {
    userID := c.Param("id")
    
    // Llamar a servicio de auth
    authData, err := callAuthService(userID)
    if err != nil {
        if isNetworkError(err) {
            errors.RespondInternalServiceDown(c, "connect-auth")
            return
        }
        if isAuthError(err) {
            errors.RespondInternalUnauthorized(c, "user session invalid")
            return
        }
        errors.RespondInternalServiceError(c, errors.InternalServerError, "auth service error", map[string]interface{}{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(200, authData)
}
```

## Integración con connect-apikey

El middleware de `connect-apikey` ha sido actualizado para usar el nuevo sistema:

```go
// Antes
c.JSON(http.StatusUnauthorized, gin.H{
    "error": "invalid or missing API key",
})

// Ahora
errors.RespondInternalUnauthorized(c, "")
```

```go
// Antes
c.JSON(http.StatusForbidden, gin.H{
    "error": "service not authorized for this endpoint",
})

// Ahora
errors.RespondInternalForbidden(c, allowedServices, serviceName)
```

## Ejemplo de Respuesta

### Unauthorized
```json
{
    "code": "INTERNAL_UNAUTHORIZED",
    "message": "invalid or missing API key",
    "service": "connect-core",
    "status": 401
}
```

### Forbidden con Detalles
```json
{
    "code": "INTERNAL_FORBIDDEN",
    "message": "service not authorized for this endpoint",
    "service": "connect-core", 
    "details": {
        "allowed_services": ["connect-auth", "connect-lobby"],
        "actual_service": "connect-rt"
    },
    "status": 403
}
```

### Not Found con Detalles
```json
{
    "code": "INTERNAL_NOT_FOUND",
    "message": "resource not found",
    "service": "connect-core",
    "details": {
        "resource": "user",
        "id": "123"
    },
    "status": 404
}
```

### Database Error
```json
{
    "code": "INTERNAL_DATABASE", 
    "message": "database operation failed",
    "service": "connect-auth",
    "details": {
        "operation": "create_user",
        "error": "duplicate key value violates unique constraint"
    },
    "status": 500
}
```

## Ventajas del Sistema

### 1. **Debugging Mejorado**
- Códigos de error específicos y consistentes
- Información contextual en el campo `details`
- Logs estructurados automáticos

### 2. **Comunicación Clara entre Servicios**
- Identificación automática del servicio que genera el error
- Respuestas consistentes en todo el ecosystem
- Información suficiente para debugging

### 3. **Detección de Servicios**
- Detección automática del servicio por headers, contexto o path
- Fallback a "connect-service" genérico

### 4. **Compatibilidad**
- No reemplaza el sistema existente de errores públicos
- Específico para comunicación interna entre servicios
- Fácil migración gradual

## Migración Gradual

1. **Fase 1**: Actualizar middleware de autenticación (✅ Completado)
2. **Fase 2**: Actualizar adaptadores internos en Connect-Core
3. **Fase 3**: Actualizar handlers internos en otros servicios
4. **Fase 4**: Implementar en nuevos endpoints

## Beneficios Esperados

- **Mejor debugging**: Información contextual clara en errores
- **Monitoreo mejorado**: Logs estructurados para análisis
- **Desarrollo más rápido**: Errores más informativos facilitan desarrollo
- **Consistencia**: Respuestas uniformes en todo el sistema
- **Mantenimiento**: Más fácil identificar y resolver problemas