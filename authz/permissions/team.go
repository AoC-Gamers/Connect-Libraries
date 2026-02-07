package permissions

// Team Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/team.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// TeamMembershipInvite invites users to the team (membership invite).
	TeamMembershipInvite uint64 = 1 << 0
	// TeamMembershipDelete removes memberships from the team.
	TeamMembershipDelete uint64 = 1 << 1

	// TeamServerAdd adds a server to the team's server list.
	TeamServerAdd uint64 = 1 << 2
	// TeamServerEdit edits an existing server in the team's server list.
	TeamServerEdit uint64 = 1 << 3
	// TeamServerDelete deletes a server from the team's server list.
	TeamServerDelete uint64 = 1 << 4

	// TeamMissionListAdd adds a new mission to the team's mission list.
	TeamMissionListAdd uint64 = 1 << 5
	// TeamMissionListEdit edits an existing mission in the team's mission list.
	TeamMissionListEdit uint64 = 1 << 6

	// TeamGamemodeListAdd adds a new gamemode to the team's gamemode list.
	TeamGamemodeListAdd uint64 = 1 << 7
	// TeamGamemodeListEdit edits an existing gamemode in the team's gamemode list.
	TeamGamemodeListEdit uint64 = 1 << 8

	// TeamInfoEdit edits team information.
	TeamInfoEdit uint64 = 1 << 9
	// TeamAnalytics views team statistics and analytics.
	TeamAnalytics uint64 = 1 << 10
	// TeamTransferOwnership transfers team ownership.
	TeamTransferOwnership uint64 = 1 << 11
	// TeamSuspend suspends or unsuspends this team.
	TeamSuspend uint64 = 1 << 12
	// TeamAuditView views audit logs for this team.
	TeamAuditView uint64 = 1 << 13
	// TeamLobbyCreate creates a lobby associated with the team.
	TeamLobbyCreate uint64 = 1 << 14

	// TeamRolesView views roles and permissions for this team.
	TeamRolesView uint64 = 1 << 15
	// TeamRolesEdit edits roles and permissions for this team.
	TeamRolesEdit uint64 = 1 << 16
)

// ============================================
// Team Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods

const (
	// PermTeamMembershipInvite is the permission key for TEAM__MEMBERSHIP_INVITE.
	PermTeamMembershipInvite = "TEAM__MEMBERSHIP_INVITE"
	// PermTeamMembershipDelete is the permission key for TEAM__MEMBERSHIP_DELETE.
	PermTeamMembershipDelete = "TEAM__MEMBERSHIP_DELETE"

	// PermTeamServerAdd is the permission key for TEAM__SERVER_ADD.
	PermTeamServerAdd = "TEAM__SERVER_ADD"
	// PermTeamServerEdit is the permission key for TEAM__SERVER_EDIT.
	PermTeamServerEdit = "TEAM__SERVER_EDIT"
	// PermTeamServerDelete is the permission key for TEAM__SERVER_DELETE.
	PermTeamServerDelete = "TEAM__SERVER_DELETE"

	// PermTeamMissionListAdd is the permission key for TEAM__MISSIONLIST_ADD.
	PermTeamMissionListAdd = "TEAM__MISSIONLIST_ADD"
	// PermTeamMissionListEdit is the permission key for TEAM__MISSIONLIST_EDIT.
	PermTeamMissionListEdit = "TEAM__MISSIONLIST_EDIT"

	// PermTeamGamemodeListAdd is the permission key for TEAM__GAMEMODELIST_ADD.
	PermTeamGamemodeListAdd = "TEAM__GAMEMODELIST_ADD"
	// PermTeamGamemodeListEdit is the permission key for TEAM__GAMEMODELIST_EDIT.
	PermTeamGamemodeListEdit = "TEAM__GAMEMODELIST_EDIT"

	// PermTeamInfoEdit is the permission key for TEAM__INFO_EDIT.
	PermTeamInfoEdit = "TEAM__INFO_EDIT"
	// PermTeamAnalytics is the permission key for TEAM__ANALYTICS.
	PermTeamAnalytics = "TEAM__ANALYTICS"
	// PermTeamTransferOwnership is the permission key for TEAM__TRANSFER_OWNERSHIP.
	PermTeamTransferOwnership = "TEAM__TRANSFER_OWNERSHIP"
	// PermTeamSuspend is the permission key for TEAM__SUSPEND.
	PermTeamSuspend = "TEAM__SUSPEND"
	// PermTeamAuditView is the permission key for TEAM__AUDIT_VIEW.
	PermTeamAuditView = "TEAM__AUDIT_VIEW"
	// PermTeamLobbyCreate is the permission key for TEAM__LOBBY_CREATE.
	PermTeamLobbyCreate = "TEAM__LOBBY_CREATE"

	// PermTeamRolesView is the permission key for TEAM__ROLES_VIEW.
	PermTeamRolesView = "TEAM__ROLES_VIEW"
	// PermTeamRolesEdit is the permission key for TEAM__ROLES_EDIT.
	PermTeamRolesEdit = "TEAM__ROLES_EDIT"
)

