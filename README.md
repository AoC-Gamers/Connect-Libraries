# Connect Libraries

[![CI](https://github.com/AoC-Gamers/connect-libraries/actions/workflows/ci.yml/badge.svg)](https://github.com/AoC-Gamers/connect-libraries/actions/workflows/ci.yml)
[![Release](https://github.com/AoC-Gamers/connect-libraries/actions/workflows/release.yml/badge.svg)](https://github.com/AoC-Gamers/connect-libraries/actions/workflows/release.yml)

🔧 Librerías compartidas y reutilizables para el ecosistema de microservicios Connect Backend.

## 📦 Librerías Disponibles

| Librería | Descripción | Versión |
|----------|-------------|---------|
| [apikey](./apikey/) | Autenticación y validación de API Keys | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [authz](./authz/) | Sistema de autorización, roles y permisos (RBAC) | ![Version](https://img.shields.io/badge/version-1.0.1-blue) |
| [errors](./errors/) | Manejo de errores estandarizado (RFC 7807) | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [middleware](./middleware/) | Middlewares HTTP para framework Chi | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [migrate](./migrate/) | Sistema de migraciones para PostgreSQL | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [nats](./nats/) | Cliente NATS con soporte JetStream | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [swagger](./swagger/) | Detección automática de Swagger/OpenAPI | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [testhelpers](./testhelpers/) | Utilidades para testing y mocks | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |

## 🚀 Instalación

### Usar en tus proyectos

```go
// go.mod
module github.com/AoC-Gamers/Connect-Auth

require (
    github.com/AoC-Gamers/connect-libraries/authz v1.0.1
    github.com/AoC-Gamers/connect-libraries/errors v1.0.0
    github.com/AoC-Gamers/connect-libraries/middleware v1.0.0
)
```

```bash
go get github.com/AoC-Gamers/connect-libraries/authz@v1.0.1
go mod tidy
```

## 📝 Versionado

Este repositorio usa **versionado independiente por biblioteca** siguiendo Semantic Versioning.

Cada biblioteca tiene su propio ciclo de versiones con tags en el formato `<librería>/v<versión>`:

- `apikey/v1.0.0`, `apikey/v1.0.1`, `apikey/v1.1.0`, ...
- `authz/v1.0.0`, `authz/v1.0.1`, ...
- `errors/v1.0.0`, `errors/v1.0.1`, `errors/v1.1.0`, ...
- etc.

### Crear nueva versión de una biblioteca

```bash
# Ejemplo: Nueva versión de authz
cd authz
# Actualizar CHANGELOG.md con los cambios
git add .
git commit -m "feat(authz): nueva funcionalidad"
git tag authz/v1.0.2
git push origin authz/v1.0.2
```

### Actualizar en proyectos

```bash
# Actualizar a una versión específica
go get github.com/AoC-Gamers/connect-libraries/authz@v1.0.1

# O usar la última versión
go get -u github.com/AoC-Gamers/connect-libraries/authz

go mod tidy
```

### Consultar versiones disponibles

```bash
# Ver todas las versiones de una biblioteca
go list -m -versions github.com/AoC-Gamers/connect-libraries/authz

# Ver tags en GitHub
git ls-remote --tags origin | grep authz
```

## 🏗️ Estructura del Repositorio

```
connect-libraries/
├── apikey/              # Autenticación API Keys
│   ├── CHANGELOG.md     # Historial de versiones
│   └── ...
├── authz/               # Autorización, roles y permisos
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
├── testhelpers/         # Utilidades de testing
│   ├── CHANGELOG.md
│   └── ...
├── .gitignore
└── README.md
```

> **Nota**: Cada biblioteca mantiene su propio CHANGELOG.md con su historial de versiones independiente.

## 📄 Licencia

AoC Gamers © 2026
