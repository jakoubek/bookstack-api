// cli/cli_test.go - unit tests for the hqcli flag parser and dispatch.
package main

import (
	"context"
	"testing"
)

func TestFlagSet_Parse(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "string flag",
			args: []string{"--url", "http://x", "--token-id", "a", "--token-secret", "b"},
			want: map[string]any{"url": "http://x", "token-id": "a", "token-secret": "b"},
		},
		{
			name: "equals form",
			args: []string{"--url=http://x", "--token-id=a", "--token-secret=b"},
			want: map[string]any{"url": "http://x", "token-id": "a", "token-secret": "b"},
		},
		{
			name: "bool flag",
			args: []string{"--all"},
			want: map[string]any{"all": true},
		},
		{
			name: "int flag",
			args: []string{"--timeout=60"},
			want: map[string]any{"timeout": 60},
		},
		{
			name: "stops at first positional",
			args: []string{"--url", "http://x", "books", "list", "--count=5"},
			want: map[string]any{"url": "http://x"},
		},
		{
			name: "unknown flag errors",
			args: []string{"--bogus"},
			wantErr: true,
		},
		{
			name: "bad int errors",
			args: []string{"--timeout=abc"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := newFlagSet("test")
			url := fs.String("url", "", "")
			id := fs.String("token-id", "", "")
			sec := fs.String("token-secret", "", "")
			all := fs.Bool("all", false, "")
			to := fs.Int("timeout", 30, "")

			err := fs.Parse(tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v, ok := tt.want["url"].(string); ok && *url != v {
				t.Errorf("url = %q, want %q", *url, v)
			}
			if v, ok := tt.want["token-id"].(string); ok && *id != v {
				t.Errorf("token-id = %q, want %q", *id, v)
			}
			if v, ok := tt.want["token-secret"].(string); ok && *sec != v {
				t.Errorf("token-secret = %q, want %q", *sec, v)
			}
			if v, ok := tt.want["all"].(bool); ok && *all != v {
				t.Errorf("all = %v, want %v", *all, v)
			}
			if v, ok := tt.want["timeout"].(int); ok && *to != v {
				t.Errorf("timeout = %d, want %d", *to, v)
			}
		})
	}
}

func TestFlagSet_Args(t *testing.T) {
	fs := newFlagSet("test")
	fs.String("url", "", "")
	fs.Parse([]string{"--url", "http://x", "books", "list", "--count=5"})
	got := fs.Args()
	if len(got) != 3 || got[0] != "books" || got[1] != "list" || got[2] != "--count=5" {
		t.Errorf("Args() = %v, want [books list --count=5]", got)
	}
}

func TestArgID(t *testing.T) {
	id, err := argID([]string{"42"})
	if err != nil || id != 42 {
		t.Errorf("argID([42]) = %d, %v, want 42, nil", id, err)
	}
	_, err = argID([]string{})
	if err == nil {
		t.Error("expected error for empty args")
	}
}

func TestExportArgs(t *testing.T) {
	id, fmt, err := exportArgs([]string{"1", "markdown"})
	if err != nil || id != 1 || fmt != "markdown" {
		t.Errorf("exportArgs = %d, %q, %v, want 1, markdown, nil", id, fmt, err)
	}
	_, _, err = exportArgs([]string{"1"})
	if err == nil {
		t.Error("expected error for missing format")
	}
}

func TestParseListOptions(t *testing.T) {
	opts := parseListOptions([]string{"--count=10", "--offset=5", "--sort=name", "--filter=visible=true"})
	if opts.Count != 10 {
		t.Errorf("Count = %d, want 10", opts.Count)
	}
	if opts.Offset != 5 {
		t.Errorf("Offset = %d, want 5", opts.Offset)
	}
	if opts.Sort != "name" {
		t.Errorf("Sort = %q, want name", opts.Sort)
	}
	if opts.Filter == nil {
		t.Fatal("Filter map is nil")
	}
	if opts.Filter["visible"] != "true" {
		t.Errorf("Filter[visible] = %q, want true", opts.Filter["visible"])
	}
}

func TestUsage(t *testing.T) {
	err := usage("books", "list", "get <id>")
	if err == nil {
		t.Error("usage should return errUsage")
	}
}

func TestDispatchUnknownGroup(t *testing.T) {
	err := dispatch(context.Background(), nil, "bogus", "", []string{}, false)
	if err == nil {
		t.Error("dispatch should error on unknown group")
	}
}

func TestOutputError(t *testing.T) {
	err := output(nil, context.Canceled)
	if err == nil {
		t.Error("output should propagate error")
	}
}

func TestExitCodeFor(t *testing.T) {
	if exitCodeFor(nil) != 0 {
		t.Error("nil should map to 0")
	}
}