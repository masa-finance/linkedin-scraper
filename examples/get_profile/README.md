# LinkedIn Profile Fetching Example

This example demonstrates how to use the LinkedIn Profile Scraper to fetch detailed profile information, including experience, education, skills, and more.

## Features Demonstrated

1. **Single Profile Fetching** - Fetch detailed profile data by public identifier
2. **Search + Profile Integration** - Search for profiles then fetch detailed data
3. **Comprehensive Data Display** - Show all profile fields in a formatted way
4. **JSON Export** - Export profile data to JSON files
5. **Error Handling** - Proper error handling and logging
6. **Rate Limiting** - Respectful delays between API calls

## Setup

### 1. Set Environment Variables

```bash
export LINKEDIN_LI_AT="your-li-at-cookie"
export LINKEDIN_CSRF_TOKEN="your-csrf-token"
export LINKEDIN_JSESSIONID="your-jsessionid"  # optional but recommended
```

### 2. Install Dependencies

```bash
cd examples/get_profile
go mod download
```

## Usage

### Basic Profile Fetching

Fetch a specific profile by public identifier:

```bash
# Use default example profile
go run main.go

# Fetch a specific profile
go run main.go john-doe-123456
```

### Output Example

The example will show comprehensive profile information:

```
🔍 LinkedIn Profile Scraper - Profile Fetching Example
====================================================

📋 Example 1: Fetching Profile by Public Identifier
---------------------------------------------------
✅ Profile fetched successfully!

👤 Basic Information:
  Name: John Doe
  Headline: Senior Software Engineer at Tech Company
  Public ID: john-doe-123456
  Profile URL: https://www.linkedin.com/in/john-doe-123456/
  URN: urn:li:fsd_profile:ACoAAABCDEF...

💼 Experience (3 entries):
  1. Senior Software Engineer at Tech Company
  2. Software Engineer at Previous Company
  3. Junior Developer at Startup
  ... and 0 more

🎓 Education (2 entries):
  1. Bachelor of Science at University
  2. High School Diploma at High School
  ... and 0 more

🛠️  Skills (15 total):
  • JavaScript (25 endorsements)
  • Python (18 endorsements)
  • React (12 endorsements)
  • Node.js (8 endorsements)
  • AWS (5 endorsements)
  ... and 10 more skills

📍 Location:
  Country: US
  Postal Code: 90210

🌐 Social Information:
  Connections: 500
  Followers: 1250
  Following: false

📊 Additional Information:
  Creator: false
  Verified: true
  Premium: false
  Memorialized: false

🔍 Example 2: Search + Profile Integration
------------------------------------------
✅ Found 3 profiles from search

📊 Fetching detailed data for profile 1: jane-smith-789
  1. Jane Smith - Software Engineer at Innovation Corp
     URL: https://www.linkedin.com/in/jane-smith-789/
     Experience: 2 entries, Education: 1 entries, Skills: 8 entries

💾 Example 3: Export Profile to JSON
------------------------------------
✅ Profile exported to john-doe-123456_profile.json

🎉 Profile fetching examples completed!
```

## Integration with Search

The example demonstrates how to combine search and profile fetching:

1. **Search for profiles** using keywords and filters
2. **Extract public identifiers** from search results
3. **Fetch detailed profile data** for each result
4. **Display comprehensive information** with proper formatting

## Profile Data Fields

The example displays all available profile fields:

### Basic Information
- Full Name, First Name, Last Name
- Headline and Summary
- Public Identifier and URN
- Profile URL

### Professional Information
- **Experience**: Job titles, companies, date ranges, descriptions
- **Education**: Schools, degrees, fields of study, activities
- **Skills**: Skill names and endorsement counts
- **Certifications**: Names, authorities, license numbers

### Personal Information
- **Location**: Country code, postal code, geo preferences
- **Profile Picture**: Display URLs and accessibility text
- **Contact**: Creator website and social links

### Social Information
- **Connections**: Connection count and mutual connections
- **Following**: Follower count and following status
- **Verification**: Verified status and creator badges

### Additional Metadata
- Creator status and premium features
- Memorialized status and temporary status
- Profile visibility and notification settings

## Error Handling

The example demonstrates proper error handling:

- **Configuration errors**: Missing environment variables
- **Authentication errors**: Invalid credentials or expired tokens
- **API errors**: Rate limiting, server errors, network issues
- **Parsing errors**: Malformed responses or unexpected data structures

## Rate Limiting

The example includes respectful rate limiting:

- **1-second delays** between profile fetching requests
- **Timeout handling** with 30-second context timeout
- **Error backoff** for failed requests

## Best Practices

1. **Set environment variables** for credentials (never hardcode)
2. **Use timeouts** for all API requests
3. **Handle errors gracefully** with proper logging
4. **Respect rate limits** with delays between requests
5. **Export data** for offline analysis when needed

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Verify your LinkedIn cookies are current
   - Check that all required environment variables are set

2. **Profile Not Found**
   - Ensure the public identifier is correct
   - Check that the profile is publicly accessible

3. **Rate Limited**
   - Increase delays between requests
   - Verify your LinkedIn account is in good standing

4. **Empty Fields**
   - Some fields may be private or not filled in by the user
   - Check the profile manually to verify expected data

### Debug Mode

To enable debug logging, modify the example to include debug output:

```go
// Add this before making API calls
log.Printf("Fetching profile: %s", publicIdentifier)
```

## Next Steps

After running this example, you might want to:

1. **Modify the search criteria** to find different types of profiles
2. **Add data persistence** to save profiles to a database
3. **Create batch processing** to handle large numbers of profiles
4. **Add data analysis** to extract insights from profile data
5. **Build a web interface** to visualize the data

## Related Examples

- `examples/search_profiles/` - Basic profile searching
- `examples/echo_api_example/` - Web API integration 