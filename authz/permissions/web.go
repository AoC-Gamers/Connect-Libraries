package permissions

// Web Permission Bitmask Constants
// Generated from Connect-Auth/seeds/permissions/web.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// WebCommunityView views community listing and details from the web.
	WebCommunityView uint64 = 1 << 0
	// WebCommunitiesAdd creates new communities.
	WebCommunitiesAdd uint64 = 1 << 2
	// WebCommunitiesEdit edits community information from the web.
	WebCommunitiesEdit uint64 = 1 << 3
	// WebCommunitiesDelete deletes communities from the web.
	WebCommunitiesDelete uint64 = 1 << 4
	// WebCommunitiesTransferOwnership transfers ownership of a community.
	WebCommunitiesTransferOwnership uint64 = 1 << 5
	// WebCommunitiesSuspend suspends or unsuspends communities.
	WebCommunitiesSuspend uint64 = 1 << 6

	// WebTeamView views team listing and details from the web.
	WebTeamView uint64 = 1 << 7
	// WebTeamsEdit edits team information from the web.
	WebTeamsEdit uint64 = 1 << 8
	// WebTeamsDelete deletes teams from the web.
	WebTeamsDelete uint64 = 1 << 9
	// WebTeamsTransferOwnership transfers ownership of teams.
	WebTeamsTransferOwnership uint64 = 1 << 10
	// WebTeamsSuspend suspends or unsuspends teams.
	WebTeamsSuspend uint64 = 1 << 11

	// WebMembershipView views membership details from the web interface.
	WebMembershipView uint64 = 1 << 12
	// WebMembershipInvite invites users to scopes (communities, teams, web) via the web interface.
	WebMembershipInvite uint64 = 1 << 13
	// WebMembershipDelete removes memberships from the web interface.
	WebMembershipDelete uint64 = 1 << 14

	// WebSanctionsView views sanctions applied to users.
	WebSanctionsView uint64 = 1 << 15
	// WebSanctionsAdd applies a sanction to a user.
	WebSanctionsAdd uint64 = 1 << 16
	// WebSanctionsEdit edits existing sanctions.
	WebSanctionsEdit uint64 = 1 << 17
	// WebSanctionsDelete removes sanctions from users.
	WebSanctionsDelete uint64 = 1 << 18
	// WebSanctionsSuspend temporarily suspends users via sanctions.
	WebSanctionsSuspend uint64 = 1 << 19

	// WebMissionView views mission listings and details from the web.
	WebMissionView uint64 = 1 << 20
	// WebMissionAdd creates new missions from the web.
	WebMissionAdd uint64 = 1 << 21
	// WebMissionEdit edits existing missions from the web.
	WebMissionEdit uint64 = 1 << 22
	// WebMissionDelete deletes missions from the web.
	WebMissionDelete uint64 = 1 << 23
	// WebMissionSuspend suspends or unsuspends missions from the web.
	WebMissionSuspend uint64 = 1 << 24

	// WebGamemodeView views gamemode listings and details from the web.
	WebGamemodeView uint64 = 1 << 25
	// WebGamemodeAdd creates new gamemodes from the web.
	WebGamemodeAdd uint64 = 1 << 26
	// WebGamemodeEdit edits existing gamemodes from the web.
	WebGamemodeEdit uint64 = 1 << 27
	// WebGamemodeDelete deletes gamemodes from the web.
	WebGamemodeDelete uint64 = 1 << 28
	// WebGamemodeSuspend suspends or unsuspends gamemodes from the web.
	WebGamemodeSuspend uint64 = 1 << 29

	// WebLobbyView views lobby listings and details from the web.
	WebLobbyView uint64 = 1 << 30
	// WebLobbyCreatePublic creates new public lobbies from the web.
	WebLobbyCreatePublic uint64 = 1 << 31
	// WebLobbyCreatePrivate creates new private lobbies from the web.
	WebLobbyCreatePrivate uint64 = 1 << 32
	// WebLobbyJoin joins lobbies from the web.
	WebLobbyJoin uint64 = 1 << 33
	// WebLobbySpectate spectates lobbies from the web.
	WebLobbySpectate uint64 = 1 << 34

	// WebViewAuditLog accesses audit logs.
	WebViewAuditLog uint64 = 1 << 35
	// WebViewMetrics accesses platform metrics and statistics.
	WebViewMetrics uint64 = 1 << 36
	// WebSettings modifies platform settings.
	WebSettings uint64 = 1 << 37

	// WebRolesView views roles and permissions.
	WebRolesView uint64 = 1 << 38
	// WebRolesEdit edits roles and permissions.
	WebRolesEdit uint64 = 1 << 39
)

