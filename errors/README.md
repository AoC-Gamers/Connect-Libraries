# Errors

**Módulo:** `github.com/AoC-Gamers/connect-libraries/errors`

## 📋 Descripción

Sistema de manejo de errores estandarizado para todos los microservicios Connect. Implementa respuestas estructuradas siguiendo RFC 7807 (Problem Details for HTTP APIs) con soporte para errores públicos (APIs cliente) e internos (comunicación entre servicios).

## 📦 Contenido

- **errors.go** - Tipos y estructuras de error principales
- **codes.go** - Códigos de error estandarizados
- **helpers.go** - Helpers para casos de uso comunes (validación, permisos, etc.)
- **internal.go** - Sistema de errores internos para comunicación entre servicios
- **types.ts** - Definiciones TypeScript para frontend
- **EXAMPLES.md** - Ejemplos de uso completos
- **INTERNAL_ERRORS_GUIDE.md** - Guía de errores internos

## � Uso

### Respuestas de Error Públicas (APIs cliente)

```go
import (
    "net/http"
    "github.com/AoC-Gamers/connect-libraries/errors"
)

func handler(w http.ResponseWriter, r *http.Request) {
    errors.RespondError(w, http.StatusBadRequest, 
        errors.CodeValidationError, 
        "invalid mission name", 
        "Mission name exceeds 128 characters",
        map[string]interface{}{
            "max_length": 128,
            "provided_length": 150,
        },
    )
}

// Helper predefinido
errors.RespondPermissionDenied(w, scopeID, "WEB", "WEB__MISSION_VIEW", false)
```

### Errores Internos (Comunicación entre servicios)

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/AoC-Gamers/connect-libraries/errors"
)

func internalHandler(c *gin.Context) {
    // Validación
    if userID == "" {
        errors.RespondInternalValidation(c, "user_id", "parameter is required")
        return
    }
    
    // Error de base de datos
    user, err := getUserFromDB(userID)
    if err != nil {
        errors.RespondInternalDatabase(c, "get_user", err)
        return
    }
    
    c.JSON(200, user)
}
errors.RespondInternalForbidden(c, []string{"connect-auth"}, "connect-rt")

// Not Found
errors.RespondNotFound(w, "mission", missionName)

// Membership Not Found
errors.RespondMembershipNotFound(w, scopeID, "WEB", userID)

// Token Expired
errors.RespondTokenExpired(w)

```

## 📋 Códigos de Error Estandarizados

### Authentication & Authorization
- `UNAUTHORIZED` - Usuario no autenticado
- `TOKEN_EXPIRED` - Token JWT expirado
- `TOKEN_INVALID` - Token JWT inválido
- `POLICY_VERSION_MISMATCH` - Versión de política desactualizada
- `PERMISSION_DENIED` - Permisos insuficientes
- `INSUFFICIENT_PERMISSIONS` - Falta de permisos específicos

### Validation & Input
- `VALIDATION_ERROR` - Error de validación de entrada
- `INVALID_INPUT` - Entrada inválida
- `MISSING_FIELD` - Campo requerido faltante

### Resources
- `NOT_FOUND` - Recurso no encontrado
- `ALREADY_EXISTS` - Recurso ya existe
- `CONFLICT` - Conflicto de estado

### Server & Database
- `INTERNAL_ERROR` - Error interno del servidor
- `DATABASE_ERROR` - Error de base de datos
- `SERVICE_UNAVAILABLE` - Servicio no disponible

## 📁 Estructura de Respuesta (RFC 7807)

```json
{
  "error": "permission denied",
  "code": "PERMISSION_DENIED",
  "status": 403,
  "detail": "User lacks required permission for this resource",
  "meta": {
    "scope_id": 1,
    "required_permission": "WEB__MISSION_VIEW",
    "has_permission": false
  }
}
```

## ⚙️ Dependencias

- `zerolog` - Logging estructurado automático
- `gin-gonic/gin` - Soporte para framework Gin (opcional)

## ⚡ Características

- ✅ RFC 7807 compliant (Problem Details)
- ✅ Códigos de error estandarizados
- ✅ Metadata extensible para debugging
- ✅ Sistema dual: errores públicos + internos
- ✅ Logging automático con zerolog
- ✅ Detección automática de servicios
- ✅ Helpers predefinidos para casos comunes
- ✅ Compatible con Gin, Chi, net/http
- ✅ TypeScript definitions incluidas

## 📚 Documentación Adicional

- Ver [EXAMPLES.md](EXAMPLES.md) para más ejemplos
- Ver [INTERNAL_ERRORS_GUIDE.md](INTERNAL_ERRORS_GUIDE.md) para errores internos

### Resource Errors
- `NOT_FOUND`
- `MEMBERSHIP_NOT_FOUND`
- `ALREADY_EXISTS`
- `CONFLICT`

### Validation Errors
- `VALIDATION_ERROR`
- `INVALID_REQUEST`
- `BAD_REQUEST`
- `MISSING_REQUIRED_FIELD`

### Server Errors
- `INTERNAL_ERROR`
- `DATABASE_ERROR`
- `SERVICE_UNAVAILABLE`

## 🔧 Integración con Servicios

### Connect-Auth

```go
import errors "github.com/AoC-Gamers/connect-libraries/errors"

func (h *AuthzHandler) GetMyCapabilities(w http.ResponseWriter, r *http.Request) {
    // ... código existente ...
    
    if !hasMembership {
        errors.RespondMembershipNotFound(w, scopeID, scopeType, userID)
        return
    }
    
    if !hasPermission {
        errors.RespondPermissionDenied(w, scopeID, scopeType, requiredPermission, true)
        return
    }
}
```

### Connect-Core

```go
import errors "github.com/AoC-Gamers/connect-libraries/errors"

func (h *MissionHandler) GetMission(w http.ResponseWriter, r *http.Request) {
    mission, err := h.service.GetByName(name)
    if err != nil {
        if errors.IsNotFound(err) {
            errors.RespondNotFound(w, "mission", name)
        } else {
            errors.RespondInternalError(w, "failed to retrieve mission")
        }
        return
    }
    
    // ... respuesta exitosa
}
```

## 🎨 Frontend Integration

TypeScript types están incluidos para el frontend:

```typescript
// Copiar de libraries/connect-errors/types.ts a tu frontend
import type { ApiErrorResponse, ErrorCode } from '@models/api-error.model';

try {
  await api.createMission(data);
} catch (error) {
  if (error.response?.data?.code === 'PERMISSION_DENIED') {
    // Manejar error de permisos específicamente
    showPermissionDeniedDialog(error.response.data);
  }
}
```

## 🧪 Testing

```go
func TestErrorResponse(t *testing.T) {
    rr := httptest.NewRecorder()
    
    errors.RespondPermissionDenied(rr, 1, "WEB", "WEB__MISSION_VIEW", false)
    
    assert.Equal(t, http.StatusForbidden, rr.Code)
    
    var response errors.ErrorResponse
    json.Unmarshal(rr.Body.Bytes(), &response)
    
    assert.Equal(t, "PERMISSION_DENIED", string(response.Code))
    assert.Equal(t, 403, response.Status)
}
```

## 📖 Referencias

- [RFC 7807 - Problem Details for HTTP APIs](https://tools.ietf.org/html/rfc7807)
- [Google API Design Guide - Errors](https://cloud.google.com/apis/design/errors)

## 🤝 Contribuir

Para agregar nuevos códigos de error o helpers:

1. Agrega el código en `codes.go`
2. Implementa el helper en `helpers.go`
3. Actualiza esta documentación
4. Agrega tests en `errors_test.go`
