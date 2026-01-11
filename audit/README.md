# Audit Library

Biblioteca compartida para gestión de auditoría en los servicios de Connect.

## Estructura

```
audit/
├── core/                      # Funcionalidad base compartida
│   ├── filters.go             # Tipos: Filters struct
│   ├── query_builder.go       # Helpers: ApplyFilters, ApplyPagination
│   ├── validators.go          # ValidateEntry, ValidateFilters
│   └── constants.go           # Constantes SQL y errores comunes
│
├── entities/                  # Definiciones por tipo de entidad
│   ├── community/             # Auditoría de comunidades
│   │   ├── actions.go         # Constantes de acciones
│   │   ├── queries.go         # Queries SQL
│   │   └── helpers.go         # Formateo de payloads
│   │
│   ├── team/                  # Auditoría de equipos
│   │   ├── actions.go
│   │   ├── queries.go
│   │   └── helpers.go
│   │
│   └── web/                   # Auditoría web/sistema
│       ├── actions.go
│       ├── queries.go
│       └── helpers.go
```

## Uso

### Connect-Core (Community + Team)

```go
import (
    auditcore "github.com/AoC-Gamers/Connect-Libraries/audit/core"
    auditcommunity "github.com/AoC-Gamers/Connect-Libraries/audit/entities/community"
)

func (r *communityAuditRepo) GetAuditEntries(ctx context.Context, filters auditcore.Filters) ([]models.CommunityAuditEntry, error) {
    query := auditcommunity.BuildSelectQuery()
    args := []interface{}{filters.ScopeID}
    
    // Aplicar filtros y paginación
    query, args = filters.ApplyFilters(query, args)
    query, args = filters.ApplyPagination(query, args)
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    // ...
}

func (r *communityAuditRepo) CreateAuditEntry(ctx context.Context, entry *models.CommunityAuditEntry) error {
    if entry == nil {
        return auditcore.ErrEntryNil
    }
    
    auditcore.EnsureTimestamp(&entry.CreatedAt)
    
    if err := auditcommunity.ValidateAction(entry.Action); err != nil {
        return err
    }
    
    query := auditcommunity.BuildInsertQuery()
    _, err := r.db.ExecContext(ctx, query, entry.ScopeID, entry.Action, entry.PerformedBy, entry.Payload, entry.CreatedAt)
    return err
}
```

### Connect-Auth (Community + Team + Web)

```go
import (
    auditcore "github.com/AoC-Gamers/Connect-Libraries/audit/core"
    auditweb "github.com/AoC-Gamers/Connect-Libraries/audit/entities/web"
)

func (r *webAuditRepo) GetAuditEntries(ctx context.Context, filters auditcore.Filters) ([]models.WebAuditEntry, error) {
    query := auditweb.BuildSelectQuery()
    args := []interface{}{}
    
    // Web audit maneja scope_id opcional
    query, args = auditweb.ApplyWebFilters(&filters, query, args)
    query, args = filters.ApplyPagination(query, args)
    
    rows, err := r.db.QueryContext(ctx, query, args...)
    // ...
}
```

### Formateo de Payloads

```go
// Community
payload := auditcommunity.FormatCreatedPayload("My Community", "ACTIVE", true)
payload := auditcommunity.FormatStatusChangePayload("ACTIVE", "SUSPENDED")
payload := auditcommunity.FormatServerPayload(123, "EU Server #1")

// Team
payload := auditteam.FormatCreatedPayload("Team Alpha", "ALPHA", "ACTIVE")
payload := auditteam.FormatOwnerTransferPayload("76561198008295809")

// Web
payload := auditweb.FormatLoginPayload("192.168.1.1", "Mozilla/5.0...")
payload := auditweb.FormatPermissionPayload("COMMUNITY__MANAGE", "COMMUNITY", 5)
```

## Extensibilidad

Para añadir un nuevo tipo de auditoría (e.g., Lobby):

1. Crear directorio `entities/lobby/`
2. Implementar `actions.go`, `queries.go`, `helpers.go`
3. Usar en Connect-Lobby importando `audit/entities/lobby`

## Ventajas

- ✅ Sin duplicación de código entre backends
- ✅ Constantes de acciones centralizadas
- ✅ Validaciones consistentes
- ✅ Fácil extensión para nuevas entidades
- ✅ Queries optimizados y reutilizables
- ✅ Tests compartidos

## Versionado

Versión actual: **v1.0.0**
