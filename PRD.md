# Product Requirements Document: bookstack-api

## 1. Projektübersicht

- **Ziel und Vision:**  
  Eine Go-Library für die Bookstack REST-API, die programmatischen Zugriff auf das Dokumentationssystem ermöglicht. Die Library wird in hqcli integriert und soll veröffentlichungsfähig sein.

- **Zielgruppe:**  
  - Primär: AI-Agenten (via hqcli mit `--json` Flag)
  - Sekundär: Go-Entwickler, die Bookstack programmatisch nutzen wollen

- **Erfolgs-Metriken:**  
  - Vollständige Abdeckung der Kern-API (Books, Pages, Search)
  - Nutzbar in hqcli für AI-gestützte Dokumentationssuche
  - Veröffentlichung als eigenständiges Go-Modul

- **Projektscope:**  
  - **In Scope:** REST-API-Client für Bookstack, Iterator-basierte Pagination, Export-Funktionen
  - **Out of Scope:** Webhooks, Caching, Admin-Funktionen (Users, Roles)

## 2. Funktionale Anforderungen

### Kern-Features

| Feature | Priorität | Version |
|---------|-----------|---------|
| Books: List, Get | P0 | v0.1 |
| Pages: List, Get | P0 | v0.1 |
| Search: All | P0 | v0.1 |
| Pages: Export (Markdown, PDF) | P1 | v0.2 |
| Chapters: List, Get | P1 | v0.2 |
| Shelves: List, Get | P1 | v0.2 |
| Pages: Create, Update | P1 | v0.3 |
| Pages: Delete | P2 | v0.4 |
| Attachments: CRUD | P2 | v0.4 |
| Comments: CRUD | P3 | v0.5 |

### User Stories

**US1: Dokumentation durchsuchen (AI-Agent)**
> Als AI-Agent möchte ich die Bookstack-Dokumentation durchsuchen können, um relevante Seiten für Benutzeranfragen zu finden.

```bash
hqcli docs search "deployment" --json
```

**US2: Seite abrufen**
> Als Benutzer möchte ich eine Dokumentationsseite anzeigen können, um deren Inhalt zu lesen.

```bash
hqcli docs page 123
hqcli docs page deployment-guide  # via Slug
```

**US3: Bücher auflisten**
> Als Benutzer möchte ich alle verfügbaren Bücher sehen, um die Dokumentationsstruktur zu verstehen.

```bash
hqcli docs books --json
```

**US4: Seite exportieren**
> Als Benutzer möchte ich eine Seite als Markdown oder PDF exportieren können.

```bash
hqcli docs page 123 --export=md > page.md
hqcli docs page 123 --export=pdf > page.pdf
```

**US5: Seite im Browser öffnen**
> Als Benutzer möchte ich eine Seite schnell im Browser öffnen können.

```bash
hqcli docs open 123
```

### Detaillierte Workflows

**Workflow: Suche und Anzeige**
```
1. Benutzer/Agent führt Suche aus
2. API gibt Liste von Treffern zurück (ID, Typ, Name, Preview)
3. Benutzer/Agent wählt Treffer aus
4. Seite wird abgerufen und angezeigt (Markdown oder JSON)
```

**Workflow: Seite bearbeiten (v0.3)**
```
1. Seite abrufen (Get)
2. Inhalt lokal bearbeiten
3. Seite aktualisieren (Update)
```

### Feature-Prioritäten

- **Must-have (v1):** List, Get, Search für Books/Pages
- **Should-have (v1):** Export Markdown/PDF, Chapters, Shelves
- **Nice-to-have (v2):** Create, Update, Delete
- **Future:** Attachments, Comments, Image Gallery

## 3. Technische Anforderungen

- **Performance-Ziele:**  
  - API-Calls < 500ms (abhängig von Netzwerk)
  - Iterator verarbeitet 10.000+ Einträge ohne Memory-Probleme

- **Concurrent User-Kapazität:**  
  Nicht zutreffend (Library, kein Server)

- **Real-time Features:**  
  Nicht zutreffend

- **Sicherheitsstandards:**  
  - Token-basierte Authentifizierung (Token ID + Secret)
  - Keine Speicherung von Credentials (Aufrufer-Verantwortung)

- **Compliance-Vorgaben:**  
  Keine speziellen

- **Plattform-Support:**  
  - Go 1.21+
  - Linux, macOS, Windows

## 4. Datenarchitektur

*Nicht zutreffend – keine eigene Datenhaltung*

### Bookstack-Hierarchie (extern)

```
Shelf (Regal)
└── Book (Buch)
    ├── Chapter (Kapitel)
    │   └── Page (Seite)
    └── Page (Seite)
```

### Datenstrukturen

```go
// Book repräsentiert ein Bookstack-Buch
type Book struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   int       `json:"created_by"`
    UpdatedBy   int       `json:"updated_by"`
}

// Page repräsentiert eine Bookstack-Seite
type Page struct {
    ID        int       `json:"id"`
    BookID    int       `json:"book_id"`
    ChapterID int       `json:"chapter_id"`
    Name      string    `json:"name"`
    Slug      string    `json:"slug"`
    HTML      string    `json:"html"`
    RawHTML   string    `json:"raw_html"`
    Markdown  string    `json:"markdown"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Chapter repräsentiert ein Bookstack-Kapitel
