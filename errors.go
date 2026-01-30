package bookstack

import (
	"errors"
	"fmt"
)

// Sentinel errors for common API error conditions.
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrRateLimited  = errors.New("rate limited")
	ErrBadRequest   = errors.New("bad request")
)

// APIError represents an error returned by the Bookstack API.
type APIError struct {
	StatusCode int    // HTTP status code
	Code       string // Error code from API response
	Message    string // Error message from API response
	Body       string // Raw response body
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("bookstack API error (status %d, code %s): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("bookstack API error (status %d): %s", e.StatusCode, e.Message)
}

// Is implements error matching for sentinel errors.
func (e *APIError) Is(target error) bool {
	switch target {
	case ErrNotFound:
		return e.StatusCode == 404
	case ErrUnauthorized:
		return e.StatusCode == 401
	case ErrForbidden:
		return e.StatusCode == 403
	case ErrRateLimited:
		return e.StatusCode == 429
	case ErrBadRequest:
		return e.StatusCode == 400
	default:
		return false
	}
}
