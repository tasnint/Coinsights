package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/tasnint/coinsights/internal/config"
	"github.com/tasnint/coinsights/internal/models"
	"github.com/tasnint/coinsights/internal/scrapers"
)

func main() {
	// Load environment variables - try multiple paths
	envPaths := []string{
		"../../.env", // From cmd/server/
		".env",       // From current dir
		"c:/Users/tanis/Downloads/GitHub Repos/Coinsights/.env", // Absolute path
	}

	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	if youtubeAPIKey == "" || youtubeAPIKey == "your_youtube_api_key_here" {
		log.Fatal("❌ YOUTUBE_API_KEY not set in .env file")
	}

	fmt.Println("🚀 Coinsights YouTube Scraper Starting...")
	fmt.Println("==========================================")

	// ================================================
	// CONFIGURATION - Edit in config/config.go
	// ================================================
	settings := config.DefaultSettings() // Or use config.AggressiveSettings() or config.LightSettings()
	queries := config.SearchQueries

	// Limit queries if MaxQueries is set
	if settings.MaxQueries > 0 && settings.MaxQueries < len(queries) {
		queries = queries[:settings.MaxQueries]
	}

	// Show configuration
	fmt.Println("\n⚙️  CONFIGURATION")
	fmt.Println("-----------------")
	fmt.Printf("📋 Total queries available: %d\n", len(config.SearchQueries))
	fmt.Printf("🔎 Queries to run:          %d\n", len(queries))
	fmt.Printf("📺 Videos per query:        %d\n", settings.VideosPerQuery)
	fmt.Printf("💬 Comments per video:      %d\n", settings.CommentsPerVideo)
	fmt.Printf("💰 Estimated quota usage:   ~%d units (out of 10,000/day)\n", settings.CalculateQuota())

	// Show queries being used
	fmt.Println("\n🔍 SEARCH QUERIES")
	fmt.Println("-----------------")
	for i, q := range queries {
		fmt.Printf("   %2d. %s\n", i+1, q)
	}

	// ========================================
	// YOUTUBE SCRAPING (Commented out to save quota while testing Gemini)
	// ========================================
	/*
	// Initialize YouTube scraper
	youtubeScraper := scrapers.NewYouTubeScraper(youtubeAPIKey)

	// Scrape YouTube
	fmt.Println("\n📺 SCRAPING YOUTUBE...")
	fmt.Println("----------------------")
	result, err := youtubeScraper.ScrapeAll(queries, settings.VideosPerQuery, settings.CommentsPerVideo)
	if err != nil {
		log.Printf("YouTube scraping error: %v", err)
	}

	// Save YouTube results to JSON file
	fmt.Println("\n💾 SAVING YOUTUBE RESULTS...")
	fmt.Println("--------------------")
	err = saveResults(result)
	if err != nil {
		log.Printf("Error saving results: %v", err)
	}

	// Print YouTube summary
	printSummary(result)
	*/
	fmt.Println("\n📺 YOUTUBE SCRAPING: Skipped (commented out to save quota)")

	// ========================================
	// GEMINI AI SEARCH (Google AI Overview)
	// ========================================
	fmt.Println("\n🤖 GEMINI AI SEARCH...")
	fmt.Println("----------------------")

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Println("⚠️  GEMINI_API_KEY not set, skipping AI search")
	} else {
		geminiScraper, err := scrapers.NewGeminiScraper()
		if err != nil {
			log.Printf("❌ Failed to create Gemini scraper: %v", err)
		} else {
			defer geminiScraper.Close()

			// Define AI search queries for Coinbase complaints from different sources
			aiQueries := []string{
				// Query 1: Reddit-focused complaints
				"coinbase user complaints and problems from reddit discussions 2024 2025",
				// Query 2: Article/website reviews and complaints
				"coinbase customer complaints reviews from news articles trustpilot bbb consumer reports",
				// Query 3: YouTube video content analysis (not comments)
				"coinbase review video analysis problems issues discussed by youtubers crypto reviewers",
			}

			ctx := context.Background()
			aiResults, err := geminiScraper.SearchMultipleQueries(ctx, aiQueries)
			if err != nil {
				log.Printf("⚠️  Gemini search error: %v", err)
			} else {
				// Save AI results
				err = saveAIResults(aiResults)
				if err != nil {
					log.Printf("Error saving AI results: %v", err)
				}

				// Print AI summary
				printAISummary(aiResults)
			}
		}
	}

	fmt.Println("\n✅ All scraping complete!")
}

