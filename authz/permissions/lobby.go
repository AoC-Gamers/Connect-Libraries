package permissions

// Lobby Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/lobby.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// Membership Management (bits 0-2)
	LOBBY__MEMBERSHIP_VIEW   uint64 = 1 << 0 // View membership details from the lobby interface
	LOBBY__MEMBERSHIP_INVITE uint64 = 1 << 1 // Invite users to scopes lobby via the lobby interface
	LOBBY__MEMBERSHIP_DELETE uint64 = 1 << 2 // Delete memberships from the lobby interface

	// Lobby Management (bits 3-6)
	LOBBY__KICK               uint64 = 1 << 3 // Remove users from the lobby
	LOBBY__MANAGE             uint64 = 1 << 4 // Full lobby management (slots, mission, server)
	LOBBY__DISBAND            uint64 = 1 << 5 // Disband the entire lobby
	LOBBY__TRANSFER_OWNERSHIP uint64 = 1 << 6 // Transfer lobby ownership to another user
)

// Lobby Permission String Constants
// Used for CheckUserPermission() API calls and string-based permission checks
const (
	PermLobbyMembershipView    = "LOBBY__MEMBERSHIP_VIEW"
	PermLobbyMembershipInvite  = "LOBBY__MEMBERSHIP_INVITE"
	PermLobbyMembershipDelete  = "LOBBY__MEMBERSHIP_DELETE"
	PermLobbyKick              = "LOBBY__KICK"
	PermLobbyManage            = "LOBBY__MANAGE"
	PermLobbyDisband           = "LOBBY__DISBAND"
	PermLobbyTransferOwnership = "LOBBY__TRANSFER_OWNERSHIP"
)

// ============================================
// Lobby Role String Constants
// ============================================
// String constants for lobby role identifiers

const (
	RoleLobbyOwnerKey = "lobby_owner"
	RoleLobbyStaffKey = "lobby_staff"
	RoleLobbyUserKey  = "lobby_user"
)

// Permission Groups - Matches seeds/permissions/lobby.json groups
var (
	// LOBBY__BASIC - Basic lobby access
	LOBBY__BASIC = LOBBY__MEMBERSHIP_VIEW |
		LOBBY__MEMBERSHIP_INVITE

	// LOBBY__STAFF - Lobby staff permissions
	LOBBY__STAFF = LOBBY__KICK |
		LOBBY__MANAGE |
		LOBBY__MEMBERSHIP_DELETE

	// LOBBY__OWNER - Full lobby control
	LOBBY__OWNER = LOBBY__DISBAND |
		LOBBY__TRANSFER_OWNERSHIP
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleLobbyUser - Basic lobby participant (LOBBY__BASIC group)
	RoleLobbyUser = LOBBY__BASIC

	// RoleLobbyStaff - Lobby staff member (LOBBY__BASIC + LOBBY__STAFF groups)
	RoleLobbyStaff = LOBBY__BASIC | LOBBY__STAFF

	// RoleLobbyOwner - Lobby owner (LOBBY__BASIC + LOBBY__STAFF + LOBBY__OWNER groups)
	RoleLobbyOwner = LOBBY__BASIC | LOBBY__STAFF | LOBBY__OWNER
)

// GetLobbyRolePermissions returns the permission bitmask for a given lobby role
func GetLobbyRolePermissions(role string) uint64 {
	switch role {
	case RoleLobbyOwnerKey:
		return RoleLobbyOwner
	case RoleLobbyStaffKey:
		return RoleLobbyStaff
	case RoleLobbyUserKey:
		return RoleLobbyUser
	default:
		return 0 // No permissions for unknown roles
	}
}

// GetLobbyRoleName returns the human-readable name of a lobby role
func GetLobbyRoleName(role string) string {
	switch role {
	case RoleLobbyOwnerKey:
		return "Lobby Owner"
	case RoleLobbyStaffKey:
		return "Lobby Staff"
	case RoleLobbyUserKey:
		return "Lobby User"
	default:
		return "Unknown"
	}
}

// IsLobbyRoleValid checks if a lobby role identifier is valid
func IsLobbyRoleValid(role string) bool {
	switch role {
	case RoleLobbyOwnerKey, RoleLobbyStaffKey, RoleLobbyUserKey:
		return true
	default:
		return false
	}
}

// LobbyPermissionNames maps each permission bit to its key name (for debugging/logging)
var LobbyPermissionNames = map[uint64]string{
	LOBBY__MEMBERSHIP_VIEW:    PermLobbyMembershipView,
	LOBBY__MEMBERSHIP_INVITE:  PermLobbyMembershipInvite,
	LOBBY__MEMBERSHIP_DELETE:  PermLobbyMembershipDelete,
	LOBBY__KICK:               PermLobbyKick,
	LOBBY__MANAGE:             PermLobbyManage,
	LOBBY__DISBAND:            PermLobbyDisband,
	LOBBY__TRANSFER_OWNERSHIP: PermLobbyTransferOwnership,
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
	if name, ok := LobbyPermissionNames[permission]; ok {
		return name
	}
	return "UNKNOWN_PERMISSION"
}

// GetAllLobbyPermissionNames returns all permission key names in a bitmask
func GetAllLobbyPermissionNames(mask uint64) []string {
	var names []string
	for perm, name := range LobbyPermissionNames {
		if mask&perm != 0 {
			names = append(names, name)
		}
	}
	return names
}
