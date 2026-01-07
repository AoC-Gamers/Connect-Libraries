package permissions

// Web Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/web.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// Community Management (bits 0-6)
	WEB__COMMUNITY_VIEW                 uint64 = 1 << 0 // View community listing and details from the web
	WEB__COMMUNITIES_ADD                uint64 = 1 << 2 // Create new communities
	WEB__COMMUNITIES_EDIT               uint64 = 1 << 3 // Edit community information from web
	WEB__COMMUNITIES_DELETE             uint64 = 1 << 4 // Delete communities from web
	WEB__COMMUNITIES_TRANSFER_OWNERSHIP uint64 = 1 << 5 // Transfer ownership of a community
	WEB__COMMUNITIES_SUSPEND            uint64 = 1 << 6 // Suspend or unsuspend communities

	// Team Management (bits 7-11)
	WEB__TEAM_VIEW                uint64 = 1 << 7  // View team listing and details from the web
	WEB__TEAMS_EDIT               uint64 = 1 << 8  // Edit team information from web
	WEB__TEAMS_DELETE             uint64 = 1 << 9  // Delete teams from web
	WEB__TEAMS_TRANSFER_OWNERSHIP uint64 = 1 << 10 // Transfer ownership of teams
	WEB__TEAMS_SUSPEND            uint64 = 1 << 11 // Suspend or unsuspend teams

	// Membership Management (bits 12-14)
	WEB__MEMBERSHIP_VIEW   uint64 = 1 << 12 // View membership details from the web interface
	WEB__MEMBERSHIP_INVITE uint64 = 1 << 13 // Invite users to scopes (communities, teams, web) via the web interface
	WEB__MEMBERSHIP_DELETE uint64 = 1 << 14 // Remove memberships from the web interface

	// Sanctions Management (bits 15-19)
	WEB__SANCTIONS_VIEW    uint64 = 1 << 15 // View sanctions applied to users
	WEB__SANCTIONS_ADD     uint64 = 1 << 16 // Apply a sanction to a user
	WEB__SANCTIONS_EDIT    uint64 = 1 << 17 // Edit existing sanctions
	WEB__SANCTIONS_DELETE  uint64 = 1 << 18 // Remove sanctions from users
	WEB__SANCTIONS_SUSPEND uint64 = 1 << 19 // Temporarily suspend users via sanctions

	// Mission Management (bits 20-24)
	WEB__MISSION_VIEW    uint64 = 1 << 20 // View mission listings and details from the web
	WEB__MISSION_ADD     uint64 = 1 << 21 // Create new missions from the web
	WEB__MISSION_EDIT    uint64 = 1 << 22 // Edit existing missions from the web
	WEB__MISSION_DELETE  uint64 = 1 << 23 // Delete missions from the web
	WEB__MISSION_SUSPEND uint64 = 1 << 24 // Suspend or unsuspend missions from the web

	// Gamemode Management (bits 25-29)
	WEB__GAMEMODE_VIEW    uint64 = 1 << 25 // View gamemode listings and details from the web
	WEB__GAMEMODE_ADD     uint64 = 1 << 26 // Create new gamemodes from the web
	WEB__GAMEMODE_EDIT    uint64 = 1 << 27 // Edit existing gamemodes from the web
	WEB__GAMEMODE_DELETE  uint64 = 1 << 28 // Delete gamemodes from the web
	WEB__GAMEMODE_SUSPEND uint64 = 1 << 29 // Suspend or unsuspend gamemodes from the web

	// Lobby Management (bits 30-34)
	WEB__LOBBY_VIEW           uint64 = 1 << 30 // View lobby listings and details from the web
	WEB__LOBBY_CREATE_PUBLIC  uint64 = 1 << 31 // Create new public lobbies from the web
	WEB__LOBBY_CREATE_PRIVATE uint64 = 1 << 32 // Create new private lobbies from the web
	WEB__LOBBY_JOIN           uint64 = 1 << 33 // Join lobbies from the web
	WEB__LOBBY_SPECTATE       uint64 = 1 << 34 // Spectate lobbies from the web

	// Platform Administration (bits 35-37)
	WEB__VIEW_AUDIT_LOG uint64 = 1 << 35 // Access audit logs
	WEB__VIEW_METRICS   uint64 = 1 << 36 // Access platform metrics and statistics
	WEB__SETTINGS       uint64 = 1 << 37 // Modify platform settings
)

