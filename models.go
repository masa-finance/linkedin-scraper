package linkedinscraper

import (
	"encoding/json"
	"fmt"
)

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

// Date represents a LinkedIn date structure
type Date struct {
	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

// DateRange represents a LinkedIn date range
type DateRange struct {
	Start *Date `json:"start,omitempty"`
	End   *Date `json:"end,omitempty"`
}

// Experience represents a work experience/position entry
type Experience struct {
	EntityURN              string              `json:"entityUrn,omitempty"`
	CompanyName            string              `json:"companyName,omitempty"`
	CompanyURN             string              `json:"companyUrn,omitempty"`
	Title                  string              `json:"title,omitempty"`
	Description            string              `json:"description,omitempty"`
	DateRange              *DateRange          `json:"dateRange,omitempty"`
	LocationName           string              `json:"locationName,omitempty"`
	MultiLocaleCompanyName []map[string]string `json:"multiLocaleCompanyName,omitempty"`
}

// Education represents an education entry
type Education struct {
	EntityURN    string     `json:"entityUrn,omitempty"`
	SchoolName   string     `json:"schoolName,omitempty"`
	SchoolURN    string     `json:"schoolUrn,omitempty"`
	DegreeName   string     `json:"degreeName,omitempty"`
	FieldOfStudy string     `json:"fieldOfStudy,omitempty"`
	DateRange    *DateRange `json:"dateRange,omitempty"`
	Description  string     `json:"description,omitempty"`
	Activities   string     `json:"activities,omitempty"`
}

// Skill represents a skill entry
type Skill struct {
	EntityURN        string `json:"entityUrn,omitempty"`
	Name             string `json:"name,omitempty"`
	EndorsementCount int    `json:"endorsementCount,omitempty"`
	EndorsedByViewer bool   `json:"endorsedByViewer,omitempty"`
}

// Certification represents a certification entry
type Certification struct {
	EntityURN     string     `json:"entityUrn,omitempty"`
	Name          string     `json:"name,omitempty"`
	Authority     string     `json:"authority,omitempty"`
	DateRange     *DateRange `json:"dateRange,omitempty"`
	LicenseNumber string     `json:"licenseNumber,omitempty"`
	URL           string     `json:"url,omitempty"`
}

// ProfileLocation represents detailed location information
type ProfileLocation struct {
	CountryCode       string `json:"countryCode,omitempty"`
	PostalCode        string `json:"postalCode,omitempty"`
	PreferredGeoPlace string `json:"preferredGeoPlace,omitempty"`
}

// ProfilePicture represents profile picture information
type ProfilePicture struct {
	DisplayImageUrn    string `json:"displayImageUrn,omitempty"`
	PhotoFilterPicture string `json:"photoFilterPicture,omitempty"`
	RootURL            string `json:"rootUrl,omitempty"`
	A11yText           string `json:"a11yText,omitempty"`
}

// ConnectionInfo represents connection and following information
type ConnectionInfo struct {
	ConnectionCount int  `json:"connectionCount,omitempty"`
	FollowerCount   int  `json:"followerCount,omitempty"`
	FollowingCount  int  `json:"followingCount,omitempty"`
	Following       bool `json:"following,omitempty"`
}

// LinkedInProfile represents the extracted information for a single LinkedIn profile.
// Extended to support both search results and detailed profile data.
type LinkedInProfile struct {
	// Basic fields (existing - for backward compatibility)
	PublicIdentifier string `json:"publicIdentifier,omitempty"` // e.g., "nic-sanchez-a8516a54"
	URN              string `json:"urn,omitempty"`              // e.g., "urn:li:fsd_profile:ACoAAAtp-4UBpQ0aZ_PeToflBoLty9BpO_CQ6-I"
	FullName         string `json:"fullName,omitempty"`         // e.g., "Nic Sanchez"
	Headline         string `json:"headline,omitempty"`         // e.g., "Investor at Bertram Capital"
	Location         string `json:"location,omitempty"`         // e.g., "San Francisco, CA"
	ProfileURL       string `json:"profileUrl,omitempty"`       // e.g., "https://www.linkedin.com/in/nic-sanchez-a8516a54?..."

	// Extended fields for detailed profile data
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Summary   string `json:"summary,omitempty"`
	Industry  string `json:"industry,omitempty"`

	// Location details
	LocationDetails *ProfileLocation `json:"locationDetails,omitempty"`

	// Professional information
	Experience     []Experience    `json:"experience,omitempty"`
	Education      []Education     `json:"education,omitempty"`
	Skills         []Skill         `json:"skills,omitempty"`
	Certifications []Certification `json:"certifications,omitempty"`

	// Profile media and presentation
	ProfilePicture     *ProfilePicture `json:"profilePicture,omitempty"`
	BackgroundImageURL string          `json:"backgroundImageUrl,omitempty"`

	// Social and verification info
	ConnectionInfo *ConnectionInfo `json:"connectionInfo,omitempty"`
	IsVerified     bool            `json:"isVerified,omitempty"`
	IsCreator      bool            `json:"isCreator,omitempty"`
	IsPremium      bool            `json:"isPremium,omitempty"`

	// Additional metadata
	IsMemorialized  bool   `json:"isMemorialized,omitempty"`
	TempStatus      string `json:"tempStatus,omitempty"`
	TempStatusEmoji string `json:"tempStatusEmoji,omitempty"`

	// Activity and engagement
	CreatorWebsite string `json:"creatorWebsite,omitempty"`

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

// FlexibleText is a custom type to handle fields that can be either a string or a TextObject
type FlexibleText string

// UnmarshalJSON implements custom unmarshaling logic for FlexibleText.
// It tries to unmarshal into a TextObject first, and falls back to a string.
func (ft *FlexibleText) UnmarshalJSON(data []byte) error {
	// 1. Try to unmarshal into a standard TextObject
	var textObj struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(data, &textObj); err == nil && textObj.Text != "" {
		*ft = FlexibleText(textObj.Text)
		return nil
	}

	// 2. Try to unmarshal into a simple string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*ft = FlexibleText(s)
		return nil
	}

	// 3. If it's a JSON null, treat it as an empty string
	if string(data) == "null" {
		*ft = ""
		return nil
	}

	return fmt.Errorf("cannot unmarshal %s into FlexibleText", string(data))
}

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
	Items    []Item        `json:"items"`
	Position int           `json:"position"`
	Image    *string       `json:"image"` // Using pointer for nullable
	Title    *FlexibleText `json:"title"` // Using pointer for nullable text object
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
	EntityURN         string       `json:"entityUrn"`         // This is the URN of the ViewModel itself
	TrackingURN       string       `json:"trackingUrn"`       // Often the URN of the underlying profile/member
	Title             FlexibleText `json:"title"`             // Maps to FullName
	PrimarySubtitle   FlexibleText `json:"primarySubtitle"`   // Maps to Headline
	SecondarySubtitle FlexibleText `json:"secondarySubtitle"` // Maps to Location
	NavigationURL     string       `json:"navigationUrl"`     // Maps to ProfileURL
	BadgeText         FlexibleText `json:"badgeText"`         // e.g., "• 2nd" for connection degree
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
	EntityURN         string        `json:"entityUrn,omitempty"`
	TrackingURN       string        `json:"trackingUrn,omitempty"`
	Title             *FlexibleText `json:"title,omitempty"`
	PrimarySubtitle   *FlexibleText `json:"primarySubtitle,omitempty"`
	SecondarySubtitle *FlexibleText `json:"secondarySubtitle,omitempty"`
	NavigationURL     string        `json:"navigationUrl,omitempty"`
	BadgeText         *FlexibleText `json:"badgeText,omitempty"`

	// Fields from Profile type
	PublicIdentifier string `json:"publicIdentifier,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Headline         string `json:"headline,omitempty"` // Note: Profile also has a headline

	// Fields from PositionResponse
	CompanyName  string             `json:"companyName,omitempty"`
	CompanyURN   string             `json:"*company,omitempty"`
	Description  string             `json:"description,omitempty"`
	DateRange    *DateRangeResponse `json:"dateRange,omitempty"`
	LocationName string             `json:"locationName,omitempty"`

	// Fields from EducationResponse
	SchoolName   string `json:"schoolName,omitempty"`
	SchoolURN    string `json:"*school,omitempty"`
	DegreeName   string `json:"degreeName,omitempty"`
	FieldOfStudy string `json:"fieldOfStudy,omitempty"`
	Activities   string `json:"activities,omitempty"`

	// Fields from Skill
	Name             string `json:"name,omitempty"`
	EndorsementCount int    `json:"endorsementCount,omitempty"`
	EndorsedByViewer bool   `json:"endorsedByViewer,omitempty"`

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

// --- Profile API Response Structures ---

// ProfileAPIResponse represents the top-level response from LinkedIn's profile API
type ProfileAPIResponse struct {
	Data     ProfileData              `json:"data"`
	Included []GenericIncludedElement `json:"included,omitempty"`
	Meta     interface{}              `json:"meta,omitempty"`
}

// ProfileData represents the data section of the profile API response
type ProfileData struct {
	Data ProfileInnerData `json:"data"`
}

// ProfileInnerData represents the inner data structure with profile collections
type ProfileInnerData struct {
	RecipeTypes                          []string                       `json:"$recipeTypes,omitempty"`
	IdentityDashProfilesByMemberIdentity IdentityDashProfilesCollection `json:"identityDashProfilesByMemberIdentity"`
	Type                                 string                         `json:"$type,omitempty"`
}

// IdentityDashProfilesCollection represents the profile collection response
type IdentityDashProfilesCollection struct {
	Elements    []string `json:"*elements,omitempty"`
	RecipeTypes []string `json:"$recipeTypes,omitempty"`
	Type        string   `json:"$type,omitempty"`
}

// ProfileResponseEntity represents a detailed profile entity from the included array
type ProfileResponseEntity struct {
	EntityURN           string                    `json:"entityUrn,omitempty"`
	FirstName           string                    `json:"firstName,omitempty"`
	LastName            string                    `json:"lastName,omitempty"`
	Headline            string                    `json:"headline,omitempty"`
	PublicIdentifier    string                    `json:"publicIdentifier,omitempty"`
	Location            *ProfileLocationResponse  `json:"location,omitempty"`
	ProfilePicture      *ProfilePictureResponse   `json:"profilePicture,omitempty"`
	ConnectionInfo      *ConnectionInfoResponse   `json:"connections,omitempty"`
	FollowingState      *FollowingStateResponse   `json:"followingState,omitempty"`
	ProfileTopPosition  *PositionsCollection      `json:"profileTopPosition,omitempty"`
	ProfileTopEducation *EducationCollection      `json:"profileTopEducation,omitempty"`
	VerificationData    *VerificationDataResponse `json:"verificationData,omitempty"`
	Creator             bool                      `json:"creator,omitempty"`
	CreatorInfo         *CreatorInfoResponse      `json:"creatorInfo,omitempty"`
	Memorialized        bool                      `json:"memorialized,omitempty"`
	TempStatus          string                    `json:"tempStatus,omitempty"`
	TempStatusEmoji     string                    `json:"tempStatusEmoji,omitempty"`
	CreatorWebsite      *TextViewModelResponse    `json:"creatorWebsite,omitempty"`
	RecipeTypes         []string                  `json:"$recipeTypes,omitempty"`
	Type                string                    `json:"$type,omitempty"`
}

// ProfileLocationResponse represents location data from API response
type ProfileLocationResponse struct {
	CountryCode       string   `json:"countryCode,omitempty"`
	PostalCode        string   `json:"postalCode,omitempty"`
	PreferredGeoPlace string   `json:"preferredGeoPlace,omitempty"`
	RecipeTypes       []string `json:"$recipeTypes,omitempty"`
	Type              string   `json:"$type,omitempty"`
}

// ProfilePictureResponse represents profile picture data from API response
type ProfilePictureResponse struct {
	DisplayImageUrn                string               `json:"displayImageUrn,omitempty"`
	DisplayImageReference          *VectorImageResponse `json:"displayImageReference,omitempty"`
	DisplayImageWithFrameReference *VectorImageResponse `json:"displayImageWithFrameReference,omitempty"`
	A11yText                       string               `json:"a11yText,omitempty"`
	FrameType                      string               `json:"frameType,omitempty"`
	IsGeneratedOrModifiedByAi      bool                 `json:"isGeneratedOrModifiedByAi,omitempty"`
	RecipeTypes                    []string             `json:"$recipeTypes,omitempty"`
	Type                           string               `json:"$type,omitempty"`
}

// VectorImageResponse represents vector image data from API response
type VectorImageResponse struct {
	RootURL           string                   `json:"rootUrl,omitempty"`
	Artifacts         []VectorArtifactResponse `json:"artifacts,omitempty"`
	DigitalMediaAsset string                   `json:"digitalmediaAsset,omitempty"`
	Attribution       string                   `json:"attribution,omitempty"`
	RecipeTypes       []string                 `json:"$recipeTypes,omitempty"`
	Type              string                   `json:"$type,omitempty"`
}

// VectorArtifactResponse represents individual image artifacts
type VectorArtifactResponse struct {
	Width                         int      `json:"width,omitempty"`
	Height                        int      `json:"height,omitempty"`
	FileIdentifyingUrlPathSegment string   `json:"fileIdentifyingUrlPathSegment,omitempty"`
	ExpiresAt                     int64    `json:"expiresAt,omitempty"`
	RecipeTypes                   []string `json:"$recipeTypes,omitempty"`
	Type                          string   `json:"$type,omitempty"`
}

// ConnectionInfoResponse represents connection data from API response
type ConnectionInfoResponse struct {
	Paging      *PagingInfoResponse `json:"paging,omitempty"`
	Elements    []string            `json:"*elements,omitempty"`
	RecipeTypes []string            `json:"$recipeTypes,omitempty"`
	Type        string              `json:"$type,omitempty"`
}

// PagingInfoResponse represents pagination information
type PagingInfoResponse struct {
	Start       int      `json:"start,omitempty"`
	Count       int      `json:"count,omitempty"`
	Total       int      `json:"total,omitempty"`
	RecipeTypes []string `json:"$recipeTypes,omitempty"`
	Type        string   `json:"$type,omitempty"`
}

// FollowingStateResponse represents following state data
type FollowingStateResponse struct {
	EntityURN     string   `json:"entityUrn,omitempty"`
	Following     bool     `json:"following,omitempty"`
	FollowerCount int64    `json:"followerCount,omitempty"`
	FolloweeCount int64    `json:"followeeCount,omitempty"`
	RecipeTypes   []string `json:"$recipeTypes,omitempty"`
	Type          string   `json:"$type,omitempty"`
}

// PositionsCollection represents a collection of position/experience data
type PositionsCollection struct {
	Paging      *PagingInfoResponse `json:"paging,omitempty"`
	Elements    []string            `json:"*elements,omitempty"`
	RecipeTypes []string            `json:"$recipeTypes,omitempty"`
	Type        string              `json:"$type,omitempty"`
}

// EducationCollection represents a collection of education data
type EducationCollection struct {
	Paging      *PagingInfoResponse `json:"paging,omitempty"`
	Elements    []string            `json:"*elements,omitempty"`
	RecipeTypes []string            `json:"$recipeTypes,omitempty"`
	Type        string              `json:"$type,omitempty"`
}

// PositionResponse represents individual position/experience data from API
type PositionResponse struct {
	EntityURN              string              `json:"entityUrn,omitempty"`
	CompanyName            string              `json:"companyName,omitempty"`
	CompanyURN             string              `json:"*company,omitempty"`
	Title                  string              `json:"title,omitempty"`
	Description            string              `json:"description,omitempty"`
	DateRange              *DateRangeResponse  `json:"dateRange,omitempty"`
	LocationName           string              `json:"locationName,omitempty"`
	MultiLocaleCompanyName []map[string]string `json:"multiLocaleCompanyName,omitempty"`
	RecipeTypes            []string            `json:"$recipeTypes,omitempty"`
	Type                   string              `json:"$type,omitempty"`
}

// EducationResponse represents individual education data from API
type EducationResponse struct {
	EntityURN    string             `json:"entityUrn,omitempty"`
	SchoolName   string             `json:"schoolName,omitempty"`
	SchoolURN    string             `json:"*school,omitempty"`
	CompanyURN   string             `json:"*company,omitempty"`
	DegreeName   string             `json:"degreeName,omitempty"`
	FieldOfStudy string             `json:"fieldOfStudy,omitempty"`
	DateRange    *DateRangeResponse `json:"dateRange,omitempty"`
	Description  string             `json:"description,omitempty"`
	Activities   string             `json:"activities,omitempty"`
	RecipeTypes  []string           `json:"$recipeTypes,omitempty"`
	Type         string             `json:"$type,omitempty"`
}

// DateRangeResponse represents date range data from API response
type DateRangeResponse struct {
	Start       *DateResponse `json:"start,omitempty"`
	End         *DateResponse `json:"end,omitempty"`
	RecipeTypes []string      `json:"$recipeTypes,omitempty"`
	Type        string        `json:"$type,omitempty"`
}

// DateResponse represents date data from API response
type DateResponse struct {
	Year        int      `json:"year,omitempty"`
	Month       int      `json:"month,omitempty"`
	Day         int      `json:"day,omitempty"`
	RecipeTypes []string `json:"$recipeTypes,omitempty"`
	Type        string   `json:"$type,omitempty"`
}

// VerificationDataResponse represents verification information
type VerificationDataResponse struct {
	VerificationState interface{} `json:"verificationState,omitempty"`
	RecipeTypes       []string    `json:"$recipeTypes,omitempty"`
	Type              string      `json:"$type,omitempty"`
}

// CreatorInfoResponse represents creator information
type CreatorInfoResponse struct {
	CreatorWebsite        *TextViewModelResponse `json:"creatorWebsite,omitempty"`
	AssociatedHashtagUrns []string               `json:"associatedHashtagUrns,omitempty"`
	CreatorPostAnalytics  interface{}            `json:"creatorPostAnalytics,omitempty"`
	RecipeTypes           []string               `json:"$recipeTypes,omitempty"`
	Type                  string                 `json:"$type,omitempty"`
}

// TextViewModelResponse represents text with formatting from API
type TextViewModelResponse struct {
	Text              string        `json:"text,omitempty"`
	TextDirection     string        `json:"textDirection,omitempty"`
	AttributesV2      []interface{} `json:"attributesV2,omitempty"`
	AccessibilityText string        `json:"accessibilityText,omitempty"`
	RecipeTypes       []string      `json:"$recipeTypes,omitempty"`
	Type              string        `json:"$type,omitempty"`
}

// --- Search API Response Structures (existing) ---
