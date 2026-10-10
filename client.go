package bookstack

// SetToken updates the client's API token credentials.
func (c *Client) SetToken(id, secret string) {
	c.tokenID = id
	c.tokenSecret = secret
}

// TokenID returns the current token ID.
func (c *Client) TokenID() string {
	return c.tokenID
}

// TokenSecret returns the current token secret.
func (c *Client) TokenSecret() string {
	return c.tokenSecret
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}