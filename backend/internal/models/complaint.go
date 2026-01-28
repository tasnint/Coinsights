package models

import "time"

// Complaint represents a user complaint or negative feedback about Coinbase
type Complaint struct {
	ID          string    `json:"id"`
	Source      string    `json:"source"`       // "youtube" or "google"
	Title       string    `json:"title"`        // Video title or search result title
	Description string    `json:"description"`  // Comment text or snippet
	URL         string    `json:"url"`          // Link to source
	Author      string    `json:"author"`       // Username or channel name
	PublishedAt time.Time `json:"published_at"` // When it was posted
	ScrapedAt   time.Time `json:"scraped_at"`   // When we found it
	Sentiment   string    `json:"sentiment"`    // "negative", "neutral", "positive"
	Category    string    `json:"category"`     // "fees", "support", "security", etc.
	Likes       int       `json:"likes"`        // Engagement metric
}

// Thumbnail represents a YouTube thumbnail image
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Thumbnails holds different thumbnail sizes
type Thumbnails struct {
	Default  *Thumbnail `json:"default,omitempty"`
	Medium   *Thumbnail `json:"medium,omitempty"`
	High     *Thumbnail `json:"high,omitempty"`
	Standard *Thumbnail `json:"standard,omitempty"`
	MaxRes   *Thumbnail `json:"maxres,omitempty"`
}

// YouTubeVideo represents a YouTube video search result
// Matches youtube#searchResult structure from YouTube Data API
type YouTubeVideo struct {
	VideoID              string     `json:"video_id"`
	ChannelID            string     `json:"channel_id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	ChannelTitle         string     `json:"channel_title"`
	PublishedAt          time.Time  `json:"published_at"`
	Thumbnails           Thumbnails `json:"thumbnails"`
	LiveBroadcastContent string     `json:"live_broadcast_content"` // "upcoming", "live", or "none"
	URL                  string     `json:"url"`
	// Statistics from videos.list API
	ViewCount    int64    `json:"view_count"`
	LikeCount    int64    `json:"like_count"`
	CommentCount int64    `json:"comment_count"`
	Duration     string   `json:"duration"` // ISO 8601 duration (e.g., "PT4M13S")
	Tags         []string `json:"tags"`     // Video tags
}

// YouTubeComment represents a comment on a YouTube video
type YouTubeComment struct {
	CommentID   string    `json:"comment_id"`
	VideoID     string    `json:"video_id"`
	AuthorName  string    `json:"author_name"`
	Text        string    `json:"text"`
	LikeCount   int       `json:"like_count"`
	PublishedAt time.Time `json:"published_at"`
}

// GoogleResult represents a Google search result
type GoogleResult struct {
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Snippet   string    `json:"snippet"`
	Source    string    `json:"source"` // Domain name
	ScrapedAt time.Time `json:"scraped_at"`
}

// ScrapeResult holds all scraped data
type ScrapeResult struct {
	Videos        []YouTubeVideo   `json:"videos"`
	Comments      []YouTubeComment `json:"comments"`
	GoogleResults []GoogleResult   `json:"google_results"`
	Complaints    []Complaint      `json:"complaints"`
	ScrapedAt     time.Time        `json:"scraped_at"`
	Query         string           `json:"query"`
}
