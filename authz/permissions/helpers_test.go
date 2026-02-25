package permissions

import "testing"

func TestGetAllPermissionNamesSorting(t *testing.T) {
	mask := CommunityRolesEdit | CommunityMembershipInvite | CommunityInfoEdit
	names := GetAllCommunityPermissionNames(mask)

	if len(names) != 3 {
		t.Fatalf("expected 3 names, got %d", len(names))
	}

	// Orden esperado por bit ascendente: invite (bit 0), info_edit (bit 9), roles_edit (bit 15)
	if names[0] != PermCommunityMembershipInvite {
		t.Fatalf("unexpected first permission: %s", names[0])
	}
	if names[1] != PermCommunityInfoEdit {
		t.Fatalf("unexpected second permission: %s", names[1])
	}
	if names[2] != PermCommunityRolesEdit {
		t.Fatalf("unexpected third permission: %s", names[2])
	}
}

func TestGetRoleNameAndValidity(t *testing.T) {
	if !IsTeamRoleValid(RoleTeamOwnerKey) {
		t.Fatalf("expected team owner role to be valid")
	}
	if IsTeamRoleValid("team_invalid") {
		t.Fatalf("did not expect invalid team role")
	}

	if GetTeamRoleName("team_invalid") != "Unknown" {
		t.Fatalf("expected unknown fallback for invalid role")
	}
}
