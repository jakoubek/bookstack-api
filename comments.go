package bookstack

import (
	"context"
	"fmt"
)

// CommentsService handles operations on comments.
type CommentsService struct {
	client *Client
}

// List returns a list of comments with optional filtering.
func (s *CommentsService) List(ctx context.Context, opts *ListOptions) ([]Comment, error) {
	var resp listResponse[Comment]
	err := s.client.do(ctx, "GET", "/api/comments"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single comment by ID.
func (s *CommentsService) Get(ctx context.Context, id int) (*Comment, error) {
	var c Comment
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/comments/%d", id), nil, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Create creates a new comment on a page.
func (s *CommentsService) Create(ctx context.Context, req *CommentCreateRequest) (*Comment, error) {
	var c Comment
	err := s.client.do(ctx, "POST", "/api/comments", req, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update updates an existing comment.
func (s *CommentsService) Update(ctx context.Context, id int, req *CommentUpdateRequest) (*Comment, error) {
	var c Comment
	err := s.client.do(ctx, "PUT", fmt.Sprintf("/api/comments/%d", id), req, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Delete deletes a comment by ID.
func (s *CommentsService) Delete(ctx context.Context, id int) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/api/comments/%d", id), nil, nil)
}
