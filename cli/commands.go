// cli/commands.go - command dispatch and implementations
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	bookstack "code.beautifulmachines.dev/jakoubek/bookstack-api"
)

var errUsage = errors.New("usage error")

// dispatch routes "group action" to the matching library service.
func dispatch(ctx context.Context, client *bookstack.Client, group, action string, args []string, all bool) error {
	switch group {
	case "books":
		return booksAction(ctx, client, action, args, all)
	case "chapters":
		return chaptersAction(ctx, client, action, args, all)
	case "pages":
		return pagesAction(ctx, client, action, args, all)
	case "shelves":
		return shelvesAction(ctx, client, action, args, all)
	case "attachments":
		return attachmentsAction(ctx, client, action, args)
	case "comments":
		return commentsAction(ctx, client, action, args)
	case "search":
		return searchAction(ctx, client, args)
	}

	if action != "" {
		fmt.Fprintf(os.Stderr, "hqcli: unknown action %q for %s\n", action, group)
		return errUsage
	}
	fmt.Fprintf(os.Stderr, "hqcli: unknown command %q\n", group)
	return errUsage
}

// --- books -----------------------------------------------------------------

func booksAction(ctx context.Context, c *bookstack.Client, action string, args []string, all bool) error {
	switch action {
	case "list":
		return listCommand(ctx, "books", args, all,
			func(ctx context.Context) (any, error) { return collectAll(ctx, c.Books.ListAll(ctx)) },
			func(ctx context.Context, args []string) (any, error) { return c.Books.List(ctx, parseListOptions(args)) })
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Books.Get(ctx, id))
	default:
		return usage("books", "list", "get <id>")
	}
}

// --- chapters ----------------------------------------------------------------

func chaptersAction(ctx context.Context, c *bookstack.Client, action string, args []string, all bool) error {
	switch action {
	case "list":
		return listCommand(ctx, "chapters", args, all,
			func(ctx context.Context) (any, error) { return collectAll(ctx, c.Chapters.ListAll(ctx)) },
			func(ctx context.Context, args []string) (any, error) { return c.Chapters.List(ctx, parseListOptions(args)) })
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Chapters.Get(ctx, id))
	default:
		return usage("chapters", "list", "get <id>")
	}
}

// --- shelves ------------------------------------------------------------------

func shelvesAction(ctx context.Context, c *bookstack.Client, action string, args []string, all bool) error {
	switch action {
	case "list":
		return listCommand(ctx, "shelves", args, all,
			func(ctx context.Context) (any, error) { return collectAll(ctx, c.Shelves.ListAll(ctx)) },
			func(ctx context.Context, args []string) (any, error) { return c.Shelves.List(ctx, parseListOptions(args)) })
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Shelves.Get(ctx, id))
	default:
		return usage("shelves", "list", "get <id>")
	}
}

// --- pages --------------------------------------------------------------------

func pagesAction(ctx context.Context, c *bookstack.Client, action string, args []string, all bool) error {
	switch action {
	case "list":
		return listCommand(ctx, "pages", args, all,
			func(ctx context.Context) (any, error) { return collectAll(ctx, c.Pages.ListAll(ctx)) },
			func(ctx context.Context, args []string) (any, error) { return c.Pages.List(ctx, parseListOptions(args)) })
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Pages.Get(ctx, id))
	case "create":
		req, err := pageCreateRequest(args)
		if err != nil {
			return err
		}
		return output(c.Pages.Create(ctx, req))
	case "update":
		id, rest, err := splitID(args)
		if err != nil {
			return err
		}
		req, err := pageUpdateRequest(rest)
		if err != nil {
			return err
		}
		return output(c.Pages.Update(ctx, id, req))
	case "delete":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return c.Pages.Delete(ctx, id)
	case "export":
		id, format, err := exportArgs(args)
		if err != nil {
			return err
		}
		var data []byte
		switch format {
		case "markdown", "md":
			data, err = c.Pages.ExportMarkdown(ctx, id)
		case "pdf":
			data, err = c.Pages.ExportPDF(ctx, id)
		default:
			return fmt.Errorf("unknown export format %q (want markdown or pdf)", format)
		}
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(data)
		return err
	default:
		return usage("pages", "list", "get <id>", "create", "update <id>", "delete <id>", "export <id> <markdown|pdf>")
	}
}

