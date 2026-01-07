# Connect Libraries

🔧 Librerías compartidas y reutilizables para el ecosistema de microservicios Connect Backend.

## 📦 Librerías Disponibles

| Librería | Descripción | Versión |
|----------|-------------|---------|
| [auth-lib](./auth-lib/) | JWT & Claims Management | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [errors](./errors/) | Standardized Error Responses (RFC 7807) | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [middleware](./middleware/) | Framework Middlewares (Gin/Chi/HTTP) | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [apikey](./apikey/) | API Key Validation & Environment Integration | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [internal](./internal/) | Internal Shared Utilities | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [migrate](./migrate/) | Database Migration Utilities | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [nats](./nats/) | NATS Connection Manager | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |
| [testhelpers](./testhelpers/) | Testing Utilities | ![Version](https://img.shields.io/badge/version-1.0.0-blue) |

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

## 🔧 Desarrollo Local

Para trabajar con las librerías localmente sin publicar cambios, usa `replace` directives:

```go
// En tu proyecto Connect-Auth/go.mod
require (
    github.com/AoC-Gamers/connect-libraries/auth-lib v1.0.0
)

// Solo para desarrollo
replace github.com/AoC-Gamers/connect-libraries/auth-lib => ../connect-libraries/auth-lib
```

**Recuerda:** Comenta o elimina las directivas `replace` antes de hacer commit en producción.

## 🏗️ Estructura del Repositorio

```
connect-libraries/
├── auth-lib/       # JWT & Claims
├── errors/          # Errores estandarizados
├── middleware/      # Middlewares HTTP
├── apikey/          # Validación API Keys
├── internal/        # Utilidades internas
├── migrate/         # Migraciones DB
├── nats/            # Cliente NATS
├── testhelpers/             # Testing utilities
├── .gitignore
└── README.md
```

## 📄 Licencia

Privado - AoC Gamers © 2026
