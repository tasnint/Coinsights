// Data shape for issues, resolutions and proofs
package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/tasnint/coinsights/internal/models"
)

// ResolutionService manages issue resolutions and their attestations
type ResolutionService struct {
	blockchain  *BlockchainService
	resolutions map[string]*models.Resolution // In-memory store (replace with DB)
	issues      map[string]*models.Issue      // In-memory store (replace with DB)
	criteria    models.ResolutionCriteria
	mu          sync.RWMutex
}

// NewResolutionService creates a new resolution service
func NewResolutionService(blockchain *BlockchainService) *ResolutionService {
	return &ResolutionService{
		blockchain:  blockchain,
		resolutions: make(map[string]*models.Resolution),
		issues:      make(map[string]*models.Issue),
		criteria:    models.DefaultResolutionCriteria(),
	}
}

// ============================================
// ISSUE MANAGEMENT
// ============================================

// CreateIssue creates a new issue being tracked
func (rs *ResolutionService) CreateIssue(issue *models.Issue) (*models.Issue, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Generate ID if not set
	if issue.ID == "" {
		issue.ID = generateID()
	}

	issue.FirstDetected = time.Now()
	issue.LastUpdated = time.Now()
	issue.Status = "active"

	rs.issues[issue.ID] = issue
	return issue, nil
}

// GetIssue retrieves an issue by ID
func (rs *ResolutionService) GetIssue(id string) (*models.Issue, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	issue, ok := rs.issues[id]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s", id)
	}
	return issue, nil
}

// ListIssues returns all tracked issues
func (rs *ResolutionService) ListIssues(status string) []*models.Issue {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	var results []*models.Issue
	for _, issue := range rs.issues {
		if status == "" || issue.Status == status {
			results = append(results, issue)
		}
	}
	return results
}

// UpdateIssue updates an existing issue
func (rs *ResolutionService) UpdateIssue(id string, update *models.Issue) (*models.Issue, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	issue, ok := rs.issues[id]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s", id)
	}

	// Update fields
	if update.ComplaintCount > 0 {
		issue.ComplaintCount = update.ComplaintCount
	}
	if update.Severity != "" {
		issue.Severity = update.Severity
	}
	if update.Status != "" {
		issue.Status = update.Status
	}
	if update.Description != "" {
		issue.Description = update.Description
	}
	issue.LastUpdated = time.Now()

	return issue, nil
}

// ============================================
// RESOLUTION MANAGEMENT
// ============================================

// CreateResolution creates a new resolution for an issue
func (rs *ResolutionService) CreateResolution(
	ctx context.Context,
	issueID string,
	evidence *models.ResolutionEvidence,
	summary string,
) (*models.Resolution, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	// Get the issue
	issue, ok := rs.issues[issueID]
	if !ok {
		return nil, fmt.Errorf("issue not found: %s", issueID)
	}

	// Calculate confidence score
	confidence := rs.calculateConfidence(evidence)

	resolution := &models.Resolution{
		ID:               generateID(),
		Exchange:         issue.Exchange,
		IssueCategory:    issue.Category,
		Summary:          summary,
		Evidence:         *evidence,
		Confidence:       confidence,
		ResolutionWindow: int(evidence.MeasurementEnd.Sub(evidence.MeasurementStart).Hours() / 24),
		Status:           "pending",
		CreatedAt:        time.Now(),
	}

	// Check if meets criteria for auto-verification
	if rs.meetsResolutionCriteria(resolution) {
		resolution.Status = "verified"
		now := time.Now()
		resolution.VerifiedAt = &now
	}

	rs.resolutions[resolution.ID] = resolution

	// Update issue status
	issue.Status = "resolved"
	issue.Resolution = resolution
	issue.LastUpdated = time.Now()

	return resolution, nil
}

// GetResolution retrieves a resolution by ID
func (rs *ResolutionService) GetResolution(id string) (*models.Resolution, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	resolution, ok := rs.resolutions[id]
	if !ok {
		return nil, fmt.Errorf("resolution not found: %s", id)
	}
	return resolution, nil
}

// ListResolutions returns all resolutions
func (rs *ResolutionService) ListResolutions(status string) []*models.Resolution {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	var results []*models.Resolution
	for _, resolution := range rs.resolutions {
		if status == "" || resolution.Status == status {
			results = append(results, resolution)
		}
	}
	return results
}

// ============================================
// ON-CHAIN ATTESTATION
// ============================================

