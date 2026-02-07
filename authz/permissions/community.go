package permissions

// ============================================
// Community Permission Bitmask Constants
// ============================================
// Generated from Connect-Auth/seeds/permissions/community.json
// Each constant represents a bit in uint64 (max 64 permissions)

const (
	// CommunityMembershipInvite invites users to the community (membership invite).
	CommunityMembershipInvite uint64 = 1 << 0
	// CommunityMembershipDelete removes memberships from the community.
	CommunityMembershipDelete uint64 = 1 << 1

	// CommunityServerAdd adds a server to the community's server list.
	CommunityServerAdd uint64 = 1 << 2
	// CommunityServerEdit edits an existing server in the community's server list.
	CommunityServerEdit uint64 = 1 << 3
	// CommunityServerDelete deletes a server from the community's server list.
	CommunityServerDelete uint64 = 1 << 4

	// CommunityMissionListAdd adds a new mission to the community's mission list.
	CommunityMissionListAdd uint64 = 1 << 5
	// CommunityMissionListEdit edits an existing mission in the community's mission list.
	CommunityMissionListEdit uint64 = 1 << 6

	// CommunityGamemodeListAdd adds a new gamemode to the community's gamemode list.
	CommunityGamemodeListAdd uint64 = 1 << 7
	// CommunityGamemodeListEdit edits an existing gamemode in the community's gamemode list.
	CommunityGamemodeListEdit uint64 = 1 << 8

	// CommunityInfoEdit edits community information.
	CommunityInfoEdit uint64 = 1 << 9
	// CommunityAnalytics views community statistics and analytics.
	CommunityAnalytics uint64 = 1 << 10
	// CommunityTransferOwnership transfers community ownership.
	CommunityTransferOwnership uint64 = 1 << 11
	// CommunitySuspend suspends or unsuspends this community.
	CommunitySuspend uint64 = 1 << 12
	// CommunityAuditView views audit logs for this community.
	CommunityAuditView uint64 = 1 << 13

	// CommunityRolesView views roles and permissions for this community.
	CommunityRolesView uint64 = 1 << 14
	// CommunityRolesEdit edits roles and permissions for this community.
	CommunityRolesEdit uint64 = 1 << 15
)

// ============================================
// Community Permission Names (for API calls)
// ============================================
// String constants for use with CheckUserPermission() and similar API methods
// These match the permission keys in Connect-Auth

const (
	// PermCommunityMembershipInvite is the permission key for COMMUNITY__MEMBERSHIP_INVITE.
	PermCommunityMembershipInvite = "COMMUNITY__MEMBERSHIP_INVITE"
	// PermCommunityMembershipDelete is the permission key for COMMUNITY__MEMBERSHIP_DELETE.
	PermCommunityMembershipDelete = "COMMUNITY__MEMBERSHIP_DELETE"

	// PermCommunityServerAdd is the permission key for COMMUNITY__SERVER_ADD.
	PermCommunityServerAdd = "COMMUNITY__SERVER_ADD"
	// PermCommunityServerEdit is the permission key for COMMUNITY__SERVER_EDIT.
	PermCommunityServerEdit = "COMMUNITY__SERVER_EDIT"
	// PermCommunityServerDelete is the permission key for COMMUNITY__SERVER_DELETE.
	PermCommunityServerDelete = "COMMUNITY__SERVER_DELETE"

	// PermCommunityMissionListAdd is the permission key for COMMUNITY__MISSIONLIST_ADD.
	PermCommunityMissionListAdd = "COMMUNITY__MISSIONLIST_ADD"
	// PermCommunityMissionListEdit is the permission key for COMMUNITY__MISSIONLIST_EDIT.
	PermCommunityMissionListEdit = "COMMUNITY__MISSIONLIST_EDIT"

	// PermCommunityGamemodeListAdd is the permission key for COMMUNITY__GAMEMODELIST_ADD.
	PermCommunityGamemodeListAdd = "COMMUNITY__GAMEMODELIST_ADD"
	// PermCommunityGamemodeListEdit is the permission key for COMMUNITY__GAMEMODELIST_EDIT.
	PermCommunityGamemodeListEdit = "COMMUNITY__GAMEMODELIST_EDIT"

	// PermCommunityInfoEdit is the permission key for COMMUNITY__INFO_EDIT.
	PermCommunityInfoEdit = "COMMUNITY__INFO_EDIT"
	// PermCommunityAnalytics is the permission key for COMMUNITY__ANALYTICS.
	PermCommunityAnalytics = "COMMUNITY__ANALYTICS"
	// PermCommunityTransferOwnership is the permission key for COMMUNITY__TRANSFER_OWNERSHIP.
	PermCommunityTransferOwnership = "COMMUNITY__TRANSFER_OWNERSHIP"
	// PermCommunitySuspend is the permission key for COMMUNITY__SUSPEND.
	PermCommunitySuspend = "COMMUNITY__SUSPEND"
	// PermCommunityAuditView is the permission key for COMMUNITY__AUDIT_VIEW.
	PermCommunityAuditView = "COMMUNITY__AUDIT_VIEW"

	// PermCommunityRolesView is the permission key for COMMUNITY__ROLES_VIEW.
	PermCommunityRolesView = "COMMUNITY__ROLES_VIEW"
	// PermCommunityRolesEdit is the permission key for COMMUNITY__ROLES_EDIT.
	PermCommunityRolesEdit = "COMMUNITY__ROLES_EDIT"
)

