# linkedin-scraper

A Go package for searching LinkedIn profiles using the Voyager API.

## Prerequisites

To use this scraper, you need to obtain a few pieces of information from your LinkedIn session in a web browser:

1.  **`li_at` Cookie**: This is your LinkedIn session cookie.
2.  **CSRF Token (`csrf-token`)**: This token is used for request verification.
3.  **`JSESSIONID` Cookie**: While not always strictly enforced for all requests in the same way as `li_at`, it's often part of the authentication context. The value typically looks like `"ajax:YOUR_SESSION_ID"`.

### How to Obtain Credentials:

1.  Log in to your LinkedIn account in a web browser (e.g., Chrome, Firefox).
2.  Open your browser's Developer Tools (usually by pressing F12 or right-clicking on the page and selecting "Inspect" or "Inspect Element").
3.  **For Cookies (`li_at`, `JSESSIONID`):**
    *   Go to the "Application" tab (in Chrome) or "Storage" tab (in Firefox).
    *   Under "Cookies" (or "Storage" -> "Cookies"), find the domain `https://www.linkedin.com`.
    *   Locate the `li_at` cookie and copy its "Value".
    *   Locate the `JSESSIONID` cookie and copy its "Value". Make sure to include the `ajax:` prefix if it's present (e.g., `ajax:1234567890123456789`).
4.  **For CSRF Token:**
    *   This token is often found in the headers of POST/PUT requests made by your browser to LinkedIn, or sometimes embedded in the page source. A common way to find it is to look at a recent authenticated request made by LinkedIn itself.
    *   Alternatively, you can sometimes find it by inspecting the page source for a hidden input field named `csrfToken` or by checking XHR request headers for `csrf-token`.
    *   *Note*: The `csrf-token` required by the Voyager API is often *different* from the `JSESSIONID` cookie value that has `ajax:` prefix. The `csrf-token` is typically a separate, longer alphanumeric string.

**Important Security Note:** These credentials provide access to your LinkedIn account. Keep them secure and do not share them publicly. For development, it's recommended to use environment variables or a `.env` file to manage these secrets.

## Installation

```bash
go get github.com/masa-finance/linkedin-scraper
```

## Example Usage

The following example demonstrates how to search for profiles. It expects your LinkedIn credentials to be set as environment variables or in a `.env` file.

1.  **Set up your credentials:**
    *   Copy the `.env.example` file located in the root of this project to a new file named `.env` in the project root:
        ```bash
        cp .env.example .env
        ```
    *   Edit the `.env` file and replace the placeholder values with your actual credentials:

        ```env
        # .env (in project root)
        LI_AT_COOKIE="YOUR_LI_AT_COOKIE_VALUE"
        CSRF_TOKEN="YOUR_CSRF_TOKEN_VALUE"
        JSESSIONID_TOKEN="ajax:YOUR_JSESSIONID_VALUE" # Ensure to include "ajax:" prefix if present
        ```
    *   The example program (`examples/search_profiles/main.go`) uses `godotenv.Load()` without a specific path. This means it will load the `.env` file if the program is run from the directory containing the `.env` file (i.e., the project root).
        ```bash
        # From the project root directory:
        go run examples/search_profiles/main.go
        ```
    *   Alternatively, you can set these as actual system environment variables.

2.  **Example program code** (`examples/search_profiles/main.go`):

    ```go
    package main

    import (
    	"context"
    	"encoding/json"
    	"fmt"
    	"log"
    	"os"

    	linkedinscraper "github.com/masa-finance/linkedin-scraper"
    	"github.com/joho/godotenv"
    )

    func main() {
    	// Load .env file. Errors are not fatal for flexibility
    	err := godotenv.Load() 
    	if err != nil {
    		log.Println("No .env file found or error loading it, relying on environment variables: ", err)
    	}

    	fmt.Println("LinkedIn Profile Scraper Example")

    	liAtCookie := os.Getenv("LI_AT_COOKIE")
    	csrfToken := os.Getenv("CSRF_TOKEN")
    	jsessionID := os.Getenv("JSESSIONID_TOKEN")

    	if liAtCookie == "" || csrfToken == "" {
    		log.Fatal("Error: LI_AT_COOKIE and CSRF_TOKEN environment variables must be set.")
    	}

    	auth := linkedinscraper.AuthCredentials{
    		LiAtCookie: liAtCookie,
    		CSRFToken:  csrfToken,
    		JSESSIONID: jsessionID,
    	}

    	cfg, err := linkedinscraper.NewConfig(auth)
    	if err != nil {
    		log.Fatalf("Error creating config: %v", err)
    	}

    	client, err := linkedinscraper.NewClient(cfg)
    	if err != nil {
    		log.Fatalf("Error creating client: %v", err)
    	}

    	fmt.Println("Searching for profiles with keyword 'software engineer'...")
    	searchArgs := linkedinscraper.ProfileSearchArgs{
    		Keywords:       "software engineer",
    		NetworkFilters: []string{"S"}, 
    		Start:          0,            
    		Count:          5,             
    	}

    	profiles, err := client.SearchProfiles(context.Background(), searchArgs)
    	if err != nil {
    		log.Fatalf("Error searching profiles: %v", err)
    	}

    	if len(profiles) == 0 {
    		fmt.Println("No profiles found or returned.")
    	} else {
    		fmt.Printf("Found %d profile(s):\n", len(profiles))
    		profilesJSON, err := json.MarshalIndent(profiles, "", "  ")
    		if err != nil {
    			log.Fatalf("Error marshalling profiles to JSON: %v", err)
    		}
    		fmt.Println(string(profilesJSON))
    	}

    	fmt.Println("\nExample finished.")
    }
    ```

This will search for profiles matching "software engineer", attempt to load credentials from `.env` or environment variables, and print the found profiles as JSON.

## Contributing

(Details to be added later - e.g., how to run tests, coding standards, etc.)

## License

This project is licensed under the [MIT License](LICENSE).