package bookstack

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestBooksService_List(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/api/books" {
			t.Errorf("path = %s, want /api/books", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"id": 1, "name": "Book One", "slug": "book-one"},
				{"id": 2, "name": "Book Two", "slug": "book-two"},
			},
			"total": 2,
		})
	})

	books, err := c.Books.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("got %d books, want 2", len(books))
	}
	if books[0].Name != "Book One" {
		t.Errorf("books[0].Name = %q, want %q", books[0].Name, "Book One")
	}
}

func TestBooksService_List_WithOptions(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("count") != "10" {
			t.Errorf("count = %q, want 10", q.Get("count"))
		}
		if q.Get("offset") != "20" {
			t.Errorf("offset = %q, want 20", q.Get("offset"))
		}
		if q.Get("sort") != "-name" {
			t.Errorf("sort = %q, want -name", q.Get("sort"))
		}
		json.NewEncoder(w).Encode(map[string]any{"data": []any{}, "total": 0})
	})

	_, err := c.Books.List(context.Background(), &ListOptions{
		Count:  10,
		Offset: 20,
		Sort:   "-name",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBooksService_Get(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/books/42" {
			t.Errorf("path = %s, want /api/books/42", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": 42, "name": "My Book", "slug": "my-book",
		})
	})

	book, err := c.Books.Get(context.Background(), 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.ID != 42 {
		t.Errorf("ID = %d, want 42", book.ID)
	}
	if book.Name != "My Book" {
		t.Errorf("Name = %q, want %q", book.Name, "My Book")
	}
}

func TestBooksService_Get_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]string{"message": "Book not found"},
		})
	})

	_, err := c.Books.Get(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound")
	}
}
