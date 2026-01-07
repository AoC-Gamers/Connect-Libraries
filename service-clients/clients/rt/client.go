package rt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AoC-Gamers/connect-libraries/core-types/endpoints"
	"github.com/AoC-Gamers/connect-libraries/core-types/errors"
	"github.com/AoC-Gamers/connect-libraries/core-types/models"
	"github.com/rs/zerolog/log"
)

// Client is a type-safe HTTP client for Connect-RT internal API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Connect-RT client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // RT should be faster
		},
	}
}

// ============================================
// PRESENCE MANAGEMENT
// ============================================

// GetUserPresence retrieves a user's presence by steamID
// Endpoint: GET /rt/internal/presence/{steamid}
// Returns nil presence if user is not online
func (c *Client) GetUserPresence(ctx context.Context, steamID string) (*models.RTUserPresence, error) {
	path := fmt.Sprintf("/rt/internal/presence/%s", steamID)

	var resp models.GetUserPresenceResponse
	err := c.doRequest(ctx, "GET", path, nil, &resp)
	if err != nil {
		// Check if it's a 404 (user not found) - this is normal
		if ierr, ok := err.(*errors.InternalError); ok && ierr.StatusCode == 404 {
			log.Debug().Str("steamid", steamID).Msg("User presence not found (offline)")
			return nil, nil // Not an error, user is just offline
		}
		return nil, fmt.Errorf("get user presence: %w", err)
	}

	return resp.Presence, nil
}

// InitializePresence initializes a user's presence after login
// Endpoint: POST /rt/internal/presence/{steamid}
func (c *Client) InitializePresence(ctx context.Context, steamID string, req models.InitializePresenceRequest) (*models.InitializePresenceResponse, error) {
	path := fmt.Sprintf("/rt/internal/presence/%s", steamID)
	var resp models.InitializePresenceResponse
	if err := c.doRequest(ctx, "POST", path, req, &resp); err != nil {
		return nil, fmt.Errorf("initialize presence: %w", err)
	}
	return &resp, nil
}

// BatchGetPresence retrieves presence for multiple users (max 100)
// Endpoint: POST /rt/internal/presence/batch
func (c *Client) BatchGetPresence(ctx context.Context, steamIDs []string) (map[string]*models.RTUserPresence, error) {
	if len(steamIDs) == 0 {
		return make(map[string]*models.RTUserPresence), nil
	}

	if len(steamIDs) > 100 {
		return nil, fmt.Errorf("maximum 100 steamids per request, got %d", len(steamIDs))
	}

	req := models.BatchGetPresenceRequest{
		SteamIDs: steamIDs,
	}

	var resp models.BatchGetPresenceResponse
	if err := c.doRequest(ctx, "POST", endpoints.RTBatchGetPresence.Path, req, &resp); err != nil {
		return nil, fmt.Errorf("batch get presence: %w", err)
	}

	return resp.Presences, nil
}

// GetOnlineUsers retrieves all currently online users
// Endpoint: GET /rt/internal/presence/online
func (c *Client) GetOnlineUsers(ctx context.Context) (*models.GetOnlineUsersResponse, error) {
	var resp models.GetOnlineUsersResponse
	if err := c.doRequest(ctx, "GET", endpoints.RTGetOnlineUsers.Path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get online users: %w", err)
	}
	return &resp, nil
}

// GetUsersByStatus retrieves users by presence status
// Endpoint: GET /rt/internal/presence/by-status/{status}
// Valid statuses: "online", "away", "busy", "offline"
func (c *Client) GetUsersByStatus(ctx context.Context, status string) (*models.GetUsersByStatusResponse, error) {
	path := fmt.Sprintf("/rt/internal/presence/by-status/%s", status)
	var resp models.GetUsersByStatusResponse
	if err := c.doRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("get users by status: %w", err)
	}
	return &resp, nil
}

// CreateNotification creates a notification for a user via RT internal API
// Endpoint: POST /rt/internal/notifications
func (c *Client) CreateNotification(ctx context.Context, req models.RTCreateNotificationRequest) (*models.RTCreateNotificationResponse, error) {
	var resp models.RTCreateNotificationResponse
	if err := c.doRequest(ctx, "POST", "/rt/internal/notifications", req, &resp); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return &resp, nil
}

