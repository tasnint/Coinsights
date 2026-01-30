// API for react dashboard to manage issues & proofs
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tasnint/coinsights/internal/models"
	"github.com/tasnint/coinsights/internal/services"
)

// BlockchainHandler handles blockchain-related API endpoints
type BlockchainHandler struct {
	resolutionService *services.ResolutionService
	blockchainService *services.BlockchainService
}

// NewBlockchainHandler creates a new blockchain handler
func NewBlockchainHandler(
	resolutionService *services.ResolutionService,
	blockchainService *services.BlockchainService,
) *BlockchainHandler {
	return &BlockchainHandler{
		resolutionService: resolutionService,
		blockchainService: blockchainService,
	}
}

// ============================================
// ISSUE ENDPOINTS
// ============================================

// CreateIssue handles POST /api/issues
func (h *BlockchainHandler) CreateIssue(w http.ResponseWriter, r *http.Request) {
	var issue models.Issue
	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	created, err := h.resolutionService.CreateIssue(&issue)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, created)
}

// GetIssue handles GET /api/issues/{id}
func (h *BlockchainHandler) GetIssue(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Issue ID required")
		return
	}

	issue, err := h.resolutionService.GetIssue(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, issue)
}

// ListIssues handles GET /api/issues
func (h *BlockchainHandler) ListIssues(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	issues := h.resolutionService.ListIssues(status)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"issues": issues,
		"count":  len(issues),
	})
}

// ============================================
// RESOLUTION ENDPOINTS
// ============================================

// CreateResolutionRequest is the request body for creating a resolution
type CreateResolutionRequest struct {
	IssueID  string                    `json:"issue_id"`
	Summary  string                    `json:"summary"`
	Evidence models.ResolutionEvidence `json:"evidence"`
}

// CreateResolution handles POST /api/resolutions
func (h *BlockchainHandler) CreateResolution(w http.ResponseWriter, r *http.Request) {
	var req CreateResolutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resolution, err := h.resolutionService.CreateResolution(
		r.Context(),
		req.IssueID,
		&req.Evidence,
		req.Summary,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, resolution)
}

// GetResolution handles GET /api/resolutions/{id}
func (h *BlockchainHandler) GetResolution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Resolution ID required")
		return
	}

	resolution, err := h.resolutionService.GetResolution(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, resolution)
}

// ListResolutions handles GET /api/resolutions
func (h *BlockchainHandler) ListResolutions(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	resolutions := h.resolutionService.ListResolutions(status)
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"resolutions": resolutions,
		"count":       len(resolutions),
	})
}

// ============================================
// ATTESTATION ENDPOINTS
// ============================================

// AttestResolution handles POST /api/attestations
func (h *BlockchainHandler) AttestResolution(w http.ResponseWriter, r *http.Request) {
	var req models.AttestationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	attestation, err := h.resolutionService.AttestResolution(r.Context(), req.ResolutionID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, models.AttestationResponse{
		Success:     true,
		Attestation: attestation,
	})
}

// VerifyAttestation handles POST /api/attestations/verify
func (h *BlockchainHandler) VerifyAttestation(w http.ResponseWriter, r *http.Request) {
	var req models.VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var response *models.VerificationResponse
	var err error

	if req.EvidenceHash != "" {
		response, err = h.resolutionService.VerifyByHash(r.Context(), req.EvidenceHash)
	} else if req.ResolutionID != "" {
		response, err = h.resolutionService.VerifyResolution(r.Context(), req.ResolutionID)
	} else {
		respondError(w, http.StatusBadRequest, "Either evidence_hash or resolution_id required")
		return
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, response)
}

// GetAttestationByResolution handles GET /api/resolutions/{id}/attestation
func (h *BlockchainHandler) GetAttestationByResolution(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Resolution ID required")
		return
	}

	resolution, err := h.resolutionService.GetResolution(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	if resolution.Attestation == nil {
		respondError(w, http.StatusNotFound, "Resolution not yet attested")
		return
	}

	respondJSON(w, http.StatusOK, resolution.Attestation)
}

// ============================================
// BLOCKCHAIN INFO ENDPOINTS
// ============================================

// GetChainInfo handles GET /api/blockchain/info
func (h *BlockchainHandler) GetChainInfo(w http.ResponseWriter, r *http.Request) {
	if h.blockchainService == nil {
		respondError(w, http.StatusServiceUnavailable, "Blockchain service not configured")
		return
	}

	chainInfo := h.blockchainService.GetChainInfo()
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"chain":            chainInfo,
		"wallet_address":   h.blockchainService.GetWalletAddress(),
		"supported_chains": models.SupportedChains(),
	})
}

// GetStats handles GET /api/blockchain/stats
func (h *BlockchainHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := h.resolutionService.GetStats()
	respondJSON(w, http.StatusOK, stats)
}

// HashEvidence handles POST /api/blockchain/hash
// Useful for pre-computing hashes before attestation
func (h *BlockchainHandler) HashEvidence(w http.ResponseWriter, r *http.Request) {
	if h.blockchainService == nil {
		respondError(w, http.StatusServiceUnavailable, "Blockchain service not configured")
		return
	}

	var evidence models.ResolutionEvidence
	if err := json.NewDecoder(r.Body).Decode(&evidence); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	hash, err := h.blockchainService.HashEvidence(&evidence)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"hash": hash,
	})
}

// ============================================
// DEMO / TEST ENDPOINTS
// ============================================

// CreateDemoIssueAndResolve handles POST /api/demo/full-workflow
// This demonstrates the complete workflow for testing
func (h *BlockchainHandler) CreateDemoIssueAndResolve(w http.ResponseWriter, r *http.Request) {
	// Step 1: Create a demo issue
	issue := &models.Issue{
		Exchange:       "coinbase",
		Category:       "withdrawal_delays",
		Title:          "Withdrawal Delays - December 2025",
		Description:    "Users reporting significant delays in withdrawal processing",
		ComplaintCount: 150,
		Severity:       "high",
	}

	createdIssue, err := h.resolutionService.CreateIssue(issue)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create issue: "+err.Error())
		return
	}

	// Step 2: Create resolution with evidence
	evidence := &models.ResolutionEvidence{
		ComplaintsBefore:    150,
		ComplaintsAfter:     22,
		PercentageDecrease:  0.85,
		SentimentShift:      0.3,
		SampleComplaints:    []string{"complaint_001", "complaint_002", "complaint_003"},
		DataSources:         []string{"youtube", "google", "reddit"},
		MeasurementStart:    time.Now().AddDate(0, 0, -14),
		MeasurementEnd:      time.Now().AddDate(0, 0, -7),
		AnalysisMethodology: "Complaint volume tracking with sentiment analysis over 7-day rolling window",
	}

	resolution, err := h.resolutionService.CreateResolution(
		r.Context(),
		createdIssue.ID,
		evidence,
		"Withdrawal delays resolved. Complaint volume decreased by 85% over 7 days. Coinbase appears to have improved their withdrawal processing infrastructure.",
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create resolution: "+err.Error())
		return
	}

	// Step 3: Compute hash (show what would be attested)
	var hash string
	if h.blockchainService != nil {
		hash, _ = h.blockchainService.HashEvidence(evidence)
	}

	// Return the complete workflow result
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"workflow_complete": true,
		"issue":             createdIssue,
		"resolution":        resolution,
		"evidence_hash":     hash,
		"next_steps": map[string]string{
			"attest": "POST /api/attestations with {resolution_id: \"" + resolution.ID + "\"}",
			"verify": "POST /api/attestations/verify with {resolution_id: \"" + resolution.ID + "\"}",
		},
	})
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
