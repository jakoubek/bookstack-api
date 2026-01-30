package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestShelvesService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/shelves" {
			t.Errorf("path = %s, want /api/shelves", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []map[string]any{{"id": 1, "name": "Shelf One"}},
			"total": 1,
		})
	})

	shelves, err := c.Shelves.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shelves) != 1 {
		t.Fatalf("got %d shelves, want 1", len(shelves))
	}
}

func TestShelvesService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/shelves/7" {
			t.Errorf("path = %s, want /api/shelves/7", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": 7, "name": "My Shelf",
		})
	})

	shelf, err := c.Shelves.Get(context.Background(), 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if shelf.Name != "My Shelf" {
		t.Errorf("Name = %q, want %q", shelf.Name, "My Shelf")
	}
}

func TestShelvesService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Shelf not found"},
		})
	})

	_, err := c.Shelves.Get(context.Background(), 999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
