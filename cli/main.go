// Command hqcli is a thin CLI wrapper around the bookstack-api library.
//
// It is designed for AI agents and shell scripts: commands map 1:1 to
// library services, output is JSON by default, and exit codes map to
// error classes so scripts can branch on them:
//
//	0  success
//	1  generic / usage error
//	2  resource not found (404)
//	3  auth error (401 / 403)
//	4  rate limited (429)
//
// Credentials are read from flags or, as a fallback, the environment:
//
//	BOOKSTACK_API_URL      instance base URL
//	BOOKSTACK_TOKEN_ID     API token ID
//	BOOKSTACK_TOKEN_SECRET API token secret
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	bookstack "code.beautifulmachines.dev/jakoubek/bookstack-api"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	// Handle help/version before flag parsing so they work without credentials.
	for _, a := range args {
		if a == "help" || a == "--help" || a == "-h" {
			printUsage(os.Stdout)
			return 0
		}
		if a == "version" || a == "--version" {
			fmt.Println(version)
			return 0
		}
	}

	fs := newFlagSet("hqcli")
	baseURL := fs.String("url", "", "BookStack base URL (default $BOOKSTACK_API_URL)")
	tokenID := fs.String("token-id", "", "API token ID (default $BOOKSTACK_TOKEN_ID)")
	secret := fs.String("token-secret", "", "API token secret (default $BOOKSTACK_TOKEN_SECRET)")
	all := fs.Bool("all", false, "list every page of results (list commands)")
	timeoutS := fs.Int("timeout", 30, "HTTP timeout in seconds")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "hqcli: %v\n", err)
		return 1
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		printUsage(os.Stderr)
		return 1
	}

	switch remaining[0] {
	case "version", "--version":
		fmt.Println(version)
		return 0
	case "help", "--help", "-h":
		printUsage(os.Stdout)
		return 0
	}

	bURL := envOr(*baseURL, "BOOKSTACK_API_URL")
	tID := envOr(*tokenID, "BOOKSTACK_TOKEN_ID")
	tSec := envOr(*secret, "BOOKSTACK_TOKEN_SECRET")
	if bURL == "" || tID == "" || tSec == "" {
		fmt.Fprintln(os.Stderr, "hqcli: missing configuration. Set --url/--token-id/--token-secret or BOOKSTACK_API_URL/BOOKSTACK_TOKEN_ID/BOOKSTACK_TOKEN_SECRET.")
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutS)*time.Second)
	defer cancel()

	client, err := bookstack.NewClient(bookstack.Config{
		BaseURL:     bURL,
		TokenID:     tID,
		TokenSecret: tSec,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "hqcli: %v\n", err)
		return 1
	}

	group, action, rest := remaining[0], "", remaining[1:]
	// search has no sub-action; the query is the first positional arg.
	if group != "search" && len(rest) > 0 {
		action, rest = rest[0], rest[1:]
	}

	err = dispatch(ctx, client, group, action, rest, *all)
	if err == nil {
		return 0
	}
	return exitCodeFor(err)
}

func envOr(v, key string) string {
	if v != "" {
		return v
	}
	return os.Getenv(key)
}