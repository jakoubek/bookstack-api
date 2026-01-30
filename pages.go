package bookstack

import (
	"context"
	"fmt"
)

// PagesService handles operations on pages.
type PagesService struct {
	client *Client
}

// List returns a list of pages with optional filtering.
func (s *PagesService) List(ctx context.Context, opts *ListOptions) ([]Page, error) {
	var resp listResponse[Page]
	err := s.client.do(ctx, "GET", "/api/pages"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single page by ID, including its content.
func (s *PagesService) Get(ctx context.Context, id int) (*Page, error) {
	var page Page
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/pages/%d", id), nil, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}
