package bookstack

import (
	"context"
	"fmt"
)

// AttachmentsService handles operations on attachments.
type AttachmentsService struct {
	client *Client
}

// List returns a list of attachments with optional filtering.
func (s *AttachmentsService) List(ctx context.Context, opts *ListOptions) ([]Attachment, error) {
	var resp listResponse[Attachment]
	err := s.client.do(ctx, "GET", "/api/attachments"+opts.queryString(), nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// Get retrieves a single attachment by ID.
func (s *AttachmentsService) Get(ctx context.Context, id int) (*Attachment, error) {
	var a Attachment
	err := s.client.do(ctx, "GET", fmt.Sprintf("/api/attachments/%d", id), nil, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Create creates a new link attachment.
func (s *AttachmentsService) Create(ctx context.Context, req *AttachmentCreateRequest) (*Attachment, error) {
	var a Attachment
	err := s.client.do(ctx, "POST", "/api/attachments", req, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Update updates an existing attachment.
func (s *AttachmentsService) Update(ctx context.Context, id int, req *AttachmentUpdateRequest) (*Attachment, error) {
	var a Attachment
	err := s.client.do(ctx, "PUT", fmt.Sprintf("/api/attachments/%d", id), req, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Delete deletes an attachment by ID.
func (s *AttachmentsService) Delete(ctx context.Context, id int) error {
	return s.client.do(ctx, "DELETE", fmt.Sprintf("/api/attachments/%d", id), nil, nil)
}
