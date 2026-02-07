package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tasnint/coinsights/internal/models"
	"google.golang.org/genai"
)

// GeminiScraper uses Gemini AI with Google Search grounding to find complaints
type GeminiScraper struct {
	client *genai.Client
	apiKey string
}

// AIOverviewResult represents the structured output from Gemini
type AIOverviewResult struct {
	Query              string               `json:"query"`
	Summary            string               `json:"summary"`
	KeyComplaints      []ExtractedComplaint `json:"key_complaints"`
	Sources            []SourceReference    `json:"sources"`
	SentimentBreakdown SentimentStats       `json:"sentiment_breakdown"`
	GeneratedAt        time.Time            `json:"generated_at"`
}

// ExtractedComplaint represents a complaint extracted by Gemini
type ExtractedComplaint struct {
	Category    string `json:"category"`    // e.g., "fees", "support", "security"
	Description string `json:"description"` // The complaint text
	Frequency   string `json:"frequency"`   // "common", "occasional", "rare"
	Platform    string `json:"platform"`    // Where it was found (reddit, twitter, etc.)
}

// SourceReference tracks where information came from
type SourceReference struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// SentimentStats holds sentiment analysis stats (using float64 for percentages)
type SentimentStats struct {
	Negative float64 `json:"negative"`
	Neutral  float64 `json:"neutral"`
	Positive float64 `json:"positive"`
}

// NewGeminiScraper creates a new Gemini-powered scraper
func NewGeminiScraper() (*GeminiScraper, error) {
	// Check for API key in environment (GEMINI_API_KEY or GOOGLE_API_KEY)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GOOGLE_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY or GOOGLE_API_KEY environment variable not set")
	}

	// Set the env var so the SDK can find it
	os.Setenv("GOOGLE_API_KEY", apiKey)

	ctx := context.Background()

	// Create client - passing nil uses GOOGLE_API_KEY from environment
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &GeminiScraper{
		client: client,
		apiKey: apiKey,
	}, nil
}

// Close closes the Gemini client (no-op for new SDK but kept for interface compatibility)
func (gs *GeminiScraper) Close() {
	// New SDK doesn't require explicit close
}

// SearchComplaintsWithAI searches for complaints using Gemini with Google Search grounding
func (gs *GeminiScraper) SearchComplaintsWithAI(ctx context.Context, query string) (*AIOverviewResult, error) {
	fmt.Printf("🤖 Searching with Gemini AI: %s\n", query)

	prompt := fmt.Sprintf(`You are a research assistant analyzing user complaints about cryptocurrency platforms.

Search the web for: "%s"

Analyze the search results and provide a comprehensive analysis in the following JSON format:
{
	"query": "%s",
	"summary": "A 2-3 sentence summary of the main complaints found",
	"key_complaints": [
		{
			"category": "category name (fees, customer_support, security, account_issues, withdrawal_problems, verification, app_bugs, other)",
			"description": "Brief description of the complaint",
			"frequency": "common/occasional/rare",
			"platform": "where this was found (reddit, twitter, trustpilot, bbb, etc.)"
		}
	],
	"sources": [
		{
			"title": "Page title",
			"url": "Full URL",
			"domain": "domain.com"
		}
	],
	"sentiment_breakdown": {
		"negative": 0,
		"neutral": 0,
		"positive": 0
	}
}

Focus on:
1. Recent complaints (within the last year if possible)
2. Common recurring issues
3. Reddit discussions, review sites, social media, and forums
4. Be objective and factual

Return ONLY valid JSON, no markdown code blocks or explanation.`, query, query)

	// Use the new SDK API with Google Search tool for grounding
	// Model: gemini-2.0-flash is recommended for speed
	modelName := "gemini-2.0-flash"

	// Create config with Google Search tool enabled
	// NOTE: Cannot use ResponseMIMEType with Google Search tool
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{GoogleSearch: &genai.GoogleSearch{}},
		},
	}

	result, err := gs.client.Models.GenerateContent(
		ctx,
		modelName,
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	// Extract text from response using the new SDK's Text() method
	responseText := result.Text()
	if responseText == "" {
		return nil, fmt.Errorf("no response from Gemini")
	}

	// Clean up the response - remove markdown code blocks if present
	responseText = cleanJSONResponse(responseText)

	// Parse the JSON response
	var aiResult AIOverviewResult
	if err := json.Unmarshal([]byte(responseText), &aiResult); err != nil {
		// If JSON parsing fails, return raw response as summary
		fmt.Printf("⚠️  JSON parsing failed, raw response: %s\n", responseText)
		return &AIOverviewResult{
			Query:       query,
			Summary:     responseText,
			GeneratedAt: time.Now(),
		}, nil
	}

	aiResult.GeneratedAt = time.Now()
	fmt.Printf("✅ Gemini found %d key complaints from %d sources\n",
		len(aiResult.KeyComplaints), len(aiResult.Sources))

	return &aiResult, nil
}

