# 🪙 Coinsights

**AI-Powered Cryptocurrency Exchange Complaint Analyzer**

Coinsights scrapes and analyzes user complaints about cryptocurrency exchanges (starting with Coinbase) from multiple sources including YouTube comments, Reddit discussions, news articles, and review sites using AI-powered search.

![Coinsights Dashboard](assets/oinsights.png)

---

## 🚀 Features

- **YouTube Scraping** - Automatically searches for complaint videos and extracts comments using YouTube Data API v3
- **Gemini AI Search** - Uses Google's Gemini AI with Google Search grounding to find and analyze complaints from:
  - Reddit discussions
  - News articles & review sites (Trustpilot, BBB, Consumer Reports)
  - YouTube video content analysis
- **Structured Output** - Categorizes complaints by type (fees, support, security, etc.)
- **Modern React Dashboard** - Clean UI to visualize and explore complaint data

---

## 📁 Project Structure

```
Coinsights/
├── backend/                 # Go backend
│   ├── cmd/server/         # Main entry point
│   ├── internal/
│   │   ├── api/            # HTTP handlers & middleware
│   │   ├── config/         # Configuration settings
│   │   ├── models/         # Data models
│   │   ├── scrapers/       # YouTube & Gemini scrapers
│   │   └── services/       # Business logic
│   └── pkg/utils/          # Utility functions
├── frontend/               # React TypeScript frontend
│   ├── public/
│   └── src/
│       ├── components/     # React components
│       ├── pages/          # Page components
│       ├── services/       # API services
│       ├── styles/         # CSS styles
│       └── types/          # TypeScript types
├── data/                   # Scraped data output
│   ├── youtube_latest_results.json
│   └── gemini_latest_results.json
└── assets/                 # Images & assets
```

---

## 🛠️ Tech Stack

### Backend
- **Go 1.24+** - Fast, compiled language
- **Colly** - Web scraping framework
- **Google Gemini AI** - AI-powered search with grounding
- **YouTube Data API v3** - Video and comment scraping

### Frontend
- **React 18** - UI framework
- **TypeScript** - Type safety
- **React Router** - Navigation
- **Lucide React** - Icons
- **Axios** - HTTP client

---

## ⚙️ Setup & Installation

### Prerequisites
- Go 1.24+
- Node.js 18+
- npm or yarn

### 1. Clone the repository
```bash
git clone https://github.com/tasnint/Coinsights.git
cd Coinsights
```

### 2. Configure environment variables
Create a `.env` file in the root directory:
```env
# Google Cloud Platform API Keys
YOUTUBE_API_KEY=your_youtube_api_key_here
GEMINI_API_KEY=your_gemini_api_key_here

# Server Configuration
PORT=8080
ENV=development
```

### 3. Get API Keys

#### YouTube Data API v3
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project
3. Enable "YouTube Data API v3"
4. Create credentials (API Key)
5. Copy the API key to `.env`

#### Gemini API
1. Go to [Google AI Studio](https://aistudio.google.com/)
2. Get an API key
3. Copy the API key to `.env`

### 4. Install & Run Backend
```bash
cd backend
go mod download
cd cmd/server
go run main.go
```

### 5. Install & Run Frontend
```bash
cd frontend
npm install
npm start
```

The frontend will be available at `http://localhost:3000`

---

## 🔧 Configuration

Edit `backend/internal/config/config.go` to customize:

```go
// Search queries for YouTube
var SearchQueries = []string{
    "coinbase problems",
    "coinbase complaints",
    "coinbase fees too high",
    // ... more queries
}

// Scraping settings
func DefaultSettings() Settings {
    return Settings{
        VideosPerQuery:   5,
        CommentsPerVideo: 20,
        MaxQueries:       25,
    }
}
```

---

## 📊 Output Data

### YouTube Results (`youtube_latest_results.json`)
```json
{
  "videos": [...],
  "comments": [...],
  "scraped_at": "2026-01-28T09:00:00Z"
}
```

### Gemini AI Results (`gemini_latest_results.json`)
```json
[
  {
    "query": "coinbase complaints reddit",
    "summary": "Users report issues with...",
    "key_complaints": [
      {
        "category": "fees",
        "description": "High transaction fees",
        "frequency": "common",
        "platform": "reddit"
      }
    ],
    "sources": [...],
    "sentiment_breakdown": {
      "negative": 0.7,
      "neutral": 0.2,
      "positive": 0.1
    }
  }
]
```

---

## 🔒 API Rate Limits

| API | Free Tier Limits |
|-----|------------------|
| YouTube Data API | 10,000 units/day |
| Gemini API | 15 requests/min, 1,500/day |

The scraper includes built-in rate limiting and retry logic to handle these limits gracefully.

---

## 📝 License

MIT License - feel free to use and modify!

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📧 Contact

**Tasnim** - [@tasnint](https://github.com/tasnint)

Project Link: [https://github.com/tasnint/Coinsights](https://github.com/tasnint/Coinsights)
