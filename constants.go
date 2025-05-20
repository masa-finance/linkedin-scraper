package linkedinscraper

const (
	VoyagerBaseURL       = "https://www.linkedin.com/voyager/api/graphql"
	DefaultSearchQueryID = "voyagerSearchDashClusters.7cdf88d3366ad02cc5a3862fb9a24085"

	AcceptHeaderValue            = "application/vnd.linkedin.normalized+json+2.1"
	AcceptEncodingHeaderValue    = "gzip, deflate, br, zstd"
	AcceptLanguageHeaderValue    = "en-GB,en-US;q=0.9,en;q=0.8"
	DefaultLiLangHeaderValue     = "en_US"
	DefaultRestliProtocolVersion = "2.0.0"
	DefaultUserAgent             = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36" // Example, make configurable
)
