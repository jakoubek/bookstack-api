# bookstack-api

[![Mirror on GitHub](https://img.shields.io/badge/mirror-GitHub-blue)](https://github.com/jakoubek/bookstack-api)
[![Go Reference](https://pkg.go.dev/badge/code.beautifulmachines.dev/jakoubek/bookstack-api.svg)](https://pkg.go.dev/code.beautifulmachines.dev/jakoubek/bookstack-api)
[![Go Report Card](https://goreportcard.com/badge/code.beautifulmachines.dev/jakoubek/bookstack-api)](https://goreportcard.com/report/code.beautifulmachines.dev/jakoubek/bookstack-api)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Go client library for the [BookStack](https://www.bookstackapp.com/) REST API. Zero external dependencies.

## Agent Quick Reference

If you are an AI agent using this library, set these env vars first:

```bash
export BOOKSTACK_API_URL="https://docs.example.com"
export BOOKSTACK_TOKEN_ID="your-token-id"
export BOOKSTACK_TOKEN_SECRET="your-token-secret"
```

Then the client pattern is:

```go
client, err := bookstack.NewClient(bookstack.Config{
    BaseURL:     os.Getenv("BOOKSTACK_API_URL"),
    TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
    TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
})
if err != nil { log.Fatal(err) }
```

Pagination uses Go 1.23+ iterators — iterate with `for x, err := range svc.ListAll(ctx)`.

## Installation

```bash
go get code.beautifulmachines.dev/jakoubek/bookstack-api
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    bookstack "code.beautifulmachines.dev/jakoubek/bookstack-api"
)

func main() {
    client, err := bookstack.NewClient(bookstack.Config{
        BaseURL:     os.Getenv("BOOKSTACK_API_URL"),
        TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
        TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List all books
    books, err := client.Books.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    for _, book := range books {
        fmt.Printf("%d: %s\n", book.ID, book.Name)
    }
}
```

## Authentication

BookStack uses token-based authentication. Create an API token in your BookStack user profile under **API Tokens**.

Set the instance URL and token as environment variables (these names are used by the examples and the live tests):

```bash
export BOOKSTACK_API_URL="https://docs.example.com"
export BOOKSTACK_TOKEN_ID="your-token-id"
export BOOKSTACK_TOKEN_SECRET="your-token-secret"
```

## Usage

### Search

```go
results, err := client.Search.Search(ctx, "deployment guide", nil)
for _, r := range results {
    fmt.Printf("[%s] %s (score: %.1f)\n", r.Type, r.Name, r.Score)
}
```

### Get a Page

```go
page, err := client.Pages.Get(ctx, 42)
fmt.Println(page.HTML)
```

### Export Page as Markdown

```go
md, err := client.Pages.ExportMarkdown(ctx, 42)
fmt.Println(string(md))
```

### Iterate All Books

Uses Go 1.23+ iterators for memory-efficient pagination:

```go
for book, err := range client.Books.ListAll(ctx) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(book.Name)
}
```

### Pagination and Filtering

```go
pages, err := client.Pages.List(ctx, &bookstack.ListOptions{
    Count:  10,
    Offset: 0,
    Sort:   "-updated_at",
    Filter: map[string]string{"book_id": "1"},
})
```

### Create and Update Pages

```go
page, err := client.Pages.Create(ctx, &bookstack.PageCreateRequest{
    BookID:   1,
    Name:     "New Page",
    Markdown: "# Hello\n\nPage content here.",
})

page, err = client.Pages.Update(ctx, page.ID, &bookstack.PageUpdateRequest{
    Markdown: "# Updated\n\nNew content.",
})
```

### Error Handling

```go
page, err := client.Pages.Get(ctx, 999)
if errors.Is(err, bookstack.ErrNotFound) {
    fmt.Println("Page not found")
} else if errors.Is(err, bookstack.ErrUnauthorized) {
    fmt.Println("Invalid credentials")
}
```

### Available Services

| Service | Operations |
|---------|-----------|
| `Books` | List, ListAll, Get |
| `Pages` | List, ListAll, Get, Create, Update, Delete, ExportMarkdown, ExportPDF |
| `Chapters` | List, ListAll, Get |
| `Shelves` | List, ListAll, Get |
| `Search` | Search |
| `Attachments` | List, Get, Create, Update, Delete |
| `Comments` | List, Get, Create, Update, Delete |

## CLI Tool (hqcli)

A standalone CLI binary for quick BookStack operations — useful for agents and automation.

### Installation

```bash
cd cli
go build -o hqcli .
# Or install globally
go install code.beautifulmachines.dev/jakoubek/bookstack-api-cli@latest
```

### Usage

```bash
# Set credentials (or use env vars: BOOKSTACK_API_URL, BOOKSTACK_TOKEN_ID, BOOKSTACK_TOKEN_SECRET)
export BOOKSTACK_API_URL="https://docs.example.com"
export BOOKSTACK_TOKEN_ID="your-token-id"
export BOOKSTACK_TOKEN_SECRET="your-token-secret"

# List books (JSON output)
hqcli books list --json

# Get a page
hqcli pages get 42 --json

# Search
hqcli search "deployment guide" --json

# Create a page
hqcli pages create --book-id=1 --name="New Page" --markdown="# Hello\n\nContent"

# Update a page
hqcli pages update 42 --name="Updated" --markdown="# New content"

# Delete a page
hqcli pages delete 42

# Upload a file attachment
hqcli attachments create --name="Report" --file=report.pdf --uploaded-to=42 --json
```

### Common Flags

| Flag | Description |
|------|-------------|
| `--json` | Machine-parseable JSON output |
| `--all` | Iterate all pages (pagination) |
| `--count N` | Limit results (default: 50) |
| `--offset N` | Pagination offset |
| `--sort FIELD` | Sort field (prefix with `-` for desc) |
| `--filter KEY=VALUE` | Filter (repeatable) |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Generic error / usage |
| 2 | Not found (404) |
| 3 | Auth error (401/403) |
| 4 | Rate limited (429) |

## File Attachments

Attachments can be created as **links** or **file uploads**:

```go
// Link attachment (JSON)
linkAtt, err := client.Attachments.Create(ctx, &bookstack.AttachmentCreateRequest{
    Name:       "Documentation",
    UploadedTo: pageID,
    Link:       "https://example.com/doc.pdf",
})

// File upload (multipart/form-data)
fileAtt, err := client.Attachments.Create(ctx, &bookstack.AttachmentCreateRequest{
    Name:       "Report",
    UploadedTo: pageID,
    File:       bytes.NewReader(fileData),
    Filename:   "report.pdf",
})
```

When `File` and `Filename` are provided, the library automatically uses multipart form-data encoding.

## Requirements

- Go 1.23+
- BookStack instance with API enabled

## Issue Tracking

This repository uses [Beads](https://github.com/gastownhall/beads) (bd) for issue tracking — issues live in the repo at `.beads/` alongside the code.

Install:

```bash
# macOS / Linux (Homebrew)
brew install beads

# Or download pre-built binary (fastest, no Go toolchain needed)
curl -L -o bd.tar.gz https://github.com/gastownhall/beads/releases/download/v1.3.1/beads_1.3.1_$(uname -s)_$(uname -m).tar.gz
tar -xzf bd.tar.gz bd
sudo mv bd /usr/local/bin/
rm bd.tar.gz

# Or build from source (requires Go)
go install github.com/gastownhall/beads/cmd/bd@latest
```

Then run `bd onboard` to get started.

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git remote
```

## Compatibility

Tested against BookStack **26.05.5**.

The library's `UserRef` type handles the dual-shape response that BookStack
returns for `created_by`, `updated_by`, and `owned_by` fields:

- **List endpoints** (`GET /api/books`, `GET /api/pages`, …): plain integer, e.g. `1`
- **Detail endpoints** (`GET /api/books/{id}`, `POST`, `PUT`): expanded object
  `{"id": 1, "name": "Aldo", "slug": "aldo"}`

Also handles nullable `chapter_id` on pages that are directly in a book (not in a chapter).

If you are using a different BookStack version, the API response shape may differ.

## License

See [LICENSE](LICENSE) file.
