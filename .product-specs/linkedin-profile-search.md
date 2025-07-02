**Project: linkedin-scraper Go Package (Phase 1: Profile Search)**

**Goal:** Create a Go package that can search for LinkedIn profiles using the Voyager API, based on keywords and necessary authentication tokens, mimicking the provided cURL request.

---

**Instruction Set:**

**Step 0: Initialize Go Module**

1.  Ensure you are in the root of your `linkedin-scraper` git repository.
2.  Initialize the Go module:
    ```bash
    go mod init github.com/masa-finance/linkedin-scraper
    ```
3.  (You can run `go mod tidy` later after adding dependencies).

**Step 1: Define Core Structures, Constants & Initial Files**

1.  **Create Directory Structure (within `linkedin-scraper`):**
    *   `/` (root directory of the package, already exists as `linkedin-scraper`)
    *   `examples/search_profiles/`
2.  **Create Initial Go Files (all files like `models.go`, `client.go` etc. will be in the root `linkedin-scraper` directory, unless specified otherwise):**
    *   `models.go`
    *   `constants.go`
    *   `config.go`
    *   `client.go`
    *   `search.go`
    *   `errors.go`
    *   `linkedinscraper.go` (main package file, e.g., `package linkedinscraper`, can be minimal initially)
    *   `examples/search_profiles/main.go`
3.  **Constants (`constants.go`):**
    *   Define: `const VoyagerBaseURL = "https://www.linkedin.com/voyager/api/graphql"`
    *   Define: `const DefaultSearchQueryID = "voyagerSearchDashClusters.7cdf88d3366ad02cc5a3862fb9a24085"`
    *   Define default header values from the cURL that are less likely to change frequently (e.g., `AcceptHeader`, `AcceptEncodingHeader`, `DefaultUserAgent`, etc. We can refine these as we go).
        ```go
        package linkedinscraper

        const (
            VoyagerBaseURL         = "https://www.linkedin.com/voyager/api/graphql"
            DefaultSearchQueryID   = "voyagerSearchDashClusters.7cdf88d3366ad02cc5a3862fb9a24085"

            AcceptHeaderValue            = "application/vnd.linkedin.normalized+json+2.1"
            AcceptEncodingHeaderValue    = "gzip, deflate, br, zstd"
            AcceptLanguageHeaderValue  = "en-GB,en-US;q=0.9,en;q=0.8"
            DefaultLiLangHeaderValue     = "en_US"
            DefaultRestliProtocolVersion = "2.0.0"
            DefaultUserAgent             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36" // Example, make configurable
        )
        ```
4.  **Configuration Struct (`config.go`):**
    *   Define: `type AuthCredentials struct { ... }`
        *   Fields:
            *   `LiAtCookie string`
            *   `CSRFToken string`
            *   `JSESSIONID string` // From the cURL example cookie: "ajax:..."
    *   Define: `type Config struct { ... }`
        *   Fields:
            *   `Auth AuthCredentials`
            *   `UserAgent string`
            *   `Referer string` // This will likely need to be dynamic based on the search
            *   `XLiPageInstance string` // From cURL, seems dynamic
            *   `XLiTrack string`      // From cURL, seems dynamic or complex
            // Add other headers from the cURL that might need to be configurable or are dynamic
            // We'll start simple and add more configurability as needed.
