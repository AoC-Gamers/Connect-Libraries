# API Key

**Módulo:** `github.com/AoC-Gamers/connect-libraries/apikey`

## 📋 Descripción

Sistema de autenticación por API Key para comunicación interna segura entre microservicios del sistema Connect. Proporciona validación automática de claves de servicio con soporte para múltiples entornos (desarrollo/producción) y logging integrado.

## 📦 Contenido

- **apikey.go** - Tipos y estructuras principales
- **validator.go** - Validador de API Keys
- **env.go** - Helpers para variables de entorno
- **config_helper.go** - Utilidades de configuración
- **init.go** - Inicialización automática

## 🔧 Uso

```go
import "github.com/AoC-Gamers/connect-libraries/apikey"

// Configuración manual
config := apikey.Config{
    APIKeys: map[string]string{
        "auth":  "secret-auth-key",
        "core":  "secret-core-key",
        "lobby": "secret-lobby-key",
        "rt":    "secret-rt-key",
    },
}

// Validar API Key
isValid := apikey.Validate(apiKey, "auth")

// Carga automática desde variables de entorno
// CONNECT_AUTH_API_KEY=xxx
// CONNECT_CORE_API_KEY=xxx
apikey.InitFromEnv()

```

## 🌍 Variables de Entorno

Configura estas variables en tu `.env`:

```bash
CONNECT_AUTH_API_KEY=secret-auth-key
CONNECT_CORE_API_KEY=secret-core-key
CONNECT_LOBBY_API_KEY=secret-lobby-key
CONNECT_RT_API_KEY=secret-rt-key
```

## ⚙️ Dependencias

- `errors` - Sistema de manejo de errores estandarizado
- `zerolog` - Logging estructurado

## ⚡ Características

- ✅ Validación de API Keys para servicios internos
- ✅ Carga automática desde variables de entorno
- ✅ Soporte para múltiples servicios Connect
- ✅ Logging integrado con zerolog
- ✅ Modo desarrollo con auto-generación de claves
- ✅ Validación estricta en producción
