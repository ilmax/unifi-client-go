// Package errors provides error types for the UniFi SDK.
package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		expected string
	}{
		{
			name: "with request ID",
			err: &APIError{
				StatusCode: 401,
				Message:    "Unauthorized",
				RequestID:  "req-123",
			},
			expected: "API error (status 401, request_id: req-123): Unauthorized",
		},
		{
			name: "without request ID",
			err: &APIError{
				StatusCode: 500,
				Message:    "Internal Server Error",
				RequestID:  "",
			},
			expected: "API error (status 500): Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("APIError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(404, "Not Found", "req-456")

	if err.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", err.StatusCode)
	}
	if err.Message != "Not Found" {
		t.Errorf("Message = %q, want %q", err.Message, "Not Found")
	}
	if err.RequestID != "req-456" {
		t.Errorf("RequestID = %q, want %q", err.RequestID, "req-456")
	}
}

func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "401 error",
			err:      NewAPIError(http.StatusUnauthorized, "Unauthorized", ""),
			expected: true,
		},
		{
			name:     "403 error",
			err:      NewAPIError(http.StatusForbidden, "Forbidden", ""),
			expected: false,
		},
		{
			name:     "wrapped 401 error",
			err:      fmt.Errorf("wrapped: %w", NewAPIError(http.StatusUnauthorized, "Unauthorized", "")),
			expected: true,
		},
		{
			name:     "non-API error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAuthenticationError(tt.err); got != tt.expected {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNotFoundError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "404 error",
			err:      NewAPIError(http.StatusNotFound, "Not Found", ""),
			expected: true,
		},
		{
			name:     "500 error",
			err:      NewAPIError(http.StatusInternalServerError, "Internal Server Error", ""),
			expected: false,
		},
		{
			name:     "wrapped 404 error",
			err:      fmt.Errorf("wrapped: %w", NewAPIError(http.StatusNotFound, "Not Found", "")),
			expected: true,
		},
		{
			name:     "non-API error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFoundError(tt.err); got != tt.expected {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "429 error",
			err:      NewAPIError(http.StatusTooManyRequests, "Too Many Requests", ""),
			expected: true,
		},
		{
			name:     "503 error",
			err:      NewAPIError(http.StatusServiceUnavailable, "Service Unavailable", ""),
			expected: false,
		},
		{
			name:     "wrapped 429 error",
			err:      fmt.Errorf("wrapped: %w", NewAPIError(http.StatusTooManyRequests, "Too Many Requests", "")),
			expected: true,
		},
		{
			name:     "non-API error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimitError(tt.err); got != tt.expected {
				t.Errorf("IsRateLimitError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCommonErrors(t *testing.T) {
	// Verify common errors are defined and not nil
	if ErrEmptyAPIKey == nil {
		t.Error("ErrEmptyAPIKey should not be nil")
	}
	if ErrInvalidInterval == nil {
		t.Error("ErrInvalidInterval should not be nil")
	}
	if ErrEmptyConfigID == nil {
		t.Error("ErrEmptyConfigID should not be nil")
	}
	if ErrEmptyHostID == nil {
		t.Error("ErrEmptyHostID should not be nil")
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "with field name",
			err: &ValidationError{
				Field:   "email",
				Message: "must be a valid email address",
			},
			expected: `validation error on field "email": must be a valid email address`,
		},
		{
			name: "without field name",
			err: &ValidationError{
				Field:   "",
				Message: "invalid input",
			},
			expected: "validation error: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("ValidationError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("username", "cannot be empty")

	if err.Field != "username" {
		t.Errorf("Field = %q, want %q", err.Field, "username")
	}
	if err.Message != "cannot be empty" {
		t.Errorf("Message = %q, want %q", err.Message, "cannot be empty")
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "validation error",
			err:      NewValidationError("field", "message"),
			expected: true,
		},
		{
			name:     "wrapped validation error",
			err:      fmt.Errorf("wrapped: %w", NewValidationError("field", "message")),
			expected: true,
		},
		{
			name:     "API error",
			err:      NewAPIError(400, "Bad Request", ""),
			expected: false,
		},
		{
			name:     "non-validation error",
			err:      errors.New("some error"),
			expected: false,
		},
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidationError(tt.err); got != tt.expected {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.expected)
			}
		})
	}
}
