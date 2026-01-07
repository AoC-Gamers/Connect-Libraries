# Authz (Authorization)

**Módulo:** `github.com/AoC-Gamers/connect-libraries/authz`

## 📋 Descripción

Biblioteca de autorización que define roles, permisos y políticas de acceso utilizadas en todo el sistema Connect. Proporciona constantes y utilidades para implementar control de acceso basado en roles (RBAC).

## 📦 Contenido

### `roles/`
Definiciones de roles y permisos del sistema:
- Roles de usuario (Admin, Moderador, Usuario, etc.)
- Permisos por recurso
- Políticas de acceso
- Jerarquías de autorización

## 🔧 Uso

```go
import (
    "github.com/AoC-Gamers/connect-libraries/authz/roles"
)

// Ejemplo: Verificar permisos
if roles.HasPermission(user.Role, roles.PermissionEditCommunity) {
    // Usuario autorizado
}
```

## ⚡ Características

- ✅ Sin dependencias externas
- ✅ Constantes de roles centralizadas
- ✅ Sistema de permisos granular
- ✅ Compatible con Casbin y otros sistemas RBAC
- ✅ Fácil extensión para nuevos roles/permisos
