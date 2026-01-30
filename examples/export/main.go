// Export example: export a page to markdown or PDF.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	bookstack "code.beautifulmachines.dev/jakoubek/bookstack-api"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <page-id> <markdown|pdf>\n", os.Args[0])
		os.Exit(1)
	}

	pageID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid page ID: %s", os.Args[1])
	}
	format := os.Args[2]

	client, err := bookstack.NewClient(bookstack.Config{
		BaseURL:     os.Getenv("BOOKSTACK_URL"),
		TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
		TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	switch format {
	case "markdown":
		data, err := client.Pages.ExportMarkdown(ctx, pageID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(data))

	case "pdf":
		data, err := client.Pages.ExportPDF(ctx, pageID)
		if err != nil {
			log.Fatal(err)
		}
		filename := fmt.Sprintf("page-%d.pdf", pageID)
		if err := os.WriteFile(filename, data, 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Saved to %s (%d bytes)\n", filename, len(data))

	default:
		log.Fatalf("Unknown format: %s (use markdown or pdf)", format)
	}
}
