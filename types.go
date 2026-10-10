package bookstack

import (
	"encoding/json"
	"fmt"
	"time"
)

// UserRef represents a BookStack user reference.
//
// BookStack returns this field in two different shapes depending on the
// endpoint:
//   - List endpoints (GET /api/books, GET /api/pages, …): a plain integer, e.g. 1
//   - Detail endpoints (GET /api/books/{id}, POST /api/books, PUT /api/books/{id}, …):
//     an expanded object: {"id": 1, "name": "Aldo", "slug": "aldo"}
//
// UnmarshalJSON handles both forms so a single struct works everywhere.
type UserRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// UnmarshalJSON accepts either a plain integer or an object.
func (u *UserRef) UnmarshalJSON(data []byte) error {
	// Try the object form first: {"id": 1, "name": "...", "slug": "..."}
	var obj struct {
		ID   *int   `json:"id"`
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal(data, &obj); err == nil && obj.ID != nil {
		u.ID = *obj.ID
		u.Name = obj.Name
		u.Slug = obj.Slug
		return nil
	}

	// Fall back to a plain integer: 1
	var id int
	if err := json.Unmarshal(data, &id); err == nil {
		u.ID = id
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into UserRef", string(data))
}

// MarshalJSON emits the integer form for compatibility with list responses.
func (u UserRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.ID)
}

// Book represents a Bookstack book.
type Book struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   UserRef   `json:"created_by"`
	UpdatedBy   UserRef   `json:"updated_by"`
	OwnedBy     UserRef   `json:"owned_by"`
}

// Page represents a Bookstack page.
type Page struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	ChapterID *int      `json:"chapter_id"` // null when page is directly in a book
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	HTML      string    `json:"html"`
	Markdown  string    `json:"markdown"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy UserRef   `json:"created_by"`
	UpdatedBy UserRef   `json:"updated_by"`
	Draft     bool      `json:"draft"`
	Revision  int       `json:"revision_count"`
	Template  bool      `json:"template"`
	OwnedBy   UserRef   `json:"owned_by"`
}

// Chapter represents a Bookstack chapter.
type Chapter struct {
	ID          int       `json:"id"`
	BookID      int       `json:"book_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   UserRef   `json:"created_by"`
	UpdatedBy   UserRef   `json:"updated_by"`
	OwnedBy     UserRef   `json:"owned_by"`
}

// Shelf represents a Bookstack shelf.
type Shelf struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   UserRef   `json:"created_by"`
	UpdatedBy   UserRef   `json:"updated_by"`
	OwnedBy     UserRef   `json:"owned_by"`
}

// PageCreateRequest contains fields for creating a new page.
type PageCreateRequest struct {
	BookID    int    `json:"book_id"`
	ChapterID *int   `json:"chapter_id,omitempty"`
	Name      string `json:"name"`
	HTML      string `json:"html,omitempty"`
	Markdown  string `json:"markdown,omitempty"`
}

// PageUpdateRequest contains fields for updating an existing page.
type PageUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	HTML     string `json:"html,omitempty"`
	Markdown string `json:"markdown,omitempty"`
}

// Attachment represents a Bookstack attachment.
type Attachment struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Extension  string    `json:"extension"`
	UploadedTo int       `json:"uploaded_to"`
	External   bool      `json:"external"`
	Order      int       `json:"order"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  UserRef   `json:"created_by"`
	UpdatedBy  UserRef   `json:"updated_by"`
	Content    string    `json:"content,omitempty"`
}

// AttachmentCreateRequest contains fields for creating an attachment.
type AttachmentCreateRequest struct {
	Name       string `json:"name"`
	UploadedTo int    `json:"uploaded_to"`
	Link       string `json:"link,omitempty"`
}

// AttachmentUpdateRequest contains fields for updating an attachment.
type AttachmentUpdateRequest struct {
	Name string `json:"name,omitempty"`
	Link string `json:"link,omitempty"`
}

// Comment represents a Bookstack comment on a page.
type Comment struct {
	ID              int       `json:"id"`
	CommentableID   int       `json:"commentable_id"`
	CommentableType string    `json:"commentable_type"`
	ParentID        int       `json:"parent_id,omitempty"`
	HTML            string    `json:"html"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedBy       UserRef   `json:"created_by"`
	UpdatedBy       UserRef   `json:"updated_by"`
}

// CommentCreateRequest contains fields for creating a comment.
// BookStack requires both page_id (for validation) and commentable_id/commentable_type.
type CommentCreateRequest struct {
	PageID          int    `json:"page_id"`
	CommentableID   int    `json:"commentable_id"`
	CommentableType string `json:"commentable_type"`
	ParentID        int    `json:"parent_id,omitempty"`
	HTML            string `json:"html"`
}

// CommentUpdateRequest contains fields for updating a comment.
type CommentUpdateRequest struct {
	HTML string `json:"html"`
}

// SearchResult represents a search result from Bookstack.
type SearchResult struct {
	Type      string  `json:"type"` // "page", "chapter", "book", or "shelf"
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	BookID    int     `json:"book_id"`    // For pages and chapters
	ChapterID int     `json:"chapter_id"` // For pages
	Preview   string  `json:"preview"`
	Score     float64 `json:"score"`
}
