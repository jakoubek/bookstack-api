package bookstack

import "time"

// Book represents a Bookstack book.
type Book struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
	OwnedBy     int       `json:"owned_by"`
}

// Page represents a Bookstack page.
type Page struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	ChapterID int       `json:"chapter_id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	HTML      string    `json:"html"`
	Markdown  string    `json:"markdown"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	Draft     bool      `json:"draft"`
	Revision  int       `json:"revision_count"`
	Template  bool      `json:"template"`
	OwnedBy   int       `json:"owned_by"`
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
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
	OwnedBy     int       `json:"owned_by"`
}

// Shelf represents a Bookstack shelf.
type Shelf struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"`
	UpdatedBy   int       `json:"updated_by"`
	OwnedBy     int       `json:"owned_by"`
}

// PageCreateRequest contains fields for creating a new page.
type PageCreateRequest struct {
	BookID    int    `json:"book_id"`
	ChapterID int    `json:"chapter_id,omitempty"`
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
	CreatedBy  int       `json:"created_by"`
	UpdatedBy  int       `json:"updated_by"`
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