// ============================================
// Web Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods

const (
	// Community Management
	PermWebCommunityView                = "WEB__COMMUNITY_VIEW"
	PermWebCommunitiesAdd               = "WEB__COMMUNITIES_ADD"
	PermWebCommunitiesEdit              = "WEB__COMMUNITIES_EDIT"
	PermWebCommunitiesDelete            = "WEB__COMMUNITIES_DELETE"
	PermWebCommunitiesTransferOwnership = "WEB__COMMUNITIES_TRANSFER_OWNERSHIP"
	PermWebCommunitiesSuspend           = "WEB__COMMUNITIES_SUSPEND"

	// Team Management
	PermWebTeamView               = "WEB__TEAM_VIEW"
	PermWebTeamsEdit              = "WEB__TEAMS_EDIT"
	PermWebTeamsDelete            = "WEB__TEAMS_DELETE"
	PermWebTeamsTransferOwnership = "WEB__TEAMS_TRANSFER_OWNERSHIP"
	PermWebTeamsSuspend           = "WEB__TEAMS_SUSPEND"

	// Membership Management
	PermWebMembershipView   = "WEB__MEMBERSHIP_VIEW"
	PermWebMembershipInvite = "WEB__MEMBERSHIP_INVITE"
	PermWebMembershipDelete = "WEB__MEMBERSHIP_DELETE"

	// Sanctions Management
	PermWebSanctionsView    = "WEB__SANCTIONS_VIEW"
	PermWebSanctionsAdd     = "WEB__SANCTIONS_ADD"
	PermWebSanctionsEdit    = "WEB__SANCTIONS_EDIT"
	PermWebSanctionsDelete  = "WEB__SANCTIONS_DELETE"
	PermWebSanctionsSuspend = "WEB__SANCTIONS_SUSPEND"

	// Mission Management
	PermWebMissionView    = "WEB__MISSION_VIEW"
	PermWebMissionAdd     = "WEB__MISSION_ADD"
	PermWebMissionEdit    = "WEB__MISSION_EDIT"
	PermWebMissionDelete  = "WEB__MISSION_DELETE"
	PermWebMissionSuspend = "WEB__MISSION_SUSPEND"

	// Gamemode Management
	PermWebGamemodeView    = "WEB__GAMEMODE_VIEW"
	PermWebGamemodeAdd     = "WEB__GAMEMODE_ADD"
	PermWebGamemodeEdit    = "WEB__GAMEMODE_EDIT"
	PermWebGamemodeDelete  = "WEB__GAMEMODE_DELETE"
	PermWebGamemodeSuspend = "WEB__GAMEMODE_SUSPEND"

	// Lobby Management
	PermWebLobbyView          = "WEB__LOBBY_VIEW"
	PermWebLobbyCreatePublic  = "WEB__LOBBY_CREATE_PUBLIC"
	PermWebLobbyCreatePrivate = "WEB__LOBBY_CREATE_PRIVATE"
	PermWebLobbyJoin          = "WEB__LOBBY_JOIN"
	PermWebLobbySpectate      = "WEB__LOBBY_SPECTATE"

	// Platform Administration
	PermWebViewAuditLog = "WEB__VIEW_AUDIT_LOG"
	PermWebViewMetrics  = "WEB__VIEW_METRICS"
	PermWebSettings     = "WEB__SETTINGS"
)

// ============================================
// Web Role String Constants
// ============================================
// String constants for web role identifiers

