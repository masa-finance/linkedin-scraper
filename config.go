package linkedinscraper

// AuthCredentials holds the necessary authentication tokens.
type AuthCredentials struct {
	LiAtCookie string
	CSRFToken  string
	JSESSIONID string // From the cURL example cookie: "ajax:..."
}

// Config holds the configuration for the LinkedIn client.
type Config struct {
	Auth            AuthCredentials
	UserAgent       string
	Referer         string // This will likely need to be dynamic based on the search
	XLiPageInstance string // From cURL, seems dynamic
	XLiTrack        string // From cURL, seems dynamic or complex
	// Add other headers from the cURL that might need to be configurable or are dynamic
	// We'll start simple and add more configurability as needed.
}
