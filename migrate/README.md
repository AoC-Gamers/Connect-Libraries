# Migrate

**Módulo:** `github.com/AoC-Gamers/connect-libraries/migrate`

## 📋 Descripción

Sistema estandarizado de migraciones de base de datos para PostgreSQL utilizado por todos los microservicios Connect. Proporciona tracking de versiones, datos iniciales opcionales y verificación post-migración.

## ✅ Prerrequisitos de desarrollo

- Go `1.24.x`
- `golangci-lint` `v2.10.1`
- `gosec` `v2.23.0`

## 📦 Contenido

- **migrator.go** - Motor principal de migraciones
- **logger.go** - Configuración de logging
- **fixtures.go** - Sistema de datos iniciales (fixtures)

## 🔧 Uso

```go
import "github.com/AoC-Gamers/connect-libraries/migrate"

// Configurar logger
migrate.SetupLogger()

// Crear configuración del migrator
config := migrate.Config{
    ServiceName:    "Connect-Auth",
    SchemaName:     "auth",
    MigrationsDir:  "migrations_sql",
    DataDir:        "data_sql",
    ApplyData:      true,
    CriticalTables: []string{"users", "roles", "permissions"},
}

// Crear y ejecutar migrator
migrator, err := migrate.New(config)
if err != nil {
    log.Fatal().Err(err).Msg("Failed to create migrator")
}

if err := migrator.Run(); err != nil {
    log.Fatal().Err(err).Msg("Migration failed")
}
```

## ⚙️ Dependencias

- `pgx/v5` - Driver PostgreSQL
- `zerolog` - Logging estructurado

## ⚡ Características

- ✅ Tracking automático de versiones (schema_migrations)
- ✅ Sistema de fixtures/datos iniciales
- ✅ Verificación de tablas críticas post-migración
- ✅ Idempotente (ejecutable múltiples veces)
- ✅ Configuración vía variables de entorno
- ✅ Logging detallado de cada paso
- ✅ Soporte para múltiples schemas
    if err := migrator.Connect(); err != nil {
        log.Fatal().Err(err).Msg("Failed to connect to database")
    }
    defer migrator.Close()

    // Ejecutar migraciones
    if err := migrator.Run(); err != nil {
        log.Fatal().Err(err).Msg("Migration failed")
    }
}
```

### 2. Variables de entorno requeridas

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_NAME=connect
POSTGRES_SSLMODE=disable

# Opcionales (con defaults)
LOG_LEVEL=info           # debug, info, warn, error
LOG_FORMAT=json          # json, console, pretty
```

### 3. Estructura de directorios

```
service/
├── cmd/
│   └── migrate/
│       └── main.go              # Tu migrator personalizado
├── migrations_sql/              # Estructura de base de datos
│   ├── 001_initial.sql
│   ├── 002_add_users.sql
│   └── 003_add_indexes.sql
└── data_sql/                    # Datos iniciales (opcional)
    ├── 002_users_data.sql       # Coincide con migración 002
    └── 003_default_roles.sql    # Coincide con migración 003
```

## Sistema de Datos Iniciales

### ¿Qué son los archivos de datos?

Los archivos en `data_sql/` contienen **datos iniciales** necesarios para el funcionamiento del sistema (no estructura).

### Convenciones

1. **Numeración**: El número debe coincidir con la migración relacionada
   - Migración: `005_missions.sql`
   - Datos: `005_missions_data.sql`, `005_missions_maps.sql`

2. **Idempotencia**: Usar `ON CONFLICT DO NOTHING`
   ```sql
   INSERT INTO table (id, name) VALUES (1, 'Admin')
   ON CONFLICT (id) DO NOTHING;
   ```

3. **Header**: Incluir información descriptiva
   ```sql
   -- =============================================
   -- DATA 005: Misiones Oficiales L4D2
   -- =============================================
   -- Descripción: Datos iniciales de misiones
   -- Migración: 005_missions.sql
   -- Idempotente: Usa ON CONFLICT DO NOTHING
   -- =============================================
   
   SET search_path TO core, public;
   
   INSERT INTO ...
   ```

### Tracking de Datos

La librería crea automáticamente una tabla `<schema>.schema_migrations_data`:

```sql
CREATE TABLE auth.schema_migrations_data (
    version VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255),
    applied_at TIMESTAMP DEFAULT NOW()
);
```

### Ejecución

Los datos se aplican **después** de las migraciones:
1. Aplica todas las migraciones (estructura)
2. Si `ApplyData = true`, aplica archivos de `data_sql/`
3. Registra datos aplicados en `schema_migrations_data`

## Tablas de Control

## Tablas de Control

La librería crea automáticamente dos tablas de tracking:

### Migraciones (estructura)
```sql
CREATE TABLE auth.schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT NOW()
);
```

### Datos iniciales
```sql
CREATE TABLE auth.schema_migrations_data (
    version VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255),
    applied_at TIMESTAMP DEFAULT NOW()
);
```

## Ejemplo Completo

Ver implementaciones en:
- `Connect-Auth/cmd/migrate/main.go`
- `Connect-Core/cmd/migrate/main.go`
- `Connect-RT/cmd/migrate/main.go`
