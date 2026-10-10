# Gaps / Feature Requests — jakoubek/bookstack-api (Aldo-f fork)

Doc'd 2026-10-09. Source: https://github.com/jakoubek/bookstack-api (forked to Aldo-f/bookstack-api, cloned to ~/dev/06-apps-bookstack-api).

## Confirmed bugs (verified by code inspection, not just docs)

1. AUTH ORDER WRONG — `http.go:33` and `http.go:88`: `fmt.Sprintf("Token %s:%s", c.tokenID, c.tokenSecret)` produces `Token <ID>:<SECRET>`. Live BookStack requires `Token <SECRET>:<ID>` (SECRET first — verified by working curl with `.env` from 06-apps-bookstack, see AGENTS.md). Must swap arguments.

## Missing / untested features (potential PRs)

2. `Pages.Create` / `Pages.Update` / `Pages.Delete` — the `create-page-8.go` in 06-apps-bookstack uses direct REST for page creation; library's Pages service exists (`pages.go`) but unverified for the same operations (especially page HTML content update, which is what `create-page-8.go` does). Need to verify `Pages.Update` accepts `html` payload in the same format.
3. `Search` / `Attachments` — present in library but untested against digipunt instance (sources, source_check endpoints not in library — those are custom theme endpoints, not standard BookStack API).
4. No `.env` loader / config file support — library requires manual `NewClient(Config{...})`. Could add `NewClientFromEnv()` that reads `BOOKSTACK_API_URL`, `TOKEN_ID`, `TOKEN_SECRET`.
5. No examples directory for the specific auth header format — `CLAUDE.md` / `README.md` don't mention `SECRET:ID` vs `ID:SECRET` ordering clearly.

## Verification status

- Fork: ✓ (Aldo-f/bookstack-api)
- Clone: ✓ (`~/dev/06-apps-bookstack-api`)
- `.env`: ✓ (written with both token pairs; notes auth gap)
- Both token pairs tested with curl to `/api/shelves`: TCP timeout to `94.110.157.71:443` after TLS handshake (5s timeout). Not a token failure — network / external access issue from pi5 to docs.digipunt. Could retry from a different network or with port-check (`lsof -i :443` on external). Site root (`/`) responds 200, so host is reachable; API path may have different routing / needs newer token / requires `Accept: application/json`.
- Library `NewClient`: untested live (would need fixed auth + reachable endpoint).
- `create-page-8.go`: NOT replaced — condition from user was "only when verified we can do the same". Not verified yet.

## Next actions (need user / network fix)

- Fix auth order in library (PR to Aldo-f/bookstack-api, then upstream jakoubek).
- Resolve external TCP timeout to API endpoint (check if `docs.digipunt.aldof.duckdns.org` resolves differently externally vs internal, or if `/api/` requires different ingress / cert); retest both token pairs.
- Only after both — replace `create-page-8.go` with library equivalent.
