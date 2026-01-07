# Auth Lib

**Módulo:** `github.com/AoC-Gamers/connect-libraries/auth-lib`

## 📋 Descripción

Biblioteca de autenticación JWT reutilizable para todos los microservicios del ecosistema Connect. Proporciona parsing, validación de tokens, extracción de claims y helpers de contexto framework-agnósticos.

## 📦 Contenido

### `jwt/`
Parsing y validación de tokens JWT:
- **parser.go** - Validación de tokens JWT
- **claims.go** - Extracción de claims (steamid, roles, permissions)
- **types.go** - Tipos y estructuras comunes

### `context/`
Helpers para manejo de contexto:
- **keys.go** - Keys estándar para context.Context

### `config/`
Configuración de autenticación:
- **policy.go** - Configuración de políticas de seguridad

### `permissions/`
Sistema de permisos y autorización

## 🔧 Uso

```go
import (
    "github.com/AoC-Gamers/connect-libraries/auth-lib/jwt"
    "github.com/AoC-Gamers/connect-libraries/auth-lib/context"
)

// Parsing y validación de JWT
claims, err := jwt.ParseAndValidate(tokenStr, jwtSecret, policyVersion)
if err != nil {
    return err
}

// Extraer información del token
steamID := claims.GetSteamID()
roles := claims.GetRoles()
permissions := claims.GetPermissions()

// Helpers de contexto
steamID := context.GetSteamIDFromContext(ctx)
```

## ⚙️ Dependencias

- `jwt/v5` - Parsing y validación de JWT tokens
- `zerolog` - Logging estructurado

## ⚡ Características

- ✅ Parsing y validación de JWT tokens
- ✅ Extracción segura de claims (steamid, roles, permissions)
- ✅ Verificación de policy version
- ✅ Context helpers framework-agnósticos
- ✅ Compatible con Gin, Chi, net/http
- ✅ Sistema de permisos granular
- ✅ Type-safe con validación automática