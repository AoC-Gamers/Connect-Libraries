package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AoC-Gamers/connect-libraries/internal/endpoints"
	"github.com/AoC-Gamers/connect-libraries/internal/errors"
	"github.com/AoC-Gamers/connect-libraries/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	missionPathFormat = "/core/internal/missions/%d"
)

// Client is a type-safe HTTP client for Connect-Core internal API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Connect-Core client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// ============================================
// USER MANAGEMENT
// ============================================

// SyncUser syncs a user from Connect-Auth after login
// Endpoint: POST /core/internal/users/{steamid}
func (c *Client) SyncUser(ctx context.Context, steamID string, req models.SyncUserRequest) (*models.SyncUserResponse, error) {
	path := fmt.Sprintf("/core/internal/users/%s", steamID)
	var resp models.SyncUserResponse
	if err := c.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, fmt.Errorf("sync user: %w", err)
	}
	return &resp, nil
}

// GetUsersBatch retrieves multiple users by their SteamIDs
// Endpoint: POST /core/internal/users/batch
func (c *Client) GetUsersBatch(ctx context.Context, steamIDs []string) (*models.BatchUsersResponse, error) {
	req := models.BatchUsersRequest{
		SteamIDs: steamIDs,
	}
	var resp models.BatchUsersResponse
	if err := c.doRequest(ctx, "POST", "/core/internal/users/batch", req, &resp); err != nil {
		return nil, fmt.Errorf("get users batch: %w", err)
	}
	return &resp, nil
}

// GetUsersBatchMini retrieves multiple users in lightweight format (only name and avatar)
// Optimized for UI cards, invitations, notifications, etc.
// Endpoint: POST /core/internal/users/batch-mini
func (c *Client) GetUsersBatchMini(ctx context.Context, steamIDs []string) (*models.BatchUsersMiniResponse, error) {
	req := models.BatchUsersMiniRequest{
		SteamIDs: steamIDs,
	}
	var resp models.BatchUsersMiniResponse
	if err := c.doRequest(ctx, "POST", "/core/internal/users/batch-mini", req, &resp); err != nil {
		return nil, fmt.Errorf("get users batch mini: %w", err)
	}
	return &resp, nil
}

// ============================================
// MISSION MANAGEMENT
// ============================================

// GetMission retrieves a mission by ID
// Endpoint: GET /core/internal/missions/{id}
func (c *Client) GetMission(ctx context.Context, missionID int64) (*models.GetMissionResponse, error) {
	path := fmt.Sprintf(missionPathFormat, missionID)
	var resp models.GetMissionResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get mission: %w", err)
	}
	return &resp, nil
}

// CreateMission creates a new mission
// Endpoint: POST /core/internal/missions
func (c *Client) CreateMission(ctx context.Context, req models.CreateMissionRequest) (*models.CreateMissionResponse, error) {
	var resp models.CreateMissionResponse
	if err := c.doRequest(ctx, "POST", endpoints.CoreCreateMission.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("create mission: %w", err)
	}
	return &resp, nil
}

// UpdateMission updates an existing mission
// Endpoint: PUT /core/internal/missions/{id}
func (c *Client) UpdateMission(ctx context.Context, missionID int64, req models.UpdateMissionRequest) (*models.UpdateMissionResponse, error) {
	path := fmt.Sprintf(missionPathFormat, missionID)
	var resp models.UpdateMissionResponse
	if err := c.doRequest(ctx, "PUT", path, req, &resp); err != nil {
		return nil, fmt.Errorf("update mission: %w", err)
	}
	return &resp, nil
}

// DeleteMission deletes a mission
// Endpoint: DELETE /core/internal/missions/{id}
func (c *Client) DeleteMission(ctx context.Context, missionID int64) (*models.DeleteMissionResponse, error) {
	path := fmt.Sprintf(missionPathFormat, missionID)
	var resp models.DeleteMissionResponse
	if err := c.doRequest(ctx, "DELETE", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("delete mission: %w", err)
	}
	return &resp, nil
}

// ============================================
// GAMEMODE MANAGEMENT
// ============================================

