package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AoC-Gamers/connect-libraries/core-types/endpoints"
	"github.com/AoC-Gamers/connect-libraries/core-types/errors"
	"github.com/AoC-Gamers/connect-libraries/core-types/models"
	"github.com/rs/zerolog/log"
)

// Client is a type-safe HTTP client for Connect-Auth internal API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Connect-Auth client
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
// PERMISSION CHECKING
// ============================================

// CheckPermission verifies if a user has a specific permission
// Endpoint: POST /auth/internal/permissions/check
func (c *Client) CheckPermission(ctx context.Context, req models.CheckPermissionRequest) (*models.CheckPermissionResponse, error) {
	var resp models.CheckPermissionResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCheckPermission.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("check permission: %w", err)
	}
	return &resp, nil
}

// GetCatalog retrieves the complete permissions catalog
// Endpoint: GET /auth/internal/catalog
func (c *Client) GetCatalog(ctx context.Context) (*models.PermissionCatalogResponse, error) {
	var resp models.PermissionCatalogResponse
	if err := c.doRequest(ctx, "GET", endpoints.AuthGetCatalog.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get catalog: %w", err)
	}
	return &resp, nil
}

// ============================================
// WEB OWNER
// ============================================

// GetWebOwner retrieves the web_owner (global owner) steamID
// Endpoint: GET /auth/internal/owner
func (c *Client) GetWebOwner(ctx context.Context) (*models.WebOwnerResponse, error) {
	var resp models.WebOwnerResponse
	if err := c.doRequest(ctx, "GET", endpoints.AuthGetWebOwner.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get web owner: %w", err)
	}
	return &resp, nil
}

// ============================================
// SCOPE MANAGEMENT
// ============================================

// CreateScope creates a new scope
//
// DEPRECATED: This endpoint no longer exists in Connect-Auth.
// The auth.scopes table was eliminated. Use AssignRole() instead to create memberships.
// Scopes are now TEXT format "TYPE:ID" constructed on-demand.
//
// Endpoint: POST /auth/internal/scopes (REMOVED)
func (c *Client) CreateScope(ctx context.Context, req models.CreateScopeRequest) (*models.CreateScopeResponse, error) {
	var resp models.CreateScopeResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCreateScope.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("create scope: %w", err)
	}
	return &resp, nil
}

// DeleteScope deletes a scope and all related data
// Endpoint: DELETE /auth/internal/scopes/{scopeId}
func (c *Client) DeleteScope(ctx context.Context, scopeID string) error {
	path := fmt.Sprintf("/auth/internal/scopes/%s", scopeID)
	if err := c.doRequest(ctx, "DELETE", path, nil, nil); err != nil {
		return fmt.Errorf("delete scope: %w", err)
	}
	return nil
}

// GetScopeStaffCount retrieves staff count for a scope
// Endpoint: GET /auth/internal/scopes/{scopeId}/staff-count
func (c *Client) GetScopeStaffCount(ctx context.Context, scopeType, entityID string) (*models.ScopeStaffCountResponse, error) {
	// Construct scopeID in format "TYPE:ID"
	scopeID := fmt.Sprintf("%s:%s", scopeType, entityID)
	path := fmt.Sprintf("/auth/internal/scopes/%s/staff-count", scopeID)
	var resp models.ScopeStaffCountResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get scope staff count: %w", err)
	}
	return &resp, nil
}

// GetScopeStaff retrieves all staff members for a scope
// Endpoint: GET /auth/internal/scopes/{scopeId}/staff
func (c *Client) GetScopeStaff(ctx context.Context, scopeID string, page, limit int) (*models.ScopeStaffResponse, error) {
	path := fmt.Sprintf("/auth/internal/scopes/%s/staff", scopeID)

	// Add pagination query params if provided
	if page > 0 || limit > 0 {
		params := make([]string, 0, 2)
		if page > 0 {
			params = append(params, fmt.Sprintf("page=%d", page))
		}
		if limit > 0 {
			params = append(params, fmt.Sprintf("limit=%d", limit))
		}
		if len(params) > 0 {
			path += "?" + params[0]
			if len(params) > 1 {
				path += "&" + params[1]
			}
		}
	}

	var resp models.ScopeStaffResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get scope staff: %w", err)
	}
	return &resp, nil
}