// --- attachments ----------------------------------------------------------------

func attachmentsAction(ctx context.Context, c *bookstack.Client, action string, args []string) error {
	switch action {
	case "list":
		return output(c.Attachments.List(ctx, parseListOptions(args)))
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Attachments.Get(ctx, id))
	case "create":
		req, err := attachmentCreateRequest(args)
		if err != nil {
			return err
		}
		return output(c.Attachments.Create(ctx, req))
	case "update":
		id, rest, err := splitID(args)
		if err != nil {
			return err
		}
		req, err := attachmentUpdateRequest(rest)
		if err != nil {
			return err
		}
		return output(c.Attachments.Update(ctx, id, req))
	case "delete":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return c.Attachments.Delete(ctx, id)
	default:
		return usage("attachments", "list", "get <id>", "create", "update <id>", "delete <id>")
	}
}

// --- comments --------------------------------------------------------------------

func commentsAction(ctx context.Context, c *bookstack.Client, action string, args []string) error {
	switch action {
	case "list":
		return output(c.Comments.List(ctx, parseListOptions(args)))
	case "get":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return output(c.Comments.Get(ctx, id))
	case "create":
		req, err := commentCreateRequest(args)
		if err != nil {
			return err
		}
		return output(c.Comments.Create(ctx, req))
	case "update":
		id, rest, err := splitID(args)
		if err != nil {
			return err
		}
		req, err := commentUpdateRequest(rest)
		if err != nil {
			return err
		}
		return output(c.Comments.Update(ctx, id, req))
	case "delete":
		id, err := argID(args)
		if err != nil {
			return err
		}
		return c.Comments.Delete(ctx, id)
	default:
		return usage("comments", "list", "get <id>", "create", "update <id>", "delete <id>")
	}
}

// --- search -------------------------------------------------------------------------

func searchAction(ctx context.Context, c *bookstack.Client, args []string) error {
	if len(args) < 1 {
		return usage("search", "<query>")
	}
	opts := parseListOptions(args[1:])
	results, err := c.Search.Search(ctx, args[0], opts)
	if err != nil {
		return err
	}
	return output(results, nil)
}

// --- helpers ---------------------------------------------------------------------------

// listCommand runs a single-page List by default, or the full ListAll iterator when --all is set.
func listCommand(ctx context.Context, name string, args []string, all bool, listAll func(context.Context) (any, error), listPage func(context.Context, []string) (any, error)) error {
	if all {
		return output(listAll(ctx))
	}
	return output(listPage(ctx, args))
}