// ============================================
// Team Role String Constants
// ============================================
// String constants for team role identifiers

const (
	// RoleTeamOwnerKey is the role key for a team owner.
	RoleTeamOwnerKey = "team_owner"
	// RoleTeamStaffKey is the role key for a team staff member.
	RoleTeamStaffKey = "team_staff"
	// RoleTeamUserKey is the role key for a team user.
	RoleTeamUserKey = "team_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/team.json groups
var (
	// TeamBasic - Basic team membership
	TeamBasic uint64 = 0

	// TeamStaff - Team staff permissions
	TeamStaff = TeamMembershipInvite |
		TeamServerAdd |
		TeamServerEdit |
		TeamServerDelete |
		TeamMissionListAdd |
		TeamMissionListEdit |
		TeamGamemodeListAdd |
		TeamGamemodeListEdit |
		TeamLobbyCreate |
		TeamRolesView

	// TeamOwner - Full team control
	TeamOwner = TeamMembershipDelete |
		TeamInfoEdit |
		TeamAnalytics |
		TeamTransferOwnership |
		TeamSuspend |
		TeamAuditView |
		TeamRolesEdit
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleTeamUser - Basic team member (TeamBasic group)
	RoleTeamUser = TeamBasic

	// RoleTeamStaff - Team staff member (TeamBasic + TeamStaff groups)
	RoleTeamStaff = TeamBasic | TeamStaff

	// RoleTeamOwner - Team owner (TeamBasic + TeamStaff + TeamOwner groups)
	RoleTeamOwner = TeamBasic | TeamStaff | TeamOwner
)

var teamRolePermissions = map[string]uint64{
	RoleTeamOwnerKey: RoleTeamOwner,
	RoleTeamStaffKey: RoleTeamStaff,
	RoleTeamUserKey:  RoleTeamUser,
}

var teamRoleNames = map[string]string{
	RoleTeamOwnerKey: "Team Owner",
	RoleTeamStaffKey: "Team Staff",
	RoleTeamUserKey:  "Team User",
}

// GetTeamRolePermissions returns the permission bitmask for a given team role
func GetTeamRolePermissions(role string) uint64 {
	return getRolePermissions(role, teamRolePermissions)
}

// GetTeamRoleName returns the human-readable name of a team role
func GetTeamRoleName(role string) string {
	return getRoleName(role, teamRoleNames)
}

// IsTeamRoleValid checks if a team role identifier is valid
func IsTeamRoleValid(role string) bool {
	return isRoleValid(role, teamRoleNames)
}

// TeamPermissionNames maps each permission bit to its key name (for debugging/logging)
var TeamPermissionNames = map[uint64]string{
	TeamMembershipInvite:  PermTeamMembershipInvite,
	TeamMembershipDelete:  PermTeamMembershipDelete,
	TeamServerAdd:         PermTeamServerAdd,
	TeamServerEdit:        PermTeamServerEdit,
	TeamServerDelete:      PermTeamServerDelete,
	TeamMissionListAdd:    PermTeamMissionListAdd,
	TeamMissionListEdit:   PermTeamMissionListEdit,
	TeamGamemodeListAdd:   PermTeamGamemodeListAdd,
	TeamGamemodeListEdit:  PermTeamGamemodeListEdit,
	TeamInfoEdit:          PermTeamInfoEdit,
	TeamAnalytics:         PermTeamAnalytics,
	TeamTransferOwnership: PermTeamTransferOwnership,
	TeamSuspend:           PermTeamSuspend,
	TeamAuditView:         PermTeamAuditView,
	TeamLobbyCreate:       PermTeamLobbyCreate,
	TeamRolesView:         PermTeamRolesView,
	TeamRolesEdit:         PermTeamRolesEdit,
}

// TeamPermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var TeamPermissionKeyToBit = map[string]uint8{
	PermTeamMembershipInvite:  0,
	PermTeamMembershipDelete:  1,
	PermTeamServerAdd:         2,
	PermTeamServerEdit:        3,
	PermTeamServerDelete:      4,
	PermTeamMissionListAdd:    5,
	PermTeamMissionListEdit:   6,
	PermTeamGamemodeListAdd:   7,
	PermTeamGamemodeListEdit:  8,
	PermTeamInfoEdit:          9,
	PermTeamAnalytics:         10,
	PermTeamTransferOwnership: 11,
	PermTeamSuspend:           12,
	PermTeamAuditView:         13,
	PermTeamLobbyCreate:       14,
	PermTeamRolesView:         15,
	PermTeamRolesEdit:         16,
}

// GetTeamPermissionName returns the key name of a single permission bit
func GetTeamPermissionName(permission uint64) string {
	return getPermissionName(permission, TeamPermissionNames)
}

// GetAllTeamPermissionNames returns all permission key names in a bitmask
func GetAllTeamPermissionNames(mask uint64) []string {
	return getAllPermissionNames(mask, TeamPermissionNames)
}
