package bookstack

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBook_JSONUnmarshal(t *testing.T) {
	data := `{
		"id": 1,
		"name": "Test Book",
		"slug": "test-book",
		"description": "A test book",
		"created_at": "2024-01-15T10:30:00.000000Z",
		"updated_at": "2024-01-16T12:00:00.000000Z",
		"created_by": 1,
		"updated_by": 2,
		"owned_by": 1
	}`

	var b Book
	if err := json.Unmarshal([]byte(data), &b); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if b.ID != 1 {
		t.Errorf("ID = %d, want 1", b.ID)
	}
	if b.Name != "Test Book" {
		t.Errorf("Name = %q, want %q", b.Name, "Test Book")
	}
	if b.CreatedAt.Year() != 2024 || b.CreatedAt.Month() != time.January || b.CreatedAt.Day() != 15 {
		t.Errorf("CreatedAt = %v, want 2024-01-15", b.CreatedAt)
	}
}

func TestPage_JSONUnmarshal(t *testing.T) {
	data := `{
		"id": 5,
		"book_id": 1,
		"chapter_id": 2,
		"name": "Test Page",
		"slug": "test-page",
		"html": "<p>Hello</p>",
		"markdown": "Hello",
		"draft": false,
		"created_at": "2024-06-01T08:00:00.000000Z",
		"updated_at": "2024-06-02T09:00:00.000000Z"
	}`

	var p Page
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.BookID != 1 {
		t.Errorf("BookID = %d, want 1", p.BookID)
	}
	if p.ChapterID != 2 {
		t.Errorf("ChapterID = %d, want 2", p.ChapterID)
	}
	if p.HTML != "<p>Hello</p>" {
		t.Errorf("HTML = %q", p.HTML)
	}
}

func TestSearchResult_JSONUnmarshal(t *testing.T) {
	data := `{
		"type": "page",
		"id": 10,
		"name": "Found Page",
		"slug": "found-page",
		"book_id": 3,
		"preview": "...match...",
		"score": 1.5
	}`

	var sr SearchResult
	if err := json.Unmarshal([]byte(data), &sr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if sr.Type != "page" {
		t.Errorf("Type = %q, want %q", sr.Type, "page")
	}
	if sr.Score != 1.5 {
		t.Errorf("Score = %f, want 1.5", sr.Score)
	}
}
