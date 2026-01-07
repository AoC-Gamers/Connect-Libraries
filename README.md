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

Este es un repositorio privado. Configura Go para acceder a repositorios privados:

### 1. Configurar GOPRIVATE

```bash
go env -w GOPRIVATE=github.com/AoC-Gamers/*
```

### 2. Configurar credenciales de Git

Opción A - HTTPS con token:
```bash
git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
```

Opción B - SSH (recomendado):
```bash
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### 3. Usar en tus proyectos

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

Este repositorio usa **versionado semántico unificado** para todas las librerías.

- **v1.0.0** - Release inicial
- **v1.1.0** - Nuevas features
- **v1.0.1** - Bug fixes

### Crear nueva versión

```bash
git tag v1.1.0
git push origin v1.1.0
```

### Actualizar en proyectos

```bash
go get github.com/AoC-Gamers/connect-libraries/auth-lib@v1.1.0
go mod tidy
```

## 🏗️ Estructura del Repositorio

```
connect-libraries/
├── apikey/              # Autenticación API Keys
├── auth-lib/            # Autenticación JWT
├── authz/               # Autorización y roles
├── core-types/          # Tipos compartidos
├── errors/              # Manejo de errores
├── middleware/          # Middlewares HTTP
├── migrate/             # Migraciones PostgreSQL
├── nats/                # Cliente NATS/JetStream
├── service-clients/     # Clientes HTTP inter-servicios
├── testhelpers/         # Utilidades de testing
├── .gitignore
└── README.md
```

## 📄 Licencia

Privado - AoC Gamers © 2026
