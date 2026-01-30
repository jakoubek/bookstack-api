package bookstack

import (
	"context"
	"fmt"
	"iter"
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

// ListAll returns an iterator over all pages, handling pagination automatically.
func (s *PagesService) ListAll(ctx context.Context) iter.Seq2[Page, error] {
	return listAll[Page](ctx, s.client, "/api/pages")
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

// ExportMarkdown exports a page as markdown.
func (s *PagesService) ExportMarkdown(ctx context.Context, id int) ([]byte, error) {
	return s.client.doRaw(ctx, "GET", fmt.Sprintf("/api/pages/%d/export/markdown", id))
}

// ExportPDF exports a page as PDF.
func (s *PagesService) ExportPDF(ctx context.Context, id int) ([]byte, error) {
	return s.client.doRaw(ctx, "GET", fmt.Sprintf("/api/pages/%d/export/pdf", id))
}
