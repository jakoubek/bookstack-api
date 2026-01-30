package bookstack

import "context"

// SearchService handles search operations.
type SearchService struct {
	client *Client
}

// Search performs a search query across all content types.
// TODO: Implement API call to GET /api/search
func (s *SearchService) Search(ctx context.Context, query string, opts *ListOptions) ([]SearchResult, error) {
	// Placeholder for future implementation
	return nil, nil
}
