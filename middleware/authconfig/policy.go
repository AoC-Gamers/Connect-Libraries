package authconfig

import (
	"os"
	"strconv"
)

// PolicyConfig holds policy version configuration
// Policy version is used to ensure JWT tokens are validated against the correct
// authorization policy. When policy changes (permissions, roles, etc.), increment
// this version to invalidate old tokens.
type PolicyConfig struct {
	PolicyVersion int
}

// DefaultPolicyVersion is the fallback version when no configuration is found
const DefaultPolicyVersion = 1

// LoadPolicyVersion reads POLICY_VERSION from environment variable
// Falls back to DefaultPolicyVersion (1) if not set or invalid
//
// Environment variable:
//
//	POLICY_VERSION - Integer value (e.g., 1, 2, 3)
//
// Returns:
//   - policy version as int
//   - always returns a valid version (minimum 1)
func LoadPolicyVersion() int {
	versionStr := os.Getenv("POLICY_VERSION")
	if versionStr == "" {
		return DefaultPolicyVersion
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil || version < 1 {
		return DefaultPolicyVersion
	}

	return version
}

// LoadPolicyVersionWithFallback reads POLICY_VERSION from environment,
// with a custom fallback list for backward compatibility
//
// Tries in order:
//  1. POLICY_VERSION (new unified variable)
//  2. First fallback variable
//  3. Second fallback variable (if provided)
//  4. DefaultPolicyVersion (1)
//
// Example:
//
//	version := LoadPolicyVersionWithFallback("AUTHZ_POLICY_VERSION", "AUTH_POLICY_VERSION")
func LoadPolicyVersionWithFallback(fallbackVars ...string) int {
	// Try primary variable first
	if version := LoadPolicyVersion(); version != DefaultPolicyVersion || os.Getenv("POLICY_VERSION") != "" {
		return version
	}

	// Try fallback variables in order
	for _, varName := range fallbackVars {
		if versionStr := os.Getenv(varName); versionStr != "" {
			if version, err := strconv.Atoi(versionStr); err == nil && version >= 1 {
				return version
			}
		}
	}

	return DefaultPolicyVersion
}

// ValidatePolicyVersion checks if token policy version matches expected version
// Returns true if versions match, false otherwise
func ValidatePolicyVersion(tokenVersion, expectedVersion int) bool {
	return tokenVersion == expectedVersion
}

// GetPolicyConfig creates a PolicyConfig from environment
func GetPolicyConfig() PolicyConfig {
	return PolicyConfig{
		PolicyVersion: LoadPolicyVersion(),
	}
}

// GetPolicyConfigWithFallback creates a PolicyConfig with fallback support
func GetPolicyConfigWithFallback(fallbackVars ...string) PolicyConfig {
	return PolicyConfig{
		PolicyVersion: LoadPolicyVersionWithFallback(fallbackVars...),
	}
}