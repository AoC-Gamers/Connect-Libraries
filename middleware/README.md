# Connect Middleware

Middlewares de autenticación y autorización reutilizables para el ecosistema Connect Backend.

## Propósito

Proporcionar middlewares estándar para autenticación JWT, validación de roles, API keys y CORS que funcionen con diferentes frameworks web (Gin, Chi, net/http).

## Características

- ✅ **Multi-framework:** Gin, Chi, net/http
- ✅ **JWT Authentication:** Middleware estándar usando connect-auth-lib
- ✅ **Role-based Authorization:** Validación de roles y permisos
- ✅ **API Key Protection:** Para endpoints internos
- ✅ **CORS Management:** Configuración flexible de CORS
- ✅ **Context Injection:** Inyección consistente de claims en contexto

## Estructura

```
connect-middleware/
├── go.mod
├── README.md
├── gin/
│   ├── auth.go       # JWT auth middleware para Gin
│   ├── roles.go      # Role validation para Gin
│   └── apikey.go     # API key middleware para Gin
├── chi/
│   ├── auth.go       # JWT auth middleware para Chi
│   ├── permissions.go # Permission validation para Chi
│   └── apikey.go     # API key middleware para Chi
├── http/
│   ├── auth.go       # JWT auth middleware para net/http
│   ├── cors.go       # CORS middleware para net/http
│   └── apikey.go     # API key middleware para net/http
└── common/
    ├── extractor.go  # Token extractors comunes
    └── validators.go # Validadores comunes
```

## Uso por Framework

### Gin (Connect-Core)
```go
import "github.com/AoC-Gamers/Connect-Backend/connect-middleware/gin"

router.Use(ginmw.RequireAuth(authConfig))
router.Use(ginmw.RequireAdmin())
```

### Chi (Connect-Auth)
```go
import "github.com/AoC-Gamers/Connect-Backend/connect-middleware/chi"

r.Use(chimw.RequireAuth(cfg))
r.Use(chimw.RequirePermission("WEB__ADMIN"))
```

### net/http (Connect-Lobby, Connect-RT)
```go
import "github.com/AoC-Gamers/Connect-Backend/connect-middleware/http"

handler = httpmw.JWTAuth(config)(handler)
handler = httpmw.RequireRoles("admin")(handler)
```