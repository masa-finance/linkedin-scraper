package linkedinscraper

import "errors"

var ErrAuthMissing = errors.New("linkedinscraper: authentication credentials (li_at, csrf_token) are missing")

var (
	ErrKeywordsMissing     = errors.New("linkedinscraper: search keywords are missing")
	ErrRequestBuildFailed  = errors.New("linkedinscraper: failed to build API request")
	ErrRequestFailed       = errors.New("linkedinscraper: API request failed") // Generic for HTTP issues
	ErrUnauthorized        = errors.New("linkedinscraper: unauthorized, check credentials or IP reputation")
	ErrRateLimited         = errors.New("linkedinscraper: rate limited by API")
	ErrResponseParseFailed = errors.New("linkedinscraper: failed to parse API response")
	ErrNoProfilesFound     = errors.New("linkedinscraper: no profiles found matching criteria") // Or handle this by returning empty slice
)
