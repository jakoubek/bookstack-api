package bookstack

import (
	"context"
	"fmt"
)

// BooksService handles operations on books.
type BooksService struct {
	client *Client
}

// List returns a list of books with optional filtering.
func (s *BooksService) List(ctx context.Context, opts *ListOptions) ([]Book, error) {
	var resp listResponse[Book]
	err := s.client.do(ctx, "GET", "/api/books"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single book by ID.
func (s *BooksService) Get(ctx context.Context, id int) (*Book, error) {
	var book Book
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/books/%d", id), nil, &book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}
