package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	linkedinscraper "github.com/masa-finance/linkedin-scraper"
)

func main() {
	// Load .env file. Errors are not fatal for flexibility (e.g. CI environments might use actual env vars)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, relying on environment variables: ", err)
	}

	fmt.Println("LinkedIn Profile Scraper Example")

	liAtCookie := os.Getenv("LI_AT_COOKIE")
	csrfToken := os.Getenv("CSRF_TOKEN")
	jsessionID := os.Getenv("JSESSIONID_TOKEN") // Ensure this matches the cURL: "ajax:..."

	if liAtCookie == "" || csrfToken == "" {
		log.Fatal("Error: LI_AT_COOKIE and CSRF_TOKEN environment variables must be set.")
	}
	if jsessionID == "" {
		log.Println("Warning: JSESSIONID_TOKEN environment variable is not set. It might be required.")
		// For the JSESSIONID from cURL, it was like "ajax:YOUR_ACTUAL_ID". Ensure this format is used if it includes a prefix.
		// If it's just the ID, that's fine. The current client.go code wraps it in quotes: JSESSIONID=\"ajax:...\"
	}

	auth := linkedinscraper.AuthCredentials{
		LiAtCookie: liAtCookie,
		CSRFToken:  csrfToken,
		JSESSIONID: jsessionID, // Pass the raw value, client.go will format it in the cookie string
	}

	cfg, err := linkedinscraper.NewConfig(auth)
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	// Optional: Set a custom User-Agent if needed, otherwise default will be used
	// cfg.UserAgent = "MyCustomUserAgent/1.0"

	client, err := linkedinscraper.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	fmt.Println("Searching for profiles with keyword 'software engineer'...")
	searchArgs := linkedinscraper.ProfileSearchArgs{
		Keywords:       "software engineer",
		NetworkFilters: []string{"S"}, // Example: "S" for Second degree connections. Use {"F", "S"} for First and Second.
		Start:          0,             // Start from the first result
		Count:          5,             // Request 5 results (Note: API actual count might differ)
	}

	profiles, err := client.SearchProfiles(context.Background(), searchArgs)
	if err != nil {
		log.Fatalf("Error searching profiles: %v", err)
	}

	if len(profiles) == 0 {
		fmt.Println("No profiles found or returned.")
	} else {
		fmt.Printf("Found %d profile(s):\n", len(profiles))
		// Pretty print the JSON output
		profilesJSON, err := json.MarshalIndent(profiles, "", "  ")
		if err != nil {
			log.Fatalf("Error marshalling profiles to JSON: %v", err)
		}
		fmt.Println(string(profilesJSON))
	}

	fmt.Println("\nExample finished.")
}
