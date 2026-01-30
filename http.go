package bookstack

import (
	"context"
	"net/http"
)

// buildRequest creates an HTTP request with proper authentication headers.
// TODO: Implement request building with Authorization header (Token <id>:<secret>)
func (c *Client) buildRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	// Placeholder for future implementation
	return nil, nil
}

// doRequest executes an HTTP request and handles the response.
// TODO: Implement response handling, error parsing, and JSON unmarshaling
func (c *Client) doRequest(ctx context.Context, req *http.Request, v interface{}) error {
	// Placeholder for future implementation
	return nil
}

// ListOptions contains common options for list operations.
type ListOptions struct {
	Count  int               // Max items per page (default 100, max 500)
	Offset int               // Offset for pagination
	Sort   string            // Sort field (e.g., "name", "-created_at")
	Filter map[string]string // Filters (e.g., {"name": "value"})
}
