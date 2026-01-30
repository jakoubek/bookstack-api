package bookstack

import (
	"net/http"
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
	Books    *BooksService
	Pages    *PagesService
	Chapters *ChaptersService
	Shelves  *ShelvesService
	Search   *SearchService
}

// NewClient creates a new Bookstack API client.
func NewClient(cfg Config) *Client {
	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	c := &Client{
		baseURL:     cfg.BaseURL,
		tokenID:     cfg.TokenID,
		tokenSecret: cfg.TokenSecret,
		httpClient:  httpClient,
	}

	// Initialize services
	c.Books = &BooksService{client: c}
	c.Pages = &PagesService{client: c}
	c.Chapters = &ChaptersService{client: c}
	c.Shelves = &ShelvesService{client: c}
	c.Search = &SearchService{client: c}

	return c
}
