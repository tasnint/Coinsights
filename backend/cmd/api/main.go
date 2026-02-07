// API server for the Coinsights frontend
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// ============================================
// DATA STRUCTURES
// ============================================

// Issue represents a tracked issue from scraped data
type Issue struct {
	ID            string   `json:"id"`
	Exchange      string   `json:"exchange"`
	Category      string   `json:"category"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	FirstDetected string   `json:"first_detected"`
	Severity      string   `json:"severity"`
	Status        string   `json:"status"`
	Count         int      `json:"count"`
	Examples      []string `json:"examples"`
}

// Resolution represents a resolved issue
type Resolution struct {
	ID            string `json:"id"`
	Exchange      string `json:"exchange"`
	IssueCategory string `json:"issue_category"`
	Summary       string `json:"summary"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// YouTubeAnalysis matches the structure from youtube_analysis.json
type YouTubeAnalysis struct {
	TotalVideos      int                     `json:"total_videos"`
	TotalComments    int                     `json:"total_comments"`
	TotalIssues      int                     `json:"total_issues"`
	Categories       map[string]CategoryData `json:"categories"`
	IssuesByCategory []CategorySummary       `json:"issues_by_category"`
	AnalyzedAt       time.Time               `json:"analyzed_at"`
}

type CategoryData struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Count    int      `json:"count"`
	Examples []string `json:"examples"`
	Severity string   `json:"severity"`
}

type CategorySummary struct {
	Category    string   `json:"category"`
	Count       int      `json:"count"`
	Percentage  float64  `json:"percentage"`
	TopExamples []string `json:"top_examples"`
}

// GeminiResult matches the structure from gemini_latest_results.json
type GeminiResult struct {
	Query              string             `json:"query"`
	Summary            string             `json:"summary"`
	KeyComplaints      []KeyComplaint     `json:"key_complaints"`
	Sources            []Source           `json:"sources"`
	SentimentBreakdown map[string]float64 `json:"sentiment_breakdown"`
	GeneratedAt        time.Time          `json:"generated_at"`
}

type KeyComplaint struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Frequency   string `json:"frequency"`
	Platform    string `json:"platform"`
}

type Source struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Domain string `json:"domain"`
}

// Stats for the dashboard
type Stats struct {
	VideosAnalyzed   int `json:"videos_analyzed"`
	CommentsAnalyzed int `json:"comments_analyzed"`
	IssuesFound      int `json:"issues_found"`
	Categories       int `json:"categories"`
}

// ============================================
// GLOBAL DATA (loaded at startup)
// ============================================

var (
	youtubeAnalysis *YouTubeAnalysis
	geminiResults   []GeminiResult
	issues          []Issue
	resolutions     []Resolution
	stats           Stats
)

