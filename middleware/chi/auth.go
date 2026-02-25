package chi

import (
	"context"
	"net/http"
	"strings"

	"github.com/AoC-Gamers/connect-libraries/middleware/authcontext"
	"github.com/AoC-Gamers/connect-libraries/middleware/authjwt"
	"github.com/rs/zerolog/log"
)

// RequireAuth middleware de autenticación JWT para Chi
func RequireAuth(config authjwt.AuthConfig) func(http.Handler) http.Handler {
	return RequireAuthWithResponder(config, nil)
}

// RequireAuthWithResponder permite inyectar un ErrorResponder personalizado
func RequireAuthWithResponder(config authjwt.AuthConfig, responder ErrorResponder) func(http.Handler) http.Handler {
	responder = ensureResponder(responder)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)

			claims, err := authjwt.ParseAndValidate(tokenStr, config.SignerMaterial, config.PolicyVersionGlobal)
			if err != nil {
				// Log del error de autenticación
				log.Error().
					Err(err).
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Int("expected_policy_version", config.PolicyVersionGlobal).
					Msg("JWT authentication failed")

				// Determinar el código de error apropiado
				if strings.Contains(err.Error(), "expired") {
					responder.TokenExpired(w)
				} else if strings.Contains(err.Error(), "policy version") {
					responder.PolicyVersionMismatch(w, 0, config.PolicyVersionGlobal)
				} else {
					responder.Unauthorized(w, "Missing or invalid authentication token")
				}
				return
			}

			// Inyectar claims en contexto usando keys estándar
			ctx := context.WithValue(r.Context(), authcontext.SteamIDKey, claims.GetSteamID())
			ctx = context.WithValue(ctx, authcontext.RoleKey, claims.GetRole())
			ctx = context.WithValue(ctx, authcontext.ClaimsKey, claims)

			// Inyectar permisos web para uso en middleware de autorización
			ctx = context.WithValue(ctx, authcontext.AllowPermissionsKey, claims.AllowPermissions)
			ctx = context.WithValue(ctx, authcontext.DenyPermissionsKey, claims.DenyPermissions)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware de autorización por roles para Chi
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return RequireRoleWithResponder(nil, allowedRoles...)
}

// RequireRoleWithResponder permite inyectar un ErrorResponder personalizado
func RequireRoleWithResponder(responder ErrorResponder, allowedRoles ...string) func(http.Handler) http.Handler {
	responder = ensureResponder(responder)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r)
			if claims == nil {
				responder.Unauthorized(w, "authentication required")
				return
			}

			// Verificar si el rol del usuario está en la lista de roles permitidos
			role := claims.GetRole()
			allowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					allowed = true
					break
				}
			}

			if !allowed {
				responder.InsufficientPermissions(w, "required roles: "+strings.Join(allowedRoles, ", "))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission was removed - use RequirePermissionBitmask instead
// This legacy function presented a security risk by not performing actual permission validation

// RequirePermissionBitmask middleware de autorización por permisos con bitmask uint64
func RequirePermissionBitmask(permission uint64) func(http.Handler) http.Handler {
	return RequirePermissionBitmaskWithResponder(permission, nil)
}

// RequirePermissionBitmaskWithResponder permite inyectar un ErrorResponder personalizado
func RequirePermissionBitmaskWithResponder(permission uint64, responder ErrorResponder) func(http.Handler) http.Handler {
	responder = ensureResponder(responder)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r)
			if claims == nil {
				responder.Unauthorized(w, "authentication required")
				return
			}

			if !claims.HasPermission(permission) {
				responder.InsufficientPermissions(w, "required permission bitmask")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdmin middleware que requiere rol de administrador
func RequireAdmin() func(http.Handler) http.Handler {
	return RequireRole("admin", "super_admin")
}

// RequireStaff middleware que requiere rol de staff o superior
func RequireStaff() func(http.Handler) http.Handler {
	return RequireRole("staff", "admin", "super_admin")
}

// OptionalAuth middleware que intenta autenticar pero no falla si no hay token
func OptionalAuth(config authjwt.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr != "" {
				claims, err := authjwt.ParseAndValidate(tokenStr, config.SignerMaterial, config.PolicyVersionGlobal)
				if err == nil {
					ctx := context.WithValue(r.Context(), authcontext.SteamIDKey, claims.GetSteamID())
					ctx = context.WithValue(ctx, authcontext.RoleKey, claims.GetRole())
					ctx = context.WithValue(ctx, authcontext.ClaimsKey, claims)

					// Inyectar permisos web para uso en middleware de autorización
					ctx = context.WithValue(ctx, authcontext.AllowPermissionsKey, claims.AllowPermissions)
					ctx = context.WithValue(ctx, authcontext.DenyPermissionsKey, claims.DenyPermissions)

					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// extractToken obtiene el JWT desde Authorization header, query param o cookie
func extractToken(r *http.Request) string {
	// 1. Authorization header (Bearer token)
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		return strings.TrimSpace(auth[7:])
	}

	// 2. Query parameter
	if token := r.URL.Query().Get("token"); token != "" {
		return token
	}

	// 3. Cookie
	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	return ""
}
