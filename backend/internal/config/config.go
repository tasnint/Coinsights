package config

// ================================================
// COINSIGHTS SCRAPER CONFIGURATION
// ================================================
// Modify these values to customize scraping behavior
// ================================================

// SearchQueries - All YouTube search queries to find Coinbase complaints
// Add, remove, or modify queries here!
var SearchQueries = []string{
	// ============================================
	// DIRECT COMPLAINT SEARCHES
	// ============================================
	"coinbase problems",
	"coinbase complaints",
	"coinbase issues",
	"coinbase bad experience",
	"coinbase terrible",
	"coinbase worst",

	// ============================================
	// CONS / DISADVANTAGES / REVIEWS
	// ============================================
	"disadvantages of coinbase",
	"cons of using coinbase",
	"why coinbase is bad",
	"coinbase review negative",
	"coinbase honest review",
	"coinbase review 2024",
	"coinbase review 2025",
	"should i use coinbase",
	"coinbase vs competitors",

	// ============================================
	// SPECIFIC PAIN POINTS - FEES
	// ============================================
	"coinbase fees too high",
	"coinbase fees explained",
	"coinbase hidden fees",
	"coinbase expensive",

	// ============================================
	// SPECIFIC PAIN POINTS - CUSTOMER SERVICE
	// ============================================
	"coinbase customer support",
	"coinbase customer service bad",
	"coinbase no response",
	"coinbase support nightmare",

	// ============================================
	// SPECIFIC PAIN POINTS - ACCOUNT ISSUES
	// ============================================
	"coinbase account locked",
	"coinbase account restricted",
	"coinbase account closed",
	"coinbase verification problems",
	"coinbase identity verification failed",

	// ============================================
	// SPECIFIC PAIN POINTS - SECURITY/TRUST
	// ============================================
	"coinbase scam",
	"coinbase security issues",
	"coinbase hacked",
	"coinbase lost money",
	"coinbase funds missing",

	// ============================================
	// SPECIFIC PAIN POINTS - WITHDRAWALS
	// ============================================
	"coinbase withdrawal problems",
	"coinbase cant withdraw",
	"coinbase withdrawal delay",
	"coinbase bank transfer issues",

	// ============================================
	// COMPARISONS (often highlight cons)
	// ============================================
	"coinbase vs kraken",
	"coinbase vs binance",
	"coinbase vs crypto.com",
	"why i left coinbase",
	"coinbase alternatives",
}

// ScraperSettings configures how much data to fetch
type ScraperSettings struct {
	VideosPerQuery   int // Number of videos to fetch per search query
	CommentsPerVideo int // Number of comments to fetch per video
	MaxQueries       int // Max number of queries to run (0 = all)
}

// DefaultSettings returns the default scraper configuration
// Calculated for ~5000 quota units/day:
// - 25 queries × 100 units = 2,500 (search)
// - 25 queries × 1 unit = 25 (videos.list batched)
// - 125 videos × 20 comments × 1 unit = 125 (commentThreads)
// Total: ~2,650 units (leaves room for retries)
func DefaultSettings() ScraperSettings {
	return ScraperSettings{
		VideosPerQuery:   5,  // 5 videos per query
		CommentsPerVideo: 20, // 20 comments per video
		MaxQueries:       25, // Run first 25 queries (out of 30+ available)
	}
}

// AggressiveSettings for maximum data collection (~5000 units)
func AggressiveSettings() ScraperSettings {
	return ScraperSettings{
		VideosPerQuery:   5,  // 5 videos per query
		CommentsPerVideo: 25, // 25 comments per video
		MaxQueries:       40, // Run 40 queries
	}
}

// LightSettings for testing or preserving quota
func LightSettings() ScraperSettings {
	return ScraperSettings{
		VideosPerQuery:   3,  // 3 videos per query
		CommentsPerVideo: 10, // 10 comments per video
		MaxQueries:       5,  // Only 5 queries
	}
}

// CalculateQuota estimates API quota usage
func (s ScraperSettings) CalculateQuota() int {
	queries := s.MaxQueries
	if queries == 0 || queries > len(SearchQueries) {
		queries = len(SearchQueries)
	}

	searchUnits := queries * 100                   // search.list = 100 units each
	videoUnits := queries * 1                      // videos.list = 1 unit (batched per query)
	commentUnits := queries * s.VideosPerQuery * 1 // commentThreads = 1 unit each

	return searchUnits + videoUnits + commentUnits
}
