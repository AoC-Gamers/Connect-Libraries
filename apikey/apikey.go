package apikey

import (
	"context"
	"net/http"
	"strings"

	zlog "github.com/rs/zerolog/log"
)

// Context keys used by the chi middleware
type ctxKey string

const (
	ctxKeyServiceName ctxKey = "service_name"
	ctxKeyAPIKey      ctxKey = "api_key"
	ctxKeyAuthType    ctxKey = "auth_type"
)

// Helper functions to set context values
func setServiceName(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, ctxKeyServiceName, serviceName)
}

func setAPIKey(ctx context.Context, apiKey string) context.Context {
	return context.WithValue(ctx, ctxKeyAPIKey, apiKey)
}

func setAuthType(ctx context.Context, authType string) context.Context {
	return context.WithValue(ctx, ctxKeyAuthType, authType)
}

// RequireAPIKey valida que la petición incluya un API Key válido
func RequireAPIKey(validator *Validator) func(http.Handler) http.Handler {
	return RequireAPIKeyWithResponder(validator, nil)
}

// RequireAPIKeyWithResponder permite inyectar un ErrorResponder personalizado
func RequireAPIKeyWithResponder(validator *Validator, responder ErrorResponder) func(http.Handler) http.Handler {
	responder = ensureResponder(responder)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg := DefaultConfig()
			key := validator.ExtractAPIKey(r, cfg)

			zlog.Debug().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Str("remote_addr", r.RemoteAddr).
				Msg("🔐 API Key validation started")

			serviceName, valid := validator.ValidateKey(key)
			if !valid {
				maskedKey := "none"
				if key != "" {
					if len(key) > 8 {
						maskedKey = key[:4] + "****" + key[len(key)-4:]
					} else {
						maskedKey = "****"
					}
				}

				zlog.Warn().
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Str("api_key_received", maskedKey).
					Str("header", cfg.HeaderName).
					Msg("❌ API Key validation failed")

				responder.Unauthorized(w, "invalid or missing API key")
				return
			}

			zlog.Debug().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Str("service", serviceName).
				Msg("✅ API Key validated successfully")

			ctx := setServiceName(r.Context(), serviceName)
			ctx = setAPIKey(ctx, key)
			ctx = setAuthType(ctx, "api_key")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireConnectAPIKey carga las API keys desde env y retorna el middleware
func RequireConnectAPIKey() func(http.Handler) http.Handler {
	return RequireConnectAPIKeyWithResponder(nil)
}

// RequireConnectAPIKeyWithResponder permite inyectar un ErrorResponder personalizado
func RequireConnectAPIKeyWithResponder(responder ErrorResponder) func(http.Handler) http.Handler {
	validator, err := LoadConnectAPIKeys()
	if err != nil {
		zlog.Warn().Err(err).Msg("⚠️ Failed to load Connect API keys from environment")
		zlog.Info().Msg("💡 Make sure your .env file contains: AUTH_API_KEY, CORE_API_KEY, LOBBY_API_KEY, RT_API_KEY")
		validator, err = LoadConnectAPIKeysPermissive()
		if err != nil {
			zlog.Warn().Err(err).Msg("⚠️ Failed to load Connect API keys in permissive mode")
			validator = NewValidator(make(map[string]string))
		}
	}

	return RequireAPIKeyWithResponder(validator, responder)
}

// RequireConnectService middleware que requiere API key de servicios Connect específicos
// Valida la API key Y verifica que el servicio esté en la lista de permitidos
func RequireConnectService(allowedServices ...string) func(http.Handler) http.Handler {
	return RequireConnectServiceWithResponder(nil, allowedServices...)
}

// RequireConnectServiceWithResponder permite inyectar un ErrorResponder personalizado
func RequireConnectServiceWithResponder(responder ErrorResponder, allowedServices ...string) func(http.Handler) http.Handler {
	responder = ensureResponder(responder)
	validator, err := LoadConnectAPIKeysPermissive()
	if err != nil {
		zlog.Warn().Err(err).Msg("⚠️ Failed to load Connect API keys")
		zlog.Warn().Msg("⚠️ Creating empty validator - ALL API KEY REQUESTS WILL FAIL")
		validator = NewValidator(make(map[string]string))
	}

	// Pre-calcular el set de servicios permitidos para O(1) lookup
	allowed := make(map[string]struct{}, len(allowedServices))
	for _, s := range allowedServices {
		allowed[s] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := validator.ExtractAPIKey(r, DefaultConfig())
			serviceName, valid := validator.ValidateKey(key)
			if !valid {
				responder.Unauthorized(w, "invalid or missing API key")
				return
			}

			// Verificar si el servicio está en la lista de permitidos
			if _, ok := allowed[serviceName]; !ok {
				responder.InsufficientPermissions(w, "service not authorized for this endpoint")
				return
			}

			ctx := r.Context()
			ctx = setServiceName(ctx, serviceName)
			ctx = setAPIKey(ctx, key)
			ctx = setAuthType(ctx, "api_key")

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuthService middleware específico para Connect-Auth
func RequireAuthService() func(http.Handler) http.Handler {
	return RequireConnectService("connect-auth")
}

// RequireCoreService middleware específico para Connect-Core
func RequireCoreService() func(http.Handler) http.Handler {
	return RequireConnectService("connect-core")
}

// RequireLobbyService middleware específico para Connect-Lobby
func RequireLobbyService() func(http.Handler) http.Handler {
	return RequireConnectService("connect-lobby")
}

// RequireRTService middleware específico para Connect-RT
func RequireRTService() func(http.Handler) http.Handler {
	return RequireConnectService("connect-rt")
}

// RequireInternalServices middleware para endpoints internos
func RequireInternalServices() func(http.Handler) http.Handler {
	return RequireInternalServicesWithResponder(nil)
}

// RequireInternalServicesWithResponder permite inyectar un ErrorResponder personalizado
func RequireInternalServicesWithResponder(responder ErrorResponder) func(http.Handler) http.Handler {
	return RequireConnectServiceWithResponder(responder, "connect-auth", "connect-core", "connect-lobby", "connect-rt")
}

// AutoAPIKeyMiddleware crea middleware de API key con configuración automática
func AutoAPIKeyMiddleware() func(http.Handler) http.Handler {
	return RequireConnectAPIKey()
}

// GetServiceNameFromContext obtiene el nombre del servicio autenticado
func GetServiceNameFromContext(r *http.Request) string {
	if v := r.Context().Value(ctxKeyServiceName); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// IsServiceAuthenticated verifica si la request tiene autenticación de servicio
func IsServiceAuthenticated(r *http.Request) bool {
	return strings.TrimSpace(GetServiceNameFromContext(r)) != ""
}

// GetConnectServiceFromContext helper que obtiene el servicio Connect autenticado
// Retorna el nombre completo del servicio (ej: "connect-core") o vacío si no es válido
func GetConnectServiceFromContext(r *http.Request) string {
	serviceName := GetServiceNameFromContext(r)
	if serviceName == "" {
		return ""
	}

	// Verificar si es un servicio Connect válido
	for _, validService := range ConnectServices {
		if serviceName == validService {
			return serviceName
		}
	}

	return ""
}

// IsConnectService verifica si la request viene de un servicio Connect válido
func IsConnectService(r *http.Request) bool {
	return GetConnectServiceFromContext(r) != ""
}

// IsServiceType verifica si la request viene de un tipo de servicio específico
// serviceType puede ser "auth", "core", "lobby", "rt" (con o sin prefijo "connect-")
func IsServiceType(r *http.Request, serviceType string) bool {
	serviceName := GetServiceNameFromContext(r)
	if serviceName == "" {
		return false
	}

	// Aceptar tanto "auth" como "connect-auth"
	return serviceName == "connect-"+serviceType || serviceName == serviceType
}

// IsAuthService verifica si la request viene de Connect-Auth
func IsAuthService(r *http.Request) bool {
	return IsServiceType(r, "auth")
}

// IsCoreService verifica si la request viene de Connect-Core
func IsCoreService(r *http.Request) bool {
	return IsServiceType(r, "core")
}

// IsLobbyService verifica si la request viene de Connect-Lobby
func IsLobbyService(r *http.Request) bool {
	return IsServiceType(r, "lobby")
}

// IsRTService verifica si la request viene de Connect-RT
func IsRTService(r *http.Request) bool {
	return IsServiceType(r, "rt")
}
