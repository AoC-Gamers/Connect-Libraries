# Core Types

**Módulo:** `github.com/AoC-Gamers/connect-libraries/core-types`

## 📋 Descripción

Biblioteca de tipos compartidos y definiciones comunes utilizadas por todos los microservicios del sistema Connect. Esta biblioteca no tiene dependencias externas y proporciona la base de tipos para la comunicación entre servicios.

## 📦 Contenido

### `endpoints/`
Constantes de rutas y endpoints HTTP para cada microservicio:
- **auth.go** - Endpoints del servicio Connect-Auth
- **core.go** - Endpoints del servicio Connect-Core
- **rt.go** - Endpoints del servicio Connect-RT

### `models/`
Modelos de datos compartidos entre servicios:
- **auth.go** - Estructuras de autenticación y usuarios
- **core.go** - Estructuras del dominio principal (comunidades, equipos, etc.)
- **rt.go** - Estructuras de tiempo real y eventos

### `errors/`
Sistema centralizado de manejo de errores con códigos consistentes y traducción i18n.

## 🔧 Uso

```go
import (
    "github.com/AoC-Gamers/connect-libraries/core-types/endpoints"
    "github.com/AoC-Gamers/connect-libraries/core-types/models"
    "github.com/AoC-Gamers/connect-libraries/core-types/errors"
)
```

## ⚡ Características

- ✅ Sin dependencias externas
- ✅ Tipos compartidos para comunicación entre servicios
- ✅ Sistema de errores estandarizado
- ✅ Constantes de endpoints centralizadas