// GetUserAllMemberships retrieves all memberships for a user
// Endpoint: GET /auth/internal/users/{steamid}/memberships
func (c *Client) GetUserAllMemberships(ctx context.Context, steamID string) (*models.UserAllMembershipsResponse, error) {
	path := fmt.Sprintf("/auth/internal/users/%s/memberships", steamID)
	var resp models.UserAllMembershipsResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get user all memberships: %w", err)
	}
	return &resp, nil
}

// ============================================
// ROLE MANAGEMENT
// ============================================

// AssignRole assigns a role to a user in a scope
// Endpoint: POST /auth/internal/roles/assign
func (c *Client) AssignRole(ctx context.Context, req models.AssignRoleRequest) (*models.AssignRoleResponse, error) {
	var resp models.AssignRoleResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthAssignRole.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("assign role: %w", err)
	}
	return &resp, nil
}

// RemoveRole removes a role from a user in a scope
// Endpoint: DELETE /auth/internal/roles/{userId}/{scopeId}
func (c *Client) RemoveRole(ctx context.Context, userID, scopeID string) (*models.RemoveRoleResponse, error) {
	path := fmt.Sprintf("/auth/internal/roles/%s/%s", userID, scopeID)
	var resp models.RemoveRoleResponse
	if err := c.doRequest(ctx, "DELETE", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("remove role: %w", err)
	}
	return &resp, nil
}

// ============================================
// MEMBERSHIP MANAGEMENT
// ============================================

// CreateMembership creates a membership for a user
// Endpoint: POST /auth/internal/memberships
func (c *Client) CreateMembership(ctx context.Context, req models.CreateMembershipRequest) (*models.CreateMembershipResponse, error) {
	var resp models.CreateMembershipResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCreateMembership.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("create membership: %w", err)
	}
	return &resp, nil
}

// CheckMembership checks if a user is a member of a scope
// Endpoint: POST /auth/internal/memberships/check
func (c *Client) CheckMembership(ctx context.Context, req models.CheckMembershipRequest) (*models.CheckMembershipResponse, error) {
	var resp models.CheckMembershipResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCheckMembership.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	return &resp, nil
}

// ============================================
// OWNERSHIP TRANSFER
// ============================================

// InitiateTransfer initiates an ownership transfer
// Endpoint: POST /auth/internal/scopes/{scopeType}/{entityId}/transfers/initiate
func (c *Client) InitiateTransfer(ctx context.Context, scopeType, entityID string, req models.InitiateTransferRequest) (*models.InitiateTransferResponse, error) {
	path := fmt.Sprintf("/auth/internal/scopes/%s/%s/transfers/initiate", scopeType, entityID)
	var resp models.InitiateTransferResponse
	if err := c.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, fmt.Errorf("initiate transfer: %w", err)
	}
	return &resp, nil
}

// GetPendingTransfer retrieves a pending ownership transfer
// Endpoint: GET /auth/internal/scopes/{scopeType}/{entityId}/transfers/pending
func (c *Client) GetPendingTransfer(ctx context.Context, scopeType, entityID string) (*models.TransferResponse, error) {
	path := fmt.Sprintf("/auth/internal/scopes/%s/%s/transfers/pending", scopeType, entityID)
	var resp models.TransferResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get pending transfer: %w", err)
	}
	return &resp, nil
}

// RespondToTransfer responds to an ownership transfer (accept/reject)
// Endpoint: POST /auth/internal/transfers/{transferId}/respond
func (c *Client) RespondToTransfer(ctx context.Context, transferID string, req models.RespondTransferRequest) (*models.RespondTransferResponse, error) {
	path := fmt.Sprintf("/auth/internal/transfers/%s/respond", transferID)
	var resp models.RespondTransferResponse
	if err := c.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, fmt.Errorf("respond to transfer: %w", err)
	}
	return &resp, nil
}

// CancelTransfer cancels an ownership transfer
// Endpoint: POST /auth/internal/transfers/{transferId}/cancel
func (c *Client) CancelTransfer(ctx context.Context, transferID string) (*models.CancelTransferResponse, error) {
	path := fmt.Sprintf("/auth/internal/transfers/%s/cancel", transferID)
	var resp models.CancelTransferResponse
	if err := c.doRequest(ctx, "POST", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("cancel transfer: %w", err)
	}
	return &resp, nil
}

