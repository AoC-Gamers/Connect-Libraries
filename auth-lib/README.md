# Connect Auth Lib

Biblioteca de autenticación reutilizable para el ecosistema Connect Backend.

## Propósito

Centralizar la lógica común de autenticación JWT, validación de tokens y manejo de contexto que es compartida entre todos los microservicios Connect (Auth, Core, Lobby, RT).

## Características

- ✅ Parsing y validación de JWT tokens
- ✅ Extracción segura de claims (steamid, roles, permissions)
- ✅ Verificación de policy version
- ✅ Context helpers framework-agnósticos
- ✅ Tipos y estructuras comunes
- ✅ Zero dependencies (solo jwt-go)

## Uso

```go
import "github.com/AoC-Gamers/Connect-Backend/connect-auth-lib/jwt"

// Validar token
claims, err := jwt.ParseAndValidate(tokenStr, secret, policyVersion)
if err != nil {
    return err
}

// Extraer información
steamID := claims.GetSteamID()
roles := claims.GetRoles()
```

## Estructura

```
connect-auth-lib/
├── go.mod
├── README.md
├── jwt/
│   ├── parser.go     # JWT parsing y validación
│   ├── claims.go     # Claims helpers
│   └── types.go      # Tipos comunes
├── context/
│   ├── keys.go       # Context keys estándar
│   └── helpers.go    # Context extraction helpers
└── config/
    └── auth.go       # Configuración común
```

## Compatibilidad

- ✅ Gin (Connect-Core)
- ✅ Chi (Connect-Auth)  
- ✅ net/http (Connect-Lobby, Connect-RT)
- ✅ Go 1.21+