package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestPagesService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pages" {
			t.Errorf("path = %s, want /api/pages", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 1, "name": "Page One", "book_id": 1},
				{"id": 2, "name": "Page Two", "book_id": 1},
			},
			"total": 2,
		})
	})

	pages, err := c.Pages.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(pages))
	}
	if pages[0].Name != "Page One" {
		t.Errorf("pages[0].Name = %q, want %q", pages[0].Name, "Page One")
	}
}

func TestPagesService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pages/5" {
			t.Errorf("path = %s, want /api/pages/5", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":      5,
			"name":    "Test Page",
			"book_id": 1,
			"html":    "<p>Content</p>",
		})
	})

	page, err := c.Pages.Get(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if page.ID != 5 {
		t.Errorf("ID = %d, want 5", page.ID)
	}
	if page.HTML != "<p>Content</p>" {
		t.Errorf("HTML = %q", page.HTML)
	}
}

func TestPagesService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Page not found"},
		})
	})

	_, err := c.Pages.Get(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound")
	}
}
