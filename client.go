package linkedinscraper

import (
	"net/http"
)

// Client is the LinkedIn API client.
type Client struct {
	httpClient *http.Client
	config     *Config
}
