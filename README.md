# Connect Libraries

🔧 Librerías compartidas y reutilizables para el ecosistema de microservicios Connect Backend.

## 📦 Librerías Disponibles

| Librería | Descripción | Versión |
|----------|-------------|---------|
| [apikey](./apikey/) | Autenticación y validación de API Keys | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [auth-lib](./auth-lib/) | Sistema de autenticación JWT y permisos | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [authz](./authz/) | Sistema de autorización y roles (RBAC) | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [core-types](./core-types/) | Tipos compartidos: endpoints, modelos y errores | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [errors](./errors/) | Manejo de errores estandarizado (RFC 7807) | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [middleware](./middleware/) | Middlewares HTTP para framework Chi | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [migrate](./migrate/) | Sistema de migraciones para PostgreSQL | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [nats](./nats/) | Cliente NATS con soporte JetStream | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [service-clients](./service-clients/) | Clientes HTTP para comunicación entre servicios | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [testhelpers](./testhelpers/) | Utilidades para testing y mocks | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |

## 🚀 Instalación

### Usar en tus proyectos

```go
// go.mod
module github.com/AoC-Gamers/Connect-Auth

require (
    github.com/AoC-Gamers/connect-libraries/auth-lib v1.0.0
    github.com/AoC-Gamers/connect-libraries/errors v1.0.0
    github.com/AoC-Gamers/connect-libraries/middleware v1.0.0
)
```

```bash
go get github.com/AoC-Gamers/connect-libraries/auth-lib@v1.0.0
go mod tidy
```

## 📝 Versionado

Este repositorio usa **versionado independiente por biblioteca** siguiendo Semantic Versioning.

Cada biblioteca tiene su propio ciclo de versiones con tags en el formato `<librería>/v<versión>`:

- `apikey/v1.0.0`, `apikey/v1.0.1`, `apikey/v1.1.0`, ...
- `auth-lib/v1.0.0`, `auth-lib/v1.0.1`, `auth-lib/v1.1.0`, ...
- `errors/v1.0.0`, `errors/v1.0.1`, `errors/v1.1.0`, ...
- etc.

### Crear nueva versión de una biblioteca

```bash
# Ejemplo: Nueva versión de auth-lib
cd auth-lib
# Actualizar CHANGELOG.md con los cambios
git add .
git commit -m "feat(auth-lib): nueva funcionalidad"
git tag auth-lib/v1.1.0
git push origin auth-lib/v1.1.0
```

### Actualizar en proyectos

```bash
# Actualizar a una versión específica
go get github.com/AoC-Gamers/connect-libraries/auth-lib@v1.1.0

# O usar la última versión
go get -u github.com/AoC-Gamers/connect-libraries/auth-lib

go mod tidy
```

### Consultar versiones disponibles

```bash
# Ver todas las versiones de una biblioteca
go list -m -versions github.com/AoC-Gamers/connect-libraries/auth-lib

# Ver tags en GitHub
git ls-remote --tags origin | grep auth-lib
```

## 🏗️ Estructura del Repositorio

```
connect-libraries/
├── apikey/              # Autenticación API Keys
│   ├── CHANGELOG.md     # Historial de versiones
│   └── ...
├── auth-lib/            # Autenticación JWT
│   ├── CHANGELOG.md
│   └── ...
├── authz/               # Autorización y roles
│   ├── CHANGELOG.md
│   └── ...
├── core-types/          # Tipos compartidos
│   ├── CHANGELOG.md
│   └── ...
├── errors/              # Manejo de errores
│   ├── CHANGELOG.md
│   └── ...
├── middleware/          # Middlewares HTTP
│   ├── CHANGELOG.md
│   └── ...
├── migrate/             # Migraciones PostgreSQL
│   ├── CHANGELOG.md
│   └── ...
├── nats/                # Cliente NATS/JetStream
│   ├── CHANGELOG.md
│   └── ...
├── service-clients/     # Clientes HTTP inter-servicios
│   ├── CHANGELOG.md
│   └── ...
├── testhelpers/         # Utilidades de testing
│   ├── CHANGELOG.md
│   └── ...
├── .gitignore
└── README.md
```

> **Nota**: Cada biblioteca mantiene su propio CHANGELOG.md con su historial de versiones independiente.

## 📄 Licencia

AoC Gamers © 2026
