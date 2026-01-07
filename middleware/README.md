# Middleware

**Módulo:** `github.com/AoC-Gamers/connect-libraries/middleware`

## 📋 Descripción

Middlewares HTTP reutilizables para autenticación, autorización y protección de APIs en todos los microservicios Connect. Proporciona middlewares estándar para JWT, validación de roles/permisos, API keys y CORS con soporte multi-framework.

## 📦 Contenido

### `chi/`
Middlewares para framework Chi (usado en Connect-Auth):
- **auth.go** - Autenticación JWT
- **permissions.go** - Validación de permisos
- **apikey.go** - Protección con API keys

## 🔧 Uso

### Con Chi (Connect-Auth)

```go
import "github.com/AoC-Gamers/connect-libraries/middleware/chi"

// Autenticación JWT
r.Use(chimw.RequireAuth(cfg))

// Validación de permisos
r.Use(chimw.RequirePermission("WEB__ADMIN"))

// Protección con API key
r.Use(chimw.RequireAPIKey(apiKeyValidator))
```

## ⚙️ Dependencias

- `auth-lib` - Para parsing y validación de JWT
- `errors` - Para respuestas de error estandarizadas
- `chi` - Framework Chi router

## ⚡ Características

- ✅ Multi-framework (Chi, con soporte futuro para Gin/net-http)
- ✅ Autenticación JWT usando auth-lib
- ✅ Validación de roles y permisos granular
- ✅ Protección de APIs internas con API keys
- ✅ Context injection consistente
- ✅ Token extraction automática (headers, cookies)
- ✅ Manejo de errores estandarizado
import "github.com/AoC-Gamers/Connect-Backend/connect-middleware/http"

handler = httpmw.JWTAuth(config)(handler)
handler = httpmw.RequireRoles("admin")(handler)
```