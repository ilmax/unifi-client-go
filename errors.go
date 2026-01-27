package unifi

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the UniFi Site Manager API.
type APIError struct {
	StatusCode int    // HTTP status code
	Message    string // Error message
	RequestID  string // Request ID if available
	Err        error  // Underlying error
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("unifi api error: status=%d, message=%s, request_id=%s", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("unifi api error: status=%d, message=%s", e.StatusCode, e.Message)
}

// Unwrap returns the underlying error.
func (e *APIError) Unwrap() error {
	return e.Err
}

// IsAuthenticationError returns true if the error is an authentication error (401 Unauthorized).
func IsAuthenticationError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsPermissionError returns true if the error is a permission error (403 Forbidden).
func IsPermissionError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsNotFoundError returns true if the error is a not found error (404 Not Found).
func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error (429 Too Many Requests).
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsServerError returns true if the error is a server error (5xx).
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// IsConnectionError returns true if the error is a connection error.
func IsConnectionError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 0 && apiErr.Err != nil
	}
	return false
}

// NewAPIError creates a new APIError with the given parameters.
func NewAPIError(statusCode int, message string, requestID string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		RequestID:  requestID,
		Err:        err,
	}
}

// ErrEmptyAPIKey is returned when an empty API key is provided.
var ErrEmptyAPIKey = errors.New("api key cannot be empty")

// ErrInvalidInterval is returned when an invalid ISP metrics interval is provided.
var ErrInvalidInterval = errors.New("invalid ISP metrics interval: must be '5m' or '1h'")

// ErrEmptyConfigID is returned when an empty config ID is provided.
var ErrEmptyConfigID = errors.New("config ID cannot be empty")
