# Permission System - Integration Guide

## Overview

El sistema de permisos está centralizado en la librería `connect-auth-lib/permissions` y usa **bitmasks (uint64)** para operaciones bitwise eficientes y **constantes string** para verificaciones de API.

## Estructura de Permisos

### Scopes (Ámbitos)

Cada scope tiene su propio conjunto de permisos:

- **WEB**: Permisos de la plataforma web (administración global)
- **LOBBY**: Permisos dentro de un lobby específico
- **COMMUNITY**: Permisos dentro de una comunidad
- **TEAM**: Permisos dentro de un equipo

### Archivos por Scope

```
connect-auth-lib/permissions/
├── web.go        # Permisos de plataforma web
├── lobby.go      # Permisos de lobby
├── community.go  # Permisos de comunidad
├── team.go       # Permisos de equipo
└── authz.go      # Helpers de autorización compartidos
```

### Dos Tipos de Constantes

Cada archivo de permisos define **dos tipos de constantes** para el mismo permiso:

1. **Bitmask Constants (uint64)**: Para operaciones bitwise locales
   - Formato: `SCOPE__PERMISSION_NAME`
   - Ejemplo: `COMMUNITY__SERVER_ADD uint64 = 1 << 2`

2. **String Constants**: Para verificaciones de API (CheckUserPermission)
   - Formato: `PermScopePermissionName`
   - Ejemplo: `PermCommunityServerAdd = "COMMUNITY__SERVER_ADD"`

3. **Role String Constants**: Para identificadores de roles
   - Formato: `RoleScopeRoleNameKey`
   - Ejemplo: `RoleCommunityOwnerKey = "community_owner"``

## Usando la Librería

### 1. Importar en tu Servicio

```go
import (
    authperms "github.com/AoC-Gamers/Connect-Backend/libraries/connect-auth-lib/permissions"
)
```

### 2. Verificar Permisos con API (Recomendado)

```go
// ✅ USAR CONSTANTES STRING para CheckUserPermission
hasPermission, err := s.authClient.CheckUserPermission(
    ctx,
    userID,
    authperms.PermCommunityServerAdd,  // ← Constante string
    "COMMUNITY",
    fmt.Sprintf("%d", communityID),
)
```

### 3. Verificar Permisos Localmente (Bitwise)

```go
// Verificar un permiso simple con bitmask
if authperms.HasPermission(userAllowWeb, authperms.WEB__MISSION_ADD) {
    // Usuario puede agregar misiones
}

// Verificar con deny mask
if authperms.CanPerformAction(allowWeb, denyWeb, authperms.WEB__MISSION_ADD) {
    // Usuario puede agregar misiones (considerando deny)
}
```

### 4. Verificar Múltiples Permisos

```go
// Usuario debe tener TODOS estos permisos
if permissions.HasAllPermissions(
    userAllowLobby,
    permissions.LOBBY__KICK,
    permissions.LOBBY__MANAGE,
) {
    // Usuario puede kickear Y gestionar lobby
}

// Usuario debe tener AL MENOS UNO de estos permisos
if permissions.HasAnyPermission(
    userAllowCommunity,
    permissions.COMMUNITY__SERVER_ADD,
    permissions.COMMUNITY__SERVER_EDIT,
) {
    // Usuario puede agregar O editar servidores
}
```

### 4. Obtener Permisos de un Rol

```go
// ✅ USAR CONSTANTES DE ROLES (Recomendado)

// Web roles
webOwnerPerms := permissions.GetRolePermissions(authperms.RoleWebOwnerKey)   // "web_owner"
webStaffPerms := permissions.GetRolePermissions(authperms.RoleWebStaffKey)   // "web_staff"
webUserPerms := permissions.GetRolePermissions(authperms.RoleWebUserKey)     // "web_user"

// Lobby roles
lobbyOwnerPerms := permissions.GetLobbyRolePermissions(authperms.RoleLobbyOwnerKey) // "lobby_owner"
lobbyStaffPerms := permissions.GetLobbyRolePermissions(authperms.RoleLobbyStaffKey) // "lobby_staff"
lobbyUserPerms := permissions.GetLobbyRolePermissions(authperms.RoleLobbyUserKey)   // "lobby_user"

// Community roles
communityOwnerPerms := permissions.GetCommunityRolePermissions(authperms.RoleCommunityOwnerKey) // "community_owner"
communityStaffPerms := permissions.GetCommunityRolePermissions(authperms.RoleCommunityStaffKey) // "community_staff"
communityUserPerms := permissions.GetCommunityRolePermissions(authperms.RoleCommunityUserKey)   // "community_user"

// Team roles
teamOwnerPerms := permissions.GetTeamRolePermissions(authperms.RoleTeamOwnerKey) // "team_owner"
teamStaffPerms := permissions.GetTeamRolePermissions(authperms.RoleTeamStaffKey) // "team_staff"
teamUserPerms := permissions.GetTeamRolePermissions(authperms.RoleTeamUserKey)   // "team_user"
```

**Ventajas de usar constantes de roles**:
- ✅ Type safety (el compilador detecta typos)
- ✅ Autocomplete en el IDE
- ✅ Refactoring seguro
- ✅ Single source of truth

## Integración en Servicios

### Connect-RT (Real-Time WebSocket)

**Caso de uso**: Autorizar acciones en tiempo real (kick, lobby management, chat)

```go
// internal/handlers/lobby_handler.go
package handlers

