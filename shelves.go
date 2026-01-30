package bookstack

import (
	"context"
	"fmt"
	"iter"
)

// ShelvesService handles operations on shelves.
type ShelvesService struct {
	client *Client
}

// List returns a list of shelves with optional filtering.
func (s *ShelvesService) List(ctx context.Context, opts *ListOptions) ([]Shelf, error) {
	var resp listResponse[Shelf]
	err := s.client.do(ctx, "GET", "/api/shelves"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListAll returns an iterator over all shelves, handling pagination automatically.
func (s *ShelvesService) ListAll(ctx context.Context) iter.Seq2[Shelf, error] {
	return listAll[Shelf](ctx, s.client, "/api/shelves")
}

// Get retrieves a single shelf by ID.
func (s *ShelvesService) Get(ctx context.Context, id int) (*Shelf, error) {
	var shelf Shelf
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/shelves/%d", id), nil, &shelf)
	if err != nil {
		return nil, err
	}
	return &shelf, nil
}
