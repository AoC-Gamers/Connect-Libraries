package authconfig

import "testing"

func TestLoadPolicyVersion(t *testing.T) {
	t.Setenv("POLICY_VERSION", "")
	if got := LoadPolicyVersion(); got != DefaultPolicyVersion {
		t.Fatalf("expected default %d, got %d", DefaultPolicyVersion, got)
	}

	t.Setenv("POLICY_VERSION", "3")
	if got := LoadPolicyVersion(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}

	t.Setenv("POLICY_VERSION", "0")
	if got := LoadPolicyVersion(); got != DefaultPolicyVersion {
		t.Fatalf("expected default for invalid value, got %d", got)
	}
}

func TestLoadPolicyVersionWithFallback(t *testing.T) {
	t.Setenv("POLICY_VERSION", "")
	t.Setenv("AUTHZ_POLICY_VERSION", "7")
	t.Setenv("AUTH_POLICY_VERSION", "9")

	if got := LoadPolicyVersionWithFallback("AUTHZ_POLICY_VERSION", "AUTH_POLICY_VERSION"); got != 7 {
		t.Fatalf("expected first fallback value 7, got %d", got)
	}

	t.Setenv("POLICY_VERSION", "5")
	if got := LoadPolicyVersionWithFallback("AUTHZ_POLICY_VERSION", "AUTH_POLICY_VERSION"); got != 5 {
		t.Fatalf("expected primary value 5, got %d", got)
	}
}

func TestValidateAndConfigs(t *testing.T) {
	if !ValidatePolicyVersion(2, 2) {
		t.Fatal("expected matching versions to be valid")
	}
	if ValidatePolicyVersion(2, 3) {
		t.Fatal("expected mismatched versions to be invalid")
	}

	t.Setenv("POLICY_VERSION", "4")
	config := GetPolicyConfig()
	if config.PolicyVersion != 4 {
		t.Fatalf("expected policy config version 4, got %d", config.PolicyVersion)
	}

	t.Setenv("POLICY_VERSION", "")
	t.Setenv("LEGACY_POLICY", "6")
	fallbackConfig := GetPolicyConfigWithFallback("LEGACY_POLICY")
	if fallbackConfig.PolicyVersion != 6 {
		t.Fatalf("expected fallback config version 6, got %d", fallbackConfig.PolicyVersion)
	}
}