func main() {
	// Load environment variables
	envPaths := []string{
		"../../.env",
		".env",
	}
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	// Load data from JSON files
	if err := loadData(); err != nil {
		log.Printf("Warning: Could not load all data: %v", err)
	}

	// Setup routes
	mux := http.NewServeMux()

	// CORS middleware wrapper
	handler := corsMiddleware(mux)

	// API endpoints
	mux.HandleFunc("GET /api/issues", handleGetIssues)
	mux.HandleFunc("GET /api/resolutions", handleGetResolutions)
	mux.HandleFunc("GET /api/stats", handleGetStats)
	mux.HandleFunc("GET /api/analysis/youtube", handleGetYouTubeAnalysis)
	mux.HandleFunc("GET /api/analysis/gemini", handleGetGeminiAnalysis)
	mux.HandleFunc("GET /health", handleHealth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Coinsights API Server starting on http://localhost:%s\n", port)
	fmt.Println("==========================================")
	fmt.Printf("📊 Loaded %d issues from scraped data\n", len(issues))
	fmt.Printf("📺 YouTube: %d videos, %d comments analyzed\n", stats.VideosAnalyzed, stats.CommentsAnalyzed)
	fmt.Println("==========================================")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}

// ============================================
// DATA LOADING
// ============================================

func loadData() error {
	dataDir := "../../data"

	// Load YouTube analysis
	ytPath := filepath.Join(dataDir, "youtube_analysis.json")
	if data, err := os.ReadFile(ytPath); err == nil {
		if err := json.Unmarshal(data, &youtubeAnalysis); err != nil {
			log.Printf("Error parsing youtube_analysis.json: %v", err)
		} else {
			// Convert categories to issues
			issues = convertCategoriesToIssues(youtubeAnalysis)
			stats.VideosAnalyzed = youtubeAnalysis.TotalVideos
			stats.CommentsAnalyzed = youtubeAnalysis.TotalComments
			stats.IssuesFound = youtubeAnalysis.TotalIssues
			stats.Categories = len(youtubeAnalysis.Categories)
		}
	} else {
		log.Printf("Could not read youtube_analysis.json: %v", err)
	}

	// Load Gemini results
	geminiPath := filepath.Join(dataDir, "gemini_latest_results.json")
	if data, err := os.ReadFile(geminiPath); err == nil {
		if err := json.Unmarshal(data, &geminiResults); err != nil {
			log.Printf("Error parsing gemini_latest_results.json: %v", err)
		} else {
			// Enrich issues with Gemini data
			enrichIssuesWithGemini(geminiResults)
		}
	} else {
		log.Printf("Could not read gemini_latest_results.json: %v", err)
	}

	// Create some demo resolutions (in a real system, these would come from a database)
	resolutions = createDemoResolutions()

	return nil
}

func convertCategoriesToIssues(analysis *YouTubeAnalysis) []Issue {
	var result []Issue

	categoryTitles := map[string]string{
		"customer_support": "Poor Customer Support Response",
		"account_locked":   "Account Locked/Frozen Without Warning",
		"fees":             "Hidden or Excessive Fees",
		"withdrawal":       "Withdrawal Delays and Blocks",
		"security":         "Security Vulnerabilities and Breaches",
		"verification":     "KYC/Verification Process Issues",
		"app_bugs":         "App Crashes and Technical Bugs",
		"deposits":         "Deposit Processing Problems",
		"general_negative": "General User Dissatisfaction",
		"trading":          "Trading Execution Issues",
	}

	categoryDescriptions := map[string]string{
		"customer_support": "Users report difficulty reaching support, long wait times, unhelpful responses, and unresolved tickets.",
		"account_locked":   "Accounts being locked or frozen without clear explanation, preventing access to funds for extended periods.",
		"fees":             "Complaints about unclear fee structures, hidden fees, high spreads, and unexpected charges on transactions.",
		"withdrawal":       "Users experiencing delays or blocks when trying to withdraw funds to external wallets or bank accounts.",
		"security":         "Concerns about unauthorized access, data breaches, SIM swap attacks, and inadequate security measures.",
		"verification":     "Tedious KYC processes, repeated document requests, verification failures, and long pending times.",
		"app_bugs":         "Application crashes, slow loading, glitches, and technical errors preventing normal usage.",
		"deposits":         "Deposits not appearing, failed bank transfers, and delays in fund availability.",
		"general_negative": "Overall negative sentiment and recommendations to avoid the platform.",
		"trading":          "Issues with order execution, slippage, spread manipulation, and inability to buy/sell.",
	}

	id := 1
	for key, cat := range analysis.Categories {
		title := categoryTitles[key]
		if title == "" {
			title = cat.Name
		}
		desc := categoryDescriptions[key]
		if desc == "" {
			desc = fmt.Sprintf("Issues related to %s", cat.Name)
		}

		status := "active"
		if cat.Count < 50 {
			status = "investigating"
		}

		// Truncate examples
		examples := cat.Examples
		if len(examples) > 3 {
			examples = examples[:3]
		}

		result = append(result, Issue{
			ID:            fmt.Sprintf("issue-%d", id),
			Exchange:      "Coinbase",
			Category:      cat.Name,
			Title:         title,
			Description:   desc,
			FirstDetected: analysis.AnalyzedAt.Format("2006-01-02"),
			Severity:      cat.Severity,
			Status:        status,
			Count:         cat.Count,
			Examples:      examples,
		})
		id++
	}

	return result
}

func enrichIssuesWithGemini(gemini []GeminiResult) {
	// Add descriptions from Gemini AI analysis
	for i := range issues {
		for _, g := range gemini {
			for _, complaint := range g.KeyComplaints {
				if matchesCategory(issues[i].Category, complaint.Category) {
					// Enrich with more detailed description
					if len(complaint.Description) > len(issues[i].Description) {
						issues[i].Description = complaint.Description
					}
				}
			}
		}
	}
}

func matchesCategory(issueCategory, complaintCategory string) bool {
	categoryMap := map[string][]string{
		"Customer Support":      {"customer_support"},
		"Account Locked/Frozen": {"account_issues", "account_locked"},
		"High Fees":             {"fees"},
		"Withdrawal Problems":   {"withdrawal_problems", "withdrawal"},
		"Security Issues":       {"security"},
		"Verification Issues":   {"verification"},
	}

	for cat, matches := range categoryMap {
		if issueCategory == cat {
			for _, m := range matches {
				if complaintCategory == m {
					return true
				}
			}
		}
	}
	return false
}

func createDemoResolutions() []Resolution {
	// In a real system, these would be fetched from the blockchain/database
	return []Resolution{
		{
			ID:            "res-1",
			Exchange:      "Coinbase",
			IssueCategory: "High Fees",
			Summary:       "Fee disclosure page updated with clearer breakdowns. Advanced Trade offers lower fees.",
			Status:        "verified",
			CreatedAt:     "2026-01-15",
		},
	}
}

// ============================================
// HANDLERS
// ============================================

func handleGetIssues(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"issues": issues,
		"count":  len(issues),
	})
}

func handleGetResolutions(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"resolutions": resolutions,
		"count":       len(resolutions),
	})
}

func handleGetStats(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, stats)
}

func handleGetYouTubeAnalysis(w http.ResponseWriter, r *http.Request) {
	if youtubeAnalysis == nil {
		respondError(w, http.StatusNotFound, "YouTube analysis not available")
		return
	}
	respondJSON(w, http.StatusOK, youtubeAnalysis)
}

func handleGetGeminiAnalysis(w http.ResponseWriter, r *http.Request) {
	if geminiResults == nil {
		respondError(w, http.StatusNotFound, "Gemini analysis not available")
		return
	}
	respondJSON(w, http.StatusOK, geminiResults)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// ============================================
// MIDDLEWARE & HELPERS
// ============================================

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
