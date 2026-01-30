// Basic example: set up client, list books, get a page.
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
		BaseURL:     os.Getenv("BOOKSTACK_URL"),
		TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
		TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List first 5 books
	books, err := client.Books.List(ctx, &bookstack.ListOptions{Count: 5})
	if err != nil {
		log.Fatal(err)
	}

	for _, book := range books {
		fmt.Printf("Book %d: %s\n", book.ID, book.Name)
	}

	// Get the first page if any books exist
	if len(books) > 0 {
		pages, err := client.Pages.List(ctx, &bookstack.ListOptions{
			Count:  1,
			Filter: map[string]string{"book_id": fmt.Sprintf("%d", books[0].ID)},
		})
		if err != nil {
			log.Fatal(err)
		}
		if len(pages) > 0 {
			page, err := client.Pages.Get(ctx, pages[0].ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("\nPage: %s\n%s\n", page.Name, page.HTML[:min(200, len(page.HTML))])
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
