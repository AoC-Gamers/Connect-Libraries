package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Connect-Core Internal API Request/Response Models

// ============================================
// USER MANAGEMENT
// ============================================

// SteamUserSnapshot represents essential Steam user data from Steam Web API
// Only includes core fields needed by Connect-Core, matching Connect-Auth's UserSnapshot
type SteamUserSnapshot struct {
	PersonaName              string `json:"personaname"`
	AvatarHash               string `json:"avatarhash"`
	ProfileURL               string `json:"profileurl"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"` // 1=Private, 2=FriendsOnly, 3=Public
}

// SyncUserRequest is the request to sync a user from Connect-Auth
type SyncUserRequest struct {
	SteamSnapshot SteamUserSnapshot `json:"steamSnapshot"`
}

// SyncUserResponse is the response after syncing a user
type SyncUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// User represents a complete user from Core
type User struct {
	SteamID       string            `json:"steamid"`
	SteamSnapshot SteamUserSnapshot `json:"steam"`
}

// BatchUsersRequest is the request to get multiple users
type BatchUsersRequest struct {
	SteamIDs []string `json:"steam_ids"`
}

// BatchUsersResponse is the response with multiple users
type BatchUsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

// SteamUserSnapshotMini - Lightweight version for UI cards (only name and avatar)
// Optimized for batch requests where full user data is not needed
type SteamUserSnapshotMini struct {
	PersonaName string `json:"personaname"`
	AvatarHash  string `json:"avatarhash"`
}

// UserMini represents a minimal user for UI cards (invitations, notifications, etc.)
type UserMini struct {
	SteamID string                `json:"steamid"`
	Steam   SteamUserSnapshotMini `json:"steam"`
}

// BatchUsersMiniRequest is the request to get multiple users (lightweight)
type BatchUsersMiniRequest struct {
	SteamIDs []string `json:"steam_ids"`
}

// BatchUsersMiniResponse is the response with multiple users (lightweight)
type BatchUsersMiniResponse struct {
	Users []UserMini `json:"users"`
	Count int        `json:"count"`
}

// ============================================
// MISSION MANAGEMENT
// ============================================

// Mission represents a game mission
type Mission struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	GamemodeID  int64     `json:"gamemodeId"`
	Difficulty  string    `json:"difficulty"`
	Status      string    `json:"status"` // "active", "suspended"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// CreateMissionRequest represents a request to create a mission
type CreateMissionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	GamemodeID  int64  `json:"gamemodeId"`
	Difficulty  string `json:"difficulty"`
}

// CreateMissionResponse is the response after creating a mission
type CreateMissionResponse struct {
	Mission Mission `json:"mission"`
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
}

// UpdateMissionRequest represents a request to update a mission
type UpdateMissionRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	GamemodeID  *int64  `json:"gamemodeId,omitempty"`
	Difficulty  *string `json:"difficulty,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// UpdateMissionResponse is the response after updating a mission
type UpdateMissionResponse struct {
	Mission Mission `json:"mission"`
	Success bool    `json:"success"`
	Message string  `json:"message,omitempty"`
}

// GetMissionResponse is the response when getting a mission
type GetMissionResponse struct {
	Mission Mission `json:"mission"`
}

// DeleteMissionResponse is the response after deleting a mission
type DeleteMissionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// GAMEMODE MANAGEMENT
// ============================================

// Gamemode represents a game mode
type Gamemode struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // "active", "suspended"
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// ListGamemodesResponse is the response when listing gamemodes
type ListGamemodesResponse struct {
	Gamemodes []Gamemode `json:"gamemodes"`
}

// CreateGamemodeRequest represents a request to create a gamemode
type CreateGamemodeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateGamemodeResponse is the response after creating a gamemode
type CreateGamemodeResponse struct {
	Gamemode Gamemode `json:"gamemode"`
	Success  bool     `json:"success"`
	Message  string   `json:"message,omitempty"`
}

// ============================================
// TEAM MANAGEMENT
// ============================================

// TeamMember represents a team member
type TeamMember struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
	AvatarURL string    `json:"avatarUrl,omitempty"`
}

// GetTeamMembersResponse is the response when getting team members
type GetTeamMembersResponse struct {
	TeamID  int64        `json:"teamId"`
	Members []TeamMember `json:"members"`
}

// TeamSettings represents team configuration from Core
type TeamSettings struct {
	TeamID          int64  `json:"team_id"`
	MaxMembers      int    `json:"max_members"`
	IsRecruiting    bool   `json:"is_recruiting"`
	InvitePolicy    string `json:"invite_policy"`    // OWNER_ONLY, STAFF_AND_OWNER, ANY_MEMBER
	JoinPolicy      string `json:"join_policy"`      // CLOSED, INVITE_ONLY, REQUEST_TO_JOIN, OPEN
	RequireApproval bool   `json:"require_approval"` // For join requests
	IsActive        bool   `json:"is_active"`
	IsSuspended     bool   `json:"is_suspended"`
}

// GetTeamSettingsResponse is the response when getting team settings
type GetTeamSettingsResponse struct {
	Settings TeamSettings `json:"settings"`
}

// ValidateTeamExistsResponse is the response when validating team existence
type ValidateTeamExistsResponse struct {
	Exists bool  `json:"exists"`
	TeamID int64 `json:"team_id,omitempty"`
}

// ============================================
// SERVER MANAGEMENT
// ============================================

// Server represents a game server
type Server struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	IPAddress      string    `json:"ipAddress"`
	Port           int       `json:"port"`
	Status         string    `json:"status"` // "online", "offline", "maintenance"
	CommunityID    int64     `json:"communityId"`
	MaxPlayers     int       `json:"maxPlayers"`
	CurrentPlayers int       `json:"currentPlayers"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// ListServersResponse is the response when listing servers
