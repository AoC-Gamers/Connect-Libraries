# NATS

**Módulo:** `github.com/AoC-Gamers/connect-libraries/nats`

## 📋 Descripción

Cliente NATS/JetStream estandarizado para comunicación asíncrona entre microservicios Connect. Proporciona configuración unificada, autenticación múltiple y reconexión automática.

## 📦 Contenido

- **config.go** - Configuración de conexión NATS
- **jetstream.go** - Helpers para JetStream
- **publisher.go** - Publisher de eventos estandarizado

## 🔧 Uso

```go
import connectnats "github.com/AoC-Gamers/connect-libraries/nats"

// Conectar con configuración por defecto
conn, err := connectnats.Connect(nil)
if err != nil {
    log.Warn().Err(err).Msg("NATS unavailable")
    return
}
defer conn.Close()

// Configuración personalizada
cfg := &connectnats.Config{
    URL:           "nats://nats-shared:4222",
    ClientID:      "connect-auth-1",
    ReconnectWait: 2,
    MaxReconnects: -1,
    Token:         "secret-token",
}
conn, err := connectnats.Connect(cfg)
```

## ⚙️ Dependencias

- `nats.go` - Cliente NATS oficial
- `zerolog` - Logging estructurado

## ⚡ Características

- ✅ Configuración unificada para todos los servicios
- ✅ Autenticación múltiple (token, credenciales, user/pass)
- ✅ TLS opcional para producción
- ✅ Reconexión automática con backoff
- ✅ Non-blocking (servicios continúan si NATS falla)
- ✅ JetStream support integrado
- ✅ Publisher de eventos estandarizado
- ✅ Logging de conexiones y errores

    Timeout:       10 * time.Second,
}

conn, err := connectnats.Connect(cfg)
```

## Variables de Entorno

La biblioteca lee automáticamente estas variables:

### Conexión Básica
- `NATS_URL` - URL del servidor (default: `nats://localhost:4222`)
- `NATS_CLIENT_ID` - Identificador del cliente (default: `connect-service`)

### Autenticación (elige una)
- `NATS_TOKEN` - Token de autenticación (recomendado para desarrollo)
- `NATS_CREDS_FILE` - Ruta al archivo .creds (recomendado para producción)
- `NATS_USER` + `NATS_PASSWORD` - Usuario y contraseña

### TLS (opcional)
- `NATS_TLS_CERT` - Certificado del cliente
- `NATS_TLS_KEY` - Llave privada del cliente
- `NATS_TLS_CA` - Certificado de la autoridad certificadora

## Ejemplo con Token (Desarrollo)

```bash
export NATS_URL="nats://nats-shared:4222"
export NATS_TOKEN="my-secret-token"
export NATS_CLIENT_ID="connect-auth-dev"
```

```go
// La biblioteca detecta automáticamente el token
conn := connectnats.MustConnect(nil) // Retorna nil si falla
if conn != nil {
    defer conn.Close()
    // Usar conexión...
}
```

## Integración en Servicios

### Connect-Auth

```go
import connectnats "github.com/AoC-Gamers/Connect-Backend/libraries/connect-nats"

func initNATSPublisher(cfg *config.Config) *natsx.Publisher {
    natsCfg := &connectnats.Config{
        URL:           cfg.NATS.URL,
        ClientID:      cfg.NATS.ClientID,
        ReconnectWait: cfg.NATS.ReconnectWait,
        MaxReconnects: cfg.NATS.MaxReconnects,
    }
    
    conn := connectnats.MustConnect(natsCfg)
    if conn == nil {
        return nil // Servicio continúa sin eventos
    }
    
    return natsx.NewFromConn(conn)
}
```

### Connect-Core, Connect-RT, Connect-Lobby

Mismo patrón - reemplazar lógica duplicada con `connectnats.Connect()`.

## Configuración Docker

En `docker-compose.shared.yml`:

```yaml
nats-shared:
  image: nats:2.10-alpine
  volumes:
    - ./nats/nats.conf.tpl:/etc/nats/nats.conf.tpl:ro
  environment:
    - NATS_TOKEN=${NATS_TOKEN}
```

En servicios que consumen NATS:

```yaml
connect-auth:
  environment:
    - NATS_URL=nats://nats-shared:4222
    - NATS_TOKEN=${NATS_TOKEN}
    - NATS_CLIENT_ID=connect-auth-dev
```

## Manejo de Errores

La biblioteca NO termina el proceso si NATS falla:

```go
// MustConnect retorna nil si falla (non-fatal)
conn := connectnats.MustConnect(cfg)
if conn == nil {
    // Servicio continúa sin eventos en tiempo real
    log.Warn().Msg("Running without NATS")
    return
}

// O manejar explícitamente
conn, err := connectnats.Connect(cfg)
if err != nil {
    // Decidir qué hacer según el servicio
    return fmt.Errorf("NATS required: %w", err)
}
```

## Beneficios

1. **DRY** - No duplicar código de conexión en cada servicio
2. **Consistencia** - Mismo comportamiento en Auth, Core, RT, Lobby
3. **Seguridad** - Configuración TLS centralizada
4. **Mantenibilidad** - Un solo lugar para actualizar lógica NATS
5. **Testing** - Más fácil mockear conexiones NATS

## Roadmap

- [ ] Soporte para NATS JetStream
- [ ] Métricas de conexión (Prometheus)
- [ ] Pool de conexiones para alto tráfico
- [ ] Configuración desde archivo YAML/JSON