const (
	RoleWebOwnerKey = "web_owner"
	RoleWebStaffKey = "web_staff"
	RoleWebUserKey  = "web_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/web.json groups
var (
	// WEB__BASIC - Basic web user permissions
	WEB__BASIC = WEB__COMMUNITY_VIEW |
		WEB__TEAM_VIEW |
		WEB__SANCTIONS_VIEW |
		WEB__MISSION_VIEW |
		WEB__GAMEMODE_VIEW |
		WEB__LOBBY_VIEW |
		WEB__LOBBY_CREATE_PUBLIC |
		WEB__LOBBY_CREATE_PRIVATE |
		WEB__LOBBY_JOIN |
		WEB__LOBBY_SPECTATE

	// WEB__STAFF - Staff management and moderation
	WEB__STAFF = WEB__TEAMS_EDIT |
		WEB__TEAMS_DELETE |
		WEB__SANCTIONS_ADD |
		WEB__SANCTIONS_EDIT |
		WEB__SANCTIONS_SUSPEND |
		WEB__MEMBERSHIP_VIEW |
		WEB__MEMBERSHIP_INVITE |
		WEB__MISSION_ADD |
		WEB__MISSION_EDIT |
		WEB__MISSION_SUSPEND |
		WEB__GAMEMODE_ADD |
		WEB__GAMEMODE_EDIT |
		WEB__GAMEMODE_SUSPEND

	// WEB__OWNER - Platform ownership (only one owner via .env)
	WEB__OWNER = WEB__COMMUNITIES_ADD |
		WEB__COMMUNITIES_EDIT |
		WEB__COMMUNITIES_DELETE |
		WEB__COMMUNITIES_TRANSFER_OWNERSHIP |
		WEB__COMMUNITIES_SUSPEND |
		WEB__TEAMS_TRANSFER_OWNERSHIP |
		WEB__TEAMS_SUSPEND |
		WEB__MEMBERSHIP_DELETE |
		WEB__SANCTIONS_DELETE |
		WEB__MISSION_DELETE |
		WEB__GAMEMODE_DELETE |
		WEB__VIEW_AUDIT_LOG |
		WEB__VIEW_METRICS |
		WEB__SETTINGS
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleWebUser - Basic web user (WEB__BASIC group)
	RoleWebUser = WEB__BASIC

	// RoleWebStaff - Staff member (WEB__BASIC + WEB__STAFF groups)
	RoleWebStaff = WEB__BASIC | WEB__STAFF

	// RoleWebOwner - Platform owner (WEB__BASIC + WEB__STAFF + WEB__OWNER groups)
	RoleWebOwner = WEB__BASIC | WEB__STAFF | WEB__OWNER
)

// GetRolePermissions returns the permission bitmask for a given role
func GetRolePermissions(role string) uint64 {
	switch role {
	case RoleWebOwnerKey:
		return RoleWebOwner
	case RoleWebStaffKey:
		return RoleWebStaff
	case RoleWebUserKey:
		return RoleWebUser
	default:
		return 0 // No permissions for unknown roles
	}
}

// GetRoleName returns the human-readable name of a role
func GetRoleName(role string) string {
	switch role {
	case RoleWebOwnerKey:
		return "Web Owner"
	case RoleWebStaffKey:
		return "Web Staff"
	case RoleWebUserKey:
		return "Web User"
	default:
		return "Unknown"
	}
}

// IsRoleValid checks if a role identifier is valid
func IsRoleValid(role string) bool {
	switch role {
	case RoleWebOwnerKey, RoleWebStaffKey, RoleWebUserKey:
		return true
	default:
		return false
	}
}

// PermissionNames maps each permission bit to its key name (for debugging/logging)
var PermissionNames = map[uint64]string{
	WEB__COMMUNITY_VIEW:                 PermWebCommunityView,
	WEB__COMMUNITIES_ADD:                PermWebCommunitiesAdd,
	WEB__COMMUNITIES_EDIT:               PermWebCommunitiesEdit,
	WEB__COMMUNITIES_DELETE:             PermWebCommunitiesDelete,
	WEB__COMMUNITIES_TRANSFER_OWNERSHIP: PermWebCommunitiesTransferOwnership,
	WEB__COMMUNITIES_SUSPEND:            PermWebCommunitiesSuspend,
	WEB__TEAM_VIEW:                      PermWebTeamView,
	WEB__TEAMS_EDIT:                     PermWebTeamsEdit,
	WEB__TEAMS_DELETE:                   PermWebTeamsDelete,
	WEB__TEAMS_TRANSFER_OWNERSHIP:       PermWebTeamsTransferOwnership,
	WEB__TEAMS_SUSPEND:                  PermWebTeamsSuspend,
	WEB__MEMBERSHIP_VIEW:                PermWebMembershipView,
	WEB__MEMBERSHIP_INVITE:              PermWebMembershipInvite,
	WEB__MEMBERSHIP_DELETE:              PermWebMembershipDelete,
	WEB__SANCTIONS_VIEW:                 PermWebSanctionsView,
	WEB__SANCTIONS_ADD:                  PermWebSanctionsAdd,
	WEB__SANCTIONS_EDIT:                 PermWebSanctionsEdit,
	WEB__SANCTIONS_DELETE:               PermWebSanctionsDelete,
	WEB__SANCTIONS_SUSPEND:              PermWebSanctionsSuspend,
	WEB__MISSION_VIEW:                   PermWebMissionView,
	WEB__MISSION_ADD:                    PermWebMissionAdd,
	WEB__MISSION_EDIT:                   PermWebMissionEdit,
	WEB__MISSION_DELETE:                 PermWebMissionDelete,
	WEB__MISSION_SUSPEND:                PermWebMissionSuspend,
	WEB__GAMEMODE_VIEW:                  PermWebGamemodeView,
	WEB__GAMEMODE_ADD:                   PermWebGamemodeAdd,
	WEB__GAMEMODE_EDIT:                  PermWebGamemodeEdit,
	WEB__GAMEMODE_DELETE:                PermWebGamemodeDelete,
	WEB__GAMEMODE_SUSPEND:               PermWebGamemodeSuspend,
	WEB__LOBBY_VIEW:                     PermWebLobbyView,
	WEB__LOBBY_CREATE_PUBLIC:            PermWebLobbyCreatePublic,
	WEB__LOBBY_CREATE_PRIVATE:           PermWebLobbyCreatePrivate,
	WEB__LOBBY_JOIN:                     PermWebLobbyJoin,
	WEB__LOBBY_SPECTATE:                 PermWebLobbySpectate,
	WEB__VIEW_AUDIT_LOG:                 PermWebViewAuditLog,
	WEB__VIEW_METRICS:                   PermWebViewMetrics,
	WEB__SETTINGS:                       PermWebSettings,
}

// PermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var PermissionKeyToBit = map[string]uint8{
	PermWebCommunityView:                0,
	PermWebCommunitiesAdd:               2,
	PermWebCommunitiesEdit:              3,
	PermWebCommunitiesDelete:            4,
	PermWebCommunitiesTransferOwnership: 5,
	PermWebCommunitiesSuspend:           6,
	PermWebTeamView:                     7,
	PermWebTeamsEdit:                    8,
	PermWebTeamsDelete:                  9,
	PermWebTeamsTransferOwnership:       10,
	PermWebTeamsSuspend:                 11,
	PermWebMembershipView:               12,
	PermWebMembershipInvite:             13,
	PermWebMembershipDelete:             14,
	PermWebSanctionsView:                15,
	PermWebSanctionsAdd:                 16,
	PermWebSanctionsEdit:                17,
	PermWebSanctionsDelete:              18,
	PermWebSanctionsSuspend:             19,
	PermWebMissionView:                  20,
	PermWebMissionAdd:                   21,
	PermWebMissionEdit:                  22,
	PermWebMissionDelete:                23,
	PermWebMissionSuspend:               24,
	PermWebGamemodeView:                 25,
	PermWebGamemodeAdd:                  26,
	PermWebGamemodeEdit:                 27,
	PermWebGamemodeDelete:               28,
	PermWebGamemodeSuspend:              29,
	PermWebLobbyView:                    30,
	PermWebLobbyCreatePublic:            31,
	PermWebLobbyCreatePrivate:           32,
	PermWebLobbyJoin:                    33,
	PermWebLobbySpectate:                34,
	PermWebViewAuditLog:                 35,
	PermWebViewMetrics:                  36,
	PermWebSettings:                     37,
}

// GetPermissionName returns the key name of a single permission bit
func GetPermissionName(permission uint64) string {
	if name, ok := PermissionNames[permission]; ok {
		return name
	}
	return "UNKNOWN_PERMISSION"
}

// GetAllPermissionNames returns all permission key names in a bitmask
func GetAllPermissionNames(mask uint64) []string {
	var names []string
	for perm, name := range PermissionNames {
		if mask&perm != 0 {
			names = append(names, name)
		}
	}
	return names
}
