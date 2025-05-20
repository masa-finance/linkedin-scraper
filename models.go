package linkedinscraper

// ProfileSearchArgs represents the arguments for initiating a profile search.
type ProfileSearchArgs struct {
	Keywords       string
	NetworkFilters []string // e.g., ["F", "O"] for 1st degree and Outside network
	Start          int
	Count          int // Added based on typical pagination and cURL example
	// Add other potential search parameters here if identified.
	// Origin string // e.g., "FACETED_SEARCH", also a potential parameter
	XLiPageInstance string // Optional: To override default placeholder
	XLiTrack        string // Optional: To override default placeholder
}

// LinkedInProfile represents the extracted information for a single LinkedIn profile.
type LinkedInProfile struct {
	PublicIdentifier string `json:"publicIdentifier,omitempty"` // e.g., "nic-sanchez-a8516a54"
	URN              string `json:"urn,omitempty"`              // e.g., "urn:li:fsd_profile:ACoAAAtp-4UBpQ0aZ_PeToflBoLty9BpO_CQ6-I"
	FullName         string `json:"fullName,omitempty"`         // e.g., "Nic Sanchez"
	Headline         string `json:"headline,omitempty"`         // e.g., "Investor at Bertram Capital"
	Location         string `json:"location,omitempty"`         // e.g., "San Francisco, CA"
	ProfileURL       string `json:"profileUrl,omitempty"`       // e.g., "https://www.linkedin.com/in/nic-sanchez-a8516a54?..."
	// Degree string `json:"degree,omitempty"` // e.g. "• 2nd", could be parsed from badgeText
}

// SearchQueryParameters represents a single key-value pair for query parameters
// within the search query.
type SearchQueryParameters struct {
	Key   string   `json:"key"`   // e.g., "network" or "resultType"
	Value []string `json:"value"` // e.g., ["F", "O"] or ["PEOPLE"]
}

// SearchQuerySubQuery represents the 'query' object within the search variables.
type SearchQuerySubQuery struct {
	Keywords                 string                  `json:"keywords"`             // e.g., "investor"
	FlagshipSearchIntent     string                  `json:"flagshipSearchIntent"` // e.g., "SEARCH_SRP"
	QueryParameters          []SearchQueryParameters `json:"queryParameters"`
	IncludeFiltersInResponse bool                    `json:"includeFiltersInResponse"` // As per cURL, this is within the query object
}

// SearchVariables represents the 'variables' object passed in the API request URL.
// This struct is designed to be marshalled into the format LinkedIn expects for its GraphQL variables.
type SearchVariables struct {
	Start  int                 `json:"start"`  // e.g., 0
	Count  int                 `json:"count"`  // e.g., 1 (present in cURL, added here)
	Origin string              `json:"origin"` // e.g., "FACETED_SEARCH"
	Query  SearchQuerySubQuery `json:"query"`
	// The 'includeFiltersInResponse' field in the user's original spec for SearchVariables
	// is not present at this top level in the provided cURL. It's inside the 'Query' sub-object.
	// If it were needed at this level, it would be:
	// IncludeFiltersInResponse bool `json:"includeFiltersInResponse,omitempty"`
}

// --- API Response Structures (to be refined in Step 4 as per your plan) ---

// TextObject is a common structure in LinkedIn's API for text fields.
type TextObject struct {
	Text string `json:"text"`
}

// APIPagingInfo holds pagination data from the API response.
type APIPagingInfo struct {
	Start int `json:"start"`
	Count int `json:"count"`
	Total int `json:"total"`
}

// ClusterMetadata holds metadata about the search cluster.
type ClusterMetadata struct {
	TotalResultCount int `json:"totalResultCount"`
	// Other metadata fields can be added here
}

// Item represents an individual item within a search result cluster.
// It often contains a URN pointing to more detailed data in the "included" section.
type Item struct {
	EntityResultURN string `json:"*entityResult"` // URN for EntityResultViewModel
	FeedbackCardURN string `json:"*feedbackCard"` // URN for FeedbackCard
	// Other types of URNs or direct data might appear here
}

