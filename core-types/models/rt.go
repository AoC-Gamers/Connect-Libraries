package models

import "time"

// Connect-RT Internal API Request/Response Models

// ============================================
// PRESENCE MANAGEMENT
// ============================================

// RTUserPresence represents a user's real-time presence
type RTUserPresence struct {
	SteamID        string     `json:"steamid"`
	Status         int        `json:"status"` // 1=Online, 2=Offline, 4=InGame, 5=Spectating (matches Connect-RT UserPlatformStatus)
	LastSeenAt     *time.Time `json:"lastSeenAt,omitempty"`
	CurrentLobbyID string     `json:"currentLobbyId,omitempty"`
	GameState      string     `json:"gameState,omitempty"` // "in_menu", "in_game", "in_lobby"
	UpdatedAt      time.Time  `json:"updatedAt"`
}

// GetUserPresenceResponse is the response when getting user presence
type GetUserPresenceResponse struct {
	Presence *RTUserPresence `json:"presence"` // nil if user has no presence
}

// InitializePresenceRequest is the request to initialize user presence
type InitializePresenceRequest struct {
	Status    string `json:"status"`    // Initial status (default: "online")
	GameState string `json:"gameState"` // Initial game state (default: "in_menu")
}

// InitializePresenceResponse is the response after initializing presence
type InitializePresenceResponse struct {
	Presence RTUserPresence `json:"presence"`
	Success  bool           `json:"success"`
	Message  string         `json:"message,omitempty"`
}

// BatchGetPresenceRequest is the request to get presence for multiple users
type BatchGetPresenceRequest struct {
	SteamIDs []string `json:"steamids"` // Max 100 steamIDs
}

// BatchGetPresenceResponse is the response with multiple user presences
type BatchGetPresenceResponse struct {
	Presences map[string]*RTUserPresence `json:"presences"` // Map: steamID -> presence (nil if no presence)
}

// GetOnlineUsersResponse is the response when getting all online users
type GetOnlineUsersResponse struct {
	Users []RTUserPresence `json:"users"`
	Count int              `json:"count"`
}

// GetUsersByStatusResponse is the response when getting users by status
type GetUsersByStatusResponse struct {
	Status string           `json:"status"`
	Users  []RTUserPresence `json:"users"`
	Count  int              `json:"count"`
}

// ============================================
// NOTIFICATIONS
// ============================================

// RTCreateNotificationRequest is the request to create a notification in RT
type RTCreateNotificationRequest struct {
	UserID   string                 `json:"user_id"`            // SteamID of the recipient
	Type     string                 `json:"type"`               // Notification type (e.g., "invitation_responded")
	Category string                 `json:"category"`           // Category: "invitation", "system", "achievement", "social"
	Title    string                 `json:"title"`              // Notification title
	Message  string                 `json:"message"`            // Notification message
	Metadata map[string]interface{} `json:"metadata,omitempty"` // Additional data (JSONB)
}

// RTCreateNotificationResponse is the response after creating a notification in RT
type RTCreateNotificationResponse struct {
	ID      string `json:"id"` // UUID of created notification
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// RTCreateEnrichedNotificationRequest is the request to create a notification with auto-enrichment
// RT will:
// 1. Fetch user data from Core if ActorSteamID is provided
// 2. Enrich metadata with SteamUserSnapshotMini
// 3. Save to database
// 4. Publish to NATS for real-time delivery
type RTCreateEnrichedNotificationRequest struct {
	UserID        string                 `json:"user_id"`                  // SteamID of the recipient
	Type          string                 `json:"type"`                     // Notification type
	Category      string                 `json:"category"`                 // Category
	Title         string                 `json:"title"`                    // Notification title
	Message       string                 `json:"message"`                  // Notification message
	ActorSteamID  string                 `json:"actor_steam_id,omitempty"` // SteamID of the user who triggered the notification (will be enriched)
	Metadata      map[string]interface{} `json:"metadata,omitempty"`       // Additional data (JSONB)
	PublishToNATS bool                   `json:"publish_to_nats"`          // If true, also publish to NATS for real-time delivery
}

// RTCreateEnrichedNotificationResponse is the response after creating enriched notification
type RTCreateEnrichedNotificationResponse struct {
	ID      string `json:"id"` // UUID of created notification
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// OWNERSHIP TRANSFER EVENTS
// ============================================

// RTOwnershipTransferredRequest is the request to handle ownership transfer event
// RT will:
// 1. Fetch scope name from Core
// 2. Fetch actor (new owner) data from Core
// 3. Create notifications for: previous owner (if exists), new owner, staff members
// 4. Enrich metadata with user snapshots
// 5. Save to database
// 6. Publish to NATS for real-time delivery
type RTOwnershipTransferredRequest struct {
	ScopeType     string  `json:"scope_type"`               // "community" | "team"
	ScopeID       int64   `json:"scope_id"`                 // ID of the scope
	NewOwner      string  `json:"new_owner"`                // SteamID of new owner
	PreviousOwner *string `json:"previous_owner,omitempty"` // SteamID of previous owner (nil if first owner)
	InviterID     string  `json:"inviter_id"`               // SteamID of who triggered the transfer
}

// RTOwnershipTransferredResponse is the response after processing ownership transfer
type RTOwnershipTransferredResponse struct {
	Success              bool   `json:"success"`
	NotificationsCreated int    `json:"notifications_created"` // Number of notifications created
	Message              string `json:"message,omitempty"`
}

// RTOwnershipDemotionRequest is the request to handle ownership demotion event
// When a user is demoted from owner to staff/member
type RTOwnershipDemotionRequest struct {
	UserSteamID  string `json:"user_steamid"`  // SteamID of user being demoted
	ScopeType    string `json:"scope_type"`    // "community" | "team"
	ScopeID      int64  `json:"scope_id"`      // ID of the scope
	PreviousRole string `json:"previous_role"` // "owner"
	NewRole      string `json:"new_role"`      // "staff" | "member"
	NewOwnerID   string `json:"new_owner_id"`  // SteamID of new owner
	InvitationID int64  `json:"invitation_id"` // Related invitation ID
}

// RTOwnershipDemotionResponse is the response after processing demotion
type RTOwnershipDemotionResponse struct {
	Success              bool   `json:"success"`
	NotificationsCreated int    `json:"notifications_created"` // Number of notifications created (usually 1)
	Message              string `json:"message,omitempty"`
}
