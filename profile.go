package linkedinscraper

import (
	"encoding/json"
	"fmt"
	"strings"
)

// parseProfileFromAPIResponse parses a ProfileAPIResponse and extracts comprehensive profile data.
func parseProfileFromAPIResponse(apiResponse *ProfileAPIResponse, publicIdentifier string) (*LinkedInProfile, error) {
	// Find the main profile entity in the included array
	var profileEntity *GenericIncludedElement

	for i, item := range apiResponse.Included {
		if item.Type == EntityTypeProfile &&
			item.PublicIdentifier == publicIdentifier {
			profileEntity = &apiResponse.Included[i]
			break
		}
	}

	if profileEntity == nil {
		return nil, fmt.Errorf("profile not found in API response for publicIdentifier: %s", publicIdentifier)
	}

	// Start building the LinkedInProfile
	profile := &LinkedInProfile{
		PublicIdentifier: profileEntity.PublicIdentifier,
		URN:              profileEntity.EntityURN,
		FirstName:        profileEntity.FirstName,
		LastName:         profileEntity.LastName,
		Headline:         profileEntity.Headline,
		ProfileURL:       fmt.Sprintf("https://www.linkedin.com/in/%s/", publicIdentifier),
	}

	// Set FullName
	profile.FullName = strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	// Parse additional profile data by finding and processing related entities
	profile.Experience = parseExperienceData(apiResponse, profileEntity.EntityURN)
	profile.Education = parseEducationData(apiResponse, profileEntity.EntityURN)
	profile.Skills = parseSkillsData(apiResponse, profileEntity.EntityURN)
	profile.LocationDetails = parseLocationData(apiResponse, profileEntity.EntityURN)
	profile.ConnectionInfo = parseConnectionData(apiResponse, profileEntity.EntityURN)
	profile.ProfilePicture = parseProfilePictureData(apiResponse, profileEntity.EntityURN)

	// Parse simple fields from the profile entity itself
	parseSimpleProfileFields(profile, apiResponse, profileEntity)

	return profile, nil
}

// parseExperienceData extracts experience/position data from the API response.
func parseExperienceData(apiResponse *ProfileAPIResponse, profileURN string) []Experience {
	var experiences []Experience
	for _, item := range apiResponse.Included {
		if item.Type == EntityTypePosition {
			experience := Experience{
				EntityURN:    item.EntityURN,
				CompanyName:  item.CompanyName,
				Description:  item.Description,
				LocationName: item.LocationName,
				CompanyURN:   item.CompanyURN,
			}
			if item.Title != nil {
				experience.Title = string(*item.Title)
			}
			if item.DateRange != nil {
				experience.DateRange = &DateRange{}
				if item.DateRange.Start != nil {
					experience.DateRange.Start = &Date{
						Year:  item.DateRange.Start.Year,
						Month: item.DateRange.Start.Month,
						Day:   item.DateRange.Start.Day,
					}
				}
				if item.DateRange.End != nil {
					experience.DateRange.End = &Date{
						Year:  item.DateRange.End.Year,
						Month: item.DateRange.End.Month,
						Day:   item.DateRange.End.Day,
					}
				}
			}
			experiences = append(experiences, experience)
		}
	}
	return experiences
}

// parseEducationData extracts education data from the API response.
func parseEducationData(apiResponse *ProfileAPIResponse, profileURN string) []Education {
	var education []Education
	for _, item := range apiResponse.Included {
		if item.Type == EntityTypeEducation {
			edu := Education{
				EntityURN:    item.EntityURN,
				SchoolName:   item.SchoolName,
				SchoolURN:    item.SchoolURN,
				DegreeName:   item.DegreeName,
				FieldOfStudy: item.FieldOfStudy,
				Description:  item.Description,
				Activities:   item.Activities,
			}
			if item.DateRange != nil {
				edu.DateRange = &DateRange{}
				if item.DateRange.Start != nil {
					edu.DateRange.Start = &Date{
						Year:  item.DateRange.Start.Year,
						Month: item.DateRange.Start.Month,
						Day:   item.DateRange.Start.Day,
					}
				}
				if item.DateRange.End != nil {
					edu.DateRange.End = &Date{
						Year:  item.DateRange.End.Year,
						Month: item.DateRange.End.Month,
						Day:   item.DateRange.End.Day,
					}
				}
			}
			education = append(education, edu)
		}
	}
	return education
}

