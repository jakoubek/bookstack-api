package bookstack

// Optional live-site verification tests.
//
// These tests exercise the library against a real BookStack instance.
// They are skipped unless the following environment variables are set:
//
//	BOOKSTACK_API_URL     e.g. https://docs.example.org
//	BOOKSTACK_TOKEN_ID    API token id (created in BookStack admin)
//	BOOKSTACK_TOKEN_SECRET  API token secret
//
// Every write operation registers a t.Cleanup that deletes the artifact
// it created, so the test leaves no trace on the instance.

import (
	"context"
	"os"
	"strings"
	"testing"
)

func liveClient(t *testing.T) *Client {
	t.Helper()
	baseURL := os.Getenv("BOOKSTACK_API_URL")
	tokenID := os.Getenv("BOOKSTACK_TOKEN_ID")
	tokenSecret := os.Getenv("BOOKSTACK_TOKEN_SECRET")
	if baseURL == "" || tokenID == "" || tokenSecret == "" {
		t.Skip("live tests skipped: set BOOKSTACK_API_URL, BOOKSTACK_TOKEN_ID, BOOKSTACK_TOKEN_SECRET")
	}
	client, err := NewClient(Config{
		BaseURL:     baseURL,
		TokenID:     tokenID,
		TokenSecret: tokenSecret,
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestVerifyLiveShelves(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	shelves, err := client.Shelves.List(ctx, nil)
	if err != nil {
		t.Fatalf("Shelves.List: %v", err)
	}
	if len(shelves) == 0 {
		t.Fatal("expected at least one shelf")
	}

	count := 0
	for range client.Shelves.ListAll(ctx) {
		count++
	}
	if count == 0 {
		t.Fatal("ListAll returned zero shelves")
	}

	shelf, err := client.Shelves.Get(ctx, shelves[0].ID)
	if err != nil {
		t.Fatalf("Shelves.Get: %v", err)
	}
	if shelf.ID != shelves[0].ID {
		t.Errorf("Shelves.Get ID = %d, want %d", shelf.ID, shelves[0].ID)
	}
}

func TestVerifyLiveBooks(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	books, err := client.Books.List(ctx, nil)
	if err != nil {
		t.Fatalf("Books.List: %v", err)
	}
	if len(books) == 0 {
		t.Fatal("expected at least one book")
	}

	count := 0
	for range client.Books.ListAll(ctx) {
		count++
	}
	if count == 0 {
		t.Fatal("Books.ListAll returned zero books")
	}

	book, err := client.Books.Get(ctx, books[0].ID)
	if err != nil {
		t.Fatalf("Books.Get: %v", err)
	}
	if book.ID != books[0].ID {
		t.Errorf("Books.Get ID = %d, want %d", book.ID, books[0].ID)
	}
	if book.CreatedBy.ID == 0 {
		t.Error("Books.Get: CreatedBy.ID is 0, expected user reference")
	}
}

func TestVerifyLiveChapters(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	chapters, err := client.Chapters.List(ctx, nil)
	if err != nil {
		t.Fatalf("Chapters.List: %v", err)
	}
	if len(chapters) == 0 {
		t.Fatal("expected at least one chapter")
	}

	ch, err := client.Chapters.Get(ctx, chapters[0].ID)
	if err != nil {
		t.Fatalf("Chapters.Get: %v", err)
	}
	if ch.ID != chapters[0].ID {
		t.Errorf("Chapters.Get ID = %d, want %d", ch.ID, chapters[0].ID)
	}
}

func TestVerifyLivePagesRead(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	pages, err := client.Pages.List(ctx, nil)
	if err != nil {
		t.Fatalf("Pages.List: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("expected at least one page")
	}

	page, err := client.Pages.Get(ctx, pages[0].ID)
	if err != nil {
		t.Fatalf("Pages.Get: %v", err)
	}
	if page.ID != pages[0].ID {
		t.Errorf("Pages.Get ID = %d, want %d", page.ID, pages[0].ID)
	}

	md, err := client.Pages.ExportMarkdown(ctx, page.ID)
	if err != nil {
		t.Fatalf("Pages.ExportMarkdown: %v", err)
	}
	if len(md) == 0 {
		t.Error("ExportMarkdown returned empty content")
	}
}

func TestVerifyLivePageWriteCycle(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	books, err := client.Books.List(ctx, nil)
	if err != nil {
		t.Fatalf("Books.List: %v", err)
	}
	if len(books) == 0 {
		t.Fatal("expected at least one book for page creation")
	}

	const name = "Live Test Page (auto-cleaned)"
	page, err := client.Pages.Create(ctx, &PageCreateRequest{
		BookID: books[0].ID,
		Name:   name,
		HTML:   "<p>Live test body</p>",
	})
	if err != nil {
		t.Fatalf("Pages.Create: %v", err)
	}
	if page.ID == 0 {
		t.Fatal("Pages.Create returned zero ID")
	}
	t.Cleanup(func() {
		if err := client.Pages.Delete(context.Background(), page.ID); err != nil {
			t.Logf("cleanup: Pages.Delete(%d): %v", page.ID, err)
		}
	})

	// Verify detail unmarshal works after Create (response is detail shape)
	if page.CreatedBy.ID == 0 {
		t.Error("Create response: CreatedBy.ID is 0, expected user reference")
	}

	// Update
	updated, err := client.Pages.Update(ctx, page.ID, &PageUpdateRequest{
		HTML: "<p>Updated body</p>",
	})
	if err != nil {
		t.Fatalf("Pages.Update: %v", err)
	}
	if !strings.Contains(updated.HTML, "Updated body") {
		t.Errorf("Update: HTML = %q, want it to contain 'Updated body'", updated.HTML)
	}

	// Get after write
	got, err := client.Pages.Get(ctx, page.ID)
	if err != nil {
		t.Fatalf("Pages.Get after write: %v", err)
	}
	if got.Name != name {
		t.Errorf("Get after write: Name = %q, want %q", got.Name, name)
	}
}

func TestVerifyLiveSearch(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	results, err := client.Search.Search(ctx, "a", nil)
	if err != nil {
		t.Fatalf("Search.Search: %v", err)
	}
	if len(results) == 0 {
		t.Log("Search returned zero results (allowed on small instances)")
	}
}

func TestVerifyLiveComments(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	pages, err := client.Pages.List(ctx, nil)
	if err != nil {
		t.Fatalf("Pages.List: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("expected at least one page for comments")
	}

	comment, err := client.Comments.Create(ctx, &CommentCreateRequest{
		PageID:          pages[0].ID,
		CommentableID:   pages[0].ID,
		CommentableType: "page",
		HTML:            "<p>Live test comment (auto-cleaned)</p>",
	})
	if err != nil {
		t.Fatalf("Comments.Create: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Comments.Delete(context.Background(), comment.ID); err != nil {
			t.Logf("cleanup: Comments.Delete(%d): %v", comment.ID, err)
		}
	})

	got, err := client.Comments.Get(ctx, comment.ID)
	if err != nil {
		t.Fatalf("Comments.Get: %v", err)
	}
	if got.CommentableID != pages[0].ID {
		t.Errorf("Comments.Get CommentableID = %d, want %d", got.CommentableID, pages[0].ID)
	}

	updated, err := client.Comments.Update(ctx, comment.ID, &CommentUpdateRequest{
		HTML: "<p>Updated comment</p>",
	})
	if err != nil {
		t.Fatalf("Comments.Update: %v", err)
	}
	// BookStack returns the comment without the html field on update;
	// verify via a subsequent GET.
	gotAfter, err := client.Comments.Get(ctx, comment.ID)
	if err != nil {
		t.Fatalf("Comments.Get after update: %v", err)
	}
	if !strings.Contains(gotAfter.HTML, "Updated comment") {
		t.Errorf("Comments.Get after update HTML = %q, want it to contain 'Updated comment'", gotAfter.HTML)
	}
	if updated.ID != comment.ID {
		t.Errorf("Comments.Update ID = %d, want %d", updated.ID, comment.ID)
	}
}

func TestVerifyLiveAttachments(t *testing.T) {
	client := liveClient(t)
	ctx := context.Background()

	pages, err := client.Pages.List(ctx, nil)
	if err != nil {
		t.Fatalf("Pages.List: %v", err)
	}
	if len(pages) == 0 {
		t.Fatal("expected at least one page for attachments")
	}

	// Create a link attachment
	att, err := client.Attachments.Create(ctx, &AttachmentCreateRequest{
		Name:       "Live test link (auto-cleaned)",
		UploadedTo: pages[0].ID,
		Link:       "https://example.com",
	})
	if err != nil {
		t.Fatalf("Attachments.Create: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Attachments.Delete(context.Background(), att.ID); err != nil {
			t.Logf("cleanup: Attachments.Delete(%d): %v", att.ID, err)
		}
	})

	if att.ID == 0 {
		t.Fatal("Attachments.Create returned zero ID")
	}

	// List and verify
	list, err := client.Attachments.List(ctx, nil)
	if err != nil {
		t.Fatalf("Attachments.List: %v", err)
	}
	found := false
	for _, a := range list {
		if a.ID == att.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Attachments.List: created attachment %d not in list", att.ID)
	}

	// Update
	updated, err := client.Attachments.Update(ctx, att.ID, &AttachmentUpdateRequest{
		Name: "Updated link name",
	})
	if err != nil {
		t.Fatalf("Attachments.Update: %v", err)
	}
	if updated.Name != "Updated link name" {
		t.Errorf("Attachments.Update Name = %q", updated.Name)
	}

	// Get
	got, err := client.Attachments.Get(ctx, att.ID)
	if err != nil {
		t.Fatalf("Attachments.Get: %v", err)
	}
	if got.Name != "Updated link name" {
		t.Errorf("Attachments.Get Name = %q", got.Name)
	}

	// Delete semantics: re-fetch after delete should fail
	if err := client.Attachments.Delete(ctx, att.ID); err != nil {
		t.Fatalf("Attachments.Delete (explicit): %v", err)
	}
	if _, err := client.Attachments.Get(ctx, att.ID); err == nil {
		t.Error("Attachments.Get after delete should fail")
	}
}
