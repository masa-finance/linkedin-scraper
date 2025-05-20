package linkedinscraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
		Count:  args.Count,       // Populate Count from args
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
	customHeaders.Set("Accept", "application/vnd.linkedin.normalized+json+2.1") // Ensure correct Accept header from cURL

	// Construct Referer URL
	// The cURL Referer is: https://www.linkedin.com/search/results/people/?keywords=investor&network=["F","O"]&origin=FACETED_SEARCH
	// Note: network parameter is a literal JSON array string, not URL encoded components.
	var refererQueryParts []string
	refererQueryParts = append(refererQueryParts, "keywords="+url.QueryEscape(args.Keywords))

	if len(args.NetworkFilters) > 0 {
		// Create the literal JSON array string for the network filter
		networkFilterString := "[\"" + strings.Join(args.NetworkFilters, "\",\"") + "\"]"
		refererQueryParts = append(refererQueryParts, "network="+networkFilterString) // Do not QueryEscape the already formatted JSON string
	}
	refererQueryParts = append(refererQueryParts, "origin=FACETED_SEARCH")

	baseURLForReferer := "https://www.linkedin.com/search/results/people/"
	fullRefererURL := baseURLForReferer + "?" + strings.Join(refererQueryParts, "&")
	customHeaders.Set("Referer", fullRefererURL)

	// Use XLiPageInstance from args if provided, otherwise use placeholder
	xLiPageInstance := "urn:li:page:d_flagship3_search_srp_people;placeholder" // Default placeholder
	if args.XLiPageInstance != "" {
		xLiPageInstance = args.XLiPageInstance
	}
	customHeaders.Set("X-Li-Page-Instance", xLiPageInstance)

	customHeaders.Set("X-Li-Pem-Metadata", "Voyager - People SRP=search-results")

	// Use XLiTrack from args if provided, otherwise use placeholder matching cURL structure
	// cURL: {"clientVersion":"1.13.35368","mpVersion":"1.13.35368","osName":"web","timezoneOffset":-7,"timezone":"America/Los_Angeles","deviceFormFactor":"DESKTOP","mpName":"voyager-web","displayDensity":2,"displayWidth":5120,"displayHeight":2880}
	xLiTrack := `{"clientVersion":"1.13.35368","mpVersion":"1.13.35368","osName":"web","timezoneOffset":-7,"timezone":"America/Los_Angeles","deviceFormFactor":"DESKTOP","mpName":"voyager-web","displayDensity":2,"displayWidth":1920,"displayHeight":1080}` // Default placeholder, using common display W/H
	if args.XLiTrack != "" {
		xLiTrack = args.XLiTrack
	}
	customHeaders.Set("X-Li-Track", xLiTrack)

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
