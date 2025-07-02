package linkedinscraper

// This package provides LinkedIn profile scraping functionality.
//
// Example usage:
//   // Create config with auth credentials
//   auth := AuthCredentials{
//       LiAtCookie: "your-li-at-cookie",
//       CSRFToken:  "your-csrf-token",
//       JSESSIONID: "your-jsessionid",
//   }
//   cfg, err := NewConfig(auth)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // Create client
//   client, err := NewClient(cfg)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   // Fetch profile
//   profile, err := client.GetProfile(ctx, "john-doe-123456")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Profile: %s - %s\n", profile.FullName, profile.Headline)
