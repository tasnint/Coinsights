package scrapers

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tasnint/coinsights/internal/models"
)

// GoogleScraper handles Google search scraping
type GoogleScraper struct {
	Collector *colly.Collector
	Delay     time.Duration
}

// NewGoogleScraper creates a new Google scraper instance
func NewGoogleScraper() *GoogleScraper {
	c := colly.NewCollector(
		colly.AllowedDomains("www.google.com", "google.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	// Rate limiting
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*google.*",
		Delay:       2 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	return &GoogleScraper{
		Collector: c,
		Delay:     2 * time.Second,
	}
}

// GoogleSearchResult holds a single Google search result (internal use)
type GoogleSearchResult struct {
	Title   string
	URL     string
	Snippet string
}

// Search performs a Google search and returns results
func (gs *GoogleScraper) Search(query string, maxResults int) ([]models.GoogleResult, error) {
	results := []models.GoogleResult{}

	// Clone collector for each search to avoid state issues
	c := gs.Collector.Clone()

	// Handle search result items
	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		if len(results) >= maxResults {
			return
		}

		title := e.ChildText("h3")
		link := e.ChildAttr("a", "href")
		snippet := e.ChildText("div.VwiC3b")

		// Filter out empty results and Google's own pages
		if title == "" || link == "" || strings.Contains(link, "google.com") {
			return
		}

		// Extract domain from URL
		domain := extractDomain(link)

		result := models.GoogleResult{
			Title:     title,
			URL:       link,
			Snippet:   snippet,
			Source:    domain,
			ScrapedAt: time.Now(),
		}
		results = append(results, result)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("⚠️  Google scraping error: %v\n", err)
	})

	// Build search URL
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s&num=%d",
		strings.ReplaceAll(query, " ", "+"),
		maxResults+10) // Request more to account for filtering

	fmt.Printf("🔍 Searching Google for: %s\n", query)

	err := c.Visit(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search Google: %w", err)
	}

	c.Wait()

	fmt.Printf("✅ Found %d Google results\n", len(results))
	return results, nil
}

// ScrapeAll searches Google for multiple queries
func (gs *GoogleScraper) ScrapeAll(queries []string, resultsPerQuery int) ([]models.GoogleResult, error) {
	allResults := []models.GoogleResult{}

	for _, query := range queries {
		results, err := gs.Search(query, resultsPerQuery)
		if err != nil {
			fmt.Printf("⚠️  Error searching for '%s': %v\n", query, err)
			continue
		}
		allResults = append(allResults, results...)

		// Be respectful with rate limiting
		time.Sleep(gs.Delay)
	}

	return allResults, nil
}

// extractDomain extracts the domain name from a URL
func extractDomain(urlStr string) string {
	// Simple extraction - remove protocol and path
	domain := urlStr

	// Remove protocol
	if strings.HasPrefix(domain, "https://") {
		domain = domain[8:]
	} else if strings.HasPrefix(domain, "http://") {
		domain = domain[7:]
	}

	// Remove path
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}

	// Remove www.
	domain = strings.TrimPrefix(domain, "www.")

	return domain
}