// parseSkillsData extracts skills data from the API response.
func parseSkillsData(apiResponse *ProfileAPIResponse, profileURN string) []Skill {
	var skills []Skill
	for _, item := range apiResponse.Included {
		if strings.Contains(item.Type, EntityTypeEndorsedSkill) { // The type can vary slightly
			skill := Skill{
				EntityURN:        item.EntityURN,
				Name:             item.Name,
				EndorsementCount: item.EndorsementCount,
				EndorsedByViewer: item.EndorsedByViewer,
			}
			skills = append(skills, skill)
		}
	}
	return skills
}

// parseLocationData extracts location information from the API response.
func parseLocationData(apiResponse *ProfileAPIResponse, profileURN string) *ProfileLocation {
	// Look for location data in the main profile entity or related entities
	for _, item := range apiResponse.Included {
		if item.Type == EntityTypeProfile &&
			item.EntityURN == profileURN {
			// Parse location from the profile entity
			// This would need to be adjusted based on actual API structure
			return &ProfileLocation{
				CountryCode: extractCountryCode(item),
			}
		}
	}

	return nil
}

// parseConnectionData extracts connection and following information.
func parseConnectionData(apiResponse *ProfileAPIResponse, profileURN string) *ConnectionInfo {
	connectionInfo := &ConnectionInfo{}

	for _, item := range apiResponse.Included {
		if strings.Contains(item.Type, EntityTypeConnection) {
			// Parse connection count from the item
			// This would need adjustment based on actual API structure
			if count, err := parseConnectionCount(item); err == nil {
				connectionInfo.ConnectionCount = count
			}
		}
		if strings.Contains(item.Type, EntityTypeFollowing) {
			// Parse follower/following information
			// This would need adjustment based on actual API structure
			if count, err := parseFollowerCount(item); err == nil {
				connectionInfo.FollowerCount = count
			}
		}
	}

	return connectionInfo
}

// parseProfilePictureData extracts profile picture information.
func parseProfilePictureData(apiResponse *ProfileAPIResponse, profileURN string) *ProfilePicture {
	for _, item := range apiResponse.Included {
		if item.Type == EntityTypeProfile &&
			item.EntityURN == profileURN {
			return &ProfilePicture{
				DisplayImageUrn: extractProfileImageURN(item),
				A11yText:        item.FirstName + " " + item.LastName,
			}
		}
	}

	return nil
}

// parseSimpleProfileFields extracts simple fields directly from the profile entity.
func parseSimpleProfileFields(profile *LinkedInProfile, apiResponse *ProfileAPIResponse, profileEntity *GenericIncludedElement) {
	// Parse creator status
	if creatorValue, exists := extractFieldFromRawJSON(profileEntity, "creator"); exists {
		if creator, ok := creatorValue.(bool); ok {
			profile.IsCreator = creator
		}
	}

	// Parse memorialized status
	if memorializedValue, exists := extractFieldFromRawJSON(profileEntity, "memorialized"); exists {
		if memorialized, ok := memorializedValue.(bool); ok {
			profile.IsMemorialized = memorialized
		}
	}

	// Parse temp status
	if tempStatusValue, exists := extractFieldFromRawJSON(profileEntity, "tempStatus"); exists {
		if tempStatus, ok := tempStatusValue.(string); ok {
			profile.TempStatus = tempStatus
		}
	}

	// Parse temp status emoji
	if tempEmojiValue, exists := extractFieldFromRawJSON(profileEntity, "tempStatusEmoji"); exists {
		if tempEmoji, ok := tempEmojiValue.(string); ok {
			profile.TempStatusEmoji = tempEmoji
		}
	}
}

// Helper functions for parsing specific data types

// extractCountryCode extracts country code from a profile entity.
func extractCountryCode(item GenericIncludedElement) string {
	// This would need to be implemented based on actual API response structure
	// For now, return empty string as placeholder
	return ""
}

