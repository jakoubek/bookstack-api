package bookstack

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("success with all fields", func(t *testing.T) {
		c, err := NewClient(Config{
			BaseURL:     "https://docs.example.com",
			TokenID:     "abc",
			TokenSecret: "xyz",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.baseURL != "https://docs.example.com" {
			t.Errorf("baseURL = %q, want %q", c.baseURL, "https://docs.example.com")
		}
	})

	t.Run("default HTTPClient", func(t *testing.T) {
		c, err := NewClient(Config{
			BaseURL:     "https://docs.example.com",
			TokenID:     "abc",
			TokenSecret: "xyz",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient == nil {
			t.Fatal("httpClient should not be nil")
		}
	})

	t.Run("custom HTTPClient", func(t *testing.T) {
		custom := &http.Client{}
		c, err := NewClient(Config{
			BaseURL:     "https://docs.example.com",
			TokenID:     "abc",
			TokenSecret: "xyz",
			HTTPClient:  custom,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.httpClient != custom {
			t.Error("expected custom HTTPClient to be used")
		}
	})

	t.Run("trailing slash stripped", func(t *testing.T) {
		c, err := NewClient(Config{
			BaseURL:     "https://docs.example.com/",
			TokenID:     "abc",
			TokenSecret: "xyz",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.baseURL != "https://docs.example.com" {
			t.Errorf("baseURL = %q, want trailing slash stripped", c.baseURL)
		}
	})

	t.Run("error on missing BaseURL", func(t *testing.T) {
		_, err := NewClient(Config{TokenID: "abc", TokenSecret: "xyz"})
		if err == nil {
			t.Fatal("expected error for missing BaseURL")
		}
	})

	t.Run("error on missing TokenID", func(t *testing.T) {
		_, err := NewClient(Config{BaseURL: "https://x.com", TokenSecret: "xyz"})
		if err == nil {
			t.Fatal("expected error for missing TokenID")
		}
	})

	t.Run("error on missing TokenSecret", func(t *testing.T) {
		_, err := NewClient(Config{BaseURL: "https://x.com", TokenID: "abc"})
		if err == nil {
			t.Fatal("expected error for missing TokenSecret")
		}
	})

	t.Run("error on all fields missing", func(t *testing.T) {
		_, err := NewClient(Config{})
		if err == nil {
			t.Fatal("expected error for all missing fields")
		}
	})
}