// ============================================
// Community Role String Constants
// ============================================
// String constants for community role identifiers

const (
	// RoleCommunityOwnerKey is the role key for a community owner.
	RoleCommunityOwnerKey = "community_owner"
	// RoleCommunityStaffKey is the role key for a community staff member.
	RoleCommunityStaffKey = "community_staff"
	// RoleCommunityUserKey is the role key for a community user.
	RoleCommunityUserKey = "community_user"
)

// ============================================
// Permission Groups
// ============================================
// Matches seeds/permissions/community.json groups
var (
	// CommunityBasic - Basic community membership (no permissions)
	CommunityBasic uint64 = 0

	// CommunityStaff - Staff management permissions
	CommunityStaff = CommunityServerAdd |
		CommunityServerEdit |
		CommunityServerDelete |
		CommunityMissionListAdd |
		CommunityMissionListEdit |
		CommunityGamemodeListAdd |
		CommunityGamemodeListEdit |
		CommunityRolesView

	// CommunityOwner - Full community control
	CommunityOwner = CommunityMembershipInvite |
		CommunityMembershipDelete |
		CommunityInfoEdit |
		CommunityAnalytics |
		CommunityTransferOwnership |
		CommunitySuspend |
		CommunityAuditView |
		CommunityRolesEdit
)

// Role Presets - Calculated bitmasks for each role
var (
	// RoleCommunityUser - Basic community member (CommunityBasic group - no permissions)
	RoleCommunityUser = CommunityBasic

	// RoleCommunityStaff - Community staff member (CommunityBasic + CommunityStaff groups)
	RoleCommunityStaff = CommunityBasic | CommunityStaff

	// RoleCommunityOwner - Community owner (CommunityBasic + CommunityStaff + CommunityOwner groups)
	RoleCommunityOwner = CommunityBasic | CommunityStaff | CommunityOwner
)

var communityRolePermissions = map[string]uint64{
	RoleCommunityOwnerKey: RoleCommunityOwner,
	RoleCommunityStaffKey: RoleCommunityStaff,
	RoleCommunityUserKey:  RoleCommunityUser,
}

var communityRoleNames = map[string]string{
	RoleCommunityOwnerKey: "Community Owner",
	RoleCommunityStaffKey: "Community Staff",
	RoleCommunityUserKey:  "Community User",
}

// GetCommunityRolePermissions returns the permission bitmask for a given community role
func GetCommunityRolePermissions(role string) uint64 {
	return getRolePermissions(role, communityRolePermissions)
}

// GetCommunityRoleName returns the human-readable name of a community role
func GetCommunityRoleName(role string) string {
	return getRoleName(role, communityRoleNames)
}

// IsCommunityRoleValid checks if a community role identifier is valid
func IsCommunityRoleValid(role string) bool {
	return isRoleValid(role, communityRoleNames)
}

// CommunityPermissionNames maps each permission bit to its key name (for debugging/logging)
var CommunityPermissionNames = map[uint64]string{
	CommunityMembershipInvite:  PermCommunityMembershipInvite,
	CommunityMembershipDelete:  PermCommunityMembershipDelete,
	CommunityServerAdd:         PermCommunityServerAdd,
	CommunityServerEdit:        PermCommunityServerEdit,
	CommunityServerDelete:      PermCommunityServerDelete,
	CommunityMissionListAdd:    PermCommunityMissionListAdd,
	CommunityMissionListEdit:   PermCommunityMissionListEdit,
	CommunityGamemodeListAdd:   PermCommunityGamemodeListAdd,
	CommunityGamemodeListEdit:  PermCommunityGamemodeListEdit,
	CommunityInfoEdit:          PermCommunityInfoEdit,
	CommunityAnalytics:         PermCommunityAnalytics,
	CommunityTransferOwnership: PermCommunityTransferOwnership,
	CommunitySuspend:           PermCommunitySuspend,
	CommunityAuditView:         PermCommunityAuditView,
	CommunityRolesView:         PermCommunityRolesView,
	CommunityRolesEdit:         PermCommunityRolesEdit,
}

// CommunityPermissionKeyToBit maps permission key names to their bit values (reverse lookup)
var CommunityPermissionKeyToBit = map[string]uint8{
	PermCommunityMembershipInvite:  0,
	PermCommunityMembershipDelete:  1,
	PermCommunityServerAdd:         2,
	PermCommunityServerEdit:        3,
	PermCommunityServerDelete:      4,
	PermCommunityMissionListAdd:    5,
	PermCommunityMissionListEdit:   6,
	PermCommunityGamemodeListAdd:   7,
	PermCommunityGamemodeListEdit:  8,
	PermCommunityInfoEdit:          9,
	PermCommunityAnalytics:         10,
	PermCommunityTransferOwnership: 11,
	PermCommunitySuspend:           12,
	PermCommunityAuditView:         13,
	PermCommunityRolesView:         14,
	PermCommunityRolesEdit:         15,
}

// GetCommunityPermissionName returns the key name of a single permission bit
func GetCommunityPermissionName(permission uint64) string {
	return getPermissionName(permission, CommunityPermissionNames)
}

// GetAllCommunityPermissionNames returns all permission key names in a bitmask
func GetAllCommunityPermissionNames(mask uint64) []string {
	return getAllPermissionNames(mask, CommunityPermissionNames)
}
