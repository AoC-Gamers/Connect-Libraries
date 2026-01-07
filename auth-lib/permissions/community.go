package permissions

// ============================================
// Community Permission Bitmask Constants
// ============================================
// Generated from Connect-Auth/seeds/permissions/community.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// Membership Management (bits 0-1)
	COMMUNITY__MEMBERSHIP_INVITE uint64 = 1 << 0 // Invite users to the community (membership invite)
	COMMUNITY__MEMBERSHIP_DELETE uint64 = 1 << 1 // Remove memberships from the community

	// Server Management (bits 2-4)
	COMMUNITY__SERVER_ADD    uint64 = 1 << 2 // Add a server to the community's server list
	COMMUNITY__SERVER_EDIT   uint64 = 1 << 3 // Edit an existing server in the community's server list
	COMMUNITY__SERVER_DELETE uint64 = 1 << 4 // Delete a server from the community's server list

	// Mission List Management (bits 5-6)
	COMMUNITY__MISSIONLIST_ADD  uint64 = 1 << 5 // Add a new mission to the community's mission list
	COMMUNITY__MISSIONLIST_EDIT uint64 = 1 << 6 // Edit an existing mission in the community's mission list

	// Gamemode List Management (bits 7-8)
	COMMUNITY__GAMEMODELIST_ADD  uint64 = 1 << 7 // Add a new gamemode to the community's gamemode list
	COMMUNITY__GAMEMODELIST_EDIT uint64 = 1 << 8 // Edit an existing gamemode in the community's gamemode list

	// Community Management (bits 9-12)
	COMMUNITY__INFO_EDIT          uint64 = 1 << 9  // Edit community information
	COMMUNITY__ANALYTICS          uint64 = 1 << 10 // View community statistics and analytics
	COMMUNITY__TRANSFER_OWNERSHIP uint64 = 1 << 11 // Transfer community ownership
	COMMUNITY__SUSPEND            uint64 = 1 << 12 // Suspend or unsuspend this community
)

// ============================================
// Community Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods
// These match the permission keys in Connect-Auth

const (
	// Membership Management
	PermCommunityMembershipInvite = "COMMUNITY__MEMBERSHIP_INVITE"
	PermCommunityMembershipDelete = "COMMUNITY__MEMBERSHIP_DELETE"

	// Server Management
	PermCommunityServerAdd    = "COMMUNITY__SERVER_ADD"
	PermCommunityServerEdit   = "COMMUNITY__SERVER_EDIT"
	PermCommunityServerDelete = "COMMUNITY__SERVER_DELETE"

	// Mission List Management
	PermCommunityMissionlistAdd  = "COMMUNITY__MISSIONLIST_ADD"
	PermCommunityMissionlistEdit = "COMMUNITY__MISSIONLIST_EDIT"

	// Gamemode List Management
	PermCommunityGamemodelistAdd  = "COMMUNITY__GAMEMODELIST_ADD"
	PermCommunityGamemodelistEdit = "COMMUNITY__GAMEMODELIST_EDIT"

	// Community Management
	PermCommunityInfoEdit          = "COMMUNITY__INFO_EDIT"
	PermCommunityAnalytics         = "COMMUNITY__ANALYTICS"
	PermCommunityTransferOwnership = "COMMUNITY__TRANSFER_OWNERSHIP"
	PermCommunitySuspend           = "COMMUNITY__SUSPEND"
)

// ============================================
// Community Role String Constants
// ============================================
// String constants for community role identifiers

