package bookstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// do executes an authenticated API request and unmarshals the response.
// method is the HTTP method, path is appended to BaseURL (e.g., "/api/books"),
// body is JSON-encoded as the request body (nil for no body),
// and result is the target for JSON unmarshaling (nil to discard response body).
func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s:%s", c.tokenID, c.tokenSecret))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
		// Try to parse error details from JSON response
		var errResp struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error.Message != "" {
			apiErr.Code = errResp.Error.Code
			apiErr.Message = errResp.Error.Message
		} else {
			apiErr.Message = http.StatusText(resp.StatusCode)
		}
		return apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

// listResponse wraps the common Bookstack list API response format.
type listResponse[T any] struct {
	Data  []T `json:"data"`
	Total int `json:"total"`
}

// ListOptions contains common options for list operations.
type ListOptions struct {
	Count  int               // Max items per page (default 100, max 500)
	Offset int               // Offset for pagination
	Sort   string            // Sort field (e.g., "name", "-created_at")
	Filter map[string]string // Filters (e.g., {"name": "value"})
}

// queryString builds a URL query string from ListOptions.
func (o *ListOptions) queryString() string {
	if o == nil {
		return ""
	}
	v := url.Values{}
	if o.Count > 0 {
		v.Set("count", strconv.Itoa(o.Count))
	}
	if o.Offset > 0 {
		v.Set("offset", strconv.Itoa(o.Offset))
	}
	if o.Sort != "" {
		v.Set("sort", o.Sort)
	}
	for key, val := range o.Filter {
		v.Set("filter["+key+"]", val)
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}
