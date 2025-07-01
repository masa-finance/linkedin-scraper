# linkedin-scraper

A Go package for searching LinkedIn profiles and fetching detailed profile data using the Voyager API.

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
    		// NetworkFilters is optional. LinkedIn's default is to search all network degrees.
    		// Valid filters: "F" (1st degree), "S" (2nd degree), "O" (Outside of Your Network).
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

### Get a Specific Profile

This example shows how to fetch detailed information for a single profile using its public identifier (the part of their profile URL, e.g., `williamhgates` from `https://www.linkedin.com/in/williamhgates/`).

This example can be found in `examples/get_profile/main.go`.

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
	// Load .env file
	_ = godotenv.Load() 

	fmt.Println("LinkedIn Profile Scraper - Fetch Profile Example")

	// Get credentials from environment
	liAtCookie := os.Getenv("LI_AT_COOKIE")
	csrfToken := os.Getenv("CSRF_TOKEN")
	jsessionID := os.Getenv("JSESSIONID_TOKEN")

	if liAtCookie == "" || csrfToken == "" {
		log.Fatal("Error: LI_AT_COOKIE and CSRF_TOKEN must be set.")
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

	// The public identifier for the profile you want to fetch
	publicIdentifier := "williamhgates"

	fmt.Printf("Fetching profile for: %s\n", publicIdentifier)
	profile, err := client.GetProfile(context.Background(), publicIdentifier)
	if err != nil {
		log.Fatalf("Error fetching profile: %v", err)
	}

	fmt.Println("Profile data fetched successfully!")
	profileJSON, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling profile to JSON: %v", err)
	}
	fmt.Println(string(profileJSON))
}
```

This will fetch the full profile for Bill Gates and print it as a JSON object.

### Available Profile Data

When using `GetProfile`, the returned `LinkedInProfile` struct is populated with rich data, including:
-   **Personal Info**: Full Name, Headline, Location, Summary, Profile Picture URL.
-   **Work Experience**: A list of positions including Company Name, Title, Date Range, and Description.
-   **Education**: A list of educational institutions attended, including School Name, Degree, and Field of Study.
-   **Skills**: A list of skills with endorsement counts.
-   **Connections**: Follower and connection counts.
-   **And more**: Industry, Certifications, etc.

## Echo API Example

This project includes a more advanced example demonstrating how to use the `linkedinscraper` package within a web API built with the [Echo framework](https://echo.labstack.com/).

The example API is located in `examples/echo_api_example/`.

### Running the Echo API Example

1.  **Set up your credentials for the example:**
    *   Navigate to the example directory:
        ```bash
        cd examples/echo_api_example
        ```
    *   Create a `.env` file in this directory by copying the main project's `.env.example` or creating a new one:
        ```env
        # examples/echo_api_example/.env
        LI_AT_COOKIE="YOUR_LI_AT_COOKIE_VALUE"
        CSRF_TOKEN="YOUR_CSRF_TOKEN_VALUE"
        JSESSIONID_TOKEN="ajax:YOUR_JSESSIONID_VALUE"
        ```
    *   Replace the placeholder values with your actual LinkedIn credentials.

2.  **Run the API server:**
    *   Ensure you are in the `examples/echo_api_example` directory.
    *   Execute the following command:
        ```bash
        go run main.go
        ```
    *   The server will start, typically on port `1323`. You should see log output indicating it's running.

3.  **Test the API:**
    *   Once the server is running, you can query the `/search/linkedin` endpoint using a tool like `curl` or your web browser.
    *   Example using `curl` (from a new terminal window):
        ```bash
        curl "http://localhost:1323/search/linkedin?keywords=golang+developer"
        ```
    *   This will send a request to the API to search for LinkedIn profiles matching "golang developer" and return the results as JSON. The example is configured to return up to 5 profiles.

This Echo example demonstrates how to integrate the scraper into a service that can be called over HTTP.

## Contributing

(Details to be added later - e.g., how to run tests, coding standards, etc.)

## License

This project is licensed under the [MIT License](LICENSE).