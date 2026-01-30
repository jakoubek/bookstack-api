package bookstack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestSearchService_Search(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/search" {
			t.Errorf("path = %s, want /api/search", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "test query" {
			t.Errorf("query = %q, want %q", r.URL.Query().Get("query"), "test query")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 1, "name": "Result One", "type": "page", "preview": "...test..."},
				{"id": 2, "name": "Result Two", "type": "book", "preview": "...test..."},
			},
			"total": 2,
		})
	})

	results, err := c.Search.Search(context.Background(), "test query", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Type != "page" {
		t.Errorf("results[0].Type = %q, want %q", results[0].Type, "page")
	}
	if results[1].Type != "book" {
		t.Errorf("results[1].Type = %q, want %q", results[1].Type, "book")
	}
}

func TestSearchService_Search_WithOptions(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("count") != "5" {
			t.Errorf("count = %q, want 5", q.Get("count"))
		}
		if q.Get("offset") != "10" {
			t.Errorf("offset = %q, want 10", q.Get("offset"))
		}
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "total": 0})
	})

	_, err := c.Search.Search(context.Background(), "test", &ListOptions{Count: 5, Offset: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
