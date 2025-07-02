package linkedinscraper_test

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	linkedinscraper "github.com/masa-finance/linkedin-scraper"
)

var _ = Describe("SearchProfiles Integration", func() {
	var (
		client *linkedinscraper.Client
		auth   linkedinscraper.AuthCredentials
	)

	BeforeEach(func() {
		// Load .env file from the project root.
		// Adjust path if your test execution directory is different or .env is elsewhere.
		err := godotenv.Load("../.env") // Assuming tests might be run from a subdir or `go test ./...`
		if err != nil {
			// If .env is in the same directory as the test (e.g. when running `ginkgo` in the package dir)
			err = godotenv.Load()
		}
		if err != nil {
			log.Println("Warning: .env file not found or error loading it. Relying on actual environment variables.", err)
		}

		liAtCookie := os.Getenv("LI_AT_COOKIE")
		csrfToken := os.Getenv("CSRF_TOKEN")
		jsessionID := os.Getenv("JSESSIONID_TOKEN")

		Expect(liAtCookie).NotTo(BeEmpty(), "LI_AT_COOKIE must be set in .env or environment variables for integration tests")
		Expect(csrfToken).NotTo(BeEmpty(), "CSRF_TOKEN must be set in .env or environment variables for integration tests")
		// JSESSIONID is sometimes optional, so we don't make its absence fatal here.

		auth = linkedinscraper.AuthCredentials{
			LiAtCookie: liAtCookie,
			CSRFToken:  csrfToken,
			JSESSIONID: jsessionID,
		}

		cfg, configErr := linkedinscraper.NewConfig(auth)
		Expect(configErr).NotTo(HaveOccurred())
		Expect(cfg).NotTo(BeNil())

		c, clientErr := linkedinscraper.NewClient(cfg)
		Expect(clientErr).NotTo(HaveOccurred())
		Expect(c).NotTo(BeNil())
		client = c
	})

	Context("when searching for 'investor' profiles", func() {
		It("should successfully retrieve and parse profiles", func() {
			// Skip("Skipping real API call test until explicitly enabled or configured. Ensure .env is set up.") // Intentionally kept enabled from previous step
			searchArgs := linkedinscraper.ProfileSearchArgs{
				Keywords:        "investor",
				NetworkFilters:  []string{"F", "O"}, // Match cURL
				Start:           0,
				Count:           1,                                                                                                                                                                                                                                             // Requesting a small number for a quick test
				XLiPageInstance: "urn:li:page:d_flagship3_search_srp_people;e+uzo8jcR7GbUNfmjxNoSw==",                                                                                                                                                                          // From cURL
				XLiTrack:        `{"clientVersion":"1.13.35368","mpVersion":"1.13.35368","osName":"web","timezoneOffset":-7,"timezone":"America/Los_Angeles","deviceFormFactor":"DESKTOP","mpName":"voyager-web","displayDensity":2,"displayWidth":5120,"displayHeight":2880}`, // From cURL, requires backticks for JSON string
			}

			profiles, err := client.SearchProfiles(context.Background(), searchArgs)

			By("checking for request errors")
			Expect(err).NotTo(HaveOccurred())

			By("checking if profiles slice is populated")
			Expect(profiles).NotTo(BeEmpty(), "Expected to find at least one investor profile")

			By("inspecting the first profile's basic details")
			if len(profiles) > 0 {
				firstProfile := profiles[0]
				Expect(firstProfile.FullName).NotTo(BeEmpty(), "First profile should have a FullName")
				Expect(firstProfile.URN).NotTo(BeEmpty(), "First profile should have a URN")
				Expect(firstProfile.Headline).NotTo(BeEmpty(), "First profile should have a Headline")
				// ProfileURL might sometimes be empty depending on the API response details for a given profile
				// Expect(firstProfile.ProfileURL).NotTo(BeEmpty(), "First profile should have a ProfileURL")
			}
		})

		It("should successfully retrieve profiles without a network filter", func() {
			searchArgs := linkedinscraper.ProfileSearchArgs{
				Keywords:        "investor",
				Start:           0,
				Count:           1,
				XLiPageInstance: "urn:li:page:d_flagship3_search_srp_people;e+uzo8jcR7GbUNfmjxNoSw==",
				XLiTrack:        `{"clientVersion":"1.13.35368","mpVersion":"1.13.35368","osName":"web","timezoneOffset":-7,"timezone":"America/Los_Angeles","deviceFormFactor":"DESKTOP","mpName":"voyager-web","displayDensity":2,"displayWidth":5120,"displayHeight":2880}`,
			}

			profiles, err := client.SearchProfiles(context.Background(), searchArgs)

			By("checking for request errors")
			Expect(err).NotTo(HaveOccurred())

			By("checking if profiles slice is populated")
			Expect(profiles).NotTo(BeEmpty(), "Expected to find at least one investor profile")

			By("inspecting the first profile's basic details")
			if len(profiles) > 0 {
				firstProfile := profiles[0]
				Expect(firstProfile.FullName).NotTo(BeEmpty(), "First profile should have a FullName")
				Expect(firstProfile.URN).NotTo(BeEmpty(), "First profile should have a URN")
				Expect(firstProfile.Headline).NotTo(BeEmpty(), "First profile should have a Headline")
			}
		})
	})
})