// ListGamemodes retrieves all gamemodes
// Endpoint: GET /core/internal/gamemodes
func (c *Client) ListGamemodes(ctx context.Context) (*models.ListGamemodesResponse, error) {
	var resp models.ListGamemodesResponse
	if err := c.doRequest(ctx, "GET", endpoints.CoreListGamemodes.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("list gamemodes: %w", err)
	}
	return &resp, nil
}

// CreateGamemode creates a new gamemode
// Endpoint: POST /core/internal/gamemodes
func (c *Client) CreateGamemode(ctx context.Context, req models.CreateGamemodeRequest) (*models.CreateGamemodeResponse, error) {
	var resp models.CreateGamemodeResponse
	if err := c.doRequest(ctx, "POST", endpoints.CoreCreateGamemode.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("create gamemode: %w", err)
	}
	return &resp, nil
}

// ============================================
// TEAM MANAGEMENT
// ============================================

// GetTeamMembers retrieves members of a team
// Endpoint: GET /core/internal/teams/{id}/members
func (c *Client) GetTeamMembers(ctx context.Context, teamID int64) (*models.GetTeamMembersResponse, error) {
	path := fmt.Sprintf("/core/internal/teams/%d/members", teamID)
	var resp models.GetTeamMembersResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}
	return &resp, nil
}

// GetTeamSettings retrieves team settings
// Endpoint: GET /core/internal/teams/{id}/settings
func (c *Client) GetTeamSettings(ctx context.Context, teamID int64) (*models.TeamSettings, error) {
	path := fmt.Sprintf("/core/internal/teams/%d/settings", teamID)
	var resp models.GetTeamSettingsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get team settings: %w", err)
	}
	return &resp.Settings, nil
}

// ValidateTeamExists checks if a team exists
// Endpoint: GET /core/internal/teams/{id}
func (c *Client) ValidateTeamExists(ctx context.Context, teamID int64) (bool, error) {
	path := fmt.Sprintf("/core/internal/teams/%d", teamID)
	var resp models.ValidateTeamExistsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		// If 404, team doesn't exist
		if internalErr, ok := err.(*errors.InternalError); ok && internalErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("validate team exists: %w", err)
	}
	return resp.Exists, nil
}

// ============================================
// SERVER MANAGEMENT
// ============================================

// ListServers retrieves all servers
// Endpoint: GET /core/internal/servers
func (c *Client) ListServers(ctx context.Context) (*models.ListServersResponse, error) {
	var resp models.ListServersResponse
	if err := c.doRequest(ctx, "GET", endpoints.CoreListServers.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("list servers: %w", err)
	}
	return &resp, nil
}

// GetServer retrieves a server by ID
// Endpoint: GET /core/internal/servers/{id}
func (c *Client) GetServer(ctx context.Context, serverID int64) (*models.GetServerResponse, error) {
	path := fmt.Sprintf("/core/internal/servers/%d", serverID)
	var resp models.GetServerResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get server: %w", err)
	}
	return &resp, nil
}

// ============================================
// SETTINGS MANAGEMENT
// ============================================

// GetSettings retrieves global settings
// Endpoint: GET /core/internal/settings
func (c *Client) GetSettings(ctx context.Context) (*models.GetSettingsResponse, error) {
	var resp models.GetSettingsResponse
	if err := c.doRequest(ctx, "GET", endpoints.CoreGetSettings.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get settings: %w", err)
	}
	return &resp, nil
}

// UpdateSettings updates global settings
// Endpoint: PUT /core/internal/settings
func (c *Client) UpdateSettings(ctx context.Context, req models.UpdateSettingsRequest) (*models.UpdateSettingsResponse, error) {
	var resp models.UpdateSettingsResponse
	if err := c.doRequest(ctx, "PUT", endpoints.CoreUpdateSettings.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("update settings: %w", err)
	}
	return &resp, nil
}

// ============================================
// HTTP CLIENT HELPERS
// ============================================

