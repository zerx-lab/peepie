package searchers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"pentagi/pkg/config"
	"pentagi/pkg/database"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/observability/langfuse"
	"pentagi/pkg/system"

	"github.com/sirupsen/logrus"
)

const (
	braveSearchURL      = "https://api.search.brave.com/res/v1/web/search"
	braveTimeout        = 30 * time.Second
	braveDefaultResults = 10
	braveMaxResults     = 20
)

// brave is a Searcher backed by the Brave Search API (https://brave.com/search/api/).
// It is a plain link-discovery engine: fast, cheap (2,000 free queries/month), and
// does not summarize — it returns the same shape of result as google/duckduckgo/searxng.
type brave struct {
	cfg *config.Config
}

// NewBrave creates a new Brave Search primitive.
func NewBrave(cfg *config.Config) Searcher {
	return &brave{cfg: cfg}
}

func (b *brave) Engine() database.SearchengineType {
	return database.SearchengineTypeBrave
}

func (b *brave) IsAvailable() bool {
	return b.apiKey() != ""
}

// Handle processes the search request from the orchestrator.
func (b *brave) Handle(ctx context.Context, req Request) (string, error) {
	if !b.IsAvailable() {
		return "", ErrNotConfigured
	}

	ctx, observation := obs.Observer.NewObservation(ctx)

	numResults := req.MaxResults
	if numResults < 1 || numResults > braveMaxResults {
		numResults = braveDefaultResults
	}

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"engine":      "brave",
		"query":       req.Query[:min(len(req.Query), 1000)],
		"num_results": numResults,
	})

	result, err := b.search(ctx, req.Query, numResults)
	if err != nil {
		observation.Event(
			langfuse.WithEventName("search engine error"),
			langfuse.WithEventInput(req.Query),
			langfuse.WithEventStatus(err.Error()),
			langfuse.WithEventLevel(langfuse.ObservationLevelWarning),
			langfuse.WithEventMetadata(langfuse.Metadata{
				"engine":      "brave",
				"query":       req.Query,
				"max_results": numResults,
				"error":       err.Error(),
			}),
		)

		obs.LogErrorOrCancel(logger, err, "failed to search in Brave")
		return "", err
	}

	return result, nil
}

func (b *brave) search(ctx context.Context, query string, maxResults int) (string, error) {
	apiURL, err := url.Parse(braveSearchURL)
	if err != nil {
		return "", Fatal(fmt.Errorf("invalid brave search URL: %w", err))
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("count", strconv.Itoa(maxResults))
	apiURL.RawQuery = params.Encode()

	client, err := system.GetHTTPClient(b.cfg)
	if err != nil {
		return "", Fatal(fmt.Errorf("failed to create http client: %w", err))
	}

	client.Timeout = braveTimeout

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return "", Fatal(fmt.Errorf("failed to create request: %w", err))
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Subscription-Token", b.apiKey())

	resp, err := client.Do(req)
	if err != nil {
		return "", Retryable(fmt.Errorf("failed to do request: %w", err), 0)
	}
	defer resp.Body.Close()

	return b.parseHTTPResponse(resp, query)
}

func (b *brave) parseHTTPResponse(resp *http.Response, query string) (string, error) {
	if resp.StatusCode != http.StatusOK {
		return "", ClassifyHTTPStatus(resp.StatusCode, "brave search failed")
	}

	var braveResp braveSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&braveResp); err != nil {
		return "", Fatal(fmt.Errorf("failed to decode response body: %w", err))
	}

	return b.formatResults(braveResp.Web.Results, query), nil
}

func (b *brave) formatResults(results []braveResult, query string) string {
	if len(results) == 0 {
		return fmt.Sprintf("# No Results Found\n\nNo results were found for query: %s", query)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("# Brave Search Results\n\n## Query: %s\n\n", query))

	for i, result := range results {
		builder.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, result.Title))

		if result.URL != "" {
			builder.WriteString(fmt.Sprintf("**URL:** [%s](%s)\n\n", result.URL, result.URL))
		}

		if result.Description != "" {
			builder.WriteString(fmt.Sprintf("**Description:** %s\n\n", stripBraveHighlightTags(result.Description)))
		}

		if result.Age != "" {
			builder.WriteString(fmt.Sprintf("**Age:** %s\n\n", result.Age))
		}

		builder.WriteString("---\n\n")
	}

	return builder.String()
}

// stripBraveHighlightTags removes the <strong>/</strong> query-highlight markup Brave
// embeds in descriptions; the markdown result is plain text, not HTML.
func stripBraveHighlightTags(s string) string {
	s = strings.ReplaceAll(s, "<strong>", "")
	s = strings.ReplaceAll(s, "</strong>", "")
	return s
}

func (b *brave) apiKey() string {
	if b.cfg == nil {
		return ""
	}

	return b.cfg.Overrides.GetString(config.CategorySearchEngines, config.KeyBraveAPIKey, b.cfg.BraveAPIKey)
}

// braveResult is a single organic web result from the Brave Search API.
type braveResult struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Age         string `json:"age"`
}

// braveWebResults holds the "web" result cluster of a Brave Search API response.
type braveWebResults struct {
	Results []braveResult `json:"results"`
}

// braveSearchResponse is the top-level Brave Search API response shape. Brave also
// returns "news", "videos", and "infobox" clusters; only organic web results are used.
type braveSearchResponse struct {
	Web braveWebResults `json:"web"`
}
