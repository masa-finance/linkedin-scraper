package linkedinscraper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is the LinkedIn API client.
type Client struct {
	httpClient *http.Client
	config     *Config
}

// NewClient creates a new LinkedIn API client.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("linkedinscraper: config cannot be nil") // Consider defining a specific error for this
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second, // Go's default http.Transport handles gzip automatically
	}

	return &Client{httpClient: httpClient, config: cfg}, nil
}

// buildGraphQLURL constructs the GraphQL request URL.
func buildGraphQLURL(baseURL string, queryID string, variables SearchVariables) (string, error) {
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return "", fmt.Errorf("failed to marshal search variables: %w", err)
	}

	encodedVariables := url.QueryEscape(string(variablesJSON))

	// Construct the full URL: baseURL + "?includeWebMetadata=true&variables=(" + encodedVariables + ")&queryId=" + queryID
	// Using fmt.Sprintf for clarity
	return fmt.Sprintf("%s?includeWebMetadata=true&variables=(%s)&queryId=%s", baseURL, encodedVariables, queryID), nil
}

// makeRequest executes an HTTP request.
func (c *Client) makeRequest(ctx context.Context, method string, urlStr string, headers http.Header, body io.Reader) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set common headers from constants and config
	req.Header.Set("Accept", AcceptHeaderValue)
	req.Header.Set("Accept-Encoding", AcceptEncodingHeaderValue)
	req.Header.Set("Accept-Language", AcceptLanguageHeaderValue)
	req.Header.Set("Csrf-Token", c.config.Auth.CSRFToken)
	req.Header.Set("X-Li-Lang", DefaultLiLangHeaderValue)
	req.Header.Set("X-Restli-Protocol-Version", DefaultRestliProtocolVersion)
	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Cookie", fmt.Sprintf("li_at=%s; JSESSIONID=\"%s\"", c.config.Auth.LiAtCookie, c.config.Auth.JSESSIONID))

	// Add any other headers passed in the headers argument
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http client failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, respBodyBytes, nil
}
