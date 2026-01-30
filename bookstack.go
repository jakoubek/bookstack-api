package bookstack

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

// Config holds configuration for the Bookstack API client.
type Config struct {
	// BaseURL is the base URL of the Bookstack instance (e.g., "https://docs.example.com")
	BaseURL string

	// TokenID is the Bookstack API token ID
	TokenID string

	// TokenSecret is the Bookstack API token secret
	TokenSecret string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, a default client with 30s timeout will be used.
	HTTPClient *http.Client
}

// Client is the main Bookstack API client.
type Client struct {
	baseURL     string
	tokenID     string
	tokenSecret string
	httpClient  *http.Client

	// Service instances
	Attachments *AttachmentsService
	Books       *BooksService
	Chapters    *ChaptersService
	Pages       *PagesService
	Search      *SearchService
	Shelves     *ShelvesService
}

// NewClient creates a new Bookstack API client.
// Returns an error if BaseURL, TokenID, or TokenSecret are empty.
func NewClient(cfg Config) (*Client, error) {
	var errs []string
	if cfg.BaseURL == "" {
		errs = append(errs, "BaseURL is required")
	}
	if cfg.TokenID == "" {
		errs = append(errs, "TokenID is required")
	}
	if cfg.TokenSecret == "" {
		errs = append(errs, "TokenSecret is required")
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "; "))
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	c := &Client{
		baseURL:     strings.TrimRight(cfg.BaseURL, "/"),
		tokenID:     cfg.TokenID,
		tokenSecret: cfg.TokenSecret,
		httpClient:  httpClient,
	}

	// Initialize services
	c.Attachments = &AttachmentsService{client: c}
	c.Books = &BooksService{client: c}
	c.Chapters = &ChaptersService{client: c}
	c.Pages = &PagesService{client: c}
	c.Search = &SearchService{client: c}
	c.Shelves = &ShelvesService{client: c}

	return c, nil
}