type Chapter struct {
    ID          int       `json:"id"`
    BookID      int       `json:"book_id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// Shelf repräsentiert ein Bookstack-Regal
type Shelf struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Slug        string    `json:"slug"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// SearchResult repräsentiert ein Suchergebnis
type SearchResult struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Slug    string `json:"slug"`
    Type    string `json:"type"` // page, chapter, book, bookshelf
    URL     string `json:"url"`
    Preview string `json:"preview"`
}
```

## 5. API & Interface-Spezifikation

### Client-Initialisierung

```go
type Config struct {
    BaseURL     string       // z.B. "https://docs.jakoubek.net"
    TokenID     string       // API Token ID
    TokenSecret string       // API Token Secret
    HTTPClient  *http.Client // optional, für Tests/Mocking
}

func NewClient(cfg Config) *Client
```

**Beispiel:**
```go
client := bookstack.NewClient(bookstack.Config{
    BaseURL:     "https://docs.jakoubek.net",
    TokenID:     os.Getenv("BOOKSTACK_TOKEN_ID"),
    TokenSecret: os.Getenv("BOOKSTACK_TOKEN_SECRET"),
})
```

### Service-Struktur

```go
type Client struct {
    Books    *BooksService
    Pages    *PagesService
    Chapters *ChaptersService
    Shelves  *ShelvesService
    Search   *SearchService
}
```

### REST-Endpoints (Bookstack API)

| Service | Methode | Endpoint | Beschreibung |
|---------|---------|----------|--------------|
| Books | List | GET /api/books | Alle Bücher |
| Books | Get | GET /api/books/{id} | Einzelnes Buch |
| Pages | List | GET /api/pages | Alle Seiten |
| Pages | Get | GET /api/pages/{id} | Einzelne Seite |
| Pages | Create | POST /api/pages | Seite erstellen |
| Pages | Update | PUT /api/pages/{id} | Seite aktualisieren |
| Pages | Delete | DELETE /api/pages/{id} | Seite löschen |
| Pages | ExportMD | GET /api/pages/{id}/export/markdown | Markdown-Export |
| Pages | ExportPDF | GET /api/pages/{id}/export/pdf | PDF-Export |
| Chapters | List | GET /api/chapters | Alle Kapitel |
| Chapters | Get | GET /api/chapters/{id} | Einzelnes Kapitel |
| Shelves | List | GET /api/shelves | Alle Regale |
| Shelves | Get | GET /api/shelves/{id} | Einzelnes Regal |
| Search | All | GET /api/search?query=... | Volltextsuche |

### Authentifizierung

```
Authorization: Token <token_id>:<token_secret>
```

### Pagination (Iterator-Pattern)

```go
// ListAll gibt einen Iterator über alle Einträge zurück
// Nutzt Go 1.23+ iter.Seq oder eigene Implementation
func (s *BooksService) ListAll(ctx context.Context) iter.Seq2[*Book, error]

// Nutzung:
for book, err := range client.Books.ListAll(ctx) {
    if err != nil {
        return err
    }
    fmt.Println(book.Name)
}
```

**Begründung:** Iterator-Pattern ist Go-idiomatisch (ab 1.23), Memory-effizient und ermöglicht frühen Abbruch.

### Rate Limiting

- Bookstack: 180 Requests/Minute (default)
- **Keine Library-interne Behandlung** – Aufrufer muss Rate-Limiting selbst handhaben
- Bei 429-Response: `ErrRateLimited` zurückgeben

## 6. Benutzeroberfläche

*Nicht zutreffend – Library ohne UI*

### hqcli-Integration (separates Projekt)

```bash
# Bücher auflisten
hqcli docs books
hqcli docs books --json

# Seiten eines Buchs
hqcli docs pages --book=<id|slug>

# Seite anzeigen
hqcli docs page <id|slug>
hqcli docs page <id> --json

# Seite exportieren
hqcli docs page <id> --export=md
hqcli docs page <id> --export=pdf

# Suche
hqcli docs search "query"
hqcli docs search "query" --json

