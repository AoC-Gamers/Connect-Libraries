package permissions

// Team Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/team.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// Membership Management (bits 0-2)
	TEAM__MEMBERSHIP_VIEW   uint64 = 1 << 0 // View membership details from the team interface
	TEAM__MEMBERSHIP_INVITE uint64 = 1 << 1 // Invite users to the team (membership invite)
	TEAM__MEMBERSHIP_DELETE uint64 = 1 << 2 // Remove memberships from the team

	// Team Management (bits 3-7)
	TEAM__INFO_EDIT          uint64 = 1 << 3 // Edit team information
	TEAM__LOBBY_CREATE       uint64 = 1 << 4 // Create a lobby associated with the team
	TEAM__TRANSFER_OWNERSHIP uint64 = 1 << 5 // Transfer team ownership
	TEAM__SUSPEND            uint64 = 1 << 6 // Suspend or unsuspend this team
	TEAM__AUDIT_VIEW         uint64 = 1 << 7 // View audit logs for this team
)

// ============================================
// Team Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods

const (
	// Membership Management
	PermTeamMembershipView   = "TEAM__MEMBERSHIP_VIEW"
	PermTeamMembershipInvite = "TEAM__MEMBERSHIP_INVITE"
	PermTeamMembershipDelete = "TEAM__MEMBERSHIP_DELETE"

	// Team Management
	PermTeamInfoEdit          = "TEAM__INFO_EDIT"
	PermTeamLobbyCreate       = "TEAM__LOBBY_CREATE"
	PermTeamTransferOwnership = "TEAM__TRANSFER_OWNERSHIP"
	PermTeamSuspend           = "TEAM__SUSPEND"
	PermTeamAuditView         = "TEAM__AUDIT_VIEW"
)

// ============================================
// Team Role String Constants
// ============================================
// String constants for team role identifiers

const (
	RoleTeamOwnerKey = "team_owner"
	RoleTeamStaffKey = "team_staff"
	RoleTeamUserKey  = "team_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/team.json groups
var (
	// TEAM__BASIC - Basic team membership
	TEAM__BASIC = TEAM__MEMBERSHIP_VIEW

	// TEAM__STAFF - Team staff permissions
	TEAM__STAFF = TEAM__MEMBERSHIP_INVITE |
		TEAM__LOBBY_CREATE

	// TEAM__OWNER - Full team control
	TEAM__OWNER = TEAM__MEMBERSHIP_DELETE |
		TEAM__INFO_EDIT |
		TEAM__TRANSFER_OWNERSHIP |
		TEAM__SUSPEND |
		TEAM__AUDIT_VIEW
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleTeamUser - Basic team member (TEAM__BASIC group)
	RoleTeamUser = TEAM__BASIC

	// RoleTeamStaff - Team staff member (TEAM__BASIC + TEAM__STAFF groups)
	RoleTeamStaff = TEAM__BASIC | TEAM__STAFF

	// RoleTeamOwner - Team owner (TEAM__BASIC + TEAM__STAFF + TEAM__OWNER groups)
	RoleTeamOwner = TEAM__BASIC | TEAM__STAFF | TEAM__OWNER
)

// GetTeamRolePermissions returns the permission bitmask for a given team role
func GetTeamRolePermissions(role string) uint64 {
	switch role {
	case RoleTeamOwnerKey:
		return RoleTeamOwner
	case RoleTeamStaffKey:
		return RoleTeamStaff
	case RoleTeamUserKey:
		return RoleTeamUser
	default:
		return 0 // No permissions for unknown roles
	}
}

// GetTeamRoleName returns the human-readable name of a team role
func GetTeamRoleName(role string) string {
	switch role {
	case RoleTeamOwnerKey:
		return "Team Owner"
	case RoleTeamStaffKey:
		return "Team Staff"
	case RoleTeamUserKey:
		return "Team User"
	default:
		return "Unknown"
	}
}

// IsTeamRoleValid checks if a team role identifier is valid
func IsTeamRoleValid(role string) bool {
	switch role {
	case RoleTeamOwnerKey, RoleTeamStaffKey, RoleTeamUserKey:
		return true
	default:
		return false
	}
}

// TeamPermissionNames maps each permission bit to its key name (for debugging/logging)
var TeamPermissionNames = map[uint64]string{
	TEAM__MEMBERSHIP_VIEW:    PermTeamMembershipView,
	TEAM__MEMBERSHIP_INVITE:  PermTeamMembershipInvite,
	TEAM__MEMBERSHIP_DELETE:  PermTeamMembershipDelete,
	TEAM__INFO_EDIT:          PermTeamInfoEdit,
	TEAM__LOBBY_CREATE:       PermTeamLobbyCreate,
	TEAM__TRANSFER_OWNERSHIP: PermTeamTransferOwnership,
	TEAM__SUSPEND:            PermTeamSuspend,
}

// TeamPermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var TeamPermissionKeyToBit = map[string]uint8{
	PermTeamMembershipView:    0,
	PermTeamMembershipInvite:  1,
	PermTeamMembershipDelete:  2,
	PermTeamInfoEdit:          3,
	PermTeamLobbyCreate:       4,
	PermTeamTransferOwnership: 5,
	PermTeamSuspend:           6,
}

// GetTeamPermissionName returns the key name of a single permission bit
func GetTeamPermissionName(permission uint64) string {
	if name, ok := TeamPermissionNames[permission]; ok {
		return name
	}
	return "UNKNOWN_PERMISSION"
}

// GetAllTeamPermissionNames returns all permission key names in a bitmask
func GetAllTeamPermissionNames(mask uint64) []string {
	var names []string
	for perm, name := range TeamPermissionNames {
		if mask&perm != 0 {
			names = append(names, name)
		}
	}
	return names
}
