# 🪙 Coinsights

**Backend-Focused Cryptocurrency Exchange Complaint Analyzer with On-Chain Verification**

> *This project demonstrates my interest in **Coinbase** and **blockchain technologies** through a backend-heavy implementation featuring Go APIs, Solidity smart contracts, and on-chain attestations.*

Coinsights scrapes and analyzes user complaints about cryptocurrency exchanges (starting with Coinbase) from multiple sources including YouTube comments, Reddit discussions, news articles, and review sites using AI-powered search. Verified resolutions are recorded on-chain for immutable proof.

---

## 🎯 Project Focus

This is a **backend-focused** project showcasing:
- **Go (Golang)** REST API development
- **Blockchain integration** with Ethereum/Base
- **Smart contract** development in Solidity
- **AI-powered data analysis** using Google's Gemini
- **Web scraping** from YouTube and other sources

The frontend is intentionally minimal - just plain text displaying issues and their resolution status.

---

## 🚀 Features

- **YouTube Scraping** - Automatically searches for complaint videos and extracts comments using YouTube Data API v3
- **Gemini AI Search** - Uses Google's Gemini AI with Google Search grounding to find and analyze complaints from:
  - Reddit discussions
  - News articles & review sites (Trustpilot, BBB, Consumer Reports)
  - YouTube video content analysis
- **Structured Output** - Categorizes complaints by type (fees, support, security, etc.)
- **⛓️ On-Chain Verification** - Tamper-proof blockchain attestations for resolved issues:
  - Evidence hashing with Keccak256
  - Smart contract on Base (Coinbase L2)
  - Public, verifiable resolution records
  - Chain-of-custody audit trail
- **Minimal Frontend** - Simple text-based display of issues and resolutions (no styling)

---

## 📁 Project Structure

```
Coinsights/
├── backend/                 # Go backend (main focus)
│   ├── cmd/server/         # Main entry point
│   ├── internal/
│   │   ├── api/            # HTTP handlers & middleware
│   │   ├── config/         # Configuration settings
│   │   ├── models/         # Data models
│   │   ├── scrapers/       # YouTube & Gemini scrapers
│   │   └── services/       # Business logic (blockchain service)
│   └── pkg/utils/          # Utility functions
├── contracts/              # Solidity smart contracts
│   └── ResolutionAttestation.sol
├── frontend/               # Minimal React frontend (plain text)
│   └── src/
│       ├── App.tsx         # Single component - displays issues/resolutions
│       └── index.tsx       # Entry point
├── data/                   # Scraped data output
└── assets/                 # Images & assets
```

---

## 🛠️ Tech Stack

### Backend (Primary Focus)
- **Go 1.24+** - Fast, compiled language for API development
- **Colly** - Web scraping framework
- **Google Gemini AI** - AI-powered search with grounding
- **YouTube Data API v3** - Video and comment scraping
- **go-ethereum** - Ethereum client library for blockchain interactions

### Blockchain
- **Solidity** - Smart contract language
- **Base (Coinbase L2)** - Deployment network (testnet)
- **Keccak256** - Evidence hashing for attestations

### Frontend (Minimal)
- **React 18** - Simple display of data
- **TypeScript** - Type safety
- No styling libraries - plain HTML/text output

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

# Blockchain Configuration (for on-chain verification)
BLOCKCHAIN_NETWORK=base_sepolia
BLOCKCHAIN_RPC_URL=https://sepolia.base.org
BLOCKCHAIN_PRIVATE_KEY=your_wallet_private_key_here
ATTESTATION_CONTRACT_ADDRESS=your_deployed_contract_address

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

