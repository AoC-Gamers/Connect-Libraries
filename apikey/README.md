# Connect API Key

Sistema de autenticación por API Key para comunicación interna entre microservicios Connect.

## Características

- ✅ **Multi-framework:** Gin, Chi, net/http
- ✅ **Environment Integration:** Carga automática desde variables de entorno (.env)
- ✅ **Connect Services:** Mapeo automático de servicios Connect (Auth, Core, Lobby, RT)  
- ✅ **Development Mode:** Auto-generación de claves para desarrollo
- ✅ **Production Ready:** Validación estricta y logging seguro
- ✅ **Granular Control:** Middleware por servicio específico
- ✅ **Observability:** Logs y debugging integrado

## Estructura

```
connect-apikey/
├── go.mod
├── README.md
├── validator.go         # Validador principal
├── env.go              # Helpers para variables de entorno
├── init.go             # Inicialización automática
├── gin/
│   ├── middleware.go   # Middleware básico para Gin
│   └── connect.go      # Helpers específicos Connect para Gin
├── chi/
│   └── middleware.go   # Middleware para Chi
└── http/
    └── middleware.go   # Middleware para net/http
```

## Uso Básico

### 🚀 Integración Automática (Recomendado)

```go
import ginapi "github.com/AoC-Gamers/Connect-Backend/libraries/connect-apikey/gin"

// Carga automática desde variables de entorno (.env)
router.Use(ginapi.RequireConnectAPIKey())

// Middleware específico por servicio
router.Use(ginapi.RequireAuthService())     // Solo Connect-Auth
router.Use(ginapi.RequireLobbyService())    // Solo Connect-Lobby
```

### 🔧 Configuración Manual

```go
import apikey "github.com/AoC-Gamers/Connect-Backend/libraries/connect-apikey"

// Configuración manual
validator := apikey.NewValidator(map[string]string{
    "connect-core-key-123": "connect-core",
    "connect-lobby-key-456": "connect-lobby",
})

// Uso con Gin
import ginapi "github.com/AoC-Gamers/Connect-Backend/libraries/connect-apikey/gin"
router.Use(ginapi.RequireAPIKey(validator))
```

### 🌍 Variables de Entorno

Configura estas variables en tu `.env`:

```bash
AUTH_API_KEY=connect-auth-internal-key
CORE_API_KEY=connect-core-internal-key  
LOBBY_API_KEY=connect-lobby-internal-key
RT_API_KEY=connect-rt-internal-key
```

## Formatos Soportados

- `X-API-Key: <key>` (header)
- `Authorization: Bearer <key>` (header)
- `api_key=<key>` (query parameter)