func saveResults(result *models.ScrapeResult) error {
	// Create data directory if it doesn't exist
	dataDir := "../../data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save to single file: youtube_latest_results.json
	filename := filepath.Join(dataDir, "youtube_latest_results.json")

	// Marshal to JSON
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ YouTube results saved to: %s\n", filename)

	return nil
}

func printSummary(result *models.ScrapeResult) {
	fmt.Println("\n📊 SCRAPE SUMMARY")
	fmt.Println("=================")
	fmt.Printf("📺 YouTube Videos:   %d\n", len(result.Videos))
	fmt.Printf("💬 YouTube Comments: %d\n", len(result.Comments))
	fmt.Printf("⏰ Scraped at:       %s\n", result.ScrapedAt.Format("2006-01-02 15:04:05"))

	// Calculate total views and engagement
	var totalViews, totalLikes int64
	for _, video := range result.Videos {
		totalViews += video.ViewCount
		totalLikes += video.LikeCount
	}
	fmt.Printf("👁️  Total Views:      %s\n", formatNumber(totalViews))
	fmt.Printf("👍 Total Likes:      %s\n", formatNumber(totalLikes))

	// Show sample results
	if len(result.Videos) > 0 {
		fmt.Println("\n📺 Sample YouTube Videos:")
		for i, video := range result.Videos {
			if i >= 3 {
				break
			}
			fmt.Printf("   %d. %s\n", i+1, video.Title)
			fmt.Printf("      Channel: %s | Views: %s | Likes: %s\n",
				video.ChannelTitle, formatNumber(video.ViewCount), formatNumber(video.LikeCount))
			fmt.Printf("      URL: %s\n", video.URL)
		}
	}

	if len(result.Comments) > 0 {
		fmt.Println("\n💬 Sample Comments:")
		for i, comment := range result.Comments {
			if i >= 3 {
				break
			}
			// Truncate long comments
			text := comment.Text
			if len(text) > 100 {
				text = text[:100] + "..."
			}
			fmt.Printf("   %d. %s: \"%s\"\n", i+1, comment.AuthorName, text)
		}
	}

	fmt.Println("\n✅ Scraping complete! Check the 'data' folder for full results.")
}

// formatNumber formats large numbers with K/M suffixes
func formatNumber(n int64) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
	if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

// saveAIResults saves Gemini AI search results to a JSON file
func saveAIResults(results []scrapers.AIOverviewResult) error {
	// Create data directory if it doesn't exist
	dataDir := "../../data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Save to single file: gemini_latest_results.json
	filename := filepath.Join(dataDir, "gemini_latest_results.json")

	// Marshal to JSON
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Gemini results saved to: %s\n", filename)

	return nil
}

// printAISummary prints a summary of AI search results
func printAISummary(results []scrapers.AIOverviewResult) {
	fmt.Println("\n🤖 GEMINI AI SEARCH SUMMARY")
	fmt.Println("============================")
	fmt.Printf("📊 Queries processed: %d\n", len(results))

	totalComplaints := 0
	totalSources := 0
	for _, result := range results {
		totalComplaints += len(result.KeyComplaints)
		totalSources += len(result.Sources)
	}
	fmt.Printf("🔍 Key complaints found: %d\n", totalComplaints)
	fmt.Printf("📚 Sources referenced: %d\n", totalSources)

	// Show summaries for each query
	for i, result := range results {
		fmt.Printf("\n📌 Query %d: \"%s\"\n", i+1, result.Query)
		if result.Summary != "" {
			summary := result.Summary
			if len(summary) > 300 {
				summary = summary[:300] + "..."
			}
			fmt.Printf("   Summary: %s\n", summary)
		}

		if len(result.KeyComplaints) > 0 {
			fmt.Println("   Top complaints:")
			for j, kc := range result.KeyComplaints {
				if j >= 5 {
					break
				}
				fmt.Printf("     • [%s] %s (from %s)\n", kc.Category, kc.Description, kc.Platform)
			}
		}
	}
}
