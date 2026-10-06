# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
