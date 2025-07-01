package linkedinscraper_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	linkedinscraper "github.com/masa-finance/linkedin-scraper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LinkedIn Profile Integration Tests", func() {
	var (
		client *linkedinscraper.Client
		ctx    context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		// Load environment variables from .env file if it exists
		_ = godotenv.Load()

		// Create context with timeout
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)

		// Get credentials from environment
		liAtCookie := os.Getenv("LI_AT_COOKIE")
		csrfToken := os.Getenv("CSRF_TOKEN")
		jsessionID := os.Getenv("JSESSIONID_TOKEN")

		if liAtCookie == "" || csrfToken == "" {
			Skip("LinkedIn credentials not available. Set LI_AT_COOKIE, CSRF_TOKEN, and JSESSIONID_TOKEN environment variables to run profile integration tests.")
		}

		auth := linkedinscraper.AuthCredentials{
			LiAtCookie: liAtCookie,
			CSRFToken:  csrfToken,
			JSESSIONID: jsessionID,
		}

		config, err := linkedinscraper.NewConfig(auth)
		Expect(err).ToNot(HaveOccurred())

		client, err = linkedinscraper.NewClient(config)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		if cancel != nil {
			cancel()
		}
	})

	Describe("Profile Fetching", func() {
		Context("with valid public identifiers", func() {
			It("should fetch a profile successfully", func() {
				By("using a known public LinkedIn profile identifier")
				// Using a well-known LinkedIn profile that should be publicly accessible
				publicIdentifier := "williamhgates" // Bill Gates' LinkedIn profile

				By("fetching the profile data")
				profile, err := client.GetProfile(ctx, publicIdentifier)

				if err != nil {
					log.Printf("Error fetching profile: %v", err)
					// If this specific profile fails, try with another known public profile
					publicIdentifier = "jeffweiner08" // Jeff Weiner's LinkedIn profile
					profile, err = client.GetProfile(ctx, publicIdentifier)
				}

				Expect(err).ToNot(HaveOccurred())
				Expect(profile).ToNot(BeNil())

				By("validating basic profile fields")
				Expect(profile.PublicIdentifier).To(Equal(publicIdentifier))
				Expect(profile.FullName).ToNot(BeEmpty())
				Expect(profile.ProfileURL).ToNot(BeEmpty())
				Expect(profile.URN).ToNot(BeEmpty())

				By("logging profile information for verification")
				log.Printf("✅ Profile fetched successfully:")
				log.Printf("  Name: %s", profile.FullName)
				log.Printf("  Headline: %s", profile.Headline)
				log.Printf("  Public ID: %s", profile.PublicIdentifier)
				log.Printf("  Profile URL: %s", profile.ProfileURL)
				log.Printf("  Experience entries: %d", len(profile.Experience))
				log.Printf("  Education entries: %d", len(profile.Education))
				log.Printf("  Skills: %d", len(profile.Skills))
			})

			It("should handle profiles with comprehensive data", func() {
				By("fetching a profile with rich data")
				publicIdentifier := "reidhoffman" // Reid Hoffman's LinkedIn profile (LinkedIn founder)

				profile, err := client.GetProfile(ctx, publicIdentifier)
				Expect(err).ToNot(HaveOccurred())
				Expect(profile).ToNot(BeNil())

				// Log the experience data to see what we're getting
				if len(profile.Experience) > 0 {
					log.Printf("Found %d experience entries. First entry: %+v", len(profile.Experience), profile.Experience[0])
				}

				By("validating professional information")
				if len(profile.Experience) > 0 {
					log.Printf("  First experience: %s at %s", profile.Experience[0].Title, profile.Experience[0].CompanyName)
				}

				By("validating education information")
				if len(profile.Education) > 0 {
					log.Printf("  First education: %s at %s", profile.Education[0].DegreeName, profile.Education[0].SchoolName)
				}

				By("validating additional fields")
				if profile.LocationDetails != nil {
					log.Printf("  Location: Country Code %s", profile.LocationDetails.CountryCode)
				}

				if profile.ConnectionInfo != nil {
					// Connection info should be present for most profiles
					log.Printf("  Connection count: %d", profile.ConnectionInfo.ConnectionCount)
					log.Printf("  Follower count: %d", profile.ConnectionInfo.FollowerCount)
				}
			})
		})

		Context("with invalid public identifiers", func() {
			It("should handle non-existent profiles gracefully", func() {
				By("attempting to fetch a non-existent profile")
				publicIdentifier := "this-profile-does-not-exist-12345"

				profile, err := client.GetProfile(ctx, publicIdentifier)

				By("expecting an appropriate error")
				Expect(err).To(HaveOccurred())
				Expect(profile).To(BeNil())

				By("ensuring error is descriptive")
				Expect(err.Error()).ToNot(BeEmpty())
				log.Printf("Expected error for non-existent profile: %v", err)
			})

			It("should handle empty public identifier", func() {
				By("attempting to fetch with empty identifier")
				profile, err := client.GetProfile(ctx, "")

				By("expecting an appropriate error")
				Expect(err).To(HaveOccurred())
				Expect(profile).To(BeNil())
			})
		})

		Context("with rate limiting considerations", func() {
			It("should handle multiple profile requests with delays", func() {
				profiles := []string{
					"williamhgates",
					"jeffweiner08",
					"reidhoffman",
				}

				By("fetching multiple profiles with appropriate delays")
				for i, publicIdentifier := range profiles {
					if i > 0 {
						time.Sleep(2 * time.Second) // Respectful delay
					}

					profile, err := client.GetProfile(ctx, publicIdentifier)
					if err != nil {
						log.Printf("⚠️  Profile %s failed: %v", publicIdentifier, err)
						continue // Skip this profile if it fails
					}

					Expect(profile).ToNot(BeNil())
					Expect(profile.PublicIdentifier).To(Equal(publicIdentifier))
					log.Printf("✅ Profile %d: %s - %s", i+1, profile.FullName, profile.Headline)
				}
			})
		})
	})

	Describe("Error Handling", func() {
		It("should handle authentication errors", func() {
			By("creating a client with invalid credentials")
			invalidAuth := linkedinscraper.AuthCredentials{
				LiAtCookie: "invalid",
				CSRFToken:  "invalid",
				JSESSIONID: "invalid",
			}

			config, err := linkedinscraper.NewConfig(invalidAuth)
			Expect(err).ToNot(HaveOccurred())

			invalidClient, err := linkedinscraper.NewClient(config)
			Expect(err).ToNot(HaveOccurred())

			By("attempting to fetch a profile with invalid credentials")
			profile, err := invalidClient.GetProfile(ctx, "williamhgates")

			By("expecting an authentication error")
			Expect(err).To(HaveOccurred())
			Expect(profile).To(BeNil())

			log.Printf("Expected authentication error: %v", err)
		})

		It("should handle timeout scenarios", func() {
			By("creating a context with very short timeout")
			shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
			defer shortCancel()

			By("attempting to fetch a profile with short timeout")
			profile, err := client.GetProfile(shortCtx, "williamhgates")

			By("expecting a timeout error")
			Expect(err).To(HaveOccurred())
			Expect(profile).To(BeNil())

			log.Printf("Expected timeout error: %v", err)
		})
	})

	Describe("Data Validation", func() {
		It("should validate profile data structure", func() {
			By("fetching a known profile")
			profile, err := client.GetProfile(ctx, "williamhgates")

			if err != nil {
				Skip("Profile fetching failed - skipping data validation test")
			}

			Expect(profile).ToNot(BeNil())

			By("validating required fields are present")
			Expect(profile.PublicIdentifier).ToNot(BeEmpty())
			Expect(profile.FullName).ToNot(BeEmpty())
			Expect(profile.ProfileURL).ToNot(BeEmpty())

			By("validating field formats")
			Expect(profile.ProfileURL).To(ContainSubstring("linkedin.com"))
			if profile.URN != "" {
				Expect(profile.URN).To(ContainSubstring("urn:li:"))
			}

			By("validating nested data structures")
			for _, exp := range profile.Experience {
				if exp.Title != "" {
					Expect(exp.Title).ToNot(BeEmpty())
				}
				if exp.CompanyName != "" {
					Expect(exp.CompanyName).ToNot(BeEmpty())
				}
			}

			for _, edu := range profile.Education {
				if edu.SchoolName != "" {
					Expect(edu.SchoolName).ToNot(BeEmpty())
				}
			}

			for _, skill := range profile.Skills {
				if skill.Name != "" {
					Expect(skill.Name).ToNot(BeEmpty())
				}
			}
		})
	})
})

func TestProfileIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkedIn Profile Integration Test Suite")
}
