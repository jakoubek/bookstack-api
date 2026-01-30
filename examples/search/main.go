// Search example: search for content and display results.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	bookstack "code.beautifulmachines.dev/jakoubek/bookstack-api"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <search query>\n", os.Args[0])
		os.Exit(1)
	}

	client, err := bookstack.NewClient(bookstack.Config{
		BaseURL:     os.Getenv("BOOKSTACK_URL"),
		TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
		TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	query := os.Args[1]

	results, err := client.Search.Search(ctx, query, &bookstack.ListOptions{Count: 10})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Search results for %q:\n\n", query)
	for i, r := range results {
		fmt.Printf("%d. [%s] %s (score: %.1f)\n", i+1, r.Type, r.Name, r.Score)
		if r.Preview != "" {
			fmt.Printf("   %s\n", r.Preview)
		}
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
	}
}