import (
    "github.com/AoC-Gamers/Connect-Backend/libraries/connect-auth-lib/permissions"
)

func (h *LobbyHandler) HandleKickPlayer(userID string, targetID string, lobbyID int64) error {
    // Obtener membership del usuario en el lobby
    membership, err := h.repo.GetLobbyMembership(userID, lobbyID)
    if err != nil {
        return err
    }

    // Verificar permiso de KICK
    if !permissions.CanPerformAction(
        membership.AllowLobby,
        membership.DenyLobby,
        permissions.LOBBY__KICK,
    ) {
        return errors.New("insufficient permissions to kick players")
    }

    // Ejecutar acción
    return h.kickPlayer(targetID, lobbyID)
}
```

### Connect-Core (HTTP REST API)

**Caso de uso**: Autorizar endpoints CRUD (missions, gamemodes, communities)

```go
// internal/api/missions/middleware.go
package missions

import (
    "github.com/AoC-Gamers/Connect-Backend/libraries/connect-auth-lib/permissions"
    "github.com/gin-gonic/gin"
)

// RequireWebPermission middleware para endpoints REST
func RequireWebPermission(requiredPerm uint64) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obtener JWT claims del contexto (inyectado por AuthMiddleware)
        allowWeb := c.GetUint64("allow_web")
        denyWeb := c.GetUint64("deny_web")

        // Verificar permiso
        if !permissions.CanPerformAction(allowWeb, denyWeb, requiredPerm) {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Uso en rutas
func SetupRoutes(r *gin.Engine) {
    missions := r.Group("/missions")
    {
        // Solo lectura - básico
        missions.GET("/", RequireWebPermission(permissions.WEB__MISSION_VIEW), ListMissions)
        
        // Crear - requiere permiso staff
        missions.POST("/", RequireWebPermission(permissions.WEB__MISSION_ADD), CreateMission)
        
        // Editar - requiere permiso staff
        missions.PUT("/:id", RequireWebPermission(permissions.WEB__MISSION_EDIT), UpdateMission)
        
        // Eliminar - requiere permiso owner
        missions.DELETE("/:id", RequireWebPermission(permissions.WEB__MISSION_DELETE), DeleteMission)
    }
}
```

## Roles y Permisos Pre-definidos

### Web Roles

| Role | Constante Bitmask | Constante String | Permisos Incluidos |
|------|-------------------|------------------|-------------------|
| web_user | `RoleWebUser` | `RoleWebUserKey` | Solo visualización y lobby básico |
| web_staff | `RoleWebStaff` | `RoleWebStaffKey` | web_user + gestión de contenido |
| web_owner | `RoleWebOwner` | `RoleWebOwnerKey` | Todos los permisos de plataforma |

### Lobby Roles

| Role | Constante Bitmask | Constante String | Permisos Incluidos |
|------|-------------------|------------------|-------------------|
| lobby_user | `RoleLobbyUser` | `RoleLobbyUserKey` | Ver y participar |
| lobby_staff | `RoleLobbyStaff` | `RoleLobbyStaffKey` | Gestión moderada |
| lobby_owner | `RoleLobbyOwner` | `RoleLobbyOwnerKey` | Control total del lobby |

### Community Roles

| Role | Constante Bitmask | Constante String | Permisos Incluidos |
|------|-------------------|------------------|-------------------|
| community_user | `RoleCommunityUser` | `RoleCommunityUserKey` | Miembro básico |
| community_staff | `RoleCommunityStaff` | `RoleCommunityStaffKey` | Gestión de servidores/misiones |
| community_owner | `RoleCommunityOwner` | `RoleCommunityOwnerKey` | Control total de comunidad |

### Team Roles

| Role | Constante Bitmask | Constante String | Permisos Incluidos |
|------|-------------------|------------------|-------------------|
| team_user | `RoleTeamUser` | `RoleTeamUserKey` | Miembro del equipo |
| team_staff | `RoleTeamStaff` | `RoleTeamStaffKey` | Gestión de membresías |
| team_owner | `RoleTeamOwner` | `RoleTeamOwnerKey` | Control total del equipo |

## Debugging de Permisos

```go
// Ver nombres de permisos en una máscara
func debugPermissions(mask uint64, scopeName string) {
    var names []string
    
    switch scopeName {
    case "web":
        names = permissions.GetAllPermissionNames(mask)
    case "lobby":
        for perm, name := range permissions.LobbyPermissionNames {
            if mask&perm != 0 {
                names = append(names, name)
            }
        }
    // Similar para community y team
    }
    
    log.Info().
        Uint64("mask", mask).
        Strs("permissions", names).
        Str("scope", scopeName).
        Msg("User permissions")
}
```

## Mejores Prácticas

1. **Siempre usa `CanPerformAction`** para considerar deny masks
2. **Usa constantes para permisos Y roles** - nunca strings hardcoded
3. **Loggea decisiones de autorización** para debugging
4. **Valida permisos en AMBOS lados**: middleware Y lógica de negocio
5. **Usa grupos de permisos** (`WEB__STAFF`, `LOBBY__OWNER`) en lugar de permisos individuales cuando asignes roles
6. **Importa con alias**: `import authperms "github.com/.../connect-auth-lib/permissions"`