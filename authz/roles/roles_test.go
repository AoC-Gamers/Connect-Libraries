package roles

import "testing"

const (
	testScopeCommunity = "COMMUNITY"
	testScopeUnknown   = "UNKNOWN"
)

func TestHasRole(t *testing.T) {
	roles := []string{COMMUNITY_OWNER, TEAM_OWNER, WEB_OWNER}

	if !HasRole(COMMUNITY_OWNER, roles) {
		t.Fatalf("expected role %s to be found", COMMUNITY_OWNER)
	}
	if HasRole("nonexistent", roles) {
		t.Fatalf("did not expect unknown role to be found")
	}
}

func TestGetRole(t *testing.T) {
	role, ok := GetRole(testScopeCommunity, COMMUNITY_OWNER)
	if !ok {
		t.Fatalf("expected role in scope %s", testScopeCommunity)
	}
	if role.Scope != testScopeCommunity {
		t.Fatalf("unexpected role scope: %s", role.Scope)
	}
	if role.Name != COMMUNITY_OWNER {
		t.Fatalf("unexpected role name: %s", role.Name)
	}

	if _, ok = GetRole(testScopeCommunity, "invalid"); ok {
		t.Fatalf("did not expect invalid role to resolve")
	}
	if _, ok = GetRole(testScopeUnknown, COMMUNITY_OWNER); ok {
		t.Fatalf("did not expect unknown scope to resolve")
	}
}

func TestGetAllRoles(t *testing.T) {
	communityRoles := GetAllRoles(testScopeCommunity)
	if len(communityRoles) == 0 {
		t.Fatalf("expected non-empty roles for scope %s", testScopeCommunity)
	}

	unknownRoles := GetAllRoles(testScopeUnknown)
	if len(unknownRoles) != 0 {
		t.Fatalf("expected empty roles for unknown scope")
	}
}
