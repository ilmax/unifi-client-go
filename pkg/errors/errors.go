// Package errors provides error types for the UniFi SDK.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Common errors
var (
	ErrEmptyAPIKey     = errors.New("API key cannot be empty")
	ErrInvalidInterval = errors.New("invalid ISP metrics interval: must be '5m' or '1h'")
	ErrEmptyConfigID   = errors.New("config ID cannot be empty")
	ErrEmptyHostID     = errors.New("host ID cannot be empty")
)

// APIError represents an error returned by the UniFi API.
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
}

// NewAPIError creates a new APIError.
func NewAPIError(statusCode int, message, requestID string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		RequestID:  requestID,
	}
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("API error (status %d, request_id: %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// IsAuthenticationError returns true if the error is an authentication error (401).
func IsAuthenticationError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsNotFoundError returns true if the error is a not found error (404).
func IsNotFoundError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsRateLimitError returns true if the error is a rate limit error (429).
func IsRateLimitError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// ValidationError represents an input validation error.
type ValidationError struct {
	Field   string
	Message string
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// IsValidationError returns true if the error is a validation error.
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}
