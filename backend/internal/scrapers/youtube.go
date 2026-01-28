package scrapers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tasnint/coinsights/internal/models"
)

// YouTubeScraper handles YouTube Data API requests
type YouTubeScraper struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

// NewYouTubeScraper creates a new YouTube scraper instance
func NewYouTubeScraper(apiKey string) *YouTubeScraper {
	return &YouTubeScraper{
		APIKey:  apiKey,
		BaseURL: "https://www.googleapis.com/youtube/v3",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ============================================
// YouTube API Response Structures
// Based on official YouTube Data API docs
// ============================================

// ThumbnailResponse represents a single thumbnail from API
type ThumbnailResponse struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ThumbnailsResponse holds all thumbnail sizes from API
type ThumbnailsResponse struct {
	Default  *ThumbnailResponse `json:"default,omitempty"`
	Medium   *ThumbnailResponse `json:"medium,omitempty"`
	High     *ThumbnailResponse `json:"high,omitempty"`
	Standard *ThumbnailResponse `json:"standard,omitempty"`
	MaxRes   *ThumbnailResponse `json:"maxres,omitempty"`
}

// SearchResultID represents the id object in search results
type SearchResultID struct {
	Kind       string `json:"kind"`       // "youtube#video", "youtube#channel", "youtube#playlist"
	VideoID    string `json:"videoId"`    // Present when kind is youtube#video
	ChannelID  string `json:"channelId"`  // Present when kind is youtube#channel
	PlaylistID string `json:"playlistId"` // Present when kind is youtube#playlist
}

// SearchResultSnippet represents the snippet object in search results
type SearchResultSnippet struct {
	PublishedAt          string             `json:"publishedAt"` // ISO 8601 datetime
	ChannelID            string             `json:"channelId"`
	Title                string             `json:"title"`
	Description          string             `json:"description"`
	Thumbnails           ThumbnailsResponse `json:"thumbnails"`
	ChannelTitle         string             `json:"channelTitle"`
	LiveBroadcastContent string             `json:"liveBroadcastContent"` // "upcoming", "live", "none"
}

// SearchResult represents a single youtube#searchResult
type SearchResult struct {
	Kind    string              `json:"kind"` // "youtube#searchResult"
	Etag    string              `json:"etag"`
	ID      SearchResultID      `json:"id"`
	Snippet SearchResultSnippet `json:"snippet"`
}

// SearchListResponse represents the response from search.list API
type SearchListResponse struct {
	Kind          string         `json:"kind"` // "youtube#searchListResponse"
	Etag          string         `json:"etag"`
	NextPageToken string         `json:"nextPageToken,omitempty"`
	PrevPageToken string         `json:"prevPageToken,omitempty"`
	PageInfo      PageInfo       `json:"pageInfo"`
	Items         []SearchResult `json:"items"`
}

// PageInfo contains paging information
type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// CommentSnippet represents a comment's snippet
type CommentSnippet struct {
	AuthorDisplayName string `json:"authorDisplayName"`
	AuthorChannelURL  string `json:"authorChannelUrl"`
	TextDisplay       string `json:"textDisplay"`
	TextOriginal      string `json:"textOriginal"`
	LikeCount         int    `json:"likeCount"`
	PublishedAt       string `json:"publishedAt"`
	UpdatedAt         string `json:"updatedAt"`
}

// TopLevelComment represents a top-level comment
type TopLevelComment struct {
	Kind    string         `json:"kind"`
	Etag    string         `json:"etag"`
	ID      string         `json:"id"`
	Snippet CommentSnippet `json:"snippet"`
}

// CommentThreadSnippet represents a comment thread's snippet
type CommentThreadSnippet struct {
	VideoID         string          `json:"videoId"`
	TopLevelComment TopLevelComment `json:"topLevelComment"`
	TotalReplyCount int             `json:"totalReplyCount"`
}

// CommentThread represents a single comment thread
type CommentThread struct {
	Kind    string               `json:"kind"`
	Etag    string               `json:"etag"`
	ID      string               `json:"id"`
	Snippet CommentThreadSnippet `json:"snippet"`
}

// CommentThreadListResponse represents the response from commentThreads.list API
type CommentThreadListResponse struct {
	Kind          string          `json:"kind"`
	Etag          string          `json:"etag"`
	NextPageToken string          `json:"nextPageToken,omitempty"`
	PageInfo      PageInfo        `json:"pageInfo"`
	Items         []CommentThread `json:"items"`
}

// ============================================
// Videos List API Response Structures
// ============================================

// VideoStatistics represents the statistics object from videos.list
type VideoStatistics struct {
	ViewCount     string `json:"viewCount"`
	LikeCount     string `json:"likeCount"`
	DislikeCount  string `json:"dislikeCount"` // Deprecated but may still appear
	FavoriteCount string `json:"favoriteCount"`
	CommentCount  string `json:"commentCount"`
}

// VideoContentDetails represents the contentDetails object from videos.list
type VideoContentDetails struct {
	Duration        string `json:"duration"`   // ISO 8601 duration (PT4M13S)
	Dimension       string `json:"dimension"`  // "2d" or "3d"
	Definition      string `json:"definition"` // "hd" or "sd"
	Caption         string `json:"caption"`    // "true" or "false"
	LicensedContent bool   `json:"licensedContent"`
}

// VideoSnippet represents the snippet object from videos.list (more detailed than search)
type VideoSnippet struct {
	PublishedAt          string             `json:"publishedAt"`
	ChannelID            string             `json:"channelId"`
	Title                string             `json:"title"`
	Description          string             `json:"description"` // Full description, not truncated
	Thumbnails           ThumbnailsResponse `json:"thumbnails"`
	ChannelTitle         string             `json:"channelTitle"`
	Tags                 []string           `json:"tags"`
	CategoryID           string             `json:"categoryId"`
	LiveBroadcastContent string             `json:"liveBroadcastContent"`
}

// VideoResource represents a single video from videos.list
type VideoResource struct {
	Kind           string              `json:"kind"` // "youtube#video"
	Etag           string              `json:"etag"`
	ID             string              `json:"id"`
	Snippet        VideoSnippet        `json:"snippet"`
	ContentDetails VideoContentDetails `json:"contentDetails"`
	Statistics     VideoStatistics     `json:"statistics"`
}

// VideoListResponse represents the response from videos.list API
type VideoListResponse struct {
	Kind          string          `json:"kind"` // "youtube#videoListResponse"
	Etag          string          `json:"etag"`
	NextPageToken string          `json:"nextPageToken,omitempty"`
	PageInfo      PageInfo        `json:"pageInfo"`
	Items         []VideoResource `json:"items"`
}

// ============================================
// API Methods
// ============================================

// SearchVideos searches for YouTube videos matching the query
// Uses: GET https://www.googleapis.com/youtube/v3/search
func (ys *YouTubeScraper) SearchVideos(query string, maxResults int) ([]models.YouTubeVideo, error) {
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("q", query)
	params.Add("type", "video") // Only return videos
	params.Add("maxResults", fmt.Sprintf("%d", maxResults))
	params.Add("order", "relevance") // Can be: date, rating, relevance, title, viewCount
	params.Add("key", ys.APIKey)

	reqURL := fmt.Sprintf("%s/search?%s", ys.BaseURL, params.Encode())

	resp, err := ys.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search videos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("YouTube API error (status %d): %s", resp.StatusCode, string(body))
	}

	var searchResp SearchListResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert API response to our model
	videos := make([]models.YouTubeVideo, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		// Only process video results
		if item.ID.VideoID == "" {
			continue
		}

		publishedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		video := models.YouTubeVideo{
			VideoID:              item.ID.VideoID,
			ChannelID:            item.Snippet.ChannelID,
			Title:                item.Snippet.Title,
			Description:          item.Snippet.Description,
			ChannelTitle:         item.Snippet.ChannelTitle,
			PublishedAt:          publishedAt,
			LiveBroadcastContent: item.Snippet.LiveBroadcastContent,
			URL:                  fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID),
			Thumbnails:           convertThumbnails(item.Snippet.Thumbnails),
		}
		videos = append(videos, video)
	}

	return videos, nil
}

