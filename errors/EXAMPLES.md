# Connect Errors - Ejemplo de Uso

Este documento muestra ejemplos prácticos de cómo usar `connect-errors` en cada microservicio.

## Connect-Auth - Permisos y Memberships

```go
package authz

import (
    "net/http"
    errors "github.com/AoC-Gamers/connect-libraries/errors"
)

func (h *AuthzHandler) GetMyCapabilities(w http.ResponseWriter, r *http.Request) {
    // ... código de autenticación ...
    
    // Validar parámetros
    scopeIDStr := r.URL.Query().Get("scopeId")
    if scopeIDStr == "" {
        errors.RespondMissingField(w, "scopeId")
        return
    }
    
    scopeID, err := strconv.ParseInt(scopeIDStr, 10, 64)
    if err != nil {
        errors.RespondInvalidFormat(w, "scopeId", "integer")
        return
    }
    
    // Obtener capabilities
    masks, err := h.authz.GetEffectiveMasks(ctx, userID, scopeID)
    if err != nil {
        // Verificar si no tiene membership
        if errors.IsNotFound(err) || err.Error() == "membership not found" {
            errors.RespondMembershipNotFound(w, scopeID, "WEB", userID)
            return
        }
        errors.RespondInternalError(w, "failed to get effective masks")
        return
    }
    
    // Respuesta exitosa
    errors.RespondJSON(w, http.StatusOK, map[string]interface{}{
        "effective": masks,
        "policy_version": h.config.Authz.PolicyVersionGlobal,
    })
}
```

## Connect-Core - CRUD de Misiones

```go
package mission

import (
    "net/http"
    errors "github.com/AoC-Gamers/connect-libraries/errors"
)

func (h *MissionHandler) CreateMission(w http.ResponseWriter, r *http.Request) {
    var req CreateMissionRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        errors.RespondBadRequest(w, "invalid JSON body")
        return
    }
    
    // Validación de campos
    validationErrors := make(map[string]string)
    
    if req.Name == "" {
        validationErrors["name"] = "required"
    } else if len(req.Name) > 128 {
        errors.RespondOutOfRange(w, "name", 1, 128, len(req.Name))
        return
    }
    
    if req.URL != "" && !isValidURL(req.URL) {
        validationErrors["url"] = "must be a valid URL"
    }
    
    if len(req.Maps) == 0 {
        validationErrors["maps"] = "at least one map is required"
    }
    
    if len(validationErrors) > 0 {
        errors.RespondValidationErrors(w, validationErrors)
        return
    }
    
    // Verificar si ya existe
    existing, _ := h.service.GetByName(req.Name)
    if existing != nil {
        errors.RespondAlreadyExists(w, "mission", req.Name)
        return
    }
    
    // Crear misión
    mission, err := h.service.Create(ctx, req)
    if err != nil {
        errors.RespondDatabaseError(w, "create mission")
        return
    }
    
    // Respuesta exitosa
    errors.RespondJSON(w, http.StatusCreated, mission)
}

func (h *MissionHandler) GetMission(w http.ResponseWriter, r *http.Request) {
    name := chi.URLParam(r, "name")
    if name == "" {
        errors.RespondMissingField(w, "name")
        return
    }
    
    mission, err := h.service.GetByName(name)
    if err != nil {
        if errors.IsNotFound(err) {
            errors.RespondNotFound(w, "mission", name)
        } else {
            errors.RespondInternalError(w, "failed to retrieve mission")
        }
        return
    }
    
    errors.RespondJSON(w, http.StatusOK, mission)
}

func (h *MissionHandler) SuspendMission(w http.ResponseWriter, r *http.Request) {
    name := chi.URLParam(r, "name")
    
    // Verificar permisos
    if !hasPermission(ctx, "WEB__MISSION_SUSPEND") {
        errors.RespondPermissionDenied(w, 1, "WEB", "WEB__MISSION_SUSPEND", true)
        return
    }
    
    // Verificar que la misión existe y no es oficial
    mission, err := h.service.GetByName(name)
    if err != nil {
        errors.RespondNotFound(w, "mission", name)
        return
    }
    
    if mission.Official {
        errors.RespondOperationNotAllowed(w, "suspend mission", "cannot suspend official missions")
        return
    }
    
    // Suspender
    if err := h.service.Suspend(ctx, name, reason); err != nil {
        errors.RespondDatabaseError(w, "suspend mission")
        return
    }
    
    errors.RespondJSON(w, http.StatusOK, map[string]string{
        "message": "mission suspended successfully",
    })
}
```

## Connect-Lobby - Rate Limiting

```go
package lobby

import (
    "net/http"
    errors "github.com/AoC-Gamers/connect-libraries/errors"
)

func (h *LobbyHandler) CreateLobby(w http.ResponseWriter, r *http.Request) {
    // Check rate limit
    userID := getUserIDFromContext(r.Context())
    if !h.rateLimiter.Allow(userID) {
        errors.RespondRateLimitExceeded(w, 5, "5 minutes", 60)
        return
    }
    
    // Check user quota
    count, _ := h.service.CountUserLobbies(userID)
    if count >= 10 {
        errors.RespondQuotaExceeded(w, "active lobbies", 10, count)
        return
    }
    
    // ... crear lobby ...
}
```

## Connect-RT - WebSocket Errors

```go
package realtime

import (
    "net/http"
    errors "github.com/AoC-Gamers/connect-libraries/errors"
)

func (h *RTHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    // Verificar token
    token := r.URL.Query().Get("token")
    if token == "" {
        errors.RespondUnauthorized(w, "token is required for WebSocket connection")
        return
    }
    
    claims, err := validateToken(token)
    if err != nil {
        if errors.IsTokenExpired(err) {
            errors.RespondTokenExpired(w)
        } else {
            errors.RespondTokenInvalid(w, "malformed token")
        }
        return
    }
    
    // Verificar servicio disponible
    if !h.isHealthy() {
        errors.RespondServiceUnavailable(w, "realtime-service")
        return
    }
    
    // ... upgrade a WebSocket ...
}
```

## Migración Gradual

Si tienes código existente con `httpresp.RespondError`, puedes migrar gradualmente:

```go
// Antes
httpresp.RespondError(w, http.StatusForbidden, "insufficient permissions")

// Después (con metadata)
errors.RespondPermissionDenied(w, scopeID, scopeType, requiredPermission, hasMembership)

// O mantener compatibilidad durante la transición
errors.RespondLegacyError(w, http.StatusForbidden, "insufficient permissions")
```

## Helpers para Detección de Errores

```go
package errors

import (
    "errors"
    "database/sql"
)

// IsNotFound verifica si un error es "not found"
func IsNotFound(err error) bool {
    return errors.Is(err, sql.ErrNoRows) || 
           err.Error() == "not found" ||
           err.Error() == "membership not found"
}

// IsTokenExpired verifica si un error es token expirado
func IsTokenExpired(err error) bool {
    return err.Error() == "token expired" ||
           err.Error() == "jwt expired"
}

// IsDatabaseError verifica si es error de BD
func IsDatabaseError(err error) bool {
    return errors.Is(err, sql.ErrConnDone) ||
           errors.Is(err, sql.ErrTxDone)
}
```
