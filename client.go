package linkedinscraper

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// buildGraphQLURL constructs the full URL for a GraphQL API request.
// It takes the base URL, query ID, and variables, then assembles them.
func buildGraphQLURL(baseURL, queryID string, variables SearchVariables) (string, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// Manually construct the variables string to match the cURL format
	// (start:0,count:1,origin:FACETED_SEARCH,query:(keywords:investor,flagshipSearchIntent:SEARCH_SRP,queryParameters:List((key:network,value:List(F,O)),(key:resultType,value:List(PEOPLE))),includeFiltersInResponse:false))
	var queryParams []string
	for _, p := range variables.Query.QueryParameters {
		// Assuming p.Value is always a list of strings for now.
		// The cURL shows List(F,O) or List(PEOPLE). We need to join them with commas.
		valueList := "List(" + stringSliceToString(p.Value, ",") + ")"
		queryParams = append(queryParams, fmt.Sprintf("(key:%s,value:%s)", p.Key, valueList))
	}
	queryParametersString := "List(" + stringSliceToString(queryParams, ",") + ")"

	// Ensure keywords are properly escaped for the URL query string part, but not for the graphql variable part
	// The variable string itself is a single query parameter value, so special characters within it are fine.
	// However, if keywords themselves contain characters like '(', ')', ',', they should be as-is per cURL.

	// Reverted: Use full variablesString including queryParameters
	variablesString := fmt.Sprintf("(start:%d,count:%d,origin:%s,query:(keywords:%s,flagshipSearchIntent:%s,queryParameters:%s,includeFiltersInResponse:%t))",
		variables.Start,
		variables.Count,
		variables.Origin,
		url.QueryEscape(variables.Query.Keywords), // URL Encode the keywords string for spaces etc.
		variables.Query.FlagshipSearchIntent,
		queryParametersString, // Reverted: Include queryParametersString
		variables.Query.IncludeFiltersInResponse,
	)

	query := parsedBaseURL.Query()
	query.Set("queryId", queryID)
	// query.Set("variables", variablesString) // Old way
	query.Set("includeWebMetadata", "true")
	// parsedBaseURL.RawQuery = query.Encode() // Old way: Encodes the whole variablesString including its parentheses

	// New way: Encode queryId and includeWebMetadata, then append raw variables string
	// This is to prevent URL-encoding of parentheses within the variablesString itself.
	// The cURL seems to pass variables=(...) with literal parentheses.
	encodedBaseQuery := query.Encode() // This will have queryId and includeWebMetadata encoded

	// Now, append the variables part more directly.
	// The variablesString itself should not be additionally URL-encoded if it's meant to be like the cURL.
	// However, the overall query string still needs to be valid.
	// The key "variables" is fine. The value is our variablesString.
	// If query.Encode() was too aggressive, we construct it piece by piece.

	// Ensure variablesString itself has its necessary internal components, but its surrounding parens are literal in the final URL.
	// This means we are treating the whole `(start:0,...false)` as a single value for the `variables` key.
	// The `url.QueryEscape` should be used for the value if it contains special chars that break URL structure (like `&`, `=`, `?`)
	// BUT, the cURL has `&variables=(...)&` - the `=` and `&` are delimiters. The `(...)` is the value.
	// The log showed `variables=%28start...%29`, meaning `query.Encode()` did encode the parens.
	// If the cURL implies those parens should NOT be encoded, then we need to add it raw.

	finalQueryString := encodedBaseQuery + "&variables=" + variablesString // Append raw variables string
	parsedBaseURL.RawQuery = finalQueryString

	return parsedBaseURL.String(), nil
}

// stringSliceToString joins a slice of strings with a separator.
// Helper function for constructing parts of the variables string.
func stringSliceToString(slice []string, sep string) string {
	return strings.Join(slice, sep)
}

// buildProfileGraphQLURL constructs the full URL for a profile GraphQL API request.
// It takes the base URL, query ID, and publicIdentifier, then assembles them.
func buildProfileGraphQLURL(baseURL, queryID, publicIdentifier string) (string, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse base URL: %w", err)
	}

	// For profile fetching, the variables format is simpler:
	// variables=(memberIdentity:{publicIdentifier})
	variablesString := fmt.Sprintf("(memberIdentity:{publicIdentifier:%s})", publicIdentifier)

	query := parsedBaseURL.Query()
	query.Set("queryId", queryID)
	query.Set("includeWebMetadata", "true")

	// Encode the base query parameters
	encodedBaseQuery := query.Encode()

	// Append the variables part with literal parentheses
	finalQueryString := encodedBaseQuery + "&variables=" + variablesString
	parsedBaseURL.RawQuery = finalQueryString

	return parsedBaseURL.String(), nil
}

