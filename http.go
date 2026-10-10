package bookstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// do executes an HTTP request against the BookStack API with the client's
// authentication header attached. The auth header format is
// "Token ***" — secret first, then id — matching the format
// expected by the BookStack API token guard.
func (c *Client) do(ctx context.Context, method, path string, body interface{}, out interface{}) error {
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return decodeAPIError(resp, respBody)
	}

	if out == nil {
		return nil
	}
	return json.Unmarshal(respBody, out)
}

// doRaw executes an HTTP request and returns the raw response body (for binary
// exports such as PDF/Markdown).
func (c *Client) doRaw(ctx context.Context, method, path string) ([]byte, error) {
	req, err := c.newRequest(ctx, method, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, decodeAPIError(resp, respBody)
	}
	return respBody, nil
}

// newRequest builds an *http.Request with the BookStack auth header attached.
//
// Auth format: "Token ***" — secret first, then id.
// This matches the BookStack ApiTokenGuard, which splits the header on ':'
// and treats the first segment as the token secret and the second as the id:
//
//	[$id, $secret] = explode(':', str_replace('Token ', '', $authToken));
func (c *Client) newRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %q: %w", path, err)
	}

	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return nil, err
	}

	// BookStack expects the token secret BEFORE the id, separated by ':'.
	// See: app/Api/ApiTokenGuard.php -> explode(':', str_replace('Token ', '', $authToken))
	req.Header.Set("Authorization", fmt.Sprintf("Token %s:%s", c.tokenSecret, c.tokenID))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// decodeAPIError parses a BookStack error response into an *APIError.
// When the body is not a JSON error payload, the HTTP status text is used as
// the message (matching the behaviour expected by existing tests).
func decodeAPIError(resp *http.Response, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}
	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error.Message != "" {
		apiErr.Code = payload.Error.Code
		apiErr.Message = payload.Error.Message
	} else {
		apiErr.Message = http.StatusText(resp.StatusCode)
	}
	return apiErr
}