// doRequest performs an HTTP request with proper error handling
func (c *Client) doRequest(ctx context.Context, method, path string, reqBody, respBody interface{}) error {
	url := c.baseURL + path

	bodyReader, err := c.prepareRequestBody(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Perform request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Msg("Connect-Core HTTP request failed")
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("path", path).
			Int("status", resp.StatusCode).
			Msg("Connect-Core failed to read response body")
		return fmt.Errorf("read response: %w", err)
	}

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(method, path, resp.StatusCode, respBodyBytes)
	}

	// Parse response body if needed
	return c.parseResponseBody(method, path, resp.StatusCode, respBodyBytes, respBody)
}

// prepareRequestBody marshals the request body to JSON
func (c *Client) prepareRequestBody(reqBody interface{}) (io.Reader, error) {
	if reqBody == nil {
		return nil, nil
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	return bytes.NewReader(bodyBytes), nil
}

// handleErrorResponse processes non-2xx HTTP responses
func (c *Client) handleErrorResponse(method, path string, statusCode int, respBodyBytes []byte) error {
	errMsg := c.extractErrorMessage(respBodyBytes)

	log.Error().
		Str("method", method).
		Str("path", path).
		Int("status", statusCode).
		Str("error", errMsg).
		Msg("Connect-Core internal API error")

	return errors.NewInternalError(statusCode, "Connect-Core", path, errMsg)
}

// extractErrorMessage attempts to parse error message from response body
func (c *Client) extractErrorMessage(respBodyBytes []byte) string {
	var errResp map[string]interface{}
	if json.Unmarshal(respBodyBytes, &errResp) == nil {
		if msg, ok := errResp["error"].(string); ok {
			return msg
		}
		if msg, ok := errResp["message"].(string); ok {
			return msg
		}
	}
	return string(respBodyBytes)
}

// parseResponseBody unmarshals the response body into the provided interface
func (c *Client) parseResponseBody(method, path string, statusCode int, respBodyBytes []byte, respBody interface{}) error {
	if respBody == nil || len(respBodyBytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBodyBytes, respBody); err != nil {
		bodyPreview := string(respBodyBytes)
		if len(bodyPreview) > 200 {
			bodyPreview = bodyPreview[:200] + "..."
		}
		log.Error().
			Err(err).
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("body_preview", bodyPreview).
			Msg("Connect-Core failed to parse response JSON")
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

// ============================================
// SETTINGS MANAGEMENT
// ============================================

// GetSetting retrieves a specific setting by entity type and key
// Endpoint: GET /core/internal/settings/{entity}/{key}
func (c *Client) GetSetting(ctx context.Context, entityType models.EntityType, key string) (*models.Setting, error) {
	path := fmt.Sprintf("/core/internal/settings/%s/%s", entityType, key)
	var resp models.GetSettingResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get setting %s/%s: %w", entityType, key, err)
	}
	return &resp.Setting, nil
}

// GetSettingsByEntity retrieves all settings for a specific entity type
// Endpoint: GET /core/internal/settings/{entity}
// Returns simple key-value pairs for easy consumption
func (c *Client) GetSettingsByEntity(ctx context.Context, entityType models.EntityType) (map[string]string, error) {
	path := fmt.Sprintf("/core/internal/settings/%s", entityType)
	var resp map[string]string
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get settings for entity %s: %w", entityType, err)
	}
	return resp, nil
}

// GetAllSettings retrieves all settings with full details
// Endpoint: GET /core/internal/settings
func (c *Client) GetAllSettings(ctx context.Context) ([]models.Setting, error) {
	path := "/core/internal/settings"
	var resp models.GetSettingsDetailResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get all settings: %w", err)
	}
	return resp.Settings, nil
}

// GetSettingValue is a convenience method that returns just the value of a setting
// Returns empty string if setting is not found
func (c *Client) GetSettingValue(ctx context.Context, entityType models.EntityType, key string) (string, error) {
	setting, err := c.GetSetting(ctx, entityType, key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

// HealthCheck performs a health check on Connect-Core
func (c *Client) HealthCheck(ctx context.Context) error {
	// Use list gamemodes endpoint as health check
	_, err := c.ListGamemodes(ctx)
	return err
}
