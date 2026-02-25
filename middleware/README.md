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

- `authjwt` (interno) - Parsing y validación de JWT
- `chi` - Framework Chi router

## ⚡ Características

- ✅ Multi-framework (Chi, con soporte futuro para Gin/net-http)
- ✅ Autenticación JWT usando authjwt interno
- ✅ Validación de roles y permisos granular
- ✅ Protección de APIs internas con API keys
- ✅ Context injection consistente
- ✅ Token extraction automática (headers, cookies)
- ✅ Manejo de errores estandarizado

## 🧩 Respuestas de error personalizadas

Puedes inyectar un `ErrorResponder` para desacoplarte de cualquier librería de errores:

```go
type MyResponder struct{}

func (MyResponder) Unauthorized(w http.ResponseWriter, detail string) {
	// usar tu librería de errores aquí
}

func (MyResponder) TokenExpired(w http.ResponseWriter) {}
func (MyResponder) PolicyVersionMismatch(w http.ResponseWriter, tokenVersion, currentVersion int) {}
func (MyResponder) InsufficientPermissions(w http.ResponseWriter, action string) {}

// Uso
r.Use(chimw.RequireAuthWithResponder(cfg, MyResponder{}))
r.Use(chimw.RequireRoleWithResponder(MyResponder{}, "admin"))
r.Use(chimw.RequirePermissionBitmaskWithResponder(perm, MyResponder{}))
```