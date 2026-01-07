package chi

import (
	"net/http"

	contextlib "github.com/AoC-Gamers/connect-libraries/auth-lib/context"
	authlib "github.com/AoC-Gamers/connect-libraries/auth-lib/jwt"
)

// GetSteamIDFromContext obtiene el SteamID del contexto de Chi
func GetSteamIDFromContext(r *http.Request) string {
	if value := r.Context().Value(contextlib.SteamIDKey); value != nil {
		if steamID, ok := value.(string); ok {
			return steamID
		}
	}
	return ""
}

// GetRoleFromContext obtiene el rol del contexto de Chi
func GetRoleFromContext(r *http.Request) string {
	if value := r.Context().Value(contextlib.RoleKey); value != nil {
		if role, ok := value.(string); ok {
			return role
		}
	}
	return ""
}

// GetClaimsFromContext obtiene los claims completos del contexto de Chi
func GetClaimsFromContext(r *http.Request) *authlib.Claims {
	if value := r.Context().Value(contextlib.ClaimsKey); value != nil {
		if claims, ok := value.(*authlib.Claims); ok {
			return claims
		}
	}
	return nil
}

// HasRole verifica si el usuario actual tiene un rol específico
func HasRole(r *http.Request, role string) bool {
	claims := GetClaimsFromContext(r)
	if claims == nil {
		return false
	}
	return claims.GetRole() == role
}

// IsAdmin verifica si el usuario actual es administrador
func IsAdmin(r *http.Request) bool {
	claims := GetClaimsFromContext(r)
	if claims == nil {
		return false
	}
	return claims.IsAdmin()
}

// IsStaff verifica si el usuario actual es staff o superior
func IsStaff(r *http.Request) bool {
	claims := GetClaimsFromContext(r)
	if claims == nil {
		return false
	}
	return claims.IsStaff()
}

// HasPermission verifica si el usuario tiene un permiso específico (bitmask)
func HasPermission(r *http.Request, permission uint64) bool {
	claims := GetClaimsFromContext(r)
	if claims == nil {
		return false
	}
	return claims.HasPermission(permission)
}