// ============================================
// Web Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods

const (
	// PermWebCommunityView is the permission key for WEB__COMMUNITY_VIEW.
	PermWebCommunityView = "WEB__COMMUNITY_VIEW"
	// PermWebCommunitiesAdd is the permission key for WEB__COMMUNITIES_ADD.
	PermWebCommunitiesAdd = "WEB__COMMUNITIES_ADD"
	// PermWebCommunitiesEdit is the permission key for WEB__COMMUNITIES_EDIT.
	PermWebCommunitiesEdit = "WEB__COMMUNITIES_EDIT"
	// PermWebCommunitiesDelete is the permission key for WEB__COMMUNITIES_DELETE.
	PermWebCommunitiesDelete = "WEB__COMMUNITIES_DELETE"
	// PermWebCommunitiesTransferOwnership is the permission key for WEB__COMMUNITIES_TRANSFER_OWNERSHIP.
	PermWebCommunitiesTransferOwnership = "WEB__COMMUNITIES_TRANSFER_OWNERSHIP"
	// PermWebCommunitiesSuspend is the permission key for WEB__COMMUNITIES_SUSPEND.
	PermWebCommunitiesSuspend = "WEB__COMMUNITIES_SUSPEND"

	// PermWebTeamView is the permission key for WEB__TEAM_VIEW.
	PermWebTeamView = "WEB__TEAM_VIEW"
	// PermWebTeamsEdit is the permission key for WEB__TEAMS_EDIT.
	PermWebTeamsEdit = "WEB__TEAMS_EDIT"
	// PermWebTeamsDelete is the permission key for WEB__TEAMS_DELETE.
	PermWebTeamsDelete = "WEB__TEAMS_DELETE"
	// PermWebTeamsTransferOwnership is the permission key for WEB__TEAMS_TRANSFER_OWNERSHIP.
	PermWebTeamsTransferOwnership = "WEB__TEAMS_TRANSFER_OWNERSHIP"
	// PermWebTeamsSuspend is the permission key for WEB__TEAMS_SUSPEND.
	PermWebTeamsSuspend = "WEB__TEAMS_SUSPEND"

	// PermWebMembershipView is the permission key for WEB__MEMBERSHIP_VIEW.
	PermWebMembershipView = "WEB__MEMBERSHIP_VIEW"
	// PermWebMembershipInvite is the permission key for WEB__MEMBERSHIP_INVITE.
	PermWebMembershipInvite = "WEB__MEMBERSHIP_INVITE"
	// PermWebMembershipDelete is the permission key for WEB__MEMBERSHIP_DELETE.
	PermWebMembershipDelete = "WEB__MEMBERSHIP_DELETE"

	// PermWebSanctionsView is the permission key for WEB__SANCTIONS_VIEW.
	PermWebSanctionsView = "WEB__SANCTIONS_VIEW"
	// PermWebSanctionsAdd is the permission key for WEB__SANCTIONS_ADD.
	PermWebSanctionsAdd = "WEB__SANCTIONS_ADD"
	// PermWebSanctionsEdit is the permission key for WEB__SANCTIONS_EDIT.
	PermWebSanctionsEdit = "WEB__SANCTIONS_EDIT"
	// PermWebSanctionsDelete is the permission key for WEB__SANCTIONS_DELETE.
	PermWebSanctionsDelete = "WEB__SANCTIONS_DELETE"
	// PermWebSanctionsSuspend is the permission key for WEB__SANCTIONS_SUSPEND.
	PermWebSanctionsSuspend = "WEB__SANCTIONS_SUSPEND"

	// PermWebMissionView is the permission key for WEB__MISSION_VIEW.
	PermWebMissionView = "WEB__MISSION_VIEW"
	// PermWebMissionAdd is the permission key for WEB__MISSION_ADD.
	PermWebMissionAdd = "WEB__MISSION_ADD"
	// PermWebMissionEdit is the permission key for WEB__MISSION_EDIT.
	PermWebMissionEdit = "WEB__MISSION_EDIT"
	// PermWebMissionDelete is the permission key for WEB__MISSION_DELETE.
	PermWebMissionDelete = "WEB__MISSION_DELETE"
	// PermWebMissionSuspend is the permission key for WEB__MISSION_SUSPEND.
	PermWebMissionSuspend = "WEB__MISSION_SUSPEND"

	// PermWebGamemodeView is the permission key for WEB__GAMEMODE_VIEW.
	PermWebGamemodeView = "WEB__GAMEMODE_VIEW"
	// PermWebGamemodeAdd is the permission key for WEB__GAMEMODE_ADD.
	PermWebGamemodeAdd = "WEB__GAMEMODE_ADD"
	// PermWebGamemodeEdit is the permission key for WEB__GAMEMODE_EDIT.
	PermWebGamemodeEdit = "WEB__GAMEMODE_EDIT"
	// PermWebGamemodeDelete is the permission key for WEB__GAMEMODE_DELETE.
	PermWebGamemodeDelete = "WEB__GAMEMODE_DELETE"
	// PermWebGamemodeSuspend is the permission key for WEB__GAMEMODE_SUSPEND.
	PermWebGamemodeSuspend = "WEB__GAMEMODE_SUSPEND"

	// PermWebLobbyView is the permission key for WEB__LOBBY_VIEW.
	PermWebLobbyView = "WEB__LOBBY_VIEW"
	// PermWebLobbyCreatePublic is the permission key for WEB__LOBBY_CREATE_PUBLIC.
	PermWebLobbyCreatePublic = "WEB__LOBBY_CREATE_PUBLIC"
	// PermWebLobbyCreatePrivate is the permission key for WEB__LOBBY_CREATE_PRIVATE.
	PermWebLobbyCreatePrivate = "WEB__LOBBY_CREATE_PRIVATE"
	// PermWebLobbyJoin is the permission key for WEB__LOBBY_JOIN.
	PermWebLobbyJoin = "WEB__LOBBY_JOIN"
	// PermWebLobbySpectate is the permission key for WEB__LOBBY_SPECTATE.
	PermWebLobbySpectate = "WEB__LOBBY_SPECTATE"

	// PermWebViewAuditLog is the permission key for WEB__VIEW_AUDIT_LOG.
	PermWebViewAuditLog = "WEB__VIEW_AUDIT_LOG"
	// PermWebViewMetrics is the permission key for WEB__VIEW_METRICS.
	PermWebViewMetrics = "WEB__VIEW_METRICS"
	// PermWebSettings is the permission key for WEB__SETTINGS.
	PermWebSettings = "WEB__SETTINGS"

	// PermWebRolesView is the permission key for WEB__ROLES_VIEW.
	PermWebRolesView = "WEB__ROLES_VIEW"
	// PermWebRolesEdit is the permission key for WEB__ROLES_EDIT.
	PermWebRolesEdit = "WEB__ROLES_EDIT"
)

