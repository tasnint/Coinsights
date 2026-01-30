package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tasnint/coinsights/internal/models"
)

// IssueCategory represents a category of complaints
type IssueCategory struct {
	Name        string   `json:"name"`
	Keywords    []string `json:"keywords"`
	Count       int      `json:"count"`
	Examples    []string `json:"examples"`
	Severity    string   `json:"severity"` // "high", "medium", "low"
}

// ExtractedIssue represents a single extracted issue
type ExtractedIssue struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"`
	Text        string    `json:"text"`
	Source      string    `json:"source"`      // "video_title", "video_description", "video_tags", "comment"
	SourceURL   string    `json:"source_url"`
	SourceTitle string    `json:"source_title"`
	Likes       int       `json:"likes"`       // For comments
	ExtractedAt time.Time `json:"extracted_at"`
}

// AnalysisResult holds the complete analysis
type AnalysisResult struct {
	TotalVideos      int                       `json:"total_videos"`
	TotalComments    int                       `json:"total_comments"`
	TotalIssues      int                       `json:"total_issues"`
	Categories       map[string]*IssueCategory `json:"categories"`
	TopIssues        []ExtractedIssue          `json:"top_issues"`
	IssuesByCategory []CategorySummary         `json:"issues_by_category"`
	AnalyzedAt       time.Time                 `json:"analyzed_at"`
}

// CategorySummary provides a summary for each category
type CategorySummary struct {
	Category   string   `json:"category"`
	Count      int      `json:"count"`
	Percentage float64  `json:"percentage"`
	TopExamples []string `json:"top_examples"`
}

// YouTubeAnalyzer analyzes YouTube scrape results
type YouTubeAnalyzer struct {
	categories map[string]*IssueCategory
	issues     []ExtractedIssue
}

// NewYouTubeAnalyzer creates a new analyzer with predefined categories
func NewYouTubeAnalyzer() *YouTubeAnalyzer {
	return &YouTubeAnalyzer{
		categories: initCategories(),
		issues:     []ExtractedIssue{},
	}
}

// initCategories sets up the complaint categories with keywords
func initCategories() map[string]*IssueCategory {
	return map[string]*IssueCategory{
		"customer_support": {
			Name: "Customer Support",
			Keywords: []string{
				"support", "customer service", "no response", "no reply", "agent", 
				"ticket", "help", "contact", "chat", "email", "phone", "waiting",
				"ignored", "unhelpful", "terrible support", "worst support",
			},
			Severity: "high",
			Examples: []string{},
		},
		"account_locked": {
			Name: "Account Locked/Frozen",
			Keywords: []string{
				"locked", "frozen", "restricted", "suspended", "blocked", "disabled",
				"can't access", "cannot access", "locked out", "freeze", "hold",
				"account closed", "account terminated", "verification hold",
			},
			Severity: "high",
			Examples: []string{},
		},
		"fees": {
			Name: "High Fees",
			Keywords: []string{
				"fees", "expensive", "high fee", "hidden fee", "spread", "commission",
				"cost", "charges", "overcharge", "rip off", "ripoff", "too much",
				"fee structure", "trading fee", "withdrawal fee",
			},
			Severity: "medium",
			Examples: []string{},
		},
		"withdrawal": {
			Name: "Withdrawal Problems",
			Keywords: []string{
				"withdraw", "withdrawal", "can't withdraw", "withdrawal pending",
				"cash out", "transfer out", "send", "move funds", "stuck funds",
				"withdrawal failed", "withdrawal delayed",
			},
			Severity: "high",
			Examples: []string{},
		},
		"security": {
			Name: "Security Issues",
			Keywords: []string{
				"hack", "hacked", "stolen", "scam", "phishing", "unauthorized",
				"security", "breach", "compromised", "fraud", "theft", "lost crypto",
				"2fa", "two factor", "sim swap",
			},
			Severity: "high",
			Examples: []string{},
		},
		"verification": {
			Name: "Verification Issues",
			Keywords: []string{
				"verification", "verify", "kyc", "identity", "id verification",
				"document", "upload", "rejected", "pending verification",
				"verification failed", "verify identity",
			},
			Severity: "medium",
			Examples: []string{},
		},
		"app_bugs": {
			Name: "App/Technical Issues",
			Keywords: []string{
				"bug", "crash", "not working", "glitch", "error", "broken",
				"app issue", "loading", "slow", "lag", "freeze", "update",
				"won't load", "won't open", "technical",
			},
			Severity: "medium",
			Examples: []string{},
		},
		"deposits": {
			Name: "Deposit Problems",
			Keywords: []string{
				"deposit", "deposit pending", "deposit missing", "deposit failed",
				"bank transfer", "wire transfer", "ach", "funds not showing",
				"money missing", "payment",
			},
			Severity: "high",
			Examples: []string{},
		},
		"trading": {
			Name: "Trading Issues",
			Keywords: []string{
				"trade", "trading", "order", "limit order", "market order",
				"execution", "slippage", "price", "spread", "liquidity",
				"can't buy", "can't sell", "order failed",
			},
			Severity: "medium",
			Examples: []string{},
		},
		"general_negative": {
			Name: "General Complaints",
			Keywords: []string{
				"terrible", "worst", "awful", "horrible", "bad", "hate",
				"never use", "avoid", "stay away", "don't use", "nightmare",
				"frustrating", "disappointed", "angry", "scam",
			},
			Severity: "low",
			Examples: []string{},
		},
	}
}

