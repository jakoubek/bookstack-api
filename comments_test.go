package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestCommentsService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []map[string]any{{"id": 1, "page_id": 5, "html": "<p>Nice</p>"}},
			"total": 1,
		})
	})

	comments, err := c.Comments.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("got %d, want 1", len(comments))
	}
}

func TestCommentsService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "page_id": 5, "html": "<p>Comment</p>",
		})
	})

	comment, err := c.Comments.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.HTML != "<p>Comment</p>" {
		t.Errorf("HTML = %q", comment.HTML)
	}
}

func TestCommentsService_Create(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s, want POST", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": 2, "page_id": 5, "html": "<p>New</p>",
		})
	})

	comment, err := c.Comments.Create(context.Background(), &CommentCreateRequest{
		PageID: 5,
		HTML:   "<p>New</p>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 2 {
		t.Errorf("ID = %d, want 2", comment.ID)
	}
}

func TestCommentsService_Update(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": 1, "html": "<p>Updated</p>",
		})
	})

	comment, err := c.Comments.Update(context.Background(), 1, &CommentUpdateRequest{
		HTML: "<p>Updated</p>",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.HTML != "<p>Updated</p>" {
		t.Errorf("HTML = %q", comment.HTML)
	}
}

func TestCommentsService_Delete(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.Comments.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommentsService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Not found"},
		})
	})

	_, err := c.Comments.Get(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
