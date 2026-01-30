# bookstack-api

A Go client library for the [BookStack](https://www.bookstackapp.com/) REST API. Zero external dependencies.

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
        BaseURL:     "https://docs.example.com",
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

Set the token ID and secret as environment variables:

```bash
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

## Available Services

| Service | Operations |
|---------|-----------|
| `Books` | List, ListAll, Get |
| `Pages` | List, ListAll, Get, Create, Update, Delete, ExportMarkdown, ExportPDF |
| `Chapters` | List, ListAll, Get |
| `Shelves` | List, ListAll, Get |
| `Search` | Search |
| `Attachments` | List, Get, Create, Update, Delete |
| `Comments` | List, Get, Create, Update, Delete |

## Requirements

- Go 1.23+
- BookStack instance with API enabled

## License

See [LICENSE](LICENSE) file.
