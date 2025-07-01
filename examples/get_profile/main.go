package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	linkedinscraper "github.com/masa-finance/linkedin-scraper"
)

func main() {
	fmt.Println("🔍 LinkedIn Profile Scraper - Profile Fetching Example")
	fmt.Println("====================================================")

	// Get configuration from environment variables
	config, err := loadConfigFromEnv()
	if err != nil {
		log.Fatalf("❌ Configuration error: %v", err)
	}

	// Create client
	client, err := linkedinscraper.NewClient(config)
	if err != nil {
		log.Fatalf("❌ Failed to create client: %v", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Example 1: Fetch a specific profile
	fmt.Println("\n📋 Example 1: Fetching Profile by Public Identifier")
	fmt.Println("---------------------------------------------------")

	publicIdentifier := getPublicIdentifierFromArgs()
	profile, err := client.GetProfile(ctx, publicIdentifier)
	if err != nil {
		log.Printf("❌ Failed to fetch profile for %s: %v", publicIdentifier, err)
	} else {
		displayProfile(profile)
	}

	// Example 2: Integration with search - find profiles and fetch detailed data
	fmt.Println("\n🔍 Example 2: Search + Profile Integration")
	fmt.Println("------------------------------------------")

	searchArgs := linkedinscraper.ProfileSearchArgs{
		Keywords:       "software engineer",
		NetworkFilters: []string{"F", "O"}, // 1st and 2nd degree connections
		Start:          0,
		Count:          3, // Fetch first 3 results
	}

	profiles, err := client.SearchProfiles(ctx, searchArgs)
	if err != nil {
		log.Printf("❌ Search failed: %v", err)
	} else {
		fmt.Printf("✅ Found %d profiles from search\n", len(profiles))

		// Fetch detailed profile data for each search result
		for i, searchProfile := range profiles {
			if searchProfile.PublicIdentifier == "" {
				fmt.Printf("⚠️  Profile %d: No public identifier available\n", i+1)
				continue
			}

			fmt.Printf("\n📊 Fetching detailed data for profile %d: %s\n", i+1, searchProfile.PublicIdentifier)

			detailedProfile, err := client.GetProfile(ctx, searchProfile.PublicIdentifier)
			if err != nil {
				log.Printf("❌ Failed to fetch detailed profile: %v", err)
				continue
			}

			displayProfileSummary(detailedProfile, i+1)

			// Add delay to respect rate limits
			time.Sleep(1 * time.Second)
		}
	}

	// Example 3: Export profile data to JSON
	if profile != nil {
		fmt.Println("\n💾 Example 3: Export Profile to JSON")
		fmt.Println("------------------------------------")

		err := exportProfileToJSON(profile, fmt.Sprintf("%s_profile.json", profile.PublicIdentifier))
		if err != nil {
			log.Printf("❌ Failed to export profile: %v", err)
		} else {
			fmt.Printf("✅ Profile exported to %s_profile.json\n", profile.PublicIdentifier)
		}
	}

	fmt.Println("\n🎉 Profile fetching examples completed!")
}

// loadConfigFromEnv loads configuration from environment variables
func loadConfigFromEnv() (*linkedinscraper.Config, error) {
	liAtCookie := os.Getenv("LINKEDIN_LI_AT")
	csrfToken := os.Getenv("LINKEDIN_CSRF_TOKEN")
	jsessionid := os.Getenv("LINKEDIN_JSESSIONID")

	if liAtCookie == "" || csrfToken == "" {
		return nil, fmt.Errorf("required environment variables missing. Please set:\n" +
			"  LINKEDIN_LI_AT=your-li-at-cookie\n" +
			"  LINKEDIN_CSRF_TOKEN=your-csrf-token\n" +
			"  LINKEDIN_JSESSIONID=your-jsessionid (optional)")
	}

	auth := linkedinscraper.AuthCredentials{
		LiAtCookie: liAtCookie,
		CSRFToken:  csrfToken,
		JSESSIONID: jsessionid,
	}

	return linkedinscraper.NewConfig(auth)
}

// getPublicIdentifierFromArgs gets public identifier from command line args or uses default
func getPublicIdentifierFromArgs() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	// Default example profile (you can change this)
	return "aluko-ayodeji-367370177"
}

// displayProfile shows comprehensive profile information
func displayProfile(profile *linkedinscraper.LinkedInProfile) {
	fmt.Printf("✅ Profile fetched successfully!\n\n")

	// Basic Information
	fmt.Printf("👤 Basic Information:\n")
	fmt.Printf("  Name: %s\n", profile.FullName)
	fmt.Printf("  Headline: %s\n", profile.Headline)
	fmt.Printf("  Public ID: %s\n", profile.PublicIdentifier)
	fmt.Printf("  Profile URL: %s\n", profile.ProfileURL)
	fmt.Printf("  URN: %s\n", profile.URN)

	// Professional Information
	if len(profile.Experience) > 0 {
		fmt.Printf("\n💼 Experience (%d entries):\n", len(profile.Experience))
		for i, exp := range profile.Experience {
			if i >= 3 { // Show only first 3 for brevity
				fmt.Printf("  ... and %d more\n", len(profile.Experience)-3)
				break
			}
			fmt.Printf("  %d. %s at %s\n", i+1, exp.Title, exp.CompanyName)
		}
	}

	if len(profile.Education) > 0 {
		fmt.Printf("\n🎓 Education (%d entries):\n", len(profile.Education))
		for i, edu := range profile.Education {
			if i >= 2 { // Show only first 2 for brevity
				fmt.Printf("  ... and %d more\n", len(profile.Education)-2)
				break
			}
			fmt.Printf("  %d. %s at %s\n", i+1, edu.DegreeName, edu.SchoolName)
		}
	}

	if len(profile.Skills) > 0 {
		fmt.Printf("\n🛠️  Skills (%d total):\n", len(profile.Skills))
		skillCount := len(profile.Skills)
		if skillCount > 5 {
			skillCount = 5 // Show only first 5
		}
		for i := 0; i < skillCount; i++ {
			fmt.Printf("  • %s", profile.Skills[i].Name)
			if profile.Skills[i].EndorsementCount > 0 {
				fmt.Printf(" (%d endorsements)", profile.Skills[i].EndorsementCount)
			}
			fmt.Printf("\n")
		}
		if len(profile.Skills) > 5 {
			fmt.Printf("  ... and %d more skills\n", len(profile.Skills)-5)
		}
	}

	// Location and Contact
	if profile.LocationDetails != nil {
		fmt.Printf("\n📍 Location:\n")
		fmt.Printf("  Country: %s\n", profile.LocationDetails.CountryCode)
		if profile.LocationDetails.PostalCode != "" {
			fmt.Printf("  Postal Code: %s\n", profile.LocationDetails.PostalCode)
		}
	}

	// Social Information
	if profile.ConnectionInfo != nil {
		fmt.Printf("\n🌐 Social Information:\n")
		if profile.ConnectionInfo.ConnectionCount > 0 {
			fmt.Printf("  Connections: %d\n", profile.ConnectionInfo.ConnectionCount)
		}
		if profile.ConnectionInfo.FollowerCount > 0 {
			fmt.Printf("  Followers: %d\n", profile.ConnectionInfo.FollowerCount)
		}
		fmt.Printf("  Following: %t\n", profile.ConnectionInfo.Following)
	}

	// Additional Metadata
	fmt.Printf("\n📊 Additional Information:\n")
	fmt.Printf("  Creator: %t\n", profile.IsCreator)
	fmt.Printf("  Verified: %t\n", profile.IsVerified)
	fmt.Printf("  Premium: %t\n", profile.IsPremium)
	fmt.Printf("  Memorialized: %t\n", profile.IsMemorialized)

	if profile.TempStatus != "" {
		fmt.Printf("  Status: %s", profile.TempStatus)
		if profile.TempStatusEmoji != "" {
			fmt.Printf(" %s", profile.TempStatusEmoji)
		}
		fmt.Printf("\n")
	}
}

// displayProfileSummary shows a brief summary of a profile
func displayProfileSummary(profile *linkedinscraper.LinkedInProfile, index int) {
	fmt.Printf("  %d. %s - %s\n", index, profile.FullName, profile.Headline)
	fmt.Printf("     URL: %s\n", profile.ProfileURL)
	fmt.Printf("     Experience: %d entries, Education: %d entries, Skills: %d entries\n",
		len(profile.Experience), len(profile.Education), len(profile.Skills))
}

// exportProfileToJSON exports a profile to a JSON file
func exportProfileToJSON(profile *linkedinscraper.LinkedInProfile, filename string) error {
	jsonData, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile to JSON: %w", err)
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
