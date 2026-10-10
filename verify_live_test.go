package bookstack

import (
	"testing"
)

// Verification script for jakoubek/bookstack-api against live BookStack.
// Run from repo root: go test -run VerifyLive ./...
//
// RESULTS (2026-10-09, pi5 / home lab):
// 1. NewClient succeeds with current .env pair (TokenID / TokenSecret set).
// 2. Auth header produced by library: fmt.Sprintf("Token %s:%s", tokenID, tokenSecret) → Token <ID>:<SECRET>
// 3. LIVE API requires: Authorization: Token <SECRET>:<ID> (SECRET first) — confirmed by AGENTS.md / working curl to docs.digipunt.aldof.duckdns.org.
// 4. GAP: library swaps ID/SECRET order in http.go line 33/88. Must fix before reliable use.
// 5. External API endpoint (/api/books) unreachable from this host (TCP timeout to 94.110.157.71:443 after TLS handshake). Could be firewall / DuckDNS / network config, not token. Could also require a specific endpoint path.
// 6. Both token pairs (current .env A and rotated B from memory) time out identically — neither is confirmed invalid; network is the blocker.

func TestVerifyLive(t *testing.T) {
	// This test documents the gap; it does not make live network calls.
	// When network is fixed, replace with real client.Books.List(ctx, nil) call.
	t.Log("Live verification blocked by external TCP timeout to docs.digipunt.aldof.duckdns.org; auth-order bug in library must be fixed before reliable use.")
}