5.  **Model Structs (`models.go`):**
    *   Define: `type ProfileSearchArgs struct { ... }`
        *   Fields:
            *   `Keywords string`
            *   `NetworkFilters []string // e.g., ["F", "O"]`
            *   `Start int`
            // Add other potential search parameters if we identify them as useful
    *   Define: `type LinkedInProfile struct { ... }`
        *   Fields: (Initial best guess based on typical search results)
            *   `PublicIdentifier string \`json:"publicIdentifier,omitempty"\`` // e.g., "john-doe-12345"
            *   `URN string \`json:"urn,omitempty"\`` // e.g., "urn:li:fsd_profile:ACoAA..."
            *   `FullName string \`json:"fullName,omitempty"\``
            *   `Headline string \`json:"headline,omitempty"\``
            *   `Location string \`json:"location,omitempty"\``
            *   `ProfileURL string \`json:"profileUrl,omitempty"\`` // Constructed or from API
    *   Define: `type SearchQueryParameters struct { ... }` (for `variables.query.queryParameters`)
        *   Fields:
            *   `Key string \`json:"key"\``
            *   `Value []string \`json:"value"\`` // e.g. value:List(F,O)
    *   Define: `type SearchQuerySubQuery struct { ... }` (for `variables.query`)
        *   Fields:
            *   `Keywords string \`json:"keywords"\``
            *   `FlagshipSearchIntent string \`json:"flagshipSearchIntent"\`` // e.g., "SEARCH_SRP"
            *   `QueryParameters []SearchQueryParameters \`json:"queryParameters"\``
            *   `IncludeFiltersInResponse bool \`json:"includeFiltersInResponse"\``
    *   Define: `type SearchVariables struct { ... }` (for `variables` in the URL)
        *   Fields:
            *   `Start int \`json:"start"\``
            *   `Origin string \`json:"origin"\`` // e.g., "FACETED_SEARCH"
            *   `Query SearchQuerySubQuery \`json:"query"\``
            *   `IncludeFiltersInResponse bool \`json:"includeFiltersInResponse"\`` // This seems redundant with the one in Query, check cURL
    *   Define: `type APIResponseData struct { ... }` (This will be a placeholder for the actual complex structure of the LinkedIn API response. We will need to inspect a real response to fill this out accurately later. For now, it can be an `interface{}` or a very generic struct).
        *   Example placeholder:
            ```go
            package linkedinscraper
            // APIPagingInfo might be part of the response
            // type APIPagingInfo struct {
            //    Start int `json:"start"`
            //    Count int `json:"count"`
            //    Total int `json:"total"`
            // }

            // APIElements will contain the actual data like profiles
            // type APIIncludedElement struct {
            //    // This will be highly dependent on the actual JSON structure
            //    // For example:
            //    // TrackingID string `json:"trackingId"`
            //    // NavigationURL string `json:"navigationUrl"`
            //    // ... fields for name, headline, etc.
            //    // We'll need to map these to our LinkedInProfile struct.
            // }

            type SearchAPIResponse struct {
                // Data interface{} `json:"data"` // Or more specific if we know parts of it
                // Included []APIIncludedElement `json:"included"` // Often profiles are in an "included" array
                // Paging APIPagingInfo `json:"paging"`
            }
            ```
            We will refine `SearchAPIResponse` and how it maps to `[]LinkedInProfile` in Step 4.

**Step 2: Implement Configuration (in `config.go`)**

