# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A Go library for the Bookstack REST API. Primary consumers are AI agents via CLI (with `--json` flag) and Go developers needing programmatic Bookstack access.

## Build & Test Commands

```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run single test
go test -run TestName ./...

# Run tests with coverage
go test -cover ./...

# Lint (if golangci-lint is installed)
golangci-lint run

# Format code
gofmt -w .
```

## Architecture

### Flat Package Structure

All code lives in the root package `bookstack`. No subpackages.

```
bookstack.go       # Client, Config, NewClient()
books.go           # BooksService
pages.go           # PagesService
chapters.go        # ChaptersService
shelves.go         # ShelvesService
search.go          # SearchService
types.go           # All data structures (Book, Page, Chapter, Shelf, SearchResult)
errors.go          # Error types (ErrNotFound, ErrUnauthorized, APIError, etc.)
http.go            # HTTP helpers, request building
iterator.go        # Pagination iterator using Go 1.23+ iter.Seq
```

### Client Pattern

Services are attached to the main `Client` struct:

```go
client := bookstack.NewClient(bookstack.Config{
    BaseURL:     "https://docs.example.com",
    TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
    TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
})

// Access services via client
books, err := client.Books.List(ctx, nil)
page, err := client.Pages.Get(ctx, 123)
```

### Pagination via Iterators

Use Go 1.23+ iterator pattern for list operations:

```go
for book, err := range client.Books.ListAll(ctx) {
    if err != nil {
        return err
    }
    // process book
}
```

### Error Handling

Use sentinel errors with `errors.Is()`:

```go
if errors.Is(err, bookstack.ErrNotFound) {
    // handle 404
}
```

`APIError` provides detailed information when needed.

## Key Design Decisions

- **Zero external dependencies** - Only Go standard library
- **Go 1.21+** required
- **No caching/rate-limiting** - Caller responsibility
- **Token auth** via `Authorization: Token <id>:<secret>` header
- **Bookstack hierarchy**: Shelf → Book → Chapter → Page

## Bookstack API Reference

- Rate limit: 180 req/min (server configurable)
- Pagination: `count` (max 500), `offset`, `sort`, `filter[field]`
