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

// NewConfig creates a new Config struct.
func NewConfig(auth AuthCredentials, userAgent ...string) (*Config, error) {
	if auth.LiAtCookie == "" || auth.CSRFToken == "" {
		return nil, ErrAuthMissing
	}

	cfg := &Config{
		Auth: auth,
	}

	if len(userAgent) > 0 && userAgent[0] != "" {
		cfg.UserAgent = userAgent[0]
	} else {
		cfg.UserAgent = DefaultUserAgent // Assumes DefaultUserAgent is defined in constants.go
	}

	// We will add more parameters like Referer, XLiPageInstance, XLiTrack later
	// as they might be dynamic or require more thought on how they're set.

	return cfg, nil
}