const (
	RoleCommunityOwnerKey = "community_owner"
	RoleCommunityStaffKey = "community_staff"
	RoleCommunityUserKey  = "community_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/community.json groups
var (
	// COMMUNITY__BASIC - Basic community membership (no permissions)
	COMMUNITY__BASIC uint64 = 0

	// COMMUNITY__STAFF - Staff management permissions
	COMMUNITY__STAFF = COMMUNITY__SERVER_ADD |
		COMMUNITY__SERVER_EDIT |
		COMMUNITY__SERVER_DELETE |
		COMMUNITY__MISSIONLIST_ADD |
		COMMUNITY__MISSIONLIST_EDIT |
		COMMUNITY__GAMEMODELIST_ADD |
		COMMUNITY__GAMEMODELIST_EDIT

	// COMMUNITY__OWNER - Full community control
	COMMUNITY__OWNER = COMMUNITY__MEMBERSHIP_INVITE |
		COMMUNITY__MEMBERSHIP_DELETE |
		COMMUNITY__INFO_EDIT |
		COMMUNITY__ANALYTICS |
		COMMUNITY__TRANSFER_OWNERSHIP |
		COMMUNITY__SUSPEND
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleCommunityUser - Basic community member (COMMUNITY__BASIC group - no permissions)
	RoleCommunityUser = COMMUNITY__BASIC

	// RoleCommunityStaff - Community staff member (COMMUNITY__BASIC + COMMUNITY__STAFF groups)
	RoleCommunityStaff = COMMUNITY__BASIC | COMMUNITY__STAFF

	// RoleCommunityOwner - Community owner (COMMUNITY__BASIC + COMMUNITY__STAFF + COMMUNITY__OWNER groups)
	RoleCommunityOwner = COMMUNITY__BASIC | COMMUNITY__STAFF | COMMUNITY__OWNER
)

// GetCommunityRolePermissions returns the permission bitmask for a given community role
func GetCommunityRolePermissions(role string) uint64 {
	switch role {
	case RoleCommunityOwnerKey:
		return RoleCommunityOwner
	case RoleCommunityStaffKey:
		return RoleCommunityStaff
	case RoleCommunityUserKey:
		return RoleCommunityUser
	default:
		return 0 // No permissions for unknown roles
	}
}

// GetCommunityRoleName returns the human-readable name of a community role
func GetCommunityRoleName(role string) string {
	switch role {
	case RoleCommunityOwnerKey:
		return "Community Owner"
	case RoleCommunityStaffKey:
		return "Community Staff"
	case RoleCommunityUserKey:
		return "Community User"
	default:
		return "Unknown"
	}
}

// IsCommunityRoleValid checks if a community role identifier is valid
func IsCommunityRoleValid(role string) bool {
	switch role {
	case RoleCommunityOwnerKey, RoleCommunityStaffKey, RoleCommunityUserKey:
		return true
	default:
		return false
	}
}

// CommunityPermissionNames maps each permission bit to its key name (for debugging/logging)
var CommunityPermissionNames = map[uint64]string{
	COMMUNITY__MEMBERSHIP_INVITE:  PermCommunityMembershipInvite,
	COMMUNITY__MEMBERSHIP_DELETE:  PermCommunityMembershipDelete,
	COMMUNITY__SERVER_ADD:         PermCommunityServerAdd,
	COMMUNITY__SERVER_EDIT:        PermCommunityServerEdit,
	COMMUNITY__SERVER_DELETE:      PermCommunityServerDelete,
	COMMUNITY__MISSIONLIST_ADD:    PermCommunityMissionlistAdd,
	COMMUNITY__MISSIONLIST_EDIT:   PermCommunityMissionlistEdit,
	COMMUNITY__GAMEMODELIST_ADD:   PermCommunityGamemodelistAdd,
	COMMUNITY__GAMEMODELIST_EDIT:  PermCommunityGamemodelistEdit,
	COMMUNITY__INFO_EDIT:          PermCommunityInfoEdit,
	COMMUNITY__ANALYTICS:          PermCommunityAnalytics,
	COMMUNITY__TRANSFER_OWNERSHIP: PermCommunityTransferOwnership,
	COMMUNITY__SUSPEND:            PermCommunitySuspend,
}

// CommunityPermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var CommunityPermissionKeyToBit = map[string]uint8{
	PermCommunityMembershipInvite:  0,
	PermCommunityMembershipDelete:  1,
	PermCommunityServerAdd:         2,
	PermCommunityServerEdit:        3,
	PermCommunityServerDelete:      4,
	PermCommunityMissionlistAdd:    5,
	PermCommunityMissionlistEdit:   6,
	PermCommunityGamemodelistAdd:   7,
	PermCommunityGamemodelistEdit:  8,
	PermCommunityInfoEdit:          9,
	PermCommunityAnalytics:         10,
	PermCommunityTransferOwnership: 11,
	PermCommunitySuspend:           12,
}

// GetCommunityPermissionName returns the key name of a single permission bit
func GetCommunityPermissionName(permission uint64) string {
	if name, ok := CommunityPermissionNames[permission]; ok {
		return name
	}
	return "UNKNOWN_PERMISSION"
}

// GetAllCommunityPermissionNames returns all permission key names in a bitmask
func GetAllCommunityPermissionNames(mask uint64) []string {
	var names []string
	for perm, name := range CommunityPermissionNames {
		if mask&perm != 0 {
			names = append(names, name)
		}
	}
	return names
}
