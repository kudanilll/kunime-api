package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kunime-api/internal/anime"
	"kunime-api/internal/config"
)

func TestRouteContracts(t *testing.T) {
	app := NewServer(config.Config{APIKey: "test-key", ScrapeBaseURL: "https://example.invalid"}, anime.NewService(nil))
	for _, tc := range []struct {
		name, method, path, key, body string
		status                        int
	}{
		{"public health", http.MethodGet, "/healthz", "", "", http.StatusOK},
		{"catalog requires a key", http.MethodGet, "/api/v1/scrape-base-url", "", "", http.StatusUnauthorized},
		{"catalog with a key", http.MethodGet, "/api/v1/scrape-base-url", "test-key", "", http.StatusOK},
		{"stream request rejects invalid JSON", http.MethodPost, "/api/v1/streams/resolve", "test-key", "{", http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("X-API-Key", tc.key)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != tc.status {
				t.Fatalf("status = %d, want %d", resp.StatusCode, tc.status)
			}
		})
	}
}
