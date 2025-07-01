package linkedinscraper

const (
	VoyagerBaseURL = "https://www.linkedin.com/voyager/api/graphql"
	// DefaultSearchQueryID is the default query ID for profile searches.
	// This was taken from a cURL command observation.
	// Example: voyagerSearchDashClusters.b1d223dcc11b2a052b967900e7388211
	// Updated based on the latest working cURL command provided by the user.
	DefaultSearchQueryID = "voyagerSearchDashClusters.7cdf88d3366ad02cc5a3862fb9a24085"

	// DefaultProfileQueryID is the default query ID for profile fetching.
	// This is used with the voyagerIdentityDashProfiles query to fetch detailed profile data.
	DefaultProfileQueryID = "voyagerIdentityDashProfiles.9d37763c7ad41b36d4fb077d1e9e6ee2"

	AcceptHeaderValue            = "application/vnd.linkedin.normalized+json+2.1"
	AcceptEncodingHeaderValue    = "gzip, deflate, br, zstd"
	AcceptLanguageHeaderValue    = "en-GB,en-US;q=0.9,en;q=0.8"
	DefaultLiLangHeaderValue     = "en_US"
	DefaultRestliProtocolVersion = "2.0.0"
	// DefaultUserAgent is the default user agent for Voyager API calls
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36"
)
