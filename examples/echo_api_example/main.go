package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	linkedinscraper "github.com/masa-finance/linkedin-scraper"
)

// SimpleMessage is a struct for JSON responses (can be used for errors or simple status)
type SimpleMessage struct {
	Message string `json:"message"`
}

func main() {
	// Load .env file from the current directory (examples/echo_api_example)
	// For this to work, you'll need to create a .env file in examples/echo_api_example
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found in ./examples/echo_api_example/. Relying on globally set environment variables.")
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/search/linkedin", searchLinkedInProfilesHandler)

	// Start server
	log.Println("Starting Echo server on :1323...")
	log.Println("Try: curl \"http://localhost:1323/search/linkedin?keywords=software+engineer\"")
	e.Logger.Fatal(e.Start(":1323"))
}

// searchLinkedInProfilesHandler handles requests to search LinkedIn profiles
func searchLinkedInProfilesHandler(c echo.Context) error {
	keywords := c.QueryParam("keywords")
	if keywords == "" {
		return c.JSON(http.StatusBadRequest, SimpleMessage{Message: "keywords query parameter is required"})
	}

	// Load credentials from environment variables
	liAtCookie := os.Getenv("LI_AT_COOKIE")
	csrfToken := os.Getenv("CSRF_TOKEN")
	jsessionID := os.Getenv("JSESSIONID_TOKEN") // Ensure this matches your .env key for JSESSIONID

	if liAtCookie == "" || csrfToken == "" || jsessionID == "" {
		log.Println("Error: Missing one or more LinkedIn API credentials in environment variables (LI_AT_COOKIE, CSRF_TOKEN, JSESSIONID_TOKEN)")
		return c.JSON(http.StatusInternalServerError, SimpleMessage{Message: "Server configuration error: API credentials missing. Ensure LI_AT_COOKIE, CSRF_TOKEN, and JSESSIONID_TOKEN are set."})
	}

	authCreds := linkedinscraper.AuthCredentials{
		LiAtCookie: liAtCookie,
		CSRFToken:  csrfToken,
		JSESSIONID: jsessionID,
		// XLiTrack and XLiPageInstance will use defaults from the linkedinscraper package
	}

	// We assume NewConfig and NewClient handle default User-Agent etc.
	// Pass empty strings for userAgent, proxyURL, customQueryID to use package defaults.
	config, err := linkedinscraper.NewConfig(authCreds, "", "", "")
	if err != nil {
		log.Printf("Error creating linkedinscraper.Config: %v", err)
		return c.JSON(http.StatusInternalServerError, SimpleMessage{Message: "Failed to initialize LinkedIn client config: " + err.Error()})
	}

	client, err := linkedinscraper.NewClient(config)
	if err != nil {
		log.Printf("Error creating linkedinscraper.Client: %v", err)
		return c.JSON(http.StatusInternalServerError, SimpleMessage{Message: "Failed to initialize LinkedIn client: " + err.Error()})
	}

	searchArgs := linkedinscraper.ProfileSearchArgs{
		Keywords: keywords,
		Count:    5, // Default to 5 results for this example, make configurable if needed
		Start:    0,
		// NetworkFilters, XLiTrack, XLiPageInstance are omitted to rely on package defaults or if not strictly needed
	}

	log.Printf("Searching LinkedIn for keywords: %s", keywords)
	profiles, err := client.SearchProfiles(context.Background(), searchArgs)
	if err != nil {
		log.Printf("Error from SearchProfiles: %v", err)
		// Check for specific, user-friendly errors from your package
		if err == linkedinscraper.ErrUnauthorized {
			return c.JSON(http.StatusUnauthorized, SimpleMessage{Message: "LinkedIn API Unauthorized: Check credentials or they might have expired."})
		}
		if err == linkedinscraper.ErrRateLimited {
			return c.JSON(http.StatusTooManyRequests, SimpleMessage{Message: "LinkedIn API rate limit hit."})
		}
		// Generic error
		return c.JSON(http.StatusInternalServerError, SimpleMessage{Message: "Failed to search LinkedIn profiles: " + err.Error()})
	}

	if len(profiles) == 0 {
		return c.JSON(http.StatusOK, SimpleMessage{Message: "No profiles found for keywords: " + keywords})
	}

	return c.JSON(http.StatusOK, profiles)
}
