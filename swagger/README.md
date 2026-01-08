# Swagger - Automatic OpenAPI Documentation

📚 Biblioteca para generación automática de documentación Swagger/OpenAPI en servicios Connect.

---

## 🎯 Características

- ✅ **Detección automática** de endpoints desde Chi Router
- ✅ **Agrupación por contexto** mediante tags configurables
- ✅ **Detección de seguridad** (JWT, API Keys) desde middlewares
- ✅ **Configuración flexible** por servicio
- ✅ **Zero annotations** - no requiere comentarios en código
- ✅ **Exportación JSON** compatible con OpenAPI

---

## 📦 Instalación

```bash
go get github.com/AoC-Gamers/connect-libraries/swagger@latest
```

---

## 🚀 Uso Básico

### 1. Crear configuración del servicio

```go
package main

import (
    "github.com/AoC-Gamers/connect-libraries/swagger"
)

func main() {
    // Configurar detector de Swagger
    swaggerConfig := swagger.DefaultConfig().
        WithServiceInfo("Connect-Auth", "1.0.0").
        WithTagRules([]swagger.TagRule{
            {PathPattern: "/teams", Tags: []string{"Teams"}},
            {PathPattern: "/roles", Tags: []string{"Roles"}},
            {PathPattern: "/memberships", Tags: []string{"Memberships"}},
            {PathPattern: "/invitations", Tags: []string{"Invitations"}},
            {PathPattern: "/sanctions", Tags: []string{"Sanctions"}},
            {PathPattern: "/auth/", Tags: []string{"Authentication"}},
        })

    // Crear detector
    detector := swagger.NewDetector(swaggerConfig)

    // Escanear router (después de definir todas las rutas)
    router := setupRouter() // tu función que crea el router
    detector.ScanRouter(router)

    // Exponer endpoint de rutas detectadas
    router.Get("/swagger/routes", detector.ServeHTTP)
}
```

### 2. Orden de Tag Rules (IMPORTANTE)

⚠️ **Las reglas se evalúan en orden - la primera coincidencia gana**

```go
// ✅ CORRECTO - Específico primero
TagRules: []swagger.TagRule{
    {PathPattern: "/teams", Tags: []string{"Teams"}},           // Específico
    {PathPattern: "/memberships", Tags: []string{"Memberships"}}, // Específico
    {PathPattern: "/auth/", Tags: []string{"Authentication"}},   // Genérico al final
}

// ❌ INCORRECTO - Genérico primero captura todo
TagRules: []swagger.TagRule{
    {PathPattern: "/auth/", Tags: []string{"Authentication"}},   // ❌ Captura /auth/teams también
    {PathPattern: "/teams", Tags: []string{"Teams"}},           // Nunca se ejecuta
}
```

---

## 📋 Ejemplos por Servicio

### Connect-Auth

```go
swaggerConfig := swagger.DefaultConfig().
    WithServiceInfo("Connect-Auth", "1.0.0").
    WithTagRules([]swagger.TagRule{
        // Recursos específicos PRIMERO
        {PathPattern: "/teams", Tags: []string{"Teams"}},
        {PathPattern: "/memberships", Tags: []string{"Memberships"}},
        {PathPattern: "/roles", Tags: []string{"Roles"}},
        {PathPattern: "/invitations", Tags: []string{"Invitations"}},
        {PathPattern: "/sanctions", Tags: []string{"Sanctions"}},
        {PathPattern: "/vip", Tags: []string{"VIP Management"}},
        {PathPattern: "/permissions", Tags: []string{"Authorization"}},
        {PathPattern: "/authz", Tags: []string{"Authorization"}},
        {PathPattern: "/me/", Tags: []string{"User Profile"}},
        {PathPattern: "/admin", Tags: []string{"Administration"}},
        {PathPattern: "/internal", Tags: []string{"Internal API"}},
        {PathPattern: "/cache", Tags: []string{"Cache Management"}},
        // Genéricos AL FINAL
        {PathPattern: "/auth/steam", Tags: []string{"Authentication"}},
        {PathPattern: "/auth/", Tags: []string{"Authentication"}},
    }).
    AddPublicPath("/auth/steam").
    AddPublicPath("/auth/steam/callback")
```

### Connect-Core

```go
swaggerConfig := swagger.DefaultConfig().
    WithServiceInfo("Connect-Core", "1.0.0").
    WithTagRules([]swagger.TagRule{
        {PathPattern: "/missions", Tags: []string{"Missions"}},
        {PathPattern: "/gamemodes", Tags: []string{"Game Modes"}},
        {PathPattern: "/communities", Tags: []string{"Communities"}},
        {PathPattern: "/teams", Tags: []string{"Teams"}},
        {PathPattern: "/users", Tags: []string{"Users"}},
        {PathPattern: "/settings", Tags: []string{"Settings"}},
        {PathPattern: "/locale", Tags: []string{"Localization"}},
        {PathPattern: "/media", Tags: []string{"Media"}},
        {PathPattern: "/internal", Tags: []string{"Internal API"}},
    }).
    AddPublicPath("/core/settings").
    AddPublicPath("/core/locale").
    AddPublicPath("/core/missions")
```

### Connect-Lobby

```go
swaggerConfig := swagger.DefaultConfig().
    WithServiceInfo("Connect-Lobby", "1.0.0").
    WithTagRules([]swagger.TagRule{
        {PathPattern: "/lobbies", Tags: []string{"Lobbies"}},
        {PathPattern: "/game-modes", Tags: []string{"Game Configuration"}},
        {PathPattern: "/missions", Tags: []string{"Game Configuration"}},
        {PathPattern: "/servers", Tags: []string{"Game Configuration"}},
    }).
    AddPublicPath("/lobby/game-modes").
    AddPublicPath("/lobby/missions")
```

