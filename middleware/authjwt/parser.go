package authjwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// ParseAndValidate extrae y valida un JWT token retornando claims tipados
func ParseAndValidate(tokenStr, secret string, expectedPolicyVersion int) (*Claims, error) {
	if tokenStr == "" {
		log.Error().Msg("[JWT] Missing token")
		return nil, ErrMissingToken
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(tok *jwt.Token) (any, error) {
		// Validar el algoritmo de firma
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error().
				Str("algorithm", tok.Method.Alg()).
				Msg("[JWT] Unexpected signing method")
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		log.Error().
			Err(err).
			Bool("token_valid", token != nil && token.Valid).
			Msg("[JWT] Failed to parse token")
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		log.Error().Msg("[JWT] Token is invalid")
		return nil, ErrInvalidToken
	}

	// Extract steamID with compatibility fallback
	steamID := extractSteamID(claims)
	if steamID == "" {
		log.Error().
			Interface("claims", claims).
			Msg("[JWT] Missing steamID in claims")
		return nil, ErrMissingSteamID
	}

	// Validate policy version if present
	if err := validatePolicyVersion(claims, expectedPolicyVersion); err != nil {
		// Extraer la versión del token para el log
		tokenPolicyVersion := -1
		if pv, ok := claims["policy_version"].(float64); ok {
			tokenPolicyVersion = int(pv)
		}
		log.Error().
			Err(err).
			Int("expected", expectedPolicyVersion).
			Int("token_version", tokenPolicyVersion).
			Msg("[JWT] Policy version mismatch")
		return nil, err
	}

	// Extract role and permissions (bitmask)
	role := extractRole(claims)
	allowPermissions := extractAllowPermissions(claims)
	denyPermissions := extractDenyPermissions(claims)
	iat := extractIat(claims)
	exp := extractExp(claims)

	return &Claims{
		SteamID:          steamID,
		PolicyVersion:    expectedPolicyVersion,
		Role:             role,
		AllowPermissions: allowPermissions,
		DenyPermissions:  denyPermissions,
		IssuedAt:         iat,
		ExpiresAt:        exp,
	}, nil
}

// extractSteamID extrae steamID de los claims
func extractSteamID(claims jwt.MapClaims) string {
	if steamID, ok := claims["steamid"].(string); ok && steamID != "" {
		return steamID
	}
	return ""
}

// extractRole extrae el rol del usuario
func extractRole(claims jwt.MapClaims) string {
	if role, ok := claims["role"].(string); ok {
		return role
	}
	return "web_user" // Default role
}

// extractAllowPermissions extrae la máscara de permisos permitidos
func extractAllowPermissions(claims jwt.MapClaims) uint64 {
	if allowPermissions, ok := claims["allow_permissions"].(float64); ok {
		return uint64(allowPermissions)
	}
	return 0
}

// extractDenyPermissions extrae la máscara de permisos denegados
func extractDenyPermissions(claims jwt.MapClaims) uint64 {
	if denyPermissions, ok := claims["deny_permissions"].(float64); ok {
		return uint64(denyPermissions)
	}
	return 0
}

// extractIat extrae el timestamp de emisión
func extractIat(claims jwt.MapClaims) int64 {
	if iat, ok := claims["iat"].(float64); ok {
		return int64(iat)
	}
	return 0
}

// extractExp extrae el timestamp de expiración
func extractExp(claims jwt.MapClaims) int64 {
	if exp, ok := claims["exp"].(float64); ok {
		return int64(exp)
	}
	return 0
}

// validatePolicyVersion verifica que la versión de política coincida
func validatePolicyVersion(claims jwt.MapClaims, expectedVersion int) error {
	policyVersion, ok := claims["policy_version"].(float64)
	if !ok {
		log.Warn().
			Int("expected", expectedVersion).
			Interface("claims", claims).
			Msg("[JWT] Policy version not found in token, allowing for backward compatibility")
		return nil
	}

	tokenVersion := int(policyVersion)
	if tokenVersion != expectedVersion {
		log.Error().
			Int("expected", expectedVersion).
			Int("token_version", tokenVersion).
			Msg("[JWT] Policy version mismatch")
		return ErrPolicyVersionMismatch
	}

	return nil
}
