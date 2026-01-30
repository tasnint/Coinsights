package models

import "time"

// ============================================
// RESOLUTION MODELS
// ============================================

// Resolution represents a resolved issue with evidence
type Resolution struct {
	ID               string             `json:"id"`
	Exchange         string             `json:"exchange"`          // "coinbase", "kraken", etc.
	IssueCategory    string             `json:"issue_category"`    // "withdrawal_delays", "support_issues", etc.
	Summary          string             `json:"summary"`           // Human-readable resolution summary
	Evidence         ResolutionEvidence `json:"evidence"`          // Structured evidence
	Confidence       float64            `json:"confidence"`        // 0.0-1.0 confidence score
	ResolutionWindow int                `json:"resolution_window"` // Days over which resolution was measured
	Status           string             `json:"status"`            // "pending", "verified", "on_chain"
	CreatedAt        time.Time          `json:"created_at"`
	VerifiedAt       *time.Time         `json:"verified_at,omitempty"`
	Attestation      *Attestation       `json:"attestation,omitempty"` // On-chain attestation (if recorded)
}

// ResolutionEvidence contains the data that gets hashed for on-chain attestation
type ResolutionEvidence struct {
	ComplaintsBefore    int       `json:"complaints_before"`   // Complaint count at start of window
	ComplaintsAfter     int       `json:"complaints_after"`    // Complaint count at end of window
	PercentageDecrease  float64   `json:"percentage_decrease"` // % drop in complaints
	SentimentShift      float64   `json:"sentiment_shift"`     // Change in avg sentiment (-1 to 1)
	SampleComplaints    []string  `json:"sample_complaints"`   // Representative complaint IDs
	DataSources         []string  `json:"data_sources"`        // Where data came from
	MeasurementStart    time.Time `json:"measurement_start"`
	MeasurementEnd      time.Time `json:"measurement_end"`
	AnalysisMethodology string    `json:"analysis_methodology"` // Brief description
}

// ResolutionCriteria defines thresholds for auto-resolution
type ResolutionCriteria struct {
	MinPercentageDecrease    float64 `json:"min_percentage_decrease"` // e.g., 0.70 (70% drop)
	MinConfidence            float64 `json:"min_confidence"`          // e.g., 0.85
	MinWindowDays            int     `json:"min_window_days"`         // e.g., 7 days
	RequirePositiveSentiment bool    `json:"require_positive_sentiment"`
}

// DefaultResolutionCriteria returns sensible defaults
func DefaultResolutionCriteria() ResolutionCriteria {
	return ResolutionCriteria{
		MinPercentageDecrease:    0.70, // 70% drop required
		MinConfidence:            0.85, // 85% confidence
		MinWindowDays:            7,    // Over 7 days
		RequirePositiveSentiment: false,
	}
}

// ============================================
// ON-CHAIN ATTESTATION MODELS
// ============================================

// Attestation represents an on-chain verification record
type Attestation struct {
	ID              uint64    `json:"id"`                      // On-chain attestation ID
	TransactionHash string    `json:"transaction_hash"`        // Ethereum tx hash
	BlockNumber     uint64    `json:"block_number"`            // Block number
	BlockTimestamp  time.Time `json:"block_timestamp"`         // Block timestamp
	ChainID         int64     `json:"chain_id"`                // Network chain ID
	ContractAddress string    `json:"contract_address"`        // Attestation contract address
	EvidenceHash    string    `json:"evidence_hash"`           // Keccak256 hash (hex)
	PreviousHash    string    `json:"previous_hash,omitempty"` // Previous attestation hash
	Attestor        string    `json:"attestor"`                // Address that submitted
	ExplorerURL     string    `json:"explorer_url"`            // Link to block explorer
	Verified        bool      `json:"verified"`                // Whether verification succeeded
}

// AttestationRequest is used to request a new attestation
type AttestationRequest struct {
	ResolutionID  string `json:"resolution_id"`
	Exchange      string `json:"exchange"`
	IssueCategory string `json:"issue_category"`
}

// AttestationResponse is returned after recording an attestation
type AttestationResponse struct {
	Success     bool         `json:"success"`
	Attestation *Attestation `json:"attestation,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// VerificationRequest is used to verify an existing attestation
type VerificationRequest struct {
	EvidenceHash string `json:"evidence_hash"` // Hash to verify
	// OR
	ResolutionID string `json:"resolution_id"` // Resolution to verify
}

// VerificationResponse is returned after verification
type VerificationResponse struct {
	Verified       bool         `json:"verified"`
	OnChain        bool         `json:"on_chain"`
	Attestation    *Attestation `json:"attestation,omitempty"`
	HashMatch      bool         `json:"hash_match"`      // Local hash matches on-chain
	TimestampValid bool         `json:"timestamp_valid"` // Timestamp is reasonable
	Message        string       `json:"message"`
}

// ============================================
// BLOCKCHAIN NETWORK CONFIGURATION
// ============================================

// ChainConfig holds configuration for a specific blockchain network
type ChainConfig struct {
	Name            string `json:"name"`
	ChainID         int64  `json:"chain_id"`
	RPCURL          string `json:"rpc_url"`
	ExplorerURL     string `json:"explorer_url"`
	ContractAddress string `json:"contract_address"`
	IsTestnet       bool   `json:"is_testnet"`
}

// SupportedChains returns configurations for supported networks
func SupportedChains() map[string]ChainConfig {
	return map[string]ChainConfig{
		"base_sepolia": {
			Name:        "Base Sepolia",
			ChainID:     84532,
			RPCURL:      "https://sepolia.base.org",
			ExplorerURL: "https://sepolia.basescan.org",
			IsTestnet:   true,
		},
		"base_mainnet": {
			Name:        "Base",
			ChainID:     8453,
			RPCURL:      "https://mainnet.base.org",
			ExplorerURL: "https://basescan.org",
			IsTestnet:   false,
		},
		"ethereum_sepolia": {
			Name:        "Ethereum Sepolia",
			ChainID:     11155111,
			RPCURL:      "https://rpc.sepolia.org",
			ExplorerURL: "https://sepolia.etherscan.io",
			IsTestnet:   true,
		},
	}
}

// ============================================
// ISSUE TRACKING MODELS
// ============================================

// Issue represents a detected issue being tracked
type Issue struct {
	ID             string       `json:"id"`
	Exchange       string       `json:"exchange"`
	Category       string       `json:"category"`
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	FirstDetected  time.Time    `json:"first_detected"`
	LastUpdated    time.Time    `json:"last_updated"`
	ComplaintCount int          `json:"complaint_count"`
	Severity       string       `json:"severity"` // "critical", "high", "medium", "low"
	Status         string       `json:"status"`   // "active", "investigating", "resolved", "verified"
	Resolution     *Resolution  `json:"resolution,omitempty"`
	Attestation    *Attestation `json:"attestation,omitempty"`
}

// IssueTimeline represents the history of an issue
type IssueTimeline struct {
	IssueID string               `json:"issue_id"`
	Events  []IssueTimelineEvent `json:"events"`
}

// IssueTimelineEvent is a single event in an issue's history
type IssueTimelineEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"event_type"` // "detected", "updated", "resolved", "attested"
	Description string    `json:"description"`
	Data        any       `json:"data,omitempty"`
}
