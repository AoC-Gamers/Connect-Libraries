package authjwt

import (
	"github.com/AoC-Gamers/connect-libraries/middleware/authconfig"
)

// AuthConfig configuración común para autenticación JWT
type AuthConfig struct {
	JWTSecret           string
	PolicyVersionGlobal int
}

// NewAuthConfig creates an AuthConfig with unified POLICY_VERSION from environment
// Uses LoadPolicyVersion() to read from POLICY_VERSION env var
func NewAuthConfig(jwtSecret string) AuthConfig {
	return AuthConfig{
		JWTSecret:           jwtSecret,
		PolicyVersionGlobal: authconfig.LoadPolicyVersion(),
	}
}

// NewAuthConfigWithFallback creates an AuthConfig with fallback support for legacy variables
// Useful during migration period to support both old and new variable names
func NewAuthConfigWithFallback(jwtSecret string, fallbackVars ...string) AuthConfig {
	return AuthConfig{
		JWTSecret:           jwtSecret,
		PolicyVersionGlobal: authconfig.LoadPolicyVersionWithFallback(fallbackVars...),
	}
}

// Claims representa los claims optimizados del JWT
// Solo contiene información de autenticación y autorización WEB
type Claims struct {
	// Identificación
	SteamID       string `json:"steamid"`
	PolicyVersion int    `json:"policy_version"`

	// Autorización WEB (único scope incluido en JWT)
	Role             string `json:"role"`              // Rol para display/logging
	AllowPermissions uint64 `json:"allow_permissions"` // Permisos permitidos (bitmask)
	DenyPermissions  uint64 `json:"deny_permissions"`  // Permisos denegados (bitmask)

	// JWT standard claims
	IssuedAt  int64 `json:"iat"` // Unix timestamp
	ExpiresAt int64 `json:"exp"` // Unix timestamp
}

// ValidationError tipos de errores de validación
type ValidationError struct {
	Type    string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Error types
var (
	ErrMissingToken          = &ValidationError{"missing_token", "missing authorization token"}
	ErrInvalidToken          = &ValidationError{"invalid_token", "invalid or expired token"}
	ErrMissingSteamID        = &ValidationError{"missing_steamid", "invalid token claims: missing steamid"}
	ErrPolicyVersionMismatch = &ValidationError{"policy_mismatch", "token policy version mismatch, please re-authenticate"}
)