// ============================================
// Web Role String Constants
// ============================================
// String constants for web role identifiers

const (
	// RoleWebOwnerKey is the role key for a web owner.
	RoleWebOwnerKey = "web_owner"
	// RoleWebStaffKey is the role key for a web staff member.
	RoleWebStaffKey = "web_staff"
	// RoleWebUserKey is the role key for a web user.
	RoleWebUserKey = "web_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/web.json groups
var (
	// WebBasic - Basic web user permissions
	WebBasic = WebCommunityView |
		WebTeamView |
		WebSanctionsView |
		WebMissionView |
		WebGamemodeView |
		WebLobbyView |
		WebLobbyCreatePublic |
		WebLobbyCreatePrivate |
		WebLobbyJoin |
		WebLobbySpectate

	// WebStaff - Staff management and moderation
	WebStaff = WebTeamsEdit |
		WebTeamsDelete |
		WebSanctionsAdd |
		WebSanctionsEdit |
		WebSanctionsSuspend |
		WebMembershipView |
		WebMembershipInvite |
		WebMissionAdd |
		WebMissionEdit |
		WebMissionSuspend |
		WebGamemodeAdd |
		WebGamemodeEdit |
		WebGamemodeSuspend |
		WebRolesView

	// WebOwner - Platform ownership (only one owner via .env)
	WebOwner = WebCommunitiesAdd |
		WebCommunitiesEdit |
		WebCommunitiesDelete |
		WebCommunitiesTransferOwnership |
		WebCommunitiesSuspend |
		WebTeamsTransferOwnership |
		WebTeamsSuspend |
		WebMembershipDelete |
		WebSanctionsDelete |
		WebMissionDelete |
		WebGamemodeDelete |
		WebViewAuditLog |
		WebViewMetrics |
		WebSettings |
		WebRolesEdit
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleWebUser - Basic web user (WebBasic group)
	RoleWebUser = WebBasic

	// RoleWebStaff - Staff member (WebBasic + WebStaff groups)
	RoleWebStaff = WebBasic | WebStaff

	// RoleWebOwner - Platform owner (WebBasic + WebStaff + WebOwner groups)
	RoleWebOwner = WebBasic | WebStaff | WebOwner
)

var webRolePermissions = map[string]uint64{
	RoleWebOwnerKey: RoleWebOwner,
	RoleWebStaffKey: RoleWebStaff,
	RoleWebUserKey:  RoleWebUser,
}

var webRoleNames = map[string]string{
	RoleWebOwnerKey: "Web Owner",
	RoleWebStaffKey: "Web Staff",
	RoleWebUserKey:  "Web User",
}

// GetRolePermissions returns the permission bitmask for a given role
func GetRolePermissions(role string) uint64 {
	return getRolePermissions(role, webRolePermissions)
}

// GetRoleName returns the human-readable name of a role
func GetRoleName(role string) string {
	return getRoleName(role, webRoleNames)
}

// IsRoleValid checks if a role identifier is valid
func IsRoleValid(role string) bool {
	return isRoleValid(role, webRoleNames)
}

// PermissionNames maps each permission bit to its key name (for debugging/logging)
var PermissionNames = map[uint64]string{
	WebCommunityView:                PermWebCommunityView,
	WebCommunitiesAdd:               PermWebCommunitiesAdd,
	WebCommunitiesEdit:              PermWebCommunitiesEdit,
	WebCommunitiesDelete:            PermWebCommunitiesDelete,
	WebCommunitiesTransferOwnership: PermWebCommunitiesTransferOwnership,
	WebCommunitiesSuspend:           PermWebCommunitiesSuspend,
	WebTeamView:                     PermWebTeamView,
	WebTeamsEdit:                    PermWebTeamsEdit,
	WebTeamsDelete:                  PermWebTeamsDelete,
	WebTeamsTransferOwnership:       PermWebTeamsTransferOwnership,
	WebTeamsSuspend:                 PermWebTeamsSuspend,
	WebMembershipView:               PermWebMembershipView,
	WebMembershipInvite:             PermWebMembershipInvite,
	WebMembershipDelete:             PermWebMembershipDelete,
	WebSanctionsView:                PermWebSanctionsView,
	WebSanctionsAdd:                 PermWebSanctionsAdd,
	WebSanctionsEdit:                PermWebSanctionsEdit,
	WebSanctionsDelete:              PermWebSanctionsDelete,
	WebSanctionsSuspend:             PermWebSanctionsSuspend,
	WebMissionView:                  PermWebMissionView,
	WebMissionAdd:                   PermWebMissionAdd,
	WebMissionEdit:                  PermWebMissionEdit,
	WebMissionDelete:                PermWebMissionDelete,
	WebMissionSuspend:               PermWebMissionSuspend,
	WebGamemodeView:                 PermWebGamemodeView,
	WebGamemodeAdd:                  PermWebGamemodeAdd,
	WebGamemodeEdit:                 PermWebGamemodeEdit,
	WebGamemodeDelete:               PermWebGamemodeDelete,
	WebGamemodeSuspend:              PermWebGamemodeSuspend,
	WebLobbyView:                    PermWebLobbyView,
	WebLobbyCreatePublic:            PermWebLobbyCreatePublic,
	WebLobbyCreatePrivate:           PermWebLobbyCreatePrivate,
	WebLobbyJoin:                    PermWebLobbyJoin,
	WebLobbySpectate:                PermWebLobbySpectate,
	WebViewAuditLog:                 PermWebViewAuditLog,
	WebViewMetrics:                  PermWebViewMetrics,
	WebSettings:                     PermWebSettings,
	WebRolesView:                    PermWebRolesView,
	WebRolesEdit:                    PermWebRolesEdit,
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
	PermWebRolesView:                    38,
	PermWebRolesEdit:                    39,
}

// GetPermissionName returns the key name of a single permission bit
func GetPermissionName(permission uint64) string {
	return getPermissionName(permission, PermissionNames)
}

// GetAllPermissionNames returns all permission key names in a bitmask
func GetAllPermissionNames(mask uint64) []string {
	return getAllPermissionNames(mask, PermissionNames)
}