#### Blockchain (for On-Chain Verification)
1. Create a wallet (e.g., MetaMask)
2. Get Base Sepolia testnet ETH from [Coinbase Faucet](https://www.coinbase.com/faucets)
3. Deploy the contract (see `contracts/README.md`)
4. Copy private key and contract address to `.env`

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

## ⛓️ Blockchain Verification

Coinsights uses blockchain as a **verification primitive** for resolved issues, not a database.

### How It Works

1. **Evidence Hashing** - When an issue is resolved, the evidence (complaint counts, sentiment data, sources) is hashed using Keccak256
2. **On-Chain Commitment** - The hash is stored on Base (Coinbase L2) via our smart contract
3. **Independent Verification** - Anyone can verify the hash exists on-chain without trusting our backend

---

### 📋 System Phases Overview

The complete data flow from scraping to blockchain verification involves **9 phases**:

| Phase | Name | Operation Type | Key Files | Description |
|:-----:|------|----------------|-----------|-------------|
| 1 | **Startup & Configuration** | Setup | `main.go`, `config.go` | Load environment, initialize scrapers |
| 2 | **YouTube Scraping** | Data Collection | `youtube.go` | Fetch videos & comments via YouTube API |
| 3 | **Gemini AI Analysis** | Data Collection | `gemini.go` | AI-powered web search with Google grounding |
| 4 | **Data Aggregation** | Processing | `main.go`, `complaint.go` | Combine & structure results |
| 5 | **Issue Detection** | Analysis | `resolution.go` | Identify complaint patterns & create issues |
| 6 | **Time-Based Monitoring** | Observation | `resolution.go` | Track complaint trends over time |
| 7 | **Resolution Creation** | Decision | `resolution.go` | Determine if criteria met, bundle evidence |
| 8 | **Blockchain Attestation** | Write (On-Chain) | `blockchain.go`, Smart Contract | Hash & record proof on Base L2 |
| 9 | **Verification** | Read (On-Chain) | `blockchain.go`, `blockchain.go` (handlers) | Prove authenticity anytime |

---

### 🔍 Phase Details

#### Phase 1: Startup & Configuration
**Files:** `backend/cmd/server/main.go`, `backend/internal/config/config.go`

The application initializes by loading environment variables and configuring scrapers:

```
Application Start
      │
      ├─→ Load .env (API keys, blockchain config)
      ├─→ Initialize YouTubeScraper with API key
      ├─→ Initialize GeminiScraper with API key
      └─→ Load search queries from config
```

**Key Functions:**
- `godotenv.Load()` - Loads environment variables
- `config.GetSearchQueries()` - Returns configured search terms
- `config.DefaultSettings()` - Returns scraping limits

---

#### Phase 2: YouTube Scraping
**Files:** `backend/internal/scrapers/youtube.go`

Fetches video metadata and comments using YouTube Data API v3:

```
YouTubeScraper.ScrapeAll()
      │
      ├─→ SearchVideos(query) ──→ YouTube search.list API
      │         │
      │         └─→ Returns: Video IDs, titles, channels
      │
      ├─→ GetVideoDetails(ids) ──→ YouTube videos.list API
      │         │
      │         └─→ Returns: View counts, descriptions, dates
      │
      └─→ GetVideoComments(id) ──→ YouTube commentThreads.list API
                │
                └─→ Returns: Comment text, likes, author
```

**Output:** `data/youtube_latest_results.json`

---

#### Phase 3: Gemini AI Analysis
**Files:** `backend/internal/scrapers/gemini.go`

Uses Google's Gemini AI with Google Search grounding for intelligent web search:

```
GeminiScraper.SearchMultipleQueries()
      │
      └─→ For each query:
            │
            ├─→ SearchComplaintsWithAI(query)
            │         │
            │         ├─→ genai.Client with GoogleSearch tool
            │         └─→ AI analyzes search results
            │
            └─→ Returns AIOverviewResult:
                  ├─ Summary (AI-generated overview)
                  ├─ KeyComplaints[] (categorized issues)
                  ├─ Sources[] (URLs with snippets)
                  └─ SentimentBreakdown (neg/neu/pos ratios)
```

**Output:** `data/gemini_latest_results.json`

---

#### Phase 4: Data Aggregation
**Files:** `backend/cmd/server/main.go`, `backend/internal/models/complaint.go`

Combines YouTube and Gemini results into unified structures:

```
┌─────────────────┐     ┌─────────────────┐
│ YouTube Results │     │  Gemini Results │
│  - Videos       │     │  - AI Summaries │
│  - Comments     │     │  - Key Issues   │
└────────┬────────┘     └────────┬────────┘
         │                       │
         └───────────┬───────────┘
                     ▼
           ┌─────────────────┐
           │  ScrapeResult   │
           │  - VideoCount   │
           │  - CommentCount │
           │  - Categories   │
           │  - Sentiment    │
           └─────────────────┘
```

**Key Types:** `ScrapeResult`, `Complaint`, `YouTubeVideo`, `YouTubeComment`

---

#### Phase 5: Issue Detection
**Files:** `backend/internal/services/resolution.go`, `backend/internal/models/resolution.go`

Analyzes aggregated data to identify trackable issues:

```
ResolutionService.CreateIssue()
      │
      ├─→ Input: Exchange, category, initial complaint count
      │
      ├─→ Generate unique Issue ID
      │
      ├─→ Store in tracked issues map
      │
      └─→ Returns: Issue{
              ID: "issue_abc123",
              Exchange: "coinbase",
              Category: "high_fees",
              Status: "active",
              InitialComplaintCount: 150
           }
```

**Example Issue:**
- Exchange: "coinbase"
- Category: "account_frozen"
- Initial Complaints: 150
- Status: "active"

---

#### Phase 6: Time-Based Monitoring
**Files:** `backend/internal/services/resolution.go`

The system monitors complaint trends over configurable time periods:

```
Time Period (e.g., 7 days)
      │
      ├─→ Day 1: 150 complaints (baseline)
      ├─→ Day 3: 120 complaints (20% decrease)
      ├─→ Day 5: 80 complaints (47% decrease)
      └─→ Day 7: 22 complaints (85% decrease) ✓ Threshold met!
```

**Resolution Criteria:**
- Complaint decrease: ≥70%
- Confidence score: ≥85%
- Minimum observation period: 7 days

---

#### Phase 7: Resolution Creation
**Files:** `backend/internal/services/resolution.go`, `backend/internal/models/resolution.go`

When criteria are met, creates a resolution with bundled evidence:

```
ResolutionService.CreateResolution()
      │
      ├─→ Validate: meetsResolutionCriteria()
      │         │
      │         ├─→ Check complaint decrease %
      │         ├─→ Check confidence score
      │         └─→ Check time period
      │
      ├─→ Bundle evidence:
      │         │
      │         └─→ ResolutionEvidence{
      │               InitialCount: 150,
      │               FinalCount: 22,
      │               DecreasePercent: 85.3,
      │               ConfidenceScore: 0.92,
      │               Sources: ["youtube", "reddit", "trustpilot"],
      │               DataPoints: [...sentiment data...]
      │             }
      │
      └─→ Returns: Resolution (ready for attestation)
```

---

#### Phase 8: Blockchain Attestation (Write Operation)
**Files:** `backend/internal/services/blockchain.go`, `contracts/ResolutionAttestation.sol`

Records the resolution proof on-chain — **happens once per resolution**:

```
ResolutionService.AttestResolution()
      │
      ├─→ BlockchainService.HashEvidence(evidence)
      │         │
      │         └─→ Keccak256(JSON(evidence)) → 0x93fa2c...b81e
      │
      ├─→ BlockchainService.RecordAttestation()
      │         │
      │         ├─→ Build transaction with ABI encoding
      │         ├─→ Sign with private key (EIP-155)
      │         ├─→ Submit to Base network
      │         └─→ Wait for receipt (confirmation)
      │
      └─→ Smart Contract executes:
            │
            ├─→ recordResolution(exchange, issueType, hash)
            ├─→ Store in attestations mapping
            ├─→ Emit ResolutionRecorded event
            └─→ Returns: Transaction ID (0xabc123...)
```

**What Gets Stored On-Chain:**
- ✅ Evidence hash (32 bytes)
- ✅ Exchange name
- ✅ Issue type
- ✅ Timestamp
- ❌ NOT the actual evidence data (too expensive)

---

#### Phase 9: Verification (Read Operation)
**Files:** `backend/internal/services/blockchain.go`, `backend/internal/api/handlers/blockchain.go`

Allows **anyone** to verify a resolution's authenticity **at any time**:

```
VerifyAttestation Request
      │
      ├─→ Input: Attestation ID + Original Evidence
      │
      ├─→ Step 1: Recalculate hash from evidence
      │         │
      │         └─→ Keccak256(evidence) → new_hash
      │
      ├─→ Step 2: Read stored hash from blockchain
      │         │
      │         └─→ contract.getAttestation(id) → stored_hash
      │
      ├─→ Step 3: Compare hashes
      │         │
      │         ├─→ new_hash == stored_hash?
      │         │         │
      │         │         ├─→ ✅ YES: Evidence is authentic
      │         │         └─→ ❌ NO: Evidence was tampered!
      │         │
      │         └─→ Also verify: timestamp, exchange match
      │
      └─→ Returns: VerificationResponse{
              Valid: true,
              OnChainHash: "0x93fa2c...",
              CalculatedHash: "0x93fa2c...",
              TransactionID: "0xabc123...",
              BlockNumber: 12345678
           }
```

**Key Difference from Phase 8:**

| Aspect | Phase 8 (Attestation) | Phase 9 (Verification) |
|--------|----------------------|------------------------|
| **Operation** | Write | Read |
| **Frequency** | Once per resolution | Unlimited times |
| **Gas Cost** | ~50,000 gas | Free (view function) |
| **Purpose** | Create proof | Prove authenticity |
| **Who calls** | Backend (automated) | Anyone (users, auditors) |

---

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Off-Chain (Backend)                      │
├─────────────────────────────────────────────────────────────┤
│  Issue Detected → Resolution Evidence → Keccak256 Hash      │
│                                              │               │
│  {complaints: 150→22, decrease: 85%...}     ▼               │
│                                    0x93fa2c...b81e          │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼ (One transaction)
┌─────────────────────────────────────────────────────────────┐
│                  On-Chain (Base L2)                          │
├─────────────────────────────────────────────────────────────┤
│  ResolutionAttestation.sol                                   │
│  ├─ recordResolution(exchange, issue, hash)                  │
│  ├─ verifyHash(hash) → bool                                  │
│  └─ getAttestation(id) → full details                        │
└─────────────────────────────────────────────────────────────┘
```

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/issues` | GET | List all tracked issues |
| `/api/issues` | POST | Create a new issue |
| `/api/resolutions` | POST | Record a resolution with evidence |
| `/api/attestations` | POST | Record resolution on-chain |
| `/api/attestations/verify` | POST | Verify hash exists on-chain |
| `/api/blockchain/info` | GET | Get chain & wallet info |

### Supported Networks

| Network | Chain ID | Status |
|---------|----------|--------|
| Base Sepolia | 84532 | ✅ Testnet |
| Base Mainnet | 8453 | ⚙️ Production-ready |
| Ethereum Sepolia | 11155111 | ✅ Testnet |

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

**Tanisha** - [@tasnint](https://github.com/tasnint)

Project Link: [https://github.com/tasnint/Coinsights](https://github.com/tasnint/Coinsights)
