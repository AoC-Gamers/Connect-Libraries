# Service Clients

**Módulo:** `github.com/AoC-Gamers/connect-libraries/service-clients`

## 📋 Descripción

Biblioteca de clientes HTTP type-safe para la comunicación entre los microservicios del sistema Connect. Proporciona una interfaz limpia y tipada para realizar llamadas entre servicios con manejo automático de errores, serialización JSON y logging.

## 📦 Contenido

### `clients/auth/`
Cliente HTTP para el servicio **Connect-Auth**:
- Gestión de usuarios y permisos
- Validación de tokens
- Operaciones de autenticación

### `clients/core/`
Cliente HTTP para el servicio **Connect-Core**:
- Operaciones de comunidades y equipos
- Gestión de jugadores y estadísticas
- Consultas del dominio principal

### `clients/rt/`
Cliente HTTP para el servicio **Connect-RT** (Real-Time):
- Eventos en tiempo real
- Notificaciones
- Websockets y streams

## 🔧 Uso

```go
import (
    authclient "github.com/AoC-Gamers/connect-libraries/service-clients/clients/auth"
    coreclient "github.com/AoC-Gamers/connect-libraries/service-clients/clients/core"
    rtclient "github.com/AoC-Gamers/connect-libraries/service-clients/clients/rt"
)

// Ejemplo: Cliente para Connect-Core
coreClient := coreclient.NewClient("http://core-service:8080", nil)
user, err := coreClient.GetUserByID(ctx, userID)
```

## ⚙️ Dependencias

- `core-types` - Para modelos de datos y endpoints
- `errors` - Para manejo de errores estandarizado

## ⚡ Características

- ✅ Type-safe con validación automática
- ✅ Manejo centralizado de errores HTTP
- ✅ Logging integrado con zerolog
- ✅ Context-aware para timeouts y cancelación
- ✅ Retry y circuit breaker ready
