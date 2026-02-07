package permissions

// Lobby Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/lobby.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// LobbyMembershipView views membership details from the lobby interface.
	LobbyMembershipView uint64 = 1 << 0
	// LobbyMembershipInvite invites users to scopes lobby via the lobby interface.
	LobbyMembershipInvite uint64 = 1 << 1
	// LobbyMembershipDelete deletes memberships from the lobby interface.
	LobbyMembershipDelete uint64 = 1 << 2

	// LobbyKick removes users from the lobby.
	LobbyKick uint64 = 1 << 3
	// LobbyManage grants full lobby management (slots, mission, server).
	LobbyManage uint64 = 1 << 4
	// LobbyDisband disbands the entire lobby.
	LobbyDisband uint64 = 1 << 5
	// LobbyTransferOwnership transfers lobby ownership to another user.
	LobbyTransferOwnership uint64 = 1 << 6
)

// Used for CheckUserPermission() API calls and string-based permission checks
const (
	// PermLobbyMembershipView is the permission key for LOBBY__MEMBERSHIP_VIEW.
	PermLobbyMembershipView = "LOBBY__MEMBERSHIP_VIEW"
	// PermLobbyMembershipInvite is the permission key for LOBBY__MEMBERSHIP_INVITE.
	PermLobbyMembershipInvite = "LOBBY__MEMBERSHIP_INVITE"
	// PermLobbyMembershipDelete is the permission key for LOBBY__MEMBERSHIP_DELETE.
	PermLobbyMembershipDelete = "LOBBY__MEMBERSHIP_DELETE"
	// PermLobbyKick is the permission key for LOBBY__KICK.
	PermLobbyKick = "LOBBY__KICK"
	// PermLobbyManage is the permission key for LOBBY__MANAGE.
	PermLobbyManage = "LOBBY__MANAGE"
	// PermLobbyDisband is the permission key for LOBBY__DISBAND.
	PermLobbyDisband = "LOBBY__DISBAND"
	// PermLobbyTransferOwnership is the permission key for LOBBY__TRANSFER_OWNERSHIP.
	PermLobbyTransferOwnership = "LOBBY__TRANSFER_OWNERSHIP"
)

// ============================================
// Lobby Role String Constants
// ============================================
// String constants for lobby role identifiers

const (
	// RoleLobbyOwnerKey is the role key for a lobby owner.
	RoleLobbyOwnerKey = "lobby_owner"
	// RoleLobbyStaffKey is the role key for a lobby staff member.
	RoleLobbyStaffKey = "lobby_staff"
	// RoleLobbyUserKey is the role key for a lobby user.
	RoleLobbyUserKey = "lobby_user"
)

// Permission Groups - Matches seeds/permissions/lobby.json groups
var (
	// LobbyBasic - Basic lobby access
	LobbyBasic = LobbyMembershipView |
		LobbyMembershipInvite

	// LobbyStaff - Lobby staff permissions
	LobbyStaff = LobbyKick |
		LobbyManage |
		LobbyMembershipDelete

	// LobbyOwner - Full lobby control
	LobbyOwner = LobbyDisband |
		LobbyTransferOwnership
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleLobbyUser - Basic lobby participant (LobbyBasic group)
	RoleLobbyUser = LobbyBasic

	// RoleLobbyStaff - Lobby staff member (LobbyBasic + LobbyStaff groups)
	RoleLobbyStaff = LobbyBasic | LobbyStaff

	// RoleLobbyOwner - Lobby owner (LobbyBasic + LobbyStaff + LobbyOwner groups)
	RoleLobbyOwner = LobbyBasic | LobbyStaff | LobbyOwner
)

var lobbyRolePermissions = map[string]uint64{
	RoleLobbyOwnerKey: RoleLobbyOwner,
	RoleLobbyStaffKey: RoleLobbyStaff,
	RoleLobbyUserKey:  RoleLobbyUser,
}

var lobbyRoleNames = map[string]string{
	RoleLobbyOwnerKey: "Lobby Owner",
	RoleLobbyStaffKey: "Lobby Staff",
	RoleLobbyUserKey:  "Lobby User",
}

// GetLobbyRolePermissions returns the permission bitmask for a given lobby role
func GetLobbyRolePermissions(role string) uint64 {
	return getRolePermissions(role, lobbyRolePermissions)
}

// GetLobbyRoleName returns the human-readable name of a lobby role
func GetLobbyRoleName(role string) string {
	return getRoleName(role, lobbyRoleNames)
}

// IsLobbyRoleValid checks if a lobby role identifier is valid
func IsLobbyRoleValid(role string) bool {
	return isRoleValid(role, lobbyRoleNames)
}

// LobbyPermissionNames maps each permission bit to its key name (for debugging/logging)
var LobbyPermissionNames = map[uint64]string{
	LobbyMembershipView:    PermLobbyMembershipView,
	LobbyMembershipInvite:  PermLobbyMembershipInvite,
	LobbyMembershipDelete:  PermLobbyMembershipDelete,
	LobbyKick:              PermLobbyKick,
	LobbyManage:            PermLobbyManage,
	LobbyDisband:           PermLobbyDisband,
	LobbyTransferOwnership: PermLobbyTransferOwnership,
}

// LobbyPermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var LobbyPermissionKeyToBit = map[string]uint8{
	PermLobbyMembershipView:    0,
	PermLobbyMembershipInvite:  1,
	PermLobbyMembershipDelete:  2,
	PermLobbyKick:              3,
	PermLobbyManage:            4,
	PermLobbyDisband:           5,
	PermLobbyTransferOwnership: 6,
}

// GetLobbyPermissionName returns the key name of a single permission bit
func GetLobbyPermissionName(permission uint64) string {
	return getPermissionName(permission, LobbyPermissionNames)
}

// GetAllLobbyPermissionNames returns all permission key names in a bitmask
func GetAllLobbyPermissionNames(mask uint64) []string {
	return getAllPermissionNames(mask, LobbyPermissionNames)
}
