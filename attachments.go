package bookstack

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
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

// Create creates a new attachment.
//
// For link attachments, set Link and leave File nil — the request is sent as JSON.
// For file attachments, set File and Filename and leave Link empty — the request
// is sent as multipart/form-data.
func (s *AttachmentsService) Create(ctx context.Context, req *AttachmentCreateRequest) (*Attachment, error) {
	var a Attachment
	if req.File != nil {
		// Multipart upload path
		body, contentType, err := s.buildMultipart(req)
		if err != nil {
			return nil, err
		}
		err = s.client.doMultipart(ctx, "POST", "/api/attachments", body, contentType, &a)
		if err != nil {
			return nil, err
		}
	} else {
		// JSON path (link attachments)
		err := s.client.do(ctx, "POST", "/api/attachments", req, &a)
		if err != nil {
			return nil, err
		}
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

// buildMultipart constructs a multipart/form-data request body for file uploads.
func (s *AttachmentsService) buildMultipart(req *AttachmentCreateRequest) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	if err := w.WriteField("name", req.Name); err != nil {
		return nil, "", fmt.Errorf("writing name field: %w", err)
	}
	if err := w.WriteField("uploaded_to", fmt.Sprintf("%d", req.UploadedTo)); err != nil {
		return nil, "", fmt.Errorf("writing uploaded_to field: %w", err)
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, req.Filename))
	h.Set("Content-Type", "application/octet-stream")

	fw, err := w.CreatePart(h)
	if err != nil {
		return nil, "", fmt.Errorf("creating file part: %w", err)
	}
	if _, err := io.Copy(fw, req.File); err != nil {
		return nil, "", fmt.Errorf("copying file data: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, "", fmt.Errorf("closing multipart writer: %w", err)
	}

	return &buf, w.FormDataContentType(), nil
}
