package bookstack

import (
	"errors"
	"fmt"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "with code",
			err:  APIError{StatusCode: 404, Code: "not_found", Message: "Page not found"},
			want: "bookstack API error (status 404, code not_found): Page not found",
		},
		{
			name: "without code",
			err:  APIError{StatusCode: 500, Message: "Internal error"},
			want: "bookstack API error (status 500): Internal error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_Is(t *testing.T) {
	tests := []struct {
		status int
		target error
		want   bool
	}{
		{400, ErrBadRequest, true},
		{401, ErrUnauthorized, true},
		{403, ErrForbidden, true},
		{404, ErrNotFound, true},
		{429, ErrRateLimited, true},
		{500, ErrNotFound, false},
		{200, ErrBadRequest, false},
	}
	for _, tt := range tests {
		apiErr := &APIError{StatusCode: tt.status}
		if got := errors.Is(apiErr, tt.target); got != tt.want {
			t.Errorf("errors.Is(APIError{%d}, %v) = %v, want %v", tt.status, tt.target, got, tt.want)
		}
	}
}

func TestAPIError_Is_Wrapped(t *testing.T) {
	inner := &APIError{StatusCode: 404, Message: "not found"}
	wrapped := fmt.Errorf("fetching page: %w", inner)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("expected errors.Is to match ErrNotFound through wrapping")
	}
}
