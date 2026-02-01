// Package http provides HTTP client utilities for the UniFi SDK.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/murasame29/unifi-client-go/pkg/config"
	"github.com/murasame29/unifi-client-go/pkg/errors"
)

// Client is the HTTP client for UniFi APIs.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	userAgent  string
}

// NewClient creates a new HTTP client from config.
func NewClient(cfg config.Config) Client {
	return Client{
		httpClient: cfg.HTTPClient,
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		userAgent:  cfg.UserAgent,
	}
}

// Get sends a GET request.
func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}

// Post sends a POST request.
func (c *Client) Post(ctx context.Context, path string, body, result interface{}) error {
	return c.do(ctx, http.MethodPost, path, body, result)
}

// Put sends a PUT request.
func (c *Client) Put(ctx context.Context, path string, body, result interface{}) error {
	return c.do(ctx, http.MethodPut, path, body, result)
}

// Delete sends a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, result interface{}) error {
	return c.do(ctx, http.MethodDelete, path, nil, result)
}

func (c *Client) do(ctx context.Context, method, path string, body, result interface{}) error {
	url := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return errors.NewAPIError(
			resp.StatusCode,
			string(respBody),
			resp.Header.Get("X-Request-Id"),
		)
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// SetBaseURL sets the base URL.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}