1.  **`NewConfig` Function:**
    *   Signature: `func NewConfig(auth AuthCredentials, userAgent ...string) (*Config, error)`
    *   Logic:
        *   Validate that `auth.LiAtCookie` and `auth.CSRFToken` are provided. Return an error if not.
        *   Initialize `Config` with provided `auth`.
        *   Set `UserAgent`: if `userAgent` is provided, use it; otherwise, use `constants.DefaultUserAgent`.
        *   (We will add more parameters like `Referer`, `XLiPageInstance`, `XLiTrack` later as they might be dynamic or require more thought on how they're set).
        *   Return `&Config{...}, nil`.

**Step 3: Implement Basic HTTP Client (in `client.go`)**

1.  **`Client` Struct:**
    *   Define: `type Client struct { ... }`
    *   Fields:
        *   `httpClient *http.Client`
        *   `config *Config`
2.  **`NewClient` Function:**
    *   Signature: `func NewClient(cfg *Config) (*Client, error)`
    *   Logic:
        *   If `cfg` is nil, return an error.
        *   Initialize a standard `httpClient` (e.g., `&http.Client{Timeout: 30 * time.Second}`). Ensure it can handle gzip automatically (Go's default `http.Transport` does).
        *   Return `&Client{httpClient: httpClient, config: cfg}, nil`.
3.  **Helper Function `buildGraphQLURL` (in `client.go` or `search.go`):**
    *   Signature: `func buildGraphQLURL(baseURL string, queryID string, variables SearchVariables) (string, error)`
    *   Logic:
        *   Marshal `variables` to a JSON string.
        *   URL encode the JSON string.
        *   Construct the full URL: `baseURL + "?includeWebMetadata=true&variables=(" + encodedVariables + ")&queryId=" + queryID`.
        *   Return the URL string and any error from marshalling/encoding.
4.  **Request Execution Method (`makeRequest`) (in `client.go`):**
    *   Signature: `func (c *Client) makeRequest(ctx context.Context, method string, urlStr string, headers http.Header, body io.Reader) (*http.Response, []byte, error)`
    *   Logic:
        *   Create `req, err := http.NewRequestWithContext(ctx, method, urlStr, body)`.
        *   Set common headers from `constants` and `c.config`:
            *   `Accept: constants.AcceptHeaderValue`
            *   `Accept-Encoding: constants.AcceptEncodingHeaderValue`
            *   `Accept-Language: constants.AcceptLanguageHeaderValue`
            *   `Csrf-Token: c.config.Auth.CSRFToken`
            *   `X-Li-Lang: constants.DefaultLiLangHeaderValue`
            *   `X-Restli-Protocol-Version: constants.DefaultRestliProtocolVersion`
            *   `User-Agent: c.config.UserAgent`
            *   Cookie: `fmt.Sprintf("li_at=%s; JSESSIONID=\"%s\"", c.config.Auth.LiAtCookie, c.config.Auth.JSESSIONID)`
            *   Add any other headers passed in the `headers` argument (this allows per-request overrides or additions like `Referer`, `X-Li-Page-Instance`, `X-Li-Track`).
        *   `resp, err := c.httpClient.Do(req)`.
        *   If error, return.
        *   `defer resp.Body.Close()`.
        *   Read `respBodyBytes, err := io.ReadAll(resp.Body)`.
        *   Return `resp, respBodyBytes, err`.

**Step 4: Implement Profile Search Logic (in `search.go`)**

1.  **`SearchProfiles` Function:**
    *   Signature: `func (c *Client) SearchProfiles(ctx context.Context, args ProfileSearchArgs) ([]LinkedInProfile, error)`
    *   Logic:
        *   **Input Validation:** Check if `c.config.Auth.LiAtCookie` and `c.config.Auth.CSRFToken` are set. Return `ErrAuthMissing` (define in `errors.go`) if not.
        *   Check if `args.Keywords` is empty. Return an error if so.
        *   **Construct `SearchVariables`:**
            *   Populate `SearchVariables` struct using `args.Keywords`, `args.NetworkFilters`, `args.Start`, and default values for `Origin`, `FlagshipSearchIntent`, etc., based on the cURL.
            ```go
            // Example:
            querySubQuery := SearchQuerySubQuery{
                Keywords:             args.Keywords,
                FlagshipSearchIntent: "SEARCH_SRP", // from cURL
                QueryParameters:      []SearchQueryParameters{},
                IncludeFiltersInResponse: false,
            }
            if len(args.NetworkFilters) > 0 {
                querySubQuery.QueryParameters = append(querySubQuery.QueryParameters, SearchQueryParameters{
                    Key:   "network",
                    Value: args.NetworkFilters, // e.g. List(F,O)
                })
            }
             // Add other fixed queryParameters from cURL like (key:resultType,value:List(PEOPLE))
            querySubQuery.QueryParameters = append(querySubQuery.QueryParameters, SearchQueryParameters{
                Key:   "resultType",
                Value: []string{"PEOPLE"},
            })


            variables := SearchVariables{
                Start:  args.Start,
                Origin: "FACETED_SEARCH", // from cURL
                Query:  querySubQuery,
                IncludeFiltersInResponse: false, // from cURL, seems to be in two places
            }
            ```
        *   **Build URL:** Call `buildGraphQLURL(constants.VoyagerBaseURL, constants.DefaultSearchQueryID, variables)`. Handle error.
        *   **Prepare Headers:**
            *   `customHeaders := http.Header{}`
            *   `customHeaders.Set("Referer", "https://www.linkedin.com/search/results/people/?keywords="+url.QueryEscape(args.Keywords)) // Example, needs to be more robust
            *   `customHeaders.Set("X-Li-Page-Instance", "urn:li:page:d_flagship3_search_srp_people;xxxx")` // This needs a way to be set or generated. For now, can be a placeholder or taken from config if made configurable.
            *   `customHeaders.Set("X-Li-Pem-Metadata", "Voyager - People SRP=search-results")`
            *   `customHeaders.Set("X-Li-Track", `{"clientVersion":"1.13.x","mpVersion":"1.13.x",...}`)` // This is complex. For initial phase, we might use a static (but valid-looking) one or make it fully configurable.
        *   **Make API Call:** Call `resp, respBodyBytes, err := c.makeRequest(ctx, http.MethodGet, requestURL, customHeaders, nil)`.
        *   **Error Handling (HTTP Status):**
            *   If `err != nil`, return it.
            *   If `resp.StatusCode != http.StatusOK`:
                *   Handle common errors (401, 403 -> `ErrUnauthorized`; 429 -> `ErrRateLimited`; etc. defined in `errors.go`).
                *   Return a generic error with status code and body.
        *   **Parse JSON Response:**
            *   `var apiResponse SearchAPIResponse` // This is the complex part
            *   `err = json.Unmarshal(respBodyBytes, &apiResponse)`
            *   If `err != nil`, return `ErrResponseParseFailed` wrapping the original error.
            *   **Extract Profiles:** Iterate through `apiResponse` (e.g., `apiResponse.Included` or `apiResponse.Data.elements` - this depends heavily on the actual response structure) and map the fields to `LinkedInProfile` structs. This is where we'll need to carefully inspect a real JSON response.
            *   *Self-correction/Refinement:* The cURL uses `application/vnd.linkedin.normalized+json+2.1`. The structure often involves a top-level `data` object, and an `included` array where entities like profiles are fully described. The `data` object might contain references (URNs) to these included elements. We will need to navigate this.
        *   Return the slice of `LinkedInProfile` and `nil` error.

**Step 5: Define Custom Errors (in `errors.go`)**

1.  **Error Variables (ensure `package linkedinscraper`):**
    *   `var ErrAuthMissing = errors.New("linkedinscraper: authentication credentials (li_at, csrf_token) are missing")`
    *   `var ErrKeywordsMissing = errors.New("linkedinscraper: search keywords are missing")`
    *   `var ErrRequestBuildFailed = errors.New("linkedinscraper: failed to build API request")`
    *   `var ErrRequestFailed = errors.New("linkedinscraper: API request failed")` // Generic for HTTP issues
    *   `var ErrUnauthorized = errors.New("linkedinscraper: unauthorized, check credentials or IP reputation")`
    *   `var ErrRateLimited = errors.New("linkedinscraper: rate limited by API")`
    *   `var ErrResponseParseFailed = errors.New("linkedinscraper: failed to parse API response")`
    *   `var ErrNoProfilesFound = errors.New("linkedinscraper: no profiles found matching criteria")` // Or handle this by returning empty slice

**Step 6: Example Usage (in `examples/search_profiles/main.go`)**

1.  **`main` Function:**
    *   Get `LI_AT_COOKIE`, `CSRF_TOKEN`, `JSESSIONID_TOKEN` from environment variables or hardcode for testing.
    *   Import the scraper: `import linkedinscraper "github.com/masa-finance/linkedin-scraper"`
    *   Create `auth := linkedinscraper.AuthCredentials{...}`.
    *   `cfg, err := linkedinscraper.NewConfig(auth)`
    *   `client, err := linkedinscraper.NewClient(cfg)`
    *   `profiles, err := client.SearchProfiles(context.Background(), linkedinscraper.ProfileSearchArgs{Keywords: "investor", NetworkFilters: []string{"F", "O"}})`
    *   Handle `err`.
    *   Print `profiles`.

**Step 7: README (in `README.md` at the project root)**

1.  Create a basic `README.md` with:
    *   Project title: `linkedin-scraper`
    *   Brief description.
    *   "How to get `li_at` and `csrf-token`" section (guiding user to browser dev tools).
    *   Installation: `go get github.com/masa-finance/linkedin-scraper`
    *   Simple usage example from Step 6.

**Step 8: Testing (Iterative - `*_test.go` files)**

1.  For each `*.go` file in the `linkedinscraper` package, create a corresponding `*_test.go` (e.g., `config_test.go`, `client_test.go`).
2.  **`config_test.go`:** Test `NewConfig` validation and default settings.
3.  **`client_test.go`:**
    *   Test `buildGraphQLURL`.
    *   Test `makeRequest` by mocking `http.Client` (e.g., using `httptest.NewServer` or an interface for the HTTP client). Verify headers, URL, method. Test error handling.
4.  **`search_test.go`:**
    *   Mock the `makeRequest` method of the `Client` or the underlying HTTP transport.
    *   Test `SearchProfiles`:
        *   Happy path: Provide a mock successful JSON response and verify profiles are parsed correctly.
        *   Error paths: Test missing auth, missing keywords, API errors (mocked 401, 429, 500), malformed JSON response.
5.  **JSON Parsing Tests:** Create specific tests with sample JSON snippets (once available) to ensure the unmarshalling into `SearchAPIResponse` and mapping to `LinkedInProfile` works as expected.

---

**Collaboration Plan:**

1.  **I (AI) will start by generating the directory structure and the initial Go files with the struct definitions, constants, and function signatures outlined in Step 1, Step 2 (NewConfig), and Step 3 (Client struct, NewClient, buildGraphQLURL, makeRequest signatures).** I'll make placeholders for complex parts like the full `SearchAPIResponse` structure. All Go files in the main package will start with `package linkedinscraper`.
2.  You (User) will review these initial files. You can help provide a sample JSON response from a successful cURL call if possible, which will greatly help in defining `SearchAPIResponse` accurately in `models.go`.
3.  Next, I'll implement the logic for Step 2 (`NewConfig`) and Step 3 (`NewClient`, `buildGraphQLURL`, and the core request sending/header logic in `makeRequest`).
4.  You'll review.
5.  Then, we'll collaboratively tackle Step 4 (`SearchProfiles` in `search.go`).
    *   I will draft the function, including variable construction, URL building, and the call to `makeRequest`.
    *   **Crucial part:** We will need to work together on the JSON parsing. Once you provide a sample JSON response (or we deduce its structure from the cURL and common LinkedIn API patterns), I will attempt to define `SearchAPIResponse` (and any nested structs) in `models.go` and write the parsing logic in `SearchProfiles` to map the API data to our `[]LinkedInProfile`. This might be iterative.
6.  After that, I will implement Step 5 (`errors.go`).
7.  Then, I will draft the example usage (Step 6) and a basic README (Step 7).
8.  Finally, we will iteratively work on tests (Step 8), starting with unit tests for config, client, and then the more complex search logic, including tests for JSON parsing once the structure is clearer.