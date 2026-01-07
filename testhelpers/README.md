# Test Helpers Library

Librería compartida de utilidades para testing reutilizable por todos los módulos del backend (Connect-Core, Connect-Auth, Connect-Lobby, Connect-RT).

## Filosofía de diseño

Esta librería usa **únicamente tipos públicos** y no depende de paquetes `internal` de ningún módulo. Esto permite que sea verdaderamente compartida y reutilizable.

## Helpers disponibles

### NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock)

Crea una base de datos mock usando sqlmock para código que usa `database/sql` directamente (sin sqlx).
El cleanup es automático vía `t.Cleanup()`.

**Ejemplo:**
```go
import th "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"

func TestMyHandler(t *testing.T) {
    db, mock := th.NewMockDB(t)
    // No need to defer db.Close() - cleanup is automatic
    
    mock.ExpectQuery("SELECT ...").WillReturnRows(...)
    
    handler := NewHandler(db)
    // ...
    
    // No need to call mock.ExpectationsWereMet() - done automatically
}
```

**Ventajas:**
- Cleanup completamente automático
- Verifica expectations automáticamente
- Reduce boilerplate en tests

### NewSQLMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock, func())

Crea una base de datos mock usando sqlmock para tests unitarios de repositorios.

**Ejemplo:**
```go
import th "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"

func TestMyRepo(t *testing.T) {
    db, mock, cleanup := th.NewSQLMock(t)
    defer cleanup()

    mock.ExpectQuery("SELECT ...").WillReturnRows(...)
    
    repo := NewRepo(db)
    result, err := repo.GetSomething(ctx)
    // ...
}
```

### NewMiniredisClient(t *testing.T) (*redis.Client, *miniredis.Miniredis, func())

Crea un servidor Redis en memoria usando miniredis y retorna un cliente estándar `go-redis`.

**Ejemplo:**
```go
import th "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"

func TestCaching(t *testing.T) {
    client, mrs, cleanup := th.NewMiniredisClient(t)
    defer cleanup()

    // Usa client como cualquier *redis.Client normal
    err := client.Set(ctx, "key", "value", 0).Err()
    // ...
}
```

**Nota para módulos con wrappers internos:** Si tu módulo tiene un wrapper interno sobre Redis (como `Connect-Core/internal/redis`), crea un adaptador en tu `internal/testhelpers` que convierta el cliente público al wrapper interno. Ver ejemplo en `Connect-Core/internal/testhelpers/helpers.go`.

### MakeAuthServer(has bool) *httptest.Server

Crea un servidor HTTP de prueba que simula el endpoint de permisos de Connect-Auth.

**Ejemplo:**
```go
import th "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"

func TestPermissions(t *testing.T) {
    // Servidor que siempre retorna hasPermission: true
    authServer := th.MakeAuthServer(true)
    defer authServer.Close()

    authClient := clients.NewAuthClient(authServer.URL, "")
    // ...
}
```

## Uso en diferentes módulos

### Opción 1: Uso directo (recomendado para nuevos tests)

Importa directamente desde `libraries/testhelpers`:

```go
import th "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"
```

### Opción 2: Wrapper interno (para módulos con tipos internos)

Si tu módulo necesita adaptar tipos públicos a wrappers internos, crea un `internal/testhelpers` que actúe como adaptador:

```go
// internal/testhelpers/helpers.go
package testhelpers

import (
    lib "github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers"
    "github.com/tu-modulo/internal/redis"
)

func NewMiniredisClient(t *testing.T) (*redis.Client, func()) {
    _, mrs, cleanup := lib.NewMiniredisClient(t)
    
    // Adapta el cliente público al wrapper interno
    cfg := &config.RedisConfig{...}
    internalClient, err := redis.NewClient(cfg)
    // ...
    
    return internalClient, cleanup
}
```

## Configuración en go.mod

Agrega la librería como dependencia:

```go.mod
require (
    github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers v0.0.0-00010101000000-000000000000
)

replace github.com/AoC-Gamers/Connect-Backend/libraries/testhelpers => ../libraries/testhelpers
```

## Ventajas de esta arquitectura

1. **Reutilización:** Una sola implementación para todos los backends
2. **Sin violación de reglas Go:** No importa paquetes `internal` de otros módulos
3. **Flexibilidad:** Cada módulo puede adaptar helpers a sus necesidades específicas
4. **Mantenimiento:** Cambios en helpers se propagan automáticamente a todos los módulos
5. **Testing más rápido:** Menos código duplicado = menos bugs en tests

## Contribuir

Al agregar nuevos helpers:
- Solo usa tipos públicos (stdlib, go-redis, sqlx, etc.)
- NO importes paquetes `internal` de ningún módulo
- Documenta el uso con ejemplos
- Asegúrate de que la función cleanup libere todos los recursos