// GetVideoComments fetches comments for a specific video
// Uses: GET https://www.googleapis.com/youtube/v3/commentThreads
func (ys *YouTubeScraper) GetVideoComments(videoID string, maxResults int) ([]models.YouTubeComment, error) {
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("videoId", videoID)
	params.Add("maxResults", fmt.Sprintf("%d", maxResults))
	params.Add("order", "relevance") // Can be: time, relevance
	params.Add("textFormat", "plainText")
	params.Add("key", ys.APIKey)

	reqURL := fmt.Sprintf("%s/commentThreads?%s", ys.BaseURL, params.Encode())

	resp, err := ys.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("YouTube API error (status %d): %s", resp.StatusCode, string(body))
	}

	var commentsResp CommentThreadListResponse
	if err := json.NewDecoder(resp.Body).Decode(&commentsResp); err != nil {
		return nil, fmt.Errorf("failed to decode comments: %w", err)
	}

	comments := make([]models.YouTubeComment, 0, len(commentsResp.Items))
	for _, item := range commentsResp.Items {
		snippet := item.Snippet.TopLevelComment.Snippet
		publishedAt, _ := time.Parse(time.RFC3339, snippet.PublishedAt)

		comment := models.YouTubeComment{
			CommentID:   item.ID,
			VideoID:     videoID,
			AuthorName:  snippet.AuthorDisplayName,
			Text:        snippet.TextOriginal,
			LikeCount:   snippet.LikeCount,
			PublishedAt: publishedAt,
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// GetVideoDetails fetches detailed information for multiple videos
// Uses: GET https://www.googleapis.com/youtube/v3/videos
// This enriches search results with stats (views, likes) and full description
func (ys *YouTubeScraper) GetVideoDetails(videoIDs []string) (map[string]*VideoResource, error) {
	if len(videoIDs) == 0 {
		return make(map[string]*VideoResource), nil
	}

	// YouTube allows up to 50 video IDs per request
	params := url.Values{}
	params.Add("part", "snippet,statistics,contentDetails")
	params.Add("id", joinStrings(videoIDs, ","))
	params.Add("key", ys.APIKey)

	reqURL := fmt.Sprintf("%s/videos?%s", ys.BaseURL, params.Encode())

	resp, err := ys.HTTPClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("YouTube API error (status %d): %s", resp.StatusCode, string(body))
	}

	var videosResp VideoListResponse
	if err := json.NewDecoder(resp.Body).Decode(&videosResp); err != nil {
		return nil, fmt.Errorf("failed to decode video details: %w", err)
	}

	// Create a map for easy lookup by video ID
	videoMap := make(map[string]*VideoResource)
	for i := range videosResp.Items {
		videoMap[videosResp.Items[i].ID] = &videosResp.Items[i]
	}

	return videoMap, nil
}

// ScrapeAll searches videos, enriches with details, and fetches comments
func (ys *YouTubeScraper) ScrapeAll(queries []string, videosPerQuery int, commentsPerVideo int) (*models.ScrapeResult, error) {
	result := &models.ScrapeResult{
		Videos:    []models.YouTubeVideo{},
		Comments:  []models.YouTubeComment{},
		ScrapedAt: time.Now(),
	}

	for _, query := range queries {
		fmt.Printf("Searching YouTube for: %s\n", query)

		videos, err := ys.SearchVideos(query, videosPerQuery)
		if err != nil {
			fmt.Printf("Error searching for '%s': %v\n", query, err)
			continue
		}
		fmt.Printf("Found %d videos\n", len(videos))

		// Collect video IDs for batch details fetch
		videoIDs := make([]string, len(videos))
		for i, v := range videos {
			videoIDs[i] = v.VideoID
		}

		// Fetch detailed stats for all videos in one API call
		fmt.Printf("Fetching video statistics...\n")
		videoDetails, err := ys.GetVideoDetails(videoIDs)
		if err != nil {
			fmt.Printf("Error fetching video details: %v\n", err)
		}

		// Enrich videos with statistics
		for i := range videos {
			if details, ok := videoDetails[videos[i].VideoID]; ok {
				videos[i].ViewCount = parseCount(details.Statistics.ViewCount)
				videos[i].LikeCount = parseCount(details.Statistics.LikeCount)
				videos[i].CommentCount = parseCount(details.Statistics.CommentCount)
				videos[i].Duration = details.ContentDetails.Duration
				videos[i].Tags = details.Snippet.Tags
				// Use full description from videos.list (not truncated)
				if details.Snippet.Description != "" {
					videos[i].Description = details.Snippet.Description
				}
			}
		}

		result.Videos = append(result.Videos, videos...)

		// Fetch comments for each video
		for _, video := range videos {
			fmt.Printf("Fetching comments for: %s\n", video.Title)

			comments, err := ys.GetVideoComments(video.VideoID, commentsPerVideo)
			if err != nil {
				fmt.Printf("Error fetching comments for %s: %v\n", video.VideoID, err)
				continue
			}

			result.Comments = append(result.Comments, comments...)
			fmt.Printf("Found %d comments\n", len(comments))

			// Rate limiting - be nice to the API
			time.Sleep(500 * time.Millisecond)
		}
	}

	return result, nil
}

// convertThumbnails converts API thumbnails to model thumbnails
func convertThumbnails(apiThumbs ThumbnailsResponse) models.Thumbnails {
	convert := func(t *ThumbnailResponse) *models.Thumbnail {
		if t == nil {
			return nil
		}
		return &models.Thumbnail{
			URL:    t.URL,
			Width:  t.Width,
			Height: t.Height,
		}
	}

	return models.Thumbnails{
		Default:  convert(apiThumbs.Default),
		Medium:   convert(apiThumbs.Medium),
		High:     convert(apiThumbs.High),
		Standard: convert(apiThumbs.Standard),
		MaxRes:   convert(apiThumbs.MaxRes),
	}
}

// joinStrings joins a slice of strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// parseCount converts a string count to int64
func parseCount(s string) int64 {
	if s == "" {
		return 0
	}
	var count int64
	fmt.Sscanf(s, "%d", &count)
	return count
}
