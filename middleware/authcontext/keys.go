package authcontext

// ContextKey define claves estándar para contexto entre frameworks
type ContextKey string

// Standard context keys usadas por todos los microservicios Connect
const (
	SteamIDKey          ContextKey = "steamid"
	RoleKey             ContextKey = "role"              // Rol único (web_user, web_admin, etc.)
	AllowPermissionsKey ContextKey = "allow_permissions" // Máscara de permisos permitidos (uint64)
	DenyPermissionsKey  ContextKey = "deny_permissions"  // Máscara de permisos denegados (uint64)
	ClaimsKey           ContextKey = "claims"            // Claims completos
)