# Im Browser öffnen
hqcli docs open <id|slug>
```

**Output-Priorität:** `--json` für AI-Agenten ist Hauptanwendungsfall

## 7. Nicht-funktionale Anforderungen

- **Verfügbarkeit:**  
  Nicht zutreffend (Library)

- **Dependencies:**  
  Nur Go-Standardbibliothek (net/http, encoding/json, etc.)

- **Backward Compatibility:**  
  Semantic Versioning (v0.x während Entwicklung, v1.x nach Stabilisierung)

- **Logging-Strategie:**  
  Keine eigene Logging – Fehler werden als `error` zurückgegeben

- **Konfiguration:**  
  Via `Config`-Struct bei Client-Erstellung

## 8. Qualitätssicherung

### Definition of Done

- [ ] Alle public APIs dokumentiert (GoDoc)
- [ ] Unit-Tests für alle Services (≥80% Coverage)
- [ ] Integration-Tests gegen Mock-Server
- [ ] Beispiel-Code in examples/
- [ ] README mit Quick-Start

### Test-Anforderungen

```go
// Unit-Tests mit Mock-HTTP-Client
func TestBooksService_List(t *testing.T) {
    server := httptest.NewServer(...)
    client := bookstack.NewClient(bookstack.Config{
        BaseURL:    server.URL,
        TokenID:    "test",
        TokenSecret: "test",
    })
    
    books, err := client.Books.List(ctx, nil)
    // assertions...
}
```

### Launch-Kriterien v1.0

- [ ] Books, Pages, Search vollständig implementiert
- [ ] Export (Markdown, PDF) funktioniert
- [ ] Dokumentation vollständig
- [ ] Keine bekannten Bugs
- [ ] hqcli-Integration getestet

## 9. Technische Implementierungshinweise

### Go-Projektstruktur

```
bookstack-api/
├── bookstack.go       # Client, Config, NewClient()
├── books.go           # BooksService
├── pages.go           # PagesService
├── chapters.go        # ChaptersService
├── shelves.go         # ShelvesService
├── search.go          # SearchService
├── types.go           # Alle Datenstrukturen
├── errors.go          # Error-Typen
├── http.go            # HTTP-Helfer, Request-Building
├── iterator.go        # Pagination-Iterator
├── bookstack_test.go  # Tests
├── README.md
├── go.mod
├── go.sum
└── examples/
    └── basic/
        └── main.go
```

### Error-Handling-Strategie

```go
// Definierte Error-Typen für häufige Fälle
var (
    ErrNotFound     = errors.New("bookstack: resource not found")
    ErrUnauthorized = errors.New("bookstack: unauthorized")
    ErrForbidden    = errors.New("bookstack: forbidden")
    ErrRateLimited  = errors.New("bookstack: rate limited")
    ErrBadRequest   = errors.New("bookstack: bad request")
)

// APIError für detaillierte Fehlerinformationen
type APIError struct {
    StatusCode int
    Code       int    `json:"code"`
    Message    string `json:"message"`
    Body       []byte // Original Response Body
}

func (e *APIError) Error() string {
    return fmt.Sprintf("bookstack: API error %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Is(target error) bool {
    switch target {
    case ErrNotFound:
        return e.StatusCode == 404
    case ErrUnauthorized:
        return e.StatusCode == 401
    case ErrForbidden:
        return e.StatusCode == 403
    case ErrRateLimited:
        return e.StatusCode == 429
    case ErrBadRequest:
        return e.StatusCode == 400
    }
    return false
}
```

**Begründung:** `errors.Is()` ermöglicht einfache Fehlerprüfung, `APIError` bietet Details wenn nötig.

### HTTP-Wrapper

```go
// Interner HTTP-Helfer
func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
    // 1. Request bauen
    // 2. Auth-Header setzen
    // 3. Request ausführen
    // 4. Response prüfen
    // 5. Bei Fehler: APIError zurückgeben
    // 6. Bei Erfolg: JSON in result unmarshalen
}
```

### Entwicklungs-Prioritäten

1. **Phase 1 (v0.1):** Foundation
   - Client-Setup, Auth, HTTP-Wrapper
   - Books.List, Books.Get
   - Pages.List, Pages.Get
   - Error-Handling

2. **Phase 2 (v0.2):** Core Features
   - Search.All
   - Pages.ExportMarkdown, Pages.ExportPDF
   - Iterator für Pagination
   - Chapters, Shelves

3. **Phase 3 (v0.3):** Write Operations
   - Pages.Create
   - Pages.Update

4. **Phase 4 (v1.0):** Release
   - Dokumentation
   - Beispiele
   - CI/CD
   - Veröffentlichung

### Potenzielle Risiken

| Risiko | Wahrscheinlichkeit | Mitigation |
|--------|-------------------|------------|
| API-Änderungen in Bookstack | Niedrig | Semantic Versioning, Tests |
| Rate-Limiting-Probleme | Mittel | Dokumentation für Aufrufer |
| Große PDF-Exports | Mittel | Streaming statt Buffer |

---

## Anhang: Bookstack API-Referenz

- **Basis-URL:** https://docs.jakoubek.net/api
- **Dokumentation:** https://docs.jakoubek.net/api/docs
- **Beispiele:** https://codeberg.org/bookstack/api-scripts
- **Rate Limit:** 180 req/min (konfigurierbar serverseitig)

### Pagination-Parameter

| Parameter | Beschreibung | Default |
|-----------|--------------|---------|
| count | Anzahl Ergebnisse | 100 (max 500) |
| offset | Start-Position | 0 |
| sort | Sortierung (+name, -created_at) | - |
| filter[field] | Filter (eq, ne, gt, lt, like) | - |