type ListServersResponse struct {
	Servers []Server `json:"servers"`
}

// GetServerResponse is the response when getting a server
type GetServerResponse struct {
	Server Server `json:"server"`
}

// ============================================
// SETTINGS MANAGEMENT
// ============================================

// Settings represents global platform settings
type Settings struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// GetSettingsResponse is the response when getting settings
type GetSettingsResponse struct {
	Settings []Settings `json:"settings"`
}

// UpdateSettingsRequest represents a request to update settings
type UpdateSettingsRequest struct {
	Settings map[string]string `json:"settings"` // key: value pairs
}

// UpdateSettingsResponse is the response after updating settings
type UpdateSettingsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// ============================================
// JSONB SERIALIZATION (sqlx/database/sql compatibility)
// ============================================

// Scan implements sql.Scanner for SteamUserSnapshot
// This allows reading JSONB fields from PostgreSQL into Go structs
func (s *SteamUserSnapshot) Scan(value interface{}) error {
	if value == nil {
		*s = SteamUserSnapshot{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SteamUserSnapshot: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, s)
}

// Value implements driver.Valuer for SteamUserSnapshot
// This allows writing Go structs as JSONB fields to PostgreSQL
func (s SteamUserSnapshot) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// ============================================
// SETTINGS MANAGEMENT
// ============================================

// EntityType represents the type of configuration entity
type EntityType string

const (
	EntityTypeConfig   EntityType = "CONFIG"
	EntityTypeGeneral  EntityType = "GENERAL"
	EntityTypeTeam     EntityType = "TEAM"
	EntityTypeLobby    EntityType = "LOBBY"
	EntityTypeSecurity EntityType = "SECURITY"
)

// SettingDataType represents the data type of a setting value
type SettingDataType string

const (
	SettingDataTypeInt    SettingDataType = "INT"
	SettingDataTypeString SettingDataType = "STRING"
	SettingDataTypeBool   SettingDataType = "BOOL"
	SettingDataTypeFloat  SettingDataType = "FLOAT"
	SettingDataTypeJSON   SettingDataType = "JSON"
	SettingDataTypeEnum   SettingDataType = "ENUM"
)

// SettingMetadata contains metadata about a setting
type SettingMetadata struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	DefaultValue  string   `json:"default_value,omitempty"`
	MinValue      *float64 `json:"min_value,omitempty"`
	MaxValue      *float64 `json:"max_value,omitempty"`
	AllowedValues []string `json:"allowed_values,omitempty"`
	DisplayOrder  int      `json:"display_order,omitempty"`
	IsVisible     bool     `json:"is_visible,omitempty"`
	RequiresAuth  bool     `json:"requires_auth,omitempty"`
	UIComponent   string   `json:"ui_component,omitempty"`
}

// Setting represents a configuration setting
type Setting struct {
	ID         int             `json:"id"`
	EntityType EntityType      `json:"entityType"`
	Key        string          `json:"key"`
	Value      string          `json:"value"`
	DataType   SettingDataType `json:"dataType"`
	Name       string          `json:"name"`
	Metadata   SettingMetadata `json:"metadata"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

// GetSettingResponse is the response when getting a single setting
type GetSettingResponse struct {
	Setting
}

// GetSettingsDetailResponse is the response when getting multiple settings with full details
type GetSettingsDetailResponse struct {
	Settings []Setting `json:"settings"`
	Total    int       `json:"total"`
}