// ClusterElement represents a cluster of search results.
type ClusterElement struct {
	Items    []Item      `json:"items"`
	Position int         `json:"position"`
	Image    *string     `json:"image"` // Using pointer for nullable
	Title    *TextObject `json:"title"` // Using pointer for nullable text object
	// Other cluster fields can be added here
}

// SearchDashClusters holds the core search results.
type SearchDashClusters struct {
	Metadata ClusterMetadata  `json:"metadata"`
	Paging   APIPagingInfo    `json:"paging"`
	Elements []ClusterElement `json:"elements"`
}

// InnerData is part of the nested "data" structure in the API response.
type InnerData struct {
	SearchDashClustersByAll SearchDashClusters `json:"searchDashClustersByAll"`
}

// RootData holds the top-level "data" object.
type RootData struct {
	InnerData InnerData `json:"data"`
}

// IncludedEntityResultViewModel represents the 'EntityResultViewModel' type found in the "included" array.
// This is a key structure for populating LinkedInProfile.
type IncludedEntityResultViewModel struct {
	EntityURN         string     `json:"entityUrn"`         // This is the URN of the ViewModel itself
	TrackingURN       string     `json:"trackingUrn"`       // Often the URN of the underlying profile/member
	Title             TextObject `json:"title"`             // Maps to FullName
	PrimarySubtitle   TextObject `json:"primarySubtitle"`   // Maps to Headline
	SecondarySubtitle TextObject `json:"secondarySubtitle"` // Maps to Location
	NavigationURL     string     `json:"navigationUrl"`     // Maps to ProfileURL
	BadgeText         TextObject `json:"badgeText"`         // e.g., "• 2nd" for connection degree
	// PublicIdentifier might be derived or sometimes present
}

// IncludedProfile represents the 'Profile' type found in the "included" array.
type IncludedProfile struct {
	EntityURN        string `json:"entityUrn"` // This is the Profile URN
	PublicIdentifier string `json:"publicIdentifier,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Headline         string `json:"headline,omitempty"`
	// Other profile-specific fields
}

// IncludedFeedbackCard represents a feedback card from the "included" array.
type IncludedFeedbackCard struct {
	EntityURN  string `json:"entityUrn"`
	TrackingId string `json:"trackingId"`
}

// IncludedLazyLoadedActions represents lazy loaded actions from the "included" array.
type IncludedLazyLoadedActions struct {
	EntityURN string `json:"entityUrn"`
}

// GenericIncludedElement is used to unmarshal items from the "included" array
// and determine their specific type using the "$type" field.
type GenericIncludedElement struct {
	Type string `json:"$type"` // e.g., "com.linkedin.voyager.dash.search.EntityResultViewModel" or "com.linkedin.voyager.dash.identity.profile.Profile"
	// Embed other fields that are common or use json.RawMessage to unmarshal specific data later
	// For simplicity, we'll assume specific unmarshalling based on $type happens after this stage.
	// The fields below are from EntityResultViewModel for direct unmarshalling if $type matches.
	EntityURN         string      `json:"entityUrn,omitempty"`
	TrackingURN       string      `json:"trackingUrn,omitempty"`
	Title             *TextObject `json:"title,omitempty"`
	PrimarySubtitle   *TextObject `json:"primarySubtitle,omitempty"`
	SecondarySubtitle *TextObject `json:"secondarySubtitle,omitempty"`
	NavigationURL     string      `json:"navigationUrl,omitempty"`
	BadgeText         *TextObject `json:"badgeText,omitempty"`

	// Fields from Profile type
	PublicIdentifier string `json:"publicIdentifier,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Headline         string `json:"headline,omitempty"` // Note: Profile also has a headline

	// Fields from FeedbackCard
	TrackingId string `json:"trackingId,omitempty"`
}

// SearchAPIResponse is the top-level structure for the entire API JSON response.
// It will be refined further in Step 4.
type SearchAPIResponse struct {
	RootData RootData                 `json:"data"`
	Included []GenericIncludedElement `json:"included"` // This will hold various types of objects
	// Meta interface{} `json:"meta"` // The meta field contains microSchema, can be added if needed
	// Extensions interface{} `json:"extensions"` // The extensions field, can be added if needed
}