// collectAll drains a ListAll iterator into a slice (JSON-friendly).
func collectAll[T any](ctx context.Context, it func(func(T, error) bool)) (any, error) {
	var out []T
	for item, err := range it {
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

// output marshals v to indented JSON on stdout; nil prints nothing.
func output(v any, err error) error {
	if err != nil {
		fmt.Fprintf(os.Stderr, "hqcli: %v\n", err)
		return err
	}
	if v == nil {
		return nil
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// exitCodeFor maps library errors to CLI exit codes.
func exitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	switch {
	case errors.Is(err, bookstack.ErrNotFound):
		return 2
	case errors.Is(err, bookstack.ErrUnauthorized), errors.Is(err, bookstack.ErrForbidden):
		return 3
	case errors.Is(err, bookstack.ErrRateLimited):
		return 4
	default:
		fmt.Fprintf(os.Stderr, "hqcli: %v\n", err)
		return 1
	}
}

func argID(args []string) (int, error) {
	if len(args) < 1 {
		return 0, fmt.Errorf("missing <id> argument")
	}
	return strconv.Atoi(args[0])
}

func splitID(args []string) (int, []string, error) {
	if len(args) < 1 {
		return 0, nil, fmt.Errorf("missing <id> argument")
	}
	id, err := strconv.Atoi(args[0])
	return id, args[1:], err
}

func exportArgs(args []string) (int, string, error) {
	if len(args) < 2 {
		return 0, "", fmt.Errorf("usage: pages export <id> <markdown|pdf>")
	}
	id, err := strconv.Atoi(args[0])
	return id, args[1], err
}

func usage(command string, actions ...string) error {
	var b strings.Builder
	b.WriteString("usage: hqcli ")
	b.WriteString(command)
	for i, a := range actions {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(a)
	}
	fmt.Fprintln(os.Stderr, b.String())
	return errUsage
}

func printUsage(w io.Writer) {
	fmt.Fprint(w, `hqcli — BookStack API CLI for agents

Usage:
  hqcli [global flags] <group> <action> [args]

Global flags:
  --url <url>            BookStack base URL (default $BOOKSTACK_API_URL)
  --token-id <id>        API token ID (default $BOOKSTACK_TOKEN_ID)
  --token-secret <sec>   API token secret (default $BOOKSTACK_TOKEN_SECRET)
  --json                 JSON output (default on; use --json=false for raw)
  --all                  list every page of results (list commands)
  --timeout <seconds>    HTTP timeout (default 30)

Commands:
  books      list [--count N --offset N --sort field] | get <id>
  chapters   list | get <id>
  shelves    list | get <id>
  pages      list | get <id> | create | update <id> | delete <id> | export <id> <markdown|pdf>
  attachments  list | get <id> | create | update <id> | delete <id>
  comments   list | get <id> | create | update <id> | delete <id>
  search     <query> [--count N --offset N]
  version

Write flags (parsed as --key=value, remaining args after <id>):
  pages create   --book-id=N --name=NAME [--html=... | --markdown=...] [--chapter-id=N]
  pages update   <id> --name=NAME [--html=... | --markdown=...]
  attachments create   --name=NAME --uploaded-to=ID [--link=URL | --file=PATH]
  attachments update   <id> --name=NAME [--link=URL]
  comments create   --page-id=ID --commentable-id=ID --commentable-type=Page --html=HTML [--parent-id=N]
  comments update   <id> --html=HTML

Exit codes:
  0 success | 1 error | 2 not found (404) | 3 auth (401/403) | 4 rate limited (429)
`)
}

// --- request builders (parse --key=value args into structs) -----------------------------------

func parseListOptions(args []string) *bookstack.ListOptions {
	opts := &bookstack.ListOptions{Filter: map[string]string{}}
	for _, a := range args {
		if strings.HasPrefix(a, "--count=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--count=")); v > 0 {
				opts.Count = v
			}
		}
		if strings.HasPrefix(a, "--offset=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--offset=")); v > 0 {
				opts.Offset = v
			}
		}
		if strings.HasPrefix(a, "--sort=") {
			opts.Sort = strings.TrimPrefix(a, "--sort=")
		}
		if strings.HasPrefix(a, "--filter=") {
			// Format: --filter=key=value
			val := strings.TrimPrefix(a, "--filter=")
			if idx := strings.Index(val, "="); idx >= 0 {
				key := val[:idx]
				opts.Filter[key] = val[idx+1:]
			}
		}
	}
	return opts
}

func pageCreateRequest(args []string) (*bookstack.PageCreateRequest, error) {
	req := &bookstack.PageCreateRequest{}
	for _, a := range args {
		if strings.HasPrefix(a, "--book-id=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--book-id=")); v > 0 {
				req.BookID = v
			}
		}
		if strings.HasPrefix(a, "--chapter-id=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--chapter-id=")); v > 0 {
				req.ChapterID = &v
			}
		}
		if strings.HasPrefix(a, "--name=") {
			req.Name = strings.TrimPrefix(a, "--name=")
		}
		if strings.HasPrefix(a, "--html=") {
			req.HTML = strings.TrimPrefix(a, "--html=")
		}
		if strings.HasPrefix(a, "--markdown=") {
			req.Markdown = strings.TrimPrefix(a, "--markdown=")
		}
	}
	if req.BookID == 0 {
		return nil, fmt.Errorf("pages create: --book-id is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("pages create: --name is required")
	}
	return req, nil
}