// GetProfile fetches a detailed LinkedIn profile by public identifier.
func (c *Client) GetProfile(ctx context.Context, publicIdentifier string) (*LinkedInProfile, error) {
	// Input Validation
	if c.config.Auth.LiAtCookie == "" || c.config.Auth.CSRFToken == "" {
		return nil, ErrAuthMissing
	}
	if publicIdentifier == "" {
		return nil, fmt.Errorf("publicIdentifier cannot be empty")
	}

	// Build URL
	requestURL, err := buildProfileGraphQLURL(VoyagerBaseURL, DefaultProfileQueryID, publicIdentifier)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestBuildFailed, err)
	}

	// Prepare Headers
	customHeaders := http.Header{}
	customHeaders.Set("Accept", AcceptHeaderValue)

	// Construct Referer URL for profile requests
	refererURL := fmt.Sprintf("https://www.linkedin.com/in/%s/", publicIdentifier)
	customHeaders.Set("Referer", refererURL)

	// Set X-Li-Page-Instance for profile pages
	xLiPageInstance := fmt.Sprintf("urn:li:page:d_flagship3_profile_view_base;%s", publicIdentifier)
	customHeaders.Set("X-Li-Page-Instance", xLiPageInstance)

	customHeaders.Set("X-Li-Pem-Metadata", "Voyager - Profile")

	// Set X-Li-Track with appropriate context for profile viewing
	xLiTrack := `{"clientVersion":"1.13.35368","mpVersion":"1.13.35368","osName":"web","timezoneOffset":-7,"timezone":"America/Los_Angeles","deviceFormFactor":"DESKTOP","mpName":"voyager-web","displayDensity":2,"displayWidth":1920,"displayHeight":1080}`
	customHeaders.Set("X-Li-Track", xLiTrack)

	// Make API Call
	resp, respBodyBytes, err := c.makeRequest(ctx, http.MethodGet, requestURL, customHeaders, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err)
	}

	// Error Handling (HTTP Status)
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized, http.StatusForbidden:
			return nil, fmt.Errorf("%w: status %d, body: %s", ErrUnauthorized, resp.StatusCode, string(respBodyBytes))
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("%w: status %d, body: %s", ErrRateLimited, resp.StatusCode, string(respBodyBytes))
		default:
			return nil, fmt.Errorf("%w: received status code %d, body: %s", ErrRequestFailed, resp.StatusCode, string(respBodyBytes))
		}
	}

	// Parse JSON Response
	var apiResponse ProfileAPIResponse
	err = json.Unmarshal(respBodyBytes, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: %v. Raw response: %s", ErrResponseParseFailed, err, string(respBodyBytes))
	}

	// Extract Profile from Response using comprehensive parsing
	profile, err := convertAPIResponseToLinkedInProfile(&apiResponse, publicIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to extract profile from response: %w", err)
	}

	return profile, nil
}

// makeRequest executes an HTTP request and returns the response and body bytes.
// It handles adding common headers like CSRF token and li_at cookie.
func (c *Client) makeRequest(ctx context.Context, method string, urlStr string, headers http.Header, body io.Reader) (*http.Response, []byte, error) {
	// log.Printf("[DEBUG] makeRequest (from Echo example context): URL: %s", urlStr) // TEMPORARY LOGGING - REMOVED
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set standard headers that are often required or good to have.
	// The Content-Type for GET requests with GraphQL variables in query params is typically not needed,
	// but if we were sending a POST with a JSON body, it would be "application/json".
	// req.Header.Set("Content-Type", "application/json") // Not for GET

	// Set User-Agent to match the cURL
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("X-Li-Lang", "en_US")
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

	// Add CSRF token and li_at cookie
	req.Header.Set("Csrf-Token", c.config.Auth.CSRFToken)
	req.Header.Set("Cookie", fmt.Sprintf("li_at=%s; JSESSIONID=\"%s\"", c.config.Auth.LiAtCookie, c.config.Auth.JSESSIONID))

	// Add any other headers passed in the headers argument
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Log all request headers before sending
	// log.Println("[DEBUG] makeRequest: All Request Headers:") // TEMPORARY LOGGING - REMOVED
	// for name, headers := range req.Header { // TEMPORARY LOGGING - REMOVED
	// 	for _, h := range headers { // TEMPORARY LOGGING - REMOVED
	// 		log.Printf("[DEBUG] makeRequest Header: %v: %v", name, h) // TEMPORARY LOGGING - REMOVED
	// 	} // TEMPORARY LOGGING - REMOVED
	// } // TEMPORARY LOGGING - REMOVED

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("http client failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	// Check if the server sent gzipped content, even if Go's client is supposed to handle it.
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return resp, nil, fmt.Errorf("failed to create gzip reader for response body: %w", err)
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	respBodyBytes, err := io.ReadAll(reader) // Read from the (potentially decompressed) reader
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, respBodyBytes, nil
}