// GetTransfer retrieves an ownership transfer by ID
// Endpoint: GET /auth/internal/transfers/{transferId}
func (c *Client) GetTransfer(ctx context.Context, transferID string) (*models.TransferResponse, error) {
	path := fmt.Sprintf("/auth/internal/transfers/%s", transferID)
	var resp models.TransferResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get transfer: %w", err)
	}
	return &resp, nil
}

// ============================================
// CACHE MANAGEMENT
// ============================================

// RefreshCache publishes a cache invalidation event
// Endpoint: POST /auth/internal/cache/refresh
func (c *Client) RefreshCache(ctx context.Context, req models.CacheRefreshRequest) (*models.CacheRefreshResponse, error) {
	var resp models.CacheRefreshResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCacheRefresh.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("refresh cache: %w", err)
	}
	return &resp, nil
}

// GetCacheStatus retrieves cache status
// Endpoint: GET /auth/internal/cache/status
func (c *Client) GetCacheStatus(ctx context.Context) (*models.CacheStatusResponse, error) {
	var resp models.CacheStatusResponse
	if err := c.doRequest(ctx, "GET", endpoints.AuthCacheStatus.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get cache status: %w", err)
	}
	return &resp, nil
}

// ============================================
// NOTIFICATIONS
// ============================================

// CreateNotification creates a notification for a user
// Endpoint: POST /auth/internal/notifications
func (c *Client) CreateNotification(ctx context.Context, req models.CreateNotificationRequest) (*models.CreateNotificationResponse, error) {
	var resp models.CreateNotificationResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthCreateNotification.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return &resp, nil
}

// NotifyTeamMatchLobby notifies team members about a lobby match
// Endpoint: POST /auth/internal/notifications/team-match-lobby
func (c *Client) NotifyTeamMatchLobby(ctx context.Context, req models.TeamMatchLobbyNotificationRequest) (*models.TeamMatchLobbyNotificationResponse, error) {
	var resp models.TeamMatchLobbyNotificationResponse
	if err := c.doRequest(ctx, "POST", endpoints.AuthTeamMatchLobbyNotification.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("notify team match lobby: %w", err)
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
			Msg("Connect-Auth HTTP request failed")
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
			Msg("Connect-Auth failed to read response body")
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
		Msg("Connect-Auth internal API error")

	return errors.NewInternalError(statusCode, "Connect-Auth", path, errMsg)
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
			Msg("Connect-Auth failed to parse response JSON")
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on Connect-Auth
func (c *Client) HealthCheck(ctx context.Context) error {
	// Use catalog endpoint as health check
	_, err := c.GetCatalog(ctx)
	return err
}

// ============================================
// MEMBERSHIP & CAPABILITIES
// ============================================

// GetMembership obtiene la información de membresía de un usuario en un scope
// Endpoint: GET /auth/internal/memberships/user/{userId}/scope/{scopeId}
func (c *Client) GetMembership(ctx context.Context, userID, scopeID string) (*models.MembershipResponse, error) {
	path := fmt.Sprintf("/auth/internal/memberships/user/%s/scope/%s", userID, scopeID)
	var resp models.MembershipResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		// Si es 404, el usuario no tiene membresía
		if isNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("get membership: %w", err)
	}
	return &resp, nil
}

// GetCapabilities obtiene los permisos efectivos de un usuario en un scope
// Endpoint: GET /auth/internal/capabilities/user/{userId}/scope/{scopeId}
func (c *Client) GetCapabilities(ctx context.Context, userID, scopeID string) (*models.CapabilitiesResponse, error) {
	path := fmt.Sprintf("/auth/internal/capabilities/user/%s/scope/%s", userID, scopeID)
	var resp models.CapabilitiesResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get capabilities: %w", err)
	}
	return &resp, nil
}

// isNotFoundError verifica si un error es 404
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	// Verificar si el error contiene el código 404
	errStr := fmt.Sprintf("%v", err)
	return strings.Contains(errStr, "404") || strings.Contains(errStr, "not found")
}
