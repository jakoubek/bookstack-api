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

func TestPagesService_ExportMarkdown(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pages/5/export/markdown" {
			t.Errorf("path = %s, want /api/pages/5/export/markdown", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("# Hello\n\nThis is markdown."))
	})

	data, err := c.Pages.ExportMarkdown(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "# Hello\n\nThis is markdown." {
		t.Errorf("got %q", string(data))
	}
}

func TestPagesService_ExportPDF(t *testing.T) {
	pdfContent := []byte("%PDF-1.4 fake content")
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/pages/5/export/pdf" {
			t.Errorf("path = %s, want /api/pages/5/export/pdf", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write(pdfContent)
	})

	data, err := c.Pages.ExportPDF(context.Background(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != string(pdfContent) {
		t.Errorf("got %d bytes, want %d", len(data), len(pdfContent))
	}
}

func TestPagesService_ExportMarkdown_NotFound(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := c.Pages.ExportMarkdown(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Error("expected ErrNotFound")
	}
}