// AttestResolution records a resolution on the blockchain
func (rs *ResolutionService) AttestResolution(ctx context.Context, resolutionID string) (*models.Attestation, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	resolution, ok := rs.resolutions[resolutionID]
	if !ok {
		return nil, fmt.Errorf("resolution not found: %s", resolutionID)
	}

	// Check if already attested
	if resolution.Attestation != nil {
		return resolution.Attestation, nil
	}

	// Check if blockchain service is available
	if rs.blockchain == nil {
		return nil, fmt.Errorf("blockchain service not configured")
	}

	// Record attestation
	attestation, err := rs.blockchain.RecordAttestation(ctx, resolution)
	if err != nil {
		return nil, fmt.Errorf("failed to record attestation: %w", err)
	}

	// Update resolution
	resolution.Attestation = attestation
	resolution.Status = "on_chain"

	// Update associated issue if exists
	for _, issue := range rs.issues {
		if issue.Resolution != nil && issue.Resolution.ID == resolutionID {
			issue.Attestation = attestation
			issue.Status = "verified"
			break
		}
	}

	return attestation, nil
}

// VerifyResolution verifies an attestation exists on-chain
func (rs *ResolutionService) VerifyResolution(ctx context.Context, resolutionID string) (*models.VerificationResponse, error) {
	resolution, err := rs.GetResolution(resolutionID)
	if err != nil {
		return nil, err
	}

	if rs.blockchain == nil {
		return nil, fmt.Errorf("blockchain service not configured")
	}

	// Hash the evidence
	evidenceHash, err := rs.blockchain.HashEvidence(&resolution.Evidence)
	if err != nil {
		return nil, fmt.Errorf("failed to hash evidence: %w", err)
	}

	// Verify on chain
	return rs.blockchain.VerifyAttestation(ctx, evidenceHash)
}

// VerifyByHash verifies an attestation by evidence hash
func (rs *ResolutionService) VerifyByHash(ctx context.Context, evidenceHash string) (*models.VerificationResponse, error) {
	if rs.blockchain == nil {
		return nil, fmt.Errorf("blockchain service not configured")
	}

	return rs.blockchain.VerifyAttestation(ctx, evidenceHash)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

// calculateConfidence calculates a confidence score for a resolution
func (rs *ResolutionService) calculateConfidence(evidence *models.ResolutionEvidence) float64 {
	confidence := 0.0

	// Base confidence from percentage decrease
	if evidence.PercentageDecrease >= 0.9 {
		confidence = 0.95
	} else if evidence.PercentageDecrease >= 0.7 {
		confidence = 0.85
	} else if evidence.PercentageDecrease >= 0.5 {
		confidence = 0.70
	} else {
		confidence = 0.50
	}

	// Bonus for positive sentiment shift
	if evidence.SentimentShift > 0.2 {
		confidence += 0.05
	}

	// Bonus for multiple data sources
	if len(evidence.DataSources) >= 3 {
		confidence += 0.03
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// meetsResolutionCriteria checks if a resolution meets auto-verification criteria
func (rs *ResolutionService) meetsResolutionCriteria(resolution *models.Resolution) bool {
	// Check percentage decrease
	if resolution.Evidence.PercentageDecrease < rs.criteria.MinPercentageDecrease {
		return false
	}

	// Check confidence
	if resolution.Confidence < rs.criteria.MinConfidence {
		return false
	}

	// Check window duration
	if resolution.ResolutionWindow < rs.criteria.MinWindowDays {
		return false
	}

	// Check sentiment if required
	if rs.criteria.RequirePositiveSentiment && resolution.Evidence.SentimentShift <= 0 {
		return false
	}

	return true
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ============================================
// STATISTICS
// ============================================

// GetStats returns resolution statistics
func (rs *ResolutionService) GetStats() map[string]interface{} {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	stats := map[string]interface{}{
		"total_issues":      len(rs.issues),
		"total_resolutions": len(rs.resolutions),
		"issues_by_status":  make(map[string]int),
		"attestation_count": 0,
	}

	issuesByStatus := stats["issues_by_status"].(map[string]int)
	attestationCount := 0

	for _, issue := range rs.issues {
		issuesByStatus[issue.Status]++
	}

	for _, resolution := range rs.resolutions {
		if resolution.Attestation != nil {
			attestationCount++
		}
	}

	stats["attestation_count"] = attestationCount

	// Get on-chain count if available
	if rs.blockchain != nil {
		if count, err := rs.blockchain.GetAttestationCount(context.Background()); err == nil {
			stats["on_chain_attestation_count"] = count
		}
	}

	return stats
}