// SearchMultipleQueries searches for multiple queries and aggregates results
func (gs *GeminiScraper) SearchMultipleQueries(ctx context.Context, queries []string) ([]AIOverviewResult, error) {
	results := []AIOverviewResult{}

	for i, query := range queries {
		// Retry logic for rate limiting
		var result *AIOverviewResult
		var err error
		maxRetries := 3

		for retry := 0; retry < maxRetries; retry++ {
			result, err = gs.SearchComplaintsWithAI(ctx, query)
			if err == nil {
				break
			}

			// Check if it's a rate limit error
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "RESOURCE_EXHAUSTED") {
				waitTime := time.Duration((retry+1)*30) * time.Second
				fmt.Printf("Rate limited, waiting %v before retry %d/%d...\n", waitTime, retry+1, maxRetries)
				time.Sleep(waitTime)
			} else {
				break // Non-rate-limit error, don't retry
			}
		}

		if err != nil {
			fmt.Printf("⚠️  Error searching '%s': %v\n", query, err)
			continue
		}
		results = append(results, *result)

		// Rate limiting between queries (10 seconds to avoid 429 errors)
		if i < len(queries)-1 {
			fmt.Println("⏳ Waiting 10 seconds before next query...")
			time.Sleep(10 * time.Second)
		}
	}

	return results, nil
}

// ConvertToComplaints converts AIOverviewResults to standard Complaint models
func ConvertToComplaints(aiResults []AIOverviewResult) []models.Complaint {
	complaints := []models.Complaint{}

	for _, result := range aiResults {
		for i, kc := range result.KeyComplaints {
			complaint := models.Complaint{
				ID:          fmt.Sprintf("gemini-%s-%d", result.GeneratedAt.Format("20060102150405"), i),
				Source:      fmt.Sprintf("gemini_search:%s", kc.Platform),
				Title:       fmt.Sprintf("[%s] %s", kc.Category, truncateString(kc.Description, 50)),
				Description: kc.Description,
				Category:    kc.Category,
				Sentiment:   "negative", // Complaints are inherently negative
				ScrapedAt:   result.GeneratedAt,
			}

			// Add URL if available from sources
			if len(result.Sources) > 0 {
				complaint.URL = result.Sources[0].URL
			}

			complaints = append(complaints, complaint)
		}
	}

	return complaints
}

// GetDefaultComplaintQueries returns default queries for Coinbase complaints
func GetDefaultComplaintQueries() []string {
	return []string{
		"user complaints regarding coinbase",
		"coinbase customer complaints reddit",
		"coinbase problems issues 2024 2025",
		"coinbase account locked complaints",
		"coinbase fees too high complaints",
		"coinbase customer support terrible reddit",
		"coinbase withdrawal problems",
		"coinbase verification issues complaints",
	}
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// cleanJSONResponse removes markdown code blocks from Gemini response
func cleanJSONResponse(response string) string {
	// Remove ```json and ``` markers if present
	response = strings.TrimSpace(response)

	// Remove leading ```json or ```
	if strings.HasPrefix(response, "```json") {
		response = strings.TrimPrefix(response, "```json")
	} else if strings.HasPrefix(response, "```") {
		response = strings.TrimPrefix(response, "```")
	}

	// Remove trailing ```
	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	return strings.TrimSpace(response)
}
