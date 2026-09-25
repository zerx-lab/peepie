package searchers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"pentagi/pkg/config"
	"pentagi/pkg/database"
)

const testBraveAPIKey = "test-brave-key"

func testBraveConfig() *config.Config {
	return &config.Config{BraveAPIKey: testBraveAPIKey}
}

func TestBraveHandle(t *testing.T) {
	var seenRequest bool
	var receivedMethod string
	var receivedToken string
	var receivedAccept string
	var receivedQuery string
	var receivedCount string

	mockMux := http.NewServeMux()
	mockMux.HandleFunc("/res/v1/web/search", func(w http.ResponseWriter, r *http.Request) {
		seenRequest = true
		receivedMethod = r.Method
		receivedToken = r.Header.Get("X-Subscription-Token")
		receivedAccept = r.Header.Get("Accept")
		receivedQuery = r.URL.Query().Get("q")
		receivedCount = r.URL.Query().Get("count")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"web": {
				"results": [
					{"title": "Test Result", "url": "https://example.com/test", "description": "A <strong>test</strong> result", "age": "3 days ago"}
				]
			}
		}`))
	})

	proxy, err := newTestProxy("api.search.brave.com", mockMux)
	if err != nil {
		t.Fatalf("failed to create proxy: %v", err)
	}
	defer proxy.Close()

	cfg := &config.Config{
		BraveAPIKey:       testBraveAPIKey,
		ProxyURL:          proxy.URL(),
		ExternalSSLCAPath: proxy.CACertPath(),
	}

	b := NewBrave(cfg)

	got, err := b.Handle(t.Context(), Request{Query: "test query", MaxResults: 5})
	if err != nil {
		t.Fatalf("Handle() unexpected error: %v", err)
	}
	if b.Engine() != database.SearchengineTypeBrave {
		t.Errorf("Engine() = %q, want %q", b.Engine(), database.SearchengineTypeBrave)
	}

	if !seenRequest {
		t.Fatal("request was not intercepted by proxy - mock handler was not called")
	}
	if receivedMethod != http.MethodGet {
		t.Errorf("request method = %q, want GET", receivedMethod)
	}
	if receivedToken != testBraveAPIKey {
		t.Errorf("X-Subscription-Token = %q, want %q", receivedToken, testBraveAPIKey)
	}
	if receivedAccept != "application/json" {
		t.Errorf("Accept = %q, want application/json", receivedAccept)
	}
	if receivedQuery != "test query" {
		t.Errorf("q = %q, want %q", receivedQuery, "test query")
	}
	if receivedCount != "5" {
		t.Errorf("count = %q, want %q", receivedCount, "5")
	}

	if !strings.Contains(got, "Test Result") {
		t.Errorf("result missing expected title: %q", got)
	}
	if !strings.Contains(got, "https://example.com/test") {
		t.Errorf("result missing expected URL: %q", got)
	}
	if !strings.Contains(got, "A test result") {
		t.Errorf("result should strip <strong> tags: %q", got)
	}
	if strings.Contains(got, "<strong>") {
		t.Errorf("result must not contain raw highlight markup: %q", got)
	}
	if !strings.Contains(got, "3 days ago") {
		t.Errorf("result missing age: %q", got)
	}
}

func TestBraveIsAvailable(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		want bool
	}{
		{name: "available with api key", cfg: testBraveConfig(), want: true},
		{name: "unavailable without api key", cfg: &config.Config{}, want: false},
		{name: "unavailable when nil config", cfg: nil, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &brave{cfg: tt.cfg}
			if got := b.IsAvailable(); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBraveHandle_NotConfigured(t *testing.T) {
	b := NewBrave(&config.Config{})
	_, err := b.Handle(context.Background(), Request{Query: "q", MaxResults: 5})
	if !errors.Is(err, ErrNotConfigured) {
		t.Errorf("Handle() error = %v, want ErrNotConfigured", err)
	}
}

func TestBraveParseHTTPResponse_StatusAndDecodeErrors(t *testing.T) {
	b := &brave{}

	t.Run("429 is retryable", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusTooManyRequests, Body: io.NopCloser(strings.NewReader(""))}
		_, err := b.parseHTTPResponse(resp, "q")
		if !IsRetryable(err) {
			t.Errorf("expected retryable error for 429, got: %v", err)
		}
	})

	t.Run("401 is fatal", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusUnauthorized, Body: io.NopCloser(strings.NewReader(""))}
		_, err := b.parseHTTPResponse(resp, "q")
		if !IsFatal(err) {
			t.Errorf("expected fatal error for 401, got: %v", err)
		}
	})

	t.Run("decode error", func(t *testing.T) {
		resp := &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{invalid json"))}
		_, err := b.parseHTTPResponse(resp, "q")
		if err == nil || !strings.Contains(err.Error(), "failed to decode response body") {
			t.Fatalf("expected decode error, got: %v", err)
		}
	})
}

func TestBraveFormatResults_NoResults(t *testing.T) {
	b := &brave{}
	result := b.formatResults(nil, "test query")
	if !strings.Contains(result, "No Results Found") {
		t.Errorf("result missing 'No Results Found': %q", result)
	}
	if !strings.Contains(result, "test query") {
		t.Errorf("result missing query: %q", result)
	}
}

func TestBraveMaxResultsClamp(t *testing.T) {
	tests := []struct {
		name       string
		maxResults int
		wantCount  string
	}{
		{name: "zero uses default", maxResults: 0, wantCount: "10"},
		{name: "negative uses default", maxResults: -1, wantCount: "10"},
		{name: "over max uses default", maxResults: 100, wantCount: "10"},
		{name: "within range is preserved", maxResults: 3, wantCount: "3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedCount string
			mockMux := http.NewServeMux()
			mockMux.HandleFunc("/res/v1/web/search", func(w http.ResponseWriter, r *http.Request) {
				receivedCount = r.URL.Query().Get("count")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"web":{"results":[]}}`))
			})

			proxy, err := newTestProxy("api.search.brave.com", mockMux)
			if err != nil {
				t.Fatalf("failed to create proxy: %v", err)
			}
			defer proxy.Close()

			cfg := &config.Config{
				BraveAPIKey:       testBraveAPIKey,
				ProxyURL:          proxy.URL(),
				ExternalSSLCAPath: proxy.CACertPath(),
			}
			b := NewBrave(cfg)

			if _, err := b.Handle(t.Context(), Request{Query: "q", MaxResults: tt.maxResults}); err != nil {
				t.Fatalf("Handle() unexpected error: %v", err)
			}
			if receivedCount != tt.wantCount {
				t.Errorf("count = %q, want %q", receivedCount, tt.wantCount)
			}
		})
	}
}