// parseConnectionCount extracts connection count from a connection entity.
func parseConnectionCount(item GenericIncludedElement) (int, error) {
	// This would need to be implemented based on actual API response structure
	return 0, fmt.Errorf("not implemented")
}

// parseFollowerCount extracts follower count from a following entity.
func parseFollowerCount(item GenericIncludedElement) (int, error) {
	// This would need to be implemented based on actual API response structure
	return 0, fmt.Errorf("not implemented")
}

// extractProfileImageURN extracts profile image URN from a profile entity.
func extractProfileImageURN(item GenericIncludedElement) string {
	// This would need to be implemented based on actual API response structure
	return ""
}

// extractFieldFromRawJSON extracts a field from the raw JSON data of an entity.
// This is a helper function to access fields that aren't in the struct.
func extractFieldFromRawJSON(item *GenericIncludedElement, fieldName string) (interface{}, bool) {
	// This would require implementing raw JSON parsing
	// For now, return false as placeholder
	return nil, false
}

// validateProfileData validates and sanitizes profile data.
func validateProfileData(profile *LinkedInProfile) error {
	if profile == nil {
		return fmt.Errorf("profile cannot be nil")
	}

	if profile.PublicIdentifier == "" {
		return fmt.Errorf("publicIdentifier is required")
	}

	// Sanitize text fields
	profile.FirstName = sanitizeTextString(profile.FirstName)
	profile.LastName = sanitizeTextString(profile.LastName)
	profile.FullName = sanitizeTextString(profile.FullName)
	profile.Headline = sanitizeTextString(profile.Headline)
	profile.Summary = sanitizeTextString(profile.Summary)

	return nil
}

// sanitizeTextString removes harmful characters and trims whitespace.
func sanitizeTextString(s string) string {
	// Remove any potentially harmful characters
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)
	return s
}

// ParseFromJSON parses a JSON string into a LinkedInProfile.
// This is useful for testing and parsing saved JSON responses.
func ParseFromJSON(jsonData []byte) (*LinkedInProfile, error) {
	var apiResponse ProfileAPIResponse
	err := json.Unmarshal(jsonData, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Extract public identifier from the response
	publicIdentifier := extractPublicIdentifierFromResponse(&apiResponse)
	if publicIdentifier == "" {
		return nil, fmt.Errorf("could not extract publicIdentifier from response")
	}

	profile, err := parseProfileFromAPIResponse(&apiResponse, publicIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profile: %w", err)
	}

	err = validateProfileData(profile)
	if err != nil {
		return nil, fmt.Errorf("profile validation failed: %w", err)
	}

	return profile, nil
}

// extractPublicIdentifierFromResponse extracts the public identifier from the API response.
func extractPublicIdentifierFromResponse(apiResponse *ProfileAPIResponse) string {
	for _, item := range apiResponse.Included {
		if item.Type == EntityTypeProfile &&
			item.PublicIdentifier != "" {
			return item.PublicIdentifier
		}
	}
	return ""
}

// Advanced parsing functions for complex nested structures

// parseVectorImage parses vector image data from the API response.
func parseVectorImage(rawData interface{}) *ProfilePicture {
	// This would implement parsing of the complex vector image structure
	// including artifacts, URLs, etc.
	return nil
}

// parseDateRange parses LinkedIn's date range format.
func parseDateRange(rawData interface{}) *DateRange {
	// This would implement parsing of LinkedIn's date structure
	// with year, month, day fields
	return nil
}

// parseTextViewModel parses LinkedIn's text view model with formatting.
func parseTextViewModel(rawData interface{}) string {
	// This would implement parsing of formatted text with attributes
	return ""
}

// convertAPIResponseToLinkedInProfile is the main conversion function used by the client.
func convertAPIResponseToLinkedInProfile(apiResponse *ProfileAPIResponse, publicIdentifier string) (*LinkedInProfile, error) {
	profile, err := parseProfileFromAPIResponse(apiResponse, publicIdentifier)
	if err != nil {
		return nil, err
	}

	err = validateProfileData(profile)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
