package models

// Connect-Auth Internal API Request/Response Models

// ============================================
// PERMISSION CHECKING
// ============================================

type CheckPermissionRequest struct {
	UserID     string `json:"userId"`     // Steam ID
	Permission string `json:"permission"` // Permission key (e.g., "WEB__COMMUNITY_VIEW")
	EntityType string `json:"entityType"` // Scope type (e.g., "WEB", "COMMUNITY", "TEAM")
	EntityID   string `json:"entityId"`   // Entity ID (e.g., "1" for web scope)
}

type CheckPermissionResponse struct {
	HasPermission bool   `json:"hasPermission"`
	Reason        string `json:"reason,omitempty"`
}

// ============================================
// WEB OWNER
// ============================================

type WebOwnerResponse struct {
	WebOwnerSteamID string `json:"webOwnerSteamId"`
}

// ============================================
// SCOPE MANAGEMENT
// ============================================

type CreateScopeRequest struct {
	EntityType string `json:"entityType"` // "COMMUNITY", "TEAM", "LOBBY"
	EntityID   string `json:"entityId"`   // e.g., "COMMUNITY_3"
	OwnerID    string `json:"ownerId"`    // Optional: Steam ID of owner
}

type CreateScopeResponse struct {
	ScopeID string `json:"scopeId"` // e.g., "COMMUNITY:COMMUNITY_3"
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ScopeStaffCountResponse struct {
	ScopeType  string `json:"scopeType"`
	EntityID   string `json:"entityId"`
	StaffCount int    `json:"staffCount"`
	OwnerCount int    `json:"ownerCount"`
	TotalStaff int    `json:"totalStaff"`
}

// ============================================
// ROLE MANAGEMENT
// ============================================

type AssignRoleRequest struct {
	UserID     string `json:"userId"`
	Role       string `json:"role"`
	EntityType string `json:"entityType"`
	EntityID   string `json:"entityId"`
}

type AssignRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type RemoveRoleResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// MEMBERSHIP MANAGEMENT
// ============================================

type CreateMembershipRequest struct {
	UserID  string `json:"userId"`
	ScopeID string `json:"scopeId"`
	Role    string `json:"role"`
	Status  string `json:"status"` // "active", "pending", "suspended"
}

type CreateMembershipResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type UpdateMembershipRequest struct {
	Role             *string `json:"role,omitempty"`
	AllowPermissions *uint64 `json:"allow_permissions,omitempty"`
	DenyPermissions  *uint64 `json:"deny_permissions,omitempty"`
	PerformedBy      string  `json:"performed_by,omitempty"`
	Reason           string  `json:"reason,omitempty"`
}

type UpdateMembershipResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type CheckMembershipRequest struct {
	UserID    string `json:"userId"`
	EntityID  string `json:"entityId"`
	ScopeType string `json:"scopeType,omitempty"` // Default: "COMMUNITY"
}

type CheckMembershipResponse struct {
	IsMember bool   `json:"isMember"`
	Message  string `json:"message,omitempty"` // e.g., "role=community_staff,status=active"
}

// ============================================
// OWNERSHIP TRANSFER
// ============================================

type InitiateTransferRequest struct {
	NewOwnerID string `json:"newOwnerId"` // Steam ID of new owner
}

type InitiateTransferResponse struct {
	TransferID string `json:"transferId"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
}

type TransferResponse struct {
	ID            string `json:"id"`
	ScopeType     string `json:"scopeType"`
	EntityID      string `json:"entityId"`
	CurrentOwner  string `json:"currentOwner"`
	ProposedOwner string `json:"proposedOwner"`
	Status        string `json:"status"` // "pending", "accepted", "rejected", "cancelled"
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type RespondTransferRequest struct {
	Action string `json:"action"` // "accept" or "reject"
}

type RespondTransferResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type CancelTransferResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// CACHE MANAGEMENT
// ============================================

type CacheRefreshRequest struct {
	Scope    string `json:"scope,omitempty"`    // Optional: specific scope to refresh
	EntityID string `json:"entityId,omitempty"` // Optional: specific entity
}

type CacheRefreshResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type CacheStatusResponse struct {
	Status     string                 `json:"status"`
	CacheSize  int                    `json:"cacheSize"`
	LastReload string                 `json:"lastReload"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// ============================================
// NOTIFICATIONS
// ============================================

type CreateNotificationRequest struct {
	UserID   string                 `json:"userId"`
	Type     string                 `json:"type"`
	Title    string                 `json:"title"`
	Message  string                 `json:"message"`
	Priority string                 `json:"priority,omitempty"` // "low", "normal", "high"
	Data     map[string]interface{} `json:"data,omitempty"`
}

type CreateNotificationResponse struct {
	NotificationID string `json:"notificationId"`
	Success        bool   `json:"success"`
	Message        string `json:"message,omitempty"`
}

type TeamMatchLobbyNotificationRequest struct {
	TeamID  string `json:"teamId"`
	LobbyID string `json:"lobbyId"`
	Message string `json:"message"`
}

type TeamMatchLobbyNotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// CATALOG
// ============================================

type PermissionCatalogResponse struct {
	Scopes map[string]ScopePermissions `json:"scopes"`
}

type ScopePermissions struct {
	Scope       string                 `json:"scope"`
	Permissions []PermissionDefinition `json:"permissions"`
	Groups      []PermissionGroupDef   `json:"groups"`
}

type PermissionDefinition struct {
	Bit         int    `json:"bit"`
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type PermissionGroupDef struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

// ============================================
// MEMBERSHIP & CAPABILITIES
// ============================================

// GetMembershipRequest - Solicitud para obtener membresía
type GetMembershipRequest struct {
	UserID  string `json:"userId"`  // Steam ID del usuario
	ScopeID string `json:"scopeId"` // Scope ID (ej: "COMMUNITY:3")
}

// MembershipResponse - Respuesta con información de membresía
type MembershipResponse struct {
	UserID   string `json:"userId"`
	ScopeID  string `json:"scopeId"`
	Role     string `json:"role"`     // community_owner, community_staff, etc.
	Status   string `json:"status"`   // active, suspended, etc.
	JoinedAt string `json:"joinedAt"` // RFC3339
}

// GetCapabilitiesRequest - Solicitud para obtener capabilities
type GetCapabilitiesRequest struct {
	UserID  string `json:"userId"`  // Steam ID del usuario
	ScopeID string `json:"scopeId"` // Scope ID (ej: "COMMUNITY:3")
}

// CapabilitiesResponse - Respuesta con capabilities del usuario
type CapabilitiesResponse struct {
	Effective map[string]uint64 `json:"effective"` // Máscara efectiva por scope type (web, community, team)
	Deny      map[string]uint64 `json:"deny"`      // Máscara de denegaciones
}

// ============================================
// SCOPE STAFF
// ============================================

// StaffMember representa un miembro del staff en un scope
type StaffMember struct {
	UserSteamID string `json:"user_steamid"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
}

// ScopeStaffResponse - Respuesta con lista de staff de un scope
type ScopeStaffResponse struct {
	Staff []StaffMember `json:"staff"`
	Count int           `json:"count"` // Total count (not paginated)

	// Pagination fields (optional)
	Page       int  `json:"page,omitempty"`
	Limit      int  `json:"limit,omitempty"`
	TotalPages int  `json:"totalPages,omitempty"`
	HasMore    bool `json:"hasMore,omitempty"`
}

// ============================================
// USER ALL MEMBERSHIPS
// ============================================

// UserMembershipItem representa una membresía individual de un usuario
type UserMembershipItem struct {
	// ScopeID eliminado - es redundante (derivable como "${ScopeType}:${EntityID}")
	ScopeType        string `json:"scope_type"`
	EntityID         string `json:"entity_id"`
	Role             string `json:"role"`
	AllowPermissions uint64 `json:"allow_permissions"`
	DenyPermissions  uint64 `json:"deny_permissions"`
	JoinedAt         string `json:"joined_at"` // RFC3339
}

// UserAllMembershipsResponse - Respuesta con todas las membresías de un usuario
type UserAllMembershipsResponse struct {
	WebRole              string               `json:"web_role"`
	WebAllowPermissions  uint64               `json:"web_allow_permissions"`
	WebDenyPermissions   uint64               `json:"web_deny_permissions"`
	CommunityMemberships []UserMembershipItem `json:"community_memberships"`
	TeamMemberships      []UserMembershipItem `json:"team_memberships"`
	TotalMemberships     int                  `json:"total_memberships"`
}
