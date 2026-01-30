package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestChaptersService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chapters" {
			t.Errorf("path = %s, want /api/chapters", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []map[string]any{{"id": 1, "name": "Ch 1", "book_id": 1}},
			"total": 1,
		})
	})

	chapters, err := c.Chapters.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chapters) != 1 {
		t.Fatalf("got %d chapters, want 1", len(chapters))
	}
}

func TestChaptersService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chapters/3" {
			t.Errorf("path = %s, want /api/chapters/3", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": 3, "name": "Chapter Three", "book_id": 1,
		})
	})

	ch, err := c.Chapters.Get(context.Background(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch.Name != "Chapter Three" {
		t.Errorf("Name = %q, want %q", ch.Name, "Chapter Three")
	}
}

func TestChaptersService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Chapter not found"},
		})
	})

	_, err := c.Chapters.Get(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
