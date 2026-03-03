package settingsruntime

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CoreHTTPClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
	baseURI *url.URL
}

func NewCoreHTTPClient(baseURL, apiKey string, timeout time.Duration) (*CoreHTTPClient, error) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return nil, errors.New("core base URL is required")
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid core base URL: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid core base URL: %s", trimmed)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("invalid core base URL scheme: %s", parsed.Scheme)
	}
	if parsed.User != nil {
		return nil, errors.New("core base URL must not include user info")
	}

	if timeout <= 0 {
		timeout = defaultHTTPTimeout
	}

	normalizedBaseURL := strings.TrimRight(trimmed, "/")
	baseURI, err := url.Parse(normalizedBaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid normalized core base URL: %w", err)
	}

	return &CoreHTTPClient{
		baseURL: normalizedBaseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: timeout},
		baseURI: baseURI,
	}, nil
}

func (c *CoreHTTPClient) Health(ctx context.Context) error {
	if c == nil {
		return errors.New("core HTTP client is nil")
	}

	healthURL, err := c.buildInternalURL("/health")
	if err != nil {
		return fmt.Errorf("build health URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
	if err != nil {
		return fmt.Errorf("build health request: %w", err)
	}

	// #nosec G704 -- URL is canonicalized and constrained in buildInternalURL.
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("core health request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("core health check returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *CoreHTTPClient) GetSettingValue(ctx context.Context, entity, key string) (string, error) {
	if c == nil {
		return "", errors.New("core HTTP client is nil")
	}

	entity = strings.TrimSpace(strings.ToUpper(entity))
	key = strings.TrimSpace(key)
	if entity == "" || key == "" {
		return "", errors.New("entity and key are required")
	}

	path := fmt.Sprintf("/core/internal/settings/%s/%s/value", url.PathEscape(entity), url.PathEscape(key))
	requestURL, err := c.buildInternalURL(path)
	if err != nil {
		return "", fmt.Errorf("build setting URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("build setting request: %w", err)
	}

	if strings.TrimSpace(c.apiKey) != "" {
		req.Header.Set("X-Internal-API-Key", c.apiKey)
	}

	// #nosec G704 -- URL is canonicalized and constrained in buildInternalURL.
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("get setting %s/%s failed: %w", entity, key, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("get setting %s/%s returned status %d", entity, key, resp.StatusCode)
	}

	var payload struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode setting %s/%s: %w", entity, key, err)
	}

	return payload.Value, nil
}

func (c *CoreHTTPClient) buildInternalURL(path string) (string, error) {
	if c == nil || c.baseURI == nil {
		return "", errors.New("core HTTP client base URL is not initialized")
	}
	if strings.Contains(path, "://") {
		return "", errors.New("absolute paths are not allowed")
	}

	joined, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return "", err
	}

	parsed, err := url.Parse(joined)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != c.baseURI.Scheme || parsed.Host != c.baseURI.Host {
		return "", errors.New("resolved URL host mismatch")
	}

	return parsed.String(), nil
}
