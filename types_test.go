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
	if p.ChapterID == nil || *p.ChapterID != 2 {
		t.Errorf("ChapterID = %v, want 2", p.ChapterID)
	}
	if p.HTML != "<p>Hello</p>" {
		t.Errorf("HTML = %q", p.HTML)
	}
}

func TestUserRef_UnmarshalFromInt(t *testing.T) {
	var u UserRef
	if err := json.Unmarshal([]byte("42"), &u); err != nil {
		t.Fatalf("unmarshal int: %v", err)
	}
	if u.ID != 42 {
		t.Errorf("ID = %d, want 42", u.ID)
	}
	if u.Name != "" {
		t.Errorf("Name = %q, want empty", u.Name)
	}
}

func TestUserRef_UnmarshalFromObject(t *testing.T) {
	var u UserRef
	data := `{"id": 7, "name": "Test User", "slug": "test-user"}`
	if err := json.Unmarshal([]byte(data), &u); err != nil {
		t.Fatalf("unmarshal object: %v", err)
	}
	if u.ID != 7 {
		t.Errorf("ID = %d, want 7", u.ID)
	}
	if u.Name != "Test User" {
		t.Errorf("Name = %q, want %q", u.Name, "Test User")
	}
	if u.Slug != "test-user" {
		t.Errorf("Slug = %q, want %q", u.Slug, "test-user")
	}
}

func TestUserRef_MarshalToInt(t *testing.T) {
	u := UserRef{ID: 5, Name: "Test", Slug: "test"}
	data, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != "5" {
		t.Errorf("marshal = %s, want 5", string(data))
	}
}

func TestPage_ChapterIDNull(t *testing.T) {
	data := `{
		"id": 5,
		"book_id": 1,
		"chapter_id": null,
		"name": "Test Page",
		"slug": "test-page",
		"html": "<p>Hello</p>",
		"markdown": "Hello",
		"created_at": "2024-06-01T08:00:00.000000Z",
		"updated_at": "2024-06-02T09:00:00.000000Z"
	}`
	var p Page
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if p.ChapterID != nil {
		t.Errorf("ChapterID = %v, want nil", p.ChapterID)
	}
}

func TestBook_DetailUnmarshal(t *testing.T) {
	data := `{
		"id": 1,
		"name": "Test Book",
		"slug": "test-book",
		"created_at": "2024-01-15T10:30:00.000000Z",
		"updated_at": "2024-01-16T12:00:00.000000Z",
		"created_by": {"id": 1, "name": "Admin", "slug": "admin"},
		"updated_by": {"id": 1, "name": "Admin", "slug": "admin"},
		"owned_by": {"id": 1, "name": "Admin", "slug": "admin"}
	}`
	var b Book
	if err := json.Unmarshal([]byte(data), &b); err != nil {
		t.Fatalf("unmarshal detail: %v", err)
	}
	if b.CreatedBy.ID != 1 {
		t.Errorf("CreatedBy.ID = %d, want 1", b.CreatedBy.ID)
	}
	if b.CreatedBy.Name != "Admin" {
		t.Errorf("CreatedBy.Name = %q, want Admin", b.CreatedBy.Name)
	}
}
