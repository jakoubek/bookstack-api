package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestAttachmentsService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/attachments" {
			t.Errorf("path = %s, want /api/attachments", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []map[string]any{{"id": 1, "name": "file.pdf", "uploaded_to": 5}},
			"total": 1,
		})
	})

	attachments, err := c.Attachments.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(attachments) != 1 {
		t.Fatalf("got %d, want 1", len(attachments))
	}
	if attachments[0].Name != "file.pdf" {
		t.Errorf("Name = %q", attachments[0].Name)
	}
}

func TestAttachmentsService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "file.pdf", "content": "https://example.com/file.pdf",
		})
	})

	a, err := c.Attachments.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Content != "https://example.com/file.pdf" {
		t.Errorf("Content = %q", a.Content)
	}
}

func TestAttachmentsService_Create(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": 2, "name": "link.pdf", "external": true,
		})
	})

	a, err := c.Attachments.Create(context.Background(), &AttachmentCreateRequest{
		Name:       "link.pdf",
		UploadedTo: 5,
		Link:       "https://example.com/link.pdf",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != 2 {
		t.Errorf("ID = %d, want 2", a.ID)
	}
}

func TestAttachmentsService_Update(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "name": "renamed.pdf",
		})
	})

	a, err := c.Attachments.Update(context.Background(), 1, &AttachmentUpdateRequest{
		Name: "renamed.pdf",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "renamed.pdf" {
		t.Errorf("Name = %q", a.Name)
	}
}

func TestAttachmentsService_Delete(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.Attachments.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAttachmentsService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Not found"},
		})
	})

	_, err := c.Attachments.Get(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
