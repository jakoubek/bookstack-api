# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- `UserRef` type now correctly unmarshals both integer and object forms of
  `created_by` / `updated_by` / `owned_by` fields returned by BookStack 26.05.5
  (detail endpoints return objects, list endpoints return integers).
- `Page.ChapterID` is now `*int` to handle `null` for pages directly in a book.
- `CommentCreateRequest` now includes both `page_id` and `commentable_id` fields
  required by BookStack validation.

### Added
- Optional live-site verification test file (`verify_live_test.go`) gated by
  `BOOKSTACK_API_URL`, `BOOKSTACK_TOKEN_ID`, `BOOKSTACK_TOKEN_SECRET` env vars.
  Tests exercise all services against a real BookStack instance with automatic
  cleanup of created resources via `t.Cleanup`.
- Husky pre-commit hook running `go vet`, `go build`, and `go test -short`.
- Compatibility section in README documenting BookStack 26.05.5 response shapes.

## [v0.1.0] - 2026-10-06

### Added
- Client with config validation and token authentication (`NewClient`)
- HTTP helpers and request building
- `BooksService`, `ChaptersService` and `ShelvesService` (List, Get, ListAll)
- `PagesService` (List, Get, Create, Update, Delete) and page export (Markdown, PDF)
- `AttachmentsService` and `CommentsService` (CRUD)
- `SearchService`
- Pagination iterator based on `iter.Seq2`
- Typed errors with `APIError` and sentinel errors
- Unit tests for error types and data types
- README with quick-start guide, GoDoc documentation and examples in `examples/`
