package authjwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestClaimsHelpers(t *testing.T) {
	const readPermission uint64 = 1 << 1
	const writePermission uint64 = 1 << 2

	claims := &Claims{
		SteamID:          "76561198000000000",
		Role:             "web_admin",
		AllowPermissions: readPermission | writePermission,
		DenyPermissions:  writePermission,
	}

	if claims.GetSteamID() == "" {
		t.Fatal("expected steam id")
	}
	if claims.GetRole() != "web_admin" {
		t.Fatalf("expected role web_admin, got %s", claims.GetRole())
	}
	if !claims.HasPermission(readPermission) {
		t.Fatal("expected read permission allowed")
	}
	if claims.HasPermission(writePermission) {
		t.Fatal("expected denied permission to override allow")
	}
	if !claims.IsAdmin() || !claims.IsStaff() || !claims.IsModerator() {
		t.Fatal("expected admin to be moderator and staff")
	}
	if claims.IsOwner() {
		t.Fatal("expected admin to not be owner")
	}
}

func TestParseAndValidate(t *testing.T) {
	secret := "secret-key"
	now := time.Now().Unix()

	if _, err := ParseAndValidate("", secret, 1); err != ErrMissingToken {
		t.Fatalf("expected ErrMissingToken, got %v", err)
	}

	if _, err := ParseAndValidate("not-a-jwt", secret, 1); err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	tokenWithoutSteamID := mustSignToken(t, secret, jwt.MapClaims{
		"policy_version": 1,
		"exp":            now + 3600,
	})
	if _, err := ParseAndValidate(tokenWithoutSteamID, secret, 1); err != ErrMissingSteamID {
		t.Fatalf("expected ErrMissingSteamID, got %v", err)
	}

	mismatchToken := mustSignToken(t, secret, jwt.MapClaims{
		"steamid":        "76561198000000000",
		"policy_version": 2,
		"exp":            now + 3600,
	})
	if _, err := ParseAndValidate(mismatchToken, secret, 1); err != ErrPolicyVersionMismatch {
		t.Fatalf("expected ErrPolicyVersionMismatch, got %v", err)
	}

	validToken := mustSignToken(t, secret, jwt.MapClaims{
		"steamid":           "76561198000000001",
		"role":              "web_owner",
		"policy_version":    3,
		"allow_permissions": float64(12),
		"deny_permissions":  float64(8),
		"iat":               float64(now),
		"exp":               float64(now + 3600),
	})

	validated, err := ParseAndValidate(validToken, secret, 3)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}

	if validated.GetSteamID() != "76561198000000001" {
		t.Fatalf("unexpected steam id %s", validated.GetSteamID())
	}
	if validated.GetRole() != "web_owner" {
		t.Fatalf("unexpected role %s", validated.GetRole())
	}
	if validated.AllowPermissions != 12 || validated.DenyPermissions != 8 {
		t.Fatalf("unexpected permissions allow=%d deny=%d", validated.AllowPermissions, validated.DenyPermissions)
	}
	if validated.PolicyVersion != 3 {
		t.Fatalf("unexpected policy version %d", validated.PolicyVersion)
	}

	legacyTokenNoPolicy := mustSignToken(t, secret, jwt.MapClaims{
		"steamid": "76561198000000002",
		"exp":     float64(now + 3600),
	})
	legacyValidated, legacyErr := ParseAndValidate(legacyTokenNoPolicy, secret, 9)
	if legacyErr != nil {
		t.Fatalf("expected compatibility for missing policy version, got %v", legacyErr)
	}
	if legacyValidated.GetRole() != "web_user" {
		t.Fatalf("expected default role web_user, got %s", legacyValidated.GetRole())
	}
}

func TestNewAuthConfig(t *testing.T) {
	t.Setenv("POLICY_VERSION", "11")
	config := NewAuthConfig("jwt-secret")

	if config.SignerMaterial != "jwt-secret" {
		t.Fatalf("unexpected signer material %s", config.SignerMaterial)
	}
	if config.PolicyVersionGlobal != 11 {
		t.Fatalf("expected policy 11, got %d", config.PolicyVersionGlobal)
	}

	t.Setenv("POLICY_VERSION", "")
	t.Setenv("AUTHZ_POLICY_VERSION", "13")
	fallback := NewAuthConfigWithFallback("jwt-secret-2", "AUTHZ_POLICY_VERSION")
	if fallback.PolicyVersionGlobal != 13 {
		t.Fatalf("expected fallback policy 13, got %d", fallback.PolicyVersionGlobal)
	}
}

func mustSignToken(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("could not sign token: %v", err)
	}
	return signed
}
