package bookstack

import (
	"context"
	"fmt"
	"iter"
)

// ChaptersService handles operations on chapters.
type ChaptersService struct {
	client *Client
}

// List returns a list of chapters with optional filtering.
func (s *ChaptersService) List(ctx context.Context, opts *ListOptions) ([]Chapter, error) {
	var resp listResponse[Chapter]
	err := s.client.do(ctx, "GET", "/api/chapters"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListAll returns an iterator over all chapters, handling pagination automatically.
func (s *ChaptersService) ListAll(ctx context.Context) iter.Seq2[Chapter, error] {
	return listAll[Chapter](ctx, s.client, "/api/chapters")
}

// Get retrieves a single chapter by ID.
func (s *ChaptersService) Get(ctx context.Context, id int) (*Chapter, error) {
	var chapter Chapter
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/chapters/%d", id), nil, &chapter)
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}