### Connect-RT

```go
swaggerConfig := swagger.DefaultConfig().
    WithServiceInfo("Connect-RT", "1.0.0").
    WithTagRules([]swagger.TagRule{
        {PathPattern: "/ws", Tags: []string{"WebSocket"}},
        {PathPattern: "/presence", Tags: []string{"Presence"}},
        {PathPattern: "/test/nats", Tags: []string{"Testing"}},
    }).
    AddPublicPath("/ws").
    AddPublicPath("/presence")
```

---

## 🔐 Detección de Seguridad

El detector identifica automáticamente el tipo de autenticación mediante middlewares:

| Middleware Pattern | Security Type | Descripción |
|-------------------|---------------|-------------|
| `RequireAuth` | `BearerAuth` | JWT Authentication |
| `JWTAuth` | `BearerAuth` | JWT Authentication |
| `RequireAPIKey` | `ApiKeyAuth` | API Key Authentication |
| `RequireInternalServices` | `ApiKeyAuth` | Service-to-Service |
| `RequireWebPermission` | `BearerAuth` | Permission-based |

### Personalizar Patrones de Seguridad

```go
config := swagger.DefaultConfig()
config.SecurityPatterns["MyCustomMiddleware"] = "BearerAuth"
```

---

## 📝 Configuración Avanzada

### Método Fluent (Recomendado)

```go
config := swagger.DefaultConfig().
    WithServiceInfo("My-Service", "2.0.0").
    WithDefaultTag("API").
    WithDefaultSecurity("BearerAuth").
    AddTagRule("/users", "Users").
    AddTagRule("/products", "Products").
    AddPublicPath("/public").
    AddSkipPath("/debug")
```

### Método Manual

```go
config := &swagger.Config{
    ServiceName: "My-Service",
    Version:     "2.0.0",
    TagRules: []swagger.TagRule{
        {PathPattern: "/users", Tags: []string{"Users"}},
        {PathPattern: "/products", Tags: []string{"Products"}},
    },
    PublicPaths: []string{"/public", "/health"},
    SkipPaths:   []string{"/debug", "/swagger"},
    DefaultTag:  "API",
    DefaultSecurity: "BearerAuth",
}
```

---

## 🔄 Endpoint de Rutas

Exponer las rutas detectadas como JSON:

```go
// Opción 1: Usar ServeHTTP del detector
router.Get("/swagger/routes", detector.ServeHTTP)

// Opción 2: Custom handler
router.Get("/swagger/routes", func(w http.ResponseWriter, r *http.Request) {
    routes := detector.GetRoutes()
    json.NewEncoder(w).Encode(routes)
})
```

**Respuesta ejemplo:**
```json
[
  {
    "method": "GET",
    "path": "/auth/teams",
    "security": ["BearerAuth"],
    "tags": ["Teams"],
    "summary": "List Teams",
    "description": "Endpoint for GET /auth/teams. Requires: JWT authentication"
  }
]
```

---

## 🎨 Convenciones de Nombres

### Tags Recomendados por Servicio

**Connect-Auth:**
- `Authentication` - Login, Steam OAuth
- `Teams` - Gestión de equipos
- `Memberships` - Membresías de usuarios
- `Roles` - Gestión de roles
- `Invitations` - Sistema de invitaciones
- `Sanctions` - Sistema de sanciones
- `Authorization` - Permisos y authz
- `User Profile` - Perfil de usuario
- `VIP Management` - Sistema VIP
- `Administration` - Endpoints admin
- `Internal API` - APIs internas
- `Cache Management` - Gestión de caché

**Connect-Core:**
- `Missions` - Misiones y mapas
- `Game Modes` - Modos de juego
- `Communities` - Comunidades
- `Teams` - Equipos
- `Users` - Usuarios
- `Settings` - Configuración
- `Localization` - Traducciones
- `Media` - Imágenes y assets
- `Internal API` - APIs internas

**Connect-Lobby:**
- `Lobbies` - Gestión de lobbies
- `Game Configuration` - Configuración de juego
- `Internal API` - APIs internas

**Connect-RT:**
- `WebSocket` - Conexiones WS
- `Presence` - Sistema de presencia
- `Testing` - Endpoints de testing

---

## 🐛 Troubleshooting

### Los tags no se asignan correctamente

**Problema:** Todos los endpoints tienen el tag "Authentication"

**Solución:** Verifica el orden de las reglas. Las reglas específicas deben ir primero:

```go
// ❌ MAL
{PathPattern: "/auth/", Tags: []string{"Authentication"}},  // Primero (captura todo)
{PathPattern: "/teams", Tags: []string{"Teams"}},          // Segundo (nunca se ejecuta)

// ✅ BIEN
{PathPattern: "/teams", Tags: []string{"Teams"}},          // Primero (específico)
{PathPattern: "/auth/", Tags: []string{"Authentication"}},  // Segundo (genérico)
```

### Endpoints públicos aparecen con seguridad

**Problema:** `/auth/steam/callback` requiere JWT incorrectamente

**Solución:** Agregar a PublicPaths:

```go
config.AddPublicPath("/auth/steam/callback")
```

### Rutas no aparecen en Swagger

**Problema:** Algunas rutas no se detectan

**Solución:** Verifica que no estén en SkipPaths:

```go
config.SkipPaths = []string{"/swagger", "/debug"} // Solo estas se omiten
```

---

## 📚 Recursos

- [OpenAPI Specification](https://swagger.io/specification/)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [Convenciones REST](https://restfulapi.net/)

---

## 🔄 Changelog

Ver [CHANGELOG.md](./CHANGELOG.md) para el historial de cambios.

---

## 📄 Licencia

Parte del ecosistema Connect Libraries - Uso interno AoC-Gamers.