func pageUpdateRequest(args []string) (*bookstack.PageUpdateRequest, error) {
	req := &bookstack.PageUpdateRequest{}
	for _, a := range args {
		if strings.HasPrefix(a, "--name=") {
			req.Name = strings.TrimPrefix(a, "--name=")
		}
		if strings.HasPrefix(a, "--html=") {
			req.HTML = strings.TrimPrefix(a, "--html=")
		}
		if strings.HasPrefix(a, "--markdown=") {
			req.Markdown = strings.TrimPrefix(a, "--markdown=")
		}
	}
	return req, nil
}

func attachmentCreateRequest(args []string) (*bookstack.AttachmentCreateRequest, error) {
	req := &bookstack.AttachmentCreateRequest{}
	var filePath string
	for _, a := range args {
		if strings.HasPrefix(a, "--name=") {
			req.Name = strings.TrimPrefix(a, "--name=")
		}
		if strings.HasPrefix(a, "--uploaded-to=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--uploaded-to=")); v > 0 {
				req.UploadedTo = v
			}
		}
		if strings.HasPrefix(a, "--link=") {
			req.Link = strings.TrimPrefix(a, "--link=")
		}
		if strings.HasPrefix(a, "--file=") {
			filePath = strings.TrimPrefix(a, "--file=")
		}
	}
	if req.Name == "" {
		return nil, fmt.Errorf("attachments create: --name is required")
	}
	if req.UploadedTo == 0 {
		return nil, fmt.Errorf("attachments create: --uploaded-to is required")
	}
	if filePath != "" {
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("open file: %w", err)
		}
		req.File = f
		req.Filename = filePath[strings.LastIndex(filePath, "/")+1:]
	}
	if req.Link == "" && req.File == nil {
		return nil, fmt.Errorf("attachments create: either --link or --file is required")
	}
	return req, nil
}

func attachmentUpdateRequest(args []string) (*bookstack.AttachmentUpdateRequest, error) {
	req := &bookstack.AttachmentUpdateRequest{}
	for _, a := range args {
		if strings.HasPrefix(a, "--name=") {
			req.Name = strings.TrimPrefix(a, "--name=")
		}
		if strings.HasPrefix(a, "--link=") {
			req.Link = strings.TrimPrefix(a, "--link=")
		}
	}
	return req, nil
}

func commentCreateRequest(args []string) (*bookstack.CommentCreateRequest, error) {
	req := &bookstack.CommentCreateRequest{}
	for _, a := range args {
		if strings.HasPrefix(a, "--page-id=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--page-id=")); v > 0 {
				req.PageID = v
			}
		}
		if strings.HasPrefix(a, "--commentable-id=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--commentable-id=")); v > 0 {
				req.CommentableID = v
			}
		}
		if strings.HasPrefix(a, "--commentable-type=") {
			req.CommentableType = strings.TrimPrefix(a, "--commentable-type=")
		}
		if strings.HasPrefix(a, "--parent-id=") {
			if v, _ := strconv.Atoi(strings.TrimPrefix(a, "--parent-id=")); v > 0 {
				req.ParentID = v
			}
		}
		if strings.HasPrefix(a, "--html=") {
			req.HTML = strings.TrimPrefix(a, "--html=")
		}
	}
	if req.PageID == 0 {
		return nil, fmt.Errorf("comments create: --page-id is required")
	}
	if req.CommentableID == 0 {
		return nil, fmt.Errorf("comments create: --commentable-id is required")
	}
	if req.CommentableType == "" {
		return nil, fmt.Errorf("comments create: --commentable-type is required")
	}
	if req.HTML == "" {
		return nil, fmt.Errorf("comments create: --html is required")
	}
	return req, nil
}

func commentUpdateRequest(args []string) (*bookstack.CommentUpdateRequest, error) {
	req := &bookstack.CommentUpdateRequest{}
	for _, a := range args {
		if strings.HasPrefix(a, "--html=") {
			req.HTML = strings.TrimPrefix(a, "--html=")
		}
	}
	if req.HTML == "" {
		return nil, fmt.Errorf("comments update: --html is required")
	}
	return req, nil
}