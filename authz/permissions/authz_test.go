package permissions

import "testing"

func TestPermissionBitmaskHelpers(t *testing.T) {
	mask := CommunityMembershipInvite | CommunityServerAdd | CommunityRolesView

	if !HasPermission(mask, CommunityServerAdd) {
		t.Fatalf("expected mask to include CommunityServerAdd")
	}
	if HasPermission(mask, CommunityRolesEdit) {
		t.Fatalf("did not expect mask to include CommunityRolesEdit")
	}

	if !HasAnyPermission(mask, CommunityRolesEdit, CommunityRolesView) {
		t.Fatalf("expected any permission check to pass")
	}
	if HasAnyPermission(mask, CommunityRolesEdit, CommunitySuspend) {
		t.Fatalf("did not expect any permission check to pass")
	}

	if !HasAllPermissions(mask, CommunityMembershipInvite, CommunityServerAdd) {
		t.Fatalf("expected all permissions check to pass")
	}
	if HasAllPermissions(mask, CommunityMembershipInvite, CommunityRolesEdit) {
		t.Fatalf("did not expect all permissions check to pass")
	}
}

func TestDenyAndEffectivePermissions(t *testing.T) {
	allow := TeamMembershipInvite | TeamServerAdd | TeamRolesView
	deny := TeamServerAdd

	effective := GetEffectivePermissions(allow, deny)
	if HasPermission(effective, TeamServerAdd) {
		t.Fatalf("expected denied permission to be removed")
	}
	if !HasPermission(effective, TeamMembershipInvite) {
		t.Fatalf("expected non-denied permission to remain")
	}

	if CanPerformAction(allow, deny, TeamServerAdd) {
		t.Fatalf("did not expect denied action to be allowed")
	}
	if !CanPerformAction(allow, deny, TeamMembershipInvite) {
		t.Fatalf("expected allowed action to pass")
	}
}

func TestRoleAndPermissionHelpers(t *testing.T) {
	if GetCommunityRolePermissions(RoleCommunityOwnerKey) == 0 {
		t.Fatalf("expected community owner role to have permissions")
	}
	if GetCommunityRolePermissions("unknown") != 0 {
		t.Fatalf("expected unknown role permissions to be zero")
	}

	if !IsCommunityRoleValid(RoleCommunityOwnerKey) {
		t.Fatalf("expected role %s to be valid", RoleCommunityOwnerKey)
	}
	if IsCommunityRoleValid("invalid") {
		t.Fatalf("expected invalid role to be rejected")
	}

	if GetCommunityRoleName("invalid") != "Unknown" {
		t.Fatalf("expected invalid role name fallback")
	}
	if GetCommunityPermissionName(CommunityMembershipInvite) != PermCommunityMembershipInvite {
		t.Fatalf("unexpected permission name for bit")
	}
	if GetCommunityPermissionName(1<<60) != "UNKNOWN_PERMISSION" {
		t.Fatalf("expected unknown permission fallback")
	}
}