// AnalyzeFile reads and analyzes a YouTube results JSON file
func (a *YouTubeAnalyzer) AnalyzeFile(filepath string) (*AnalysisResult, error) {
	// Read the file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var result models.ScrapeResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("📊 Analyzing %d videos and %d comments...\n", len(result.Videos), len(result.Comments))

	// Analyze videos
	for _, video := range result.Videos {
		a.analyzeVideo(video)
	}

	// Analyze comments
	for _, comment := range result.Comments {
		a.analyzeComment(comment, result.Videos)
	}

	// Build result
	return a.buildResult(len(result.Videos), len(result.Comments)), nil
}

// analyzeVideo extracts issues from a video's title, description, and tags
func (a *YouTubeAnalyzer) analyzeVideo(video models.YouTubeVideo) {
	// Analyze title
	if issues := a.findIssuesInText(video.Title); len(issues) > 0 {
		for _, category := range issues {
			a.addIssue(ExtractedIssue{
				Category:    category,
				Text:        video.Title,
				Source:      "video_title",
				SourceURL:   video.URL,
				SourceTitle: video.Title,
			})
		}
	}

	// Analyze description (first 500 chars)
	desc := video.Description
	if len(desc) > 500 {
		desc = desc[:500]
	}
	if issues := a.findIssuesInText(desc); len(issues) > 0 {
		for _, category := range issues {
			a.addIssue(ExtractedIssue{
				Category:    category,
				Text:        desc,
				Source:      "video_description",
				SourceURL:   video.URL,
				SourceTitle: video.Title,
			})
		}
	}

	// Analyze tags
	tagText := strings.Join(video.Tags, " ")
	if issues := a.findIssuesInText(tagText); len(issues) > 0 {
		for _, category := range issues {
			a.addIssue(ExtractedIssue{
				Category:    category,
				Text:        tagText,
				Source:      "video_tags",
				SourceURL:   video.URL,
				SourceTitle: video.Title,
			})
		}
	}
}

// analyzeComment extracts issues from a comment
func (a *YouTubeAnalyzer) analyzeComment(comment models.YouTubeComment, videos []models.YouTubeVideo) {
	if issues := a.findIssuesInText(comment.Text); len(issues) > 0 {
		// Find the video this comment belongs to
		var videoURL, videoTitle string
		for _, v := range videos {
			if v.VideoID == comment.VideoID {
				videoURL = v.URL
				videoTitle = v.Title
				break
			}
		}

		for _, category := range issues {
			a.addIssue(ExtractedIssue{
				Category:    category,
				Text:        comment.Text,
				Source:      "comment",
				SourceURL:   videoURL,
				SourceTitle: videoTitle,
				Likes:       comment.LikeCount,
			})
		}
	}
}

