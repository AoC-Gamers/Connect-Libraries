package chi

import (
	"context"
	"net/http"
	"strings"

	contextlib "github.com/AoC-Gamers/connect-libraries/auth-lib/context"
	authlib "github.com/AoC-Gamers/connect-libraries/auth-lib/jwt"
	errors "github.com/AoC-Gamers/connect-libraries/errors"
	"github.com/rs/zerolog/log"
)

// RequireAuth middleware de autenticación JWT para Chi usando connect-auth-lib
func RequireAuth(config authlib.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)

			claims, err := authlib.ParseAndValidate(tokenStr, config.JWTSecret, config.PolicyVersionGlobal)
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
					errors.RespondTokenExpired(w)
				} else if strings.Contains(err.Error(), "policy version") {
					errors.RespondPolicyVersionMismatch(w, 0, config.PolicyVersionGlobal)
				} else {
					errors.RespondUnauthorized(w, "Missing or invalid authentication token")
				}
				return
			}

			// Inyectar claims en contexto usando keys estándar
			ctx := context.WithValue(r.Context(), contextlib.SteamIDKey, claims.GetSteamID())
			ctx = context.WithValue(ctx, contextlib.RoleKey, claims.GetRole())
			ctx = context.WithValue(ctx, contextlib.ClaimsKey, claims)

			// Inyectar permisos web para uso en middleware de autorización
			ctx = context.WithValue(ctx, contextlib.AllowPermissionsKey, claims.AllowPermissions)
			ctx = context.WithValue(ctx, contextlib.DenyPermissionsKey, claims.DenyPermissions)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware de autorización por roles para Chi
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r)
			if claims == nil {
				errors.RespondUnauthorized(w, "authentication required")
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
				errors.RespondInsufficientPermissions(w, "required roles: "+strings.Join(allowedRoles, ", "))
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
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaimsFromContext(r)
			if claims == nil {
				errors.RespondUnauthorized(w, "authentication required")
				return
			}

			if !claims.HasPermission(permission) {
				errors.RespondInsufficientPermissions(w, "required permission bitmask")
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
func OptionalAuth(config authlib.AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr != "" {
				claims, err := authlib.ParseAndValidate(tokenStr, config.JWTSecret, config.PolicyVersionGlobal)
				if err == nil {
					ctx := context.WithValue(r.Context(), contextlib.SteamIDKey, claims.GetSteamID())
					ctx = context.WithValue(ctx, contextlib.RoleKey, claims.GetRole())
					ctx = context.WithValue(ctx, contextlib.ClaimsKey, claims)

					// Inyectar permisos web para uso en middleware de autorización
					ctx = context.WithValue(ctx, contextlib.AllowPermissionsKey, claims.AllowPermissions)
					ctx = context.WithValue(ctx, contextlib.DenyPermissionsKey, claims.DenyPermissions)

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
