import React, { useState } from 'react';
import { Issue, Resolution } from '../../types';
import { attestationsApi } from '../../services/api';
import AttestationBadge from './AttestationBadge';
import './IssueCard.css';

interface IssueCardProps {
  issue: Issue;
  onAttestationCreated?: (issue: Issue) => void;
}

/**
 * IssueCard displays an issue with its resolution and attestation status.
 * Provides actions to attest resolved issues on-chain.
 */
const IssueCard: React.FC<IssueCardProps> = ({ issue, onAttestationCreated }) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [expanded, setExpanded] = useState(false);

  const handleAttest = async () => {
    if (!issue.resolution) return;
    
    setLoading(true);
    setError(null);

    try {
      const response = await attestationsApi.create(issue.resolution.id);
      if (response.data.success && onAttestationCreated) {
        // Update the issue with attestation
        const updatedIssue = {
          ...issue,
          status: 'verified' as const,
          attestation: response.data.attestation,
          resolution: {
            ...issue.resolution,
            attestation: response.data.attestation,
            status: 'on_chain' as const,
          },
        };
        onAttestationCreated(updatedIssue);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create attestation');
    } finally {
      setLoading(false);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return '#EF4444';
      case 'high': return '#F97316';
      case 'medium': return '#F59E0B';
      case 'low': return '#10B981';
      default: return '#64748B';
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

  return (
    <div className={`issue-card ${issue.status}`}>
      {/* Header */}
      <div className="card-header">
        <div className="header-left">
          <span 
            className="severity-badge"
            style={{ backgroundColor: getSeverityColor(issue.severity) }}
          >
            {issue.severity.toUpperCase()}
          </span>
          <span className="exchange-badge">{issue.exchange}</span>
          <span className="category-badge">{issue.category.replace(/_/g, ' ')}</span>
        </div>
        <AttestationBadge 
          attestation={issue.attestation || issue.resolution?.attestation} 
          status={issue.status}
          compact
        />
      </div>

      {/* Title & Description */}
      <h3 className="card-title">{issue.title}</h3>
      <p className="card-description">{issue.description}</p>

      {/* Metrics */}
      <div className="card-metrics">
        <div className="metric">
          <span className="metric-value">{issue.complaint_count}</span>
          <span className="metric-label">Complaints</span>
        </div>
        <div className="metric">
          <span className="metric-value">{formatDate(issue.first_detected)}</span>
          <span className="metric-label">Detected</span>
        </div>
        <div className="metric">
          <span className="metric-value">{formatDate(issue.last_updated)}</span>
          <span className="metric-label">Updated</span>
        </div>
      </div>

      {/* Resolution Section */}
      {issue.resolution && (
        <div className="resolution-section">
          <button 
            className="expand-button"
            onClick={() => setExpanded(!expanded)}
          >
            {expanded ? '▼' : '▶'} Resolution Details
          </button>

          {expanded && (
            <div className="resolution-content">
              <p className="resolution-summary">{issue.resolution.summary}</p>
              
              <div className="evidence-grid">
                <div className="evidence-item">
                  <span className="evidence-label">Complaints Before</span>
                  <span className="evidence-value">
                    {issue.resolution.evidence.complaints_before}
                  </span>
                </div>
                <div className="evidence-item">
                  <span className="evidence-label">Complaints After</span>
                  <span className="evidence-value">
                    {issue.resolution.evidence.complaints_after}
                  </span>
                </div>
                <div className="evidence-item highlight">
                  <span className="evidence-label">Decrease</span>
                  <span className="evidence-value">
                    {(issue.resolution.evidence.percentage_decrease * 100).toFixed(0)}%
                  </span>
                </div>
                <div className="evidence-item">
                  <span className="evidence-label">Confidence</span>
                  <span className="evidence-value">
                    {(issue.resolution.confidence * 100).toFixed(0)}%
                  </span>
                </div>
              </div>

              <div className="data-sources">
                <span className="sources-label">Data Sources:</span>
                {issue.resolution.evidence.data_sources.map((source, i) => (
                  <span key={i} className="source-tag">{source}</span>
                ))}
              </div>

              <div className="measurement-window">
                {formatDate(issue.resolution.evidence.measurement_start)} → {' '}
                {formatDate(issue.resolution.evidence.measurement_end)}
                {' '}({issue.resolution.resolution_window} days)
              </div>
            </div>
          )}
        </div>
      )}

      {/* Actions */}
      <div className="card-actions">
        {issue.resolution && !issue.attestation && !issue.resolution.attestation && (
          <button 
            className="attest-button"
            onClick={handleAttest}
            disabled={loading}
          >
            {loading ? 'Recording...' : '⛓️ Record On-Chain'}
          </button>
        )}
        
        {(issue.attestation || issue.resolution?.attestation) && (
          <a 
            href={(issue.attestation || issue.resolution?.attestation)?.explorer_url}
            target="_blank"
            rel="noopener noreferrer"
            className="explorer-button"
          >
            🔗 View on Explorer
          </a>
        )}
      </div>

      {error && <div className="card-error">{error}</div>}
    </div>
  );
};

export default IssueCard;