// findIssuesInText searches text for issue keywords and returns matching categories
func (a *YouTubeAnalyzer) findIssuesInText(text string) []string {
	textLower := strings.ToLower(text)
	foundCategories := []string{}

	for categoryName, category := range a.categories {
		for _, keyword := range category.Keywords {
			// Use word boundary matching for more accuracy
			pattern := `\b` + regexp.QuoteMeta(strings.ToLower(keyword)) + `\b`
			if matched, _ := regexp.MatchString(pattern, textLower); matched {
				foundCategories = append(foundCategories, categoryName)
				break // One match per category is enough
			}
		}
	}

	return foundCategories
}

// addIssue adds an issue and updates category counts
func (a *YouTubeAnalyzer) addIssue(issue ExtractedIssue) {
	issue.ID = fmt.Sprintf("issue_%d", len(a.issues)+1)
	issue.ExtractedAt = time.Now()
	a.issues = append(a.issues, issue)

	// Update category
	if cat, exists := a.categories[issue.Category]; exists {
		cat.Count++
		// Keep top 5 examples
		if len(cat.Examples) < 5 {
			// Truncate long text
			example := issue.Text
			if len(example) > 150 {
				example = example[:150] + "..."
			}
			cat.Examples = append(cat.Examples, example)
		}
	}
}

// buildResult compiles the final analysis result
func (a *YouTubeAnalyzer) buildResult(videoCount, commentCount int) *AnalysisResult {
	result := &AnalysisResult{
		TotalVideos:   videoCount,
		TotalComments: commentCount,
		TotalIssues:   len(a.issues),
		Categories:    a.categories,
		AnalyzedAt:    time.Now(),
	}

	// Build category summaries sorted by count
	summaries := []CategorySummary{}
	for name, cat := range a.categories {
		if cat.Count > 0 {
			percentage := 0.0
			if len(a.issues) > 0 {
				percentage = float64(cat.Count) / float64(len(a.issues)) * 100
			}
			summaries = append(summaries, CategorySummary{
				Category:    name,
				Count:       cat.Count,
				Percentage:  percentage,
				TopExamples: cat.Examples,
			})
		}
	}

	// Sort by count descending
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Count > summaries[j].Count
	})
	result.IssuesByCategory = summaries

	// Get top issues (comments with most likes)
	sort.Slice(a.issues, func(i, j int) bool {
		return a.issues[i].Likes > a.issues[j].Likes
	})

	// Top 20 issues
	topCount := 20
	if len(a.issues) < topCount {
		topCount = len(a.issues)
	}
	result.TopIssues = a.issues[:topCount]

	return result
}

// PrintSummary prints a human-readable summary
func (a *YouTubeAnalyzer) PrintSummary(result *AnalysisResult) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 YOUTUBE COMPLAINT ANALYSIS SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	
	fmt.Printf("\n📺 Videos Analyzed:    %d\n", result.TotalVideos)
	fmt.Printf("💬 Comments Analyzed:  %d\n", result.TotalComments)
	fmt.Printf("🔍 Issues Identified:  %d\n", result.TotalIssues)
	
	fmt.Println("\n📈 ISSUES BY CATEGORY (sorted by frequency)")
	fmt.Println(strings.Repeat("-", 50))
	
	for i, summary := range result.IssuesByCategory {
		if i >= 10 {
			break
		}
		bar := strings.Repeat("█", int(summary.Percentage/5))
		fmt.Printf("%-20s %4d (%5.1f%%) %s\n", 
			a.categories[summary.Category].Name, 
			summary.Count, 
			summary.Percentage,
			bar)
	}

	fmt.Println("\n🔥 TOP COMPLAINTS (by engagement)")
	fmt.Println(strings.Repeat("-", 50))
	
	for i, issue := range result.TopIssues {
		if i >= 5 {
			break
		}
		text := issue.Text
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		fmt.Printf("%d. [%s] (👍 %d likes)\n   \"%s\"\n\n", 
			i+1, 
			a.categories[issue.Category].Name,
			issue.Likes,
			text)
	}
}

// SaveResults saves the analysis to a JSON file
func SaveAnalysisResults(result *AnalysisResult, filepath string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✅ Analysis saved to: %s\n", filepath)
	return nil
}
