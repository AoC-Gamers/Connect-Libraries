# Connect Errors Library

**Standardized Error Responses for Connect Backend Services**

Biblioteca compartida que implementa respuestas de error estructuradas siguiendo RFC 7807 (Problem Details for HTTP APIs) para todos los microservicios del ecosistema Connect Backend.

## 🎯 Características

- ✅ Respuestas de error estructuradas y consistentes
- ✅ Códigos de error estandarizados  
- ✅ Metadata extensible para debugging
- ✅ Compatible con RFC 7807
- ✅ Helpers para casos de uso comunes
- ✅ **Sistema de errores internos para comunicación entre servicios**
- ✅ **Logging automático con zerolog**
- ✅ **Detección automática de servicios**

## 📦 Instalación

```go
import "github.com/AoC-Gamers/Connect-Backend/libraries/connect-errors"
```

## 🚀 Uso Básico

### Respuestas de Error para APIs Públicas

```go
import (
    "net/http"
    errors "github.com/AoC-Gamers/Connect-Backend/libraries/connect-errors"
)

func handler(w http.ResponseWriter, r *http.Request) {
    errors.RespondError(w, http.StatusBadRequest, errors.CodeValidationError, 
        "invalid mission name", 
        "Mission name exceeds 128 characters",
        map[string]interface{}{
            "max_length": 128,
            "provided_length": 150,
        },
    )
}
```

### Errores Internos entre Servicios (NUEVO)

```go
import (
    "github.com/gin-gonic/gin"
    errors "github.com/AoC-Gamers/Connect-Backend/libraries/connect-errors"
)

func internalHandler(c *gin.Context) {
    userID := c.Param("id")
    
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
```

### Helpers Predefinidos

```go
// Permission Denied (APIs públicas)
errors.RespondPermissionDenied(w, scopeID, "WEB", "WEB__MISSION_VIEW", false)

// Service Forbidden (APIs internas)
errors.RespondInternalForbidden(c, []string{"connect-auth"}, "connect-rt")

// Not Found
errors.RespondNotFound(w, "mission", missionName)

// Membership Not Found
errors.RespondMembershipNotFound(w, scopeID, "WEB", userID)

// Token Expired
errors.RespondTokenExpired(w)

// Policy Version Mismatch
errors.RespondPolicyVersionMismatch(w, tokenVersion, currentVersion)

// Internal Server Error (APIs públicas)
errors.RespondInternalError(w, "database connection failed")
```

## 📁 Estructura del Proyecto

- `codes.go` - Códigos de error estandarizados
- `errors.go` - Funciones base y estructura ErrorResponse (RFC 7807)
- `helpers.go` - Helpers para APIs públicas y casos de uso comunes
- `internal.go` - **Sistema de errores internos para comunicación entre servicios**
- `INTERNAL_ERRORS_GUIDE.md` - **Guía completa del sistema de errores internos**

## 📋 Estructura de Respuesta

Todas las respuestas de error siguen este formato:

```json
{
  "error": "short error message",
  "code": "ERROR_CODE",
  "status": 403,
  "detail": "Detailed explanation of what went wrong",
  "meta": {
    "scope_id": 1,
    "required_permission": "WEB__MISSION_VIEW"
  }
}
```

## 🏷️ Códigos de Error Disponibles

### Authentication & Authorization
- `UNAUTHORIZED`
- `TOKEN_EXPIRED`
- `TOKEN_INVALID`
- `POLICY_VERSION_MISMATCH`
- `PERMISSION_DENIED`
- `INSUFFICIENT_PERMISSIONS`

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
import errors "github.com/AoC-Gamers/Connect-Backend/libraries/connect-errors"

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
import errors "github.com/AoC-Gamers/Connect-Backend/libraries/connect-errors"

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