// CreateEnrichedNotification creates an enriched notification via RT internal API
// RT will fetch actor data from Core, enrich metadata, save to DB, and publish to NATS
// Endpoint: POST /rt/internal/notifications/enriched
func (c *Client) CreateEnrichedNotification(ctx context.Context, req models.RTCreateEnrichedNotificationRequest) (*models.RTCreateEnrichedNotificationResponse, error) {
	var resp models.RTCreateEnrichedNotificationResponse
	if err := c.doRequest(ctx, "POST", "/rt/internal/notifications/enriched", req, &resp); err != nil {
		return nil, fmt.Errorf("create enriched notification: %w", err)
	}
	return &resp, nil
}

// CreateOwnershipTransferredNotification handles ownership transfer event via RT internal API
// RT will fetch scope name, create notifications for affected users, and publish to NATS
// Endpoint: POST /rt/internal/events/ownership/transferred
func (c *Client) CreateOwnershipTransferredNotification(ctx context.Context, req models.RTOwnershipTransferredRequest) (*models.RTOwnershipTransferredResponse, error) {
	var resp models.RTOwnershipTransferredResponse
	if err := c.doRequest(ctx, "POST", "/rt/internal/events/ownership/transferred", req, &resp); err != nil {
		return nil, fmt.Errorf("create ownership transferred notification: %w", err)
	}
	return &resp, nil
}

// CreateOwnershipDemotionNotification handles ownership demotion event via RT internal API
// RT will create notification for the demoted user
// Endpoint: POST /rt/internal/events/ownership/demotion
func (c *Client) CreateOwnershipDemotionNotification(ctx context.Context, req models.RTOwnershipDemotionRequest) (*models.RTOwnershipDemotionResponse, error) {
	var resp models.RTOwnershipDemotionResponse
	if err := c.doRequest(ctx, "POST", "/rt/internal/events/ownership/demotion", req, &resp); err != nil {
		return nil, fmt.Errorf("create ownership demotion notification: %w", err)
	}
	return &resp, nil
}

// ============================================
// HTTP CLIENT HELPERS
// ============================================

// doRequest performs an HTTP request with proper error handling
func (c *Client) doRequest(ctx context.Context, method, path string, reqBody, respBody interface{}) error {
	url := c.baseURL + path

	bodyReader, err := c.marshalRequestBody(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	c.setRequestHeaders(req, reqBody)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("url", url).
			Msg("Connect-RT HTTP request failed")
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("method", method).
			Str("path", path).
			Int("status", resp.StatusCode).
			Msg("Connect-RT failed to read response body")
		return fmt.Errorf("read response: %w", err)
	}

	if err := c.checkStatusCode(resp.StatusCode, method, path, respBodyBytes); err != nil {
		return err
	}

	return c.parseResponseBody(respBody, respBodyBytes, method, path, resp.StatusCode)
}

// marshalRequestBody marshals the request body to JSON
func (c *Client) marshalRequestBody(reqBody interface{}) (io.Reader, error) {
	if reqBody == nil {
		return nil, nil
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	return bytes.NewReader(bodyBytes), nil
}

// setRequestHeaders sets the required headers for the request
func (c *Client) setRequestHeaders(req *http.Request, reqBody interface{}) {
	req.Header.Set("X-API-Key", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
}

// checkStatusCode validates the HTTP status code and returns an error if unsuccessful
func (c *Client) checkStatusCode(statusCode int, method, path string, respBodyBytes []byte) error {
	if statusCode >= 200 && statusCode < 300 {
		return nil
	}

	errMsg := c.extractErrorMessage(respBodyBytes)

	log.Error().
		Str("method", method).
		Str("path", path).
		Int("status", statusCode).
		Str("error", errMsg).
		Msg("Connect-RT internal API error")

	return errors.NewInternalError(statusCode, "Connect-RT", path, errMsg)
}

// extractErrorMessage extracts error message from response body
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
func (c *Client) parseResponseBody(respBody interface{}, respBodyBytes []byte, method, path string, statusCode int) error {
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
			Msg("Connect-RT failed to parse response JSON")
		return fmt.Errorf("unmarshal response: %w", err)
	}

	return nil
}

// HealthCheck performs a health check on Connect-RT
func (c *Client) HealthCheck(ctx context.Context) error {
	// Use get online users endpoint as health check
	_, err := c.GetOnlineUsers(ctx)
	return err
}
