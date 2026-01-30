package bookstack

import (
	"context"
	"net/url"
	"strconv"
)

// SearchService handles search operations.
type SearchService struct {
	client *Client
}

// Search performs a full-text search query across all content types.
// The query parameter uses Bookstack's search syntax.
func (s *SearchService) Search(ctx context.Context, query string, opts *ListOptions) ([]SearchResult, error) {
	v := url.Values{}
	v.Set("query", query)
	if opts != nil {
		if opts.Count > 0 {
			v.Set("count", strconv.Itoa(opts.Count))
		}
		if opts.Offset > 0 {
			v.Set("offset", strconv.Itoa(opts.Offset))
		}
	}

	var resp listResponse[SearchResult]
	err := s.client.do(ctx, "GET", "/api/search?"+v.Encode(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
