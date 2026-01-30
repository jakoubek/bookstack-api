package bookstack

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestBooksService_ListAll(t *testing.T) {
	callCount := 0
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		callCount++
		offset := r.URL.Query().Get("offset")
		switch offset {
		case "0":
			json.NewEncoder(w).Encode(map[string]any{
				"data":  []map[string]any{{"id": 1, "name": "Book 1"}, {"id": 2, "name": "Book 2"}},
				"total": 3,
			})
		case "2":
			json.NewEncoder(w).Encode(map[string]any{
				"data":  []map[string]any{{"id": 3, "name": "Book 3"}},
				"total": 3,
			})
		default:
			t.Errorf("unexpected offset %q", offset)
		}
	})

	var books []Book
	for book, err := range c.Books.ListAll(context.Background()) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		books = append(books, book)
	}

	if len(books) != 3 {
		t.Fatalf("got %d books, want 3", len(books))
	}
	if callCount != 2 {
		t.Errorf("API called %d times, want 2", callCount)
	}
}

func TestBooksService_ListAll_EarlyBreak(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []map[string]any{{"id": 1}, {"id": 2}, {"id": 3}},
			"total": 100,
		})
	})

	count := 0
	for _, err := range c.Books.ListAll(context.Background()) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("got %d items, want 2 (early break)", count)
	}
}

func TestBooksService_ListAll_Error(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	})

	for _, err := range c.Books.ListAll(context.Background()) {
		if err == nil {
			t.Fatal("expected error")
		}
		return // got the error, done
	}
	t.Fatal("iterator should have yielded an error")
}

func TestBooksService_ListAll_Empty(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data":  []any{},
			"total": 0,
		})
	})

	count := 0
	for _, err := range c.Books.ListAll(context.Background()) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		count++
	}
	if count != 0 {
		t.Errorf("got %d items, want 0", count)
	}
}
