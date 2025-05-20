package linkedinscraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SearchProfiles searches for LinkedIn profiles based on the provided arguments.
func (c *Client) SearchProfiles(ctx context.Context, args ProfileSearchArgs) ([]LinkedInProfile, error) {
	// Input Validation
	if c.config.Auth.LiAtCookie == "" || c.config.Auth.CSRFToken == "" {
		return nil, ErrAuthMissing
	}
	if args.Keywords == "" {
		return nil, ErrKeywordsMissing
	}

	// Construct SearchVariables
	querySubQuery := SearchQuerySubQuery{
		Keywords:                 args.Keywords,
		FlagshipSearchIntent:     "SEARCH_SRP", // from cURL
		QueryParameters:          []SearchQueryParameters{},
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
	}

	// Build URL
	requestURL, err := buildGraphQLURL(VoyagerBaseURL, DefaultSearchQueryID, variables)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestBuildFailed, err) // Wrap ErrRequestBuildFailed
	}

	// Prepare Headers
	customHeaders := http.Header{}
	// Example Referer, needs to be more robust or configurable
	customHeaders.Set("Referer", "https://www.linkedin.com/search/results/people/?keywords="+url.QueryEscape(args.Keywords))
	// X-Li-Page-Instance and X-Li-Track are complex and might need dynamic generation or configuration.
	// For now, using placeholders or values that might be common/static enough for initial testing.
	customHeaders.Set("X-Li-Page-Instance", "urn:li:page:d_flagship3_search_srp_people;placeholder") // Placeholder
	customHeaders.Set("X-Li-Pem-Metadata", "Voyager - People SRP=search-results")
	// A simplified or static X-Li-Track. This is highly likely to need adjustment.
	customHeaders.Set("X-Li-Track", `{"clientVersion":"1.13.x","mpVersion":"1.13.x","pageKey":"p_flagship3_search_srp_people","traceId":"placeholderTraceId"}`)

	// Make API Call
	resp, respBodyBytes, err := c.makeRequest(ctx, http.MethodGet, requestURL, customHeaders, nil)
	if err != nil {
		// It might be beneficial to inspect the error type if makeRequest returns a wrapped error
		// that could indicate a more specific issue (e.g., context canceled, network error before HTTP execution)
		return nil, fmt.Errorf("%w: %v", ErrRequestFailed, err) // Wrap ErrRequestFailed
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
	var apiResponse SearchAPIResponse
	err = json.Unmarshal(respBodyBytes, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("%w: %v. Raw response: %s", ErrResponseParseFailed, err, string(respBodyBytes))
	}

	// Extract Profiles
	var profiles []LinkedInProfile
	profileDataMap := make(map[string]IncludedProfile) // To store IncludedProfile data by URN for enrichment

	// First pass: collect all IncludedProfile data
	for _, item := range apiResponse.Included {
		if item.Type == "com.linkedin.voyager.dash.identity.profile.Profile" {
			// Check for nil pointers before dereferencing, though fields are not pointers in IncludedProfile itself based on current models.go
			// However, item itself could represent a partially unmarshalled element if not all fields were present.
			// For simplicity, we'll assume direct field access is safe if Type matches.
			profileDataMap[item.EntityURN] = IncludedProfile{
				EntityURN:        item.EntityURN,
				PublicIdentifier: item.PublicIdentifier,
				FirstName:        item.FirstName,
				LastName:         item.LastName,
				Headline:         item.Headline,
			}
		}
	}

	// Second pass: build LinkedInProfile from EntityResultViewModel, enriching with Profile data
	for _, item := range apiResponse.Included {
		if item.Type == "com.linkedin.voyager.dash.search.EntityResultViewModel" {
			if item.Title == nil || item.PrimarySubtitle == nil || item.SecondarySubtitle == nil {
				// Skip if essential fields are missing to avoid nil pointer dereference
				// Consider logging this case if robust error handling/reporting is needed
				continue
			}

			profile := LinkedInProfile{
				URN:        item.TrackingURN, // TrackingURN from EntityResultViewModel is often the profile URN
				FullName:   item.Title.Text,
				Headline:   item.PrimarySubtitle.Text,
				Location:   item.SecondarySubtitle.Text,
				ProfileURL: item.NavigationURL,
				// PublicIdentifier can come from EntityResultViewModel itself or be enriched
			}

			// Attempt to get PublicIdentifier directly from EntityResultViewModel's own PublicIdentifier field if it exists and is populated
			if item.PublicIdentifier != "" {
				profile.PublicIdentifier = item.PublicIdentifier
			}

			// Enrich with data from IncludedProfile if available, prioritizing already set publicIdentifier
			if linkedProfileData, ok := profileDataMap[item.TrackingURN]; ok {
				if profile.PublicIdentifier == "" && linkedProfileData.PublicIdentifier != "" {
					profile.PublicIdentifier = linkedProfileData.PublicIdentifier
				}
				// Potentially update other fields if EntityResultViewModel's were less complete, e.g. headline
				// For now, we primarily use EntityResultViewModel and supplement publicId
			}

			// If PublicIdentifier is still empty, and URN looks like a profile URN,
			// we might be able to derive it, but this is often unreliable.
			// Example: urn:li:fsd_profile:ACoAAAtp-4UBpQ0aZ_PeToflBoLty9BpO_CQ6-I
			// Public ID can sometimes be part of another field or require a separate lookup/parsing strategy if not directly available.
			// For now, we rely on it being present in either EntityResultViewModel or IncludedProfile.

			profiles = append(profiles, profile)
		}
	}

	if len(profiles) == 0 {
		// Depending on requirements, could return ErrNoProfilesFound or empty slice.
		// The current error definition notes "Or handle this by returning empty slice"
		// For now, let's stick to returning an empty slice if no profiles were parsed,
		// as the API call itself might have been successful but yielded no relevant entities.
		// If an error like ErrNoProfilesFound is desired, it should be returned here.
		return []LinkedInProfile{}, nil
	}

	return profiles, nil
}
