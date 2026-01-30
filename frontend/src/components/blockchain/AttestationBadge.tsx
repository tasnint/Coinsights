import React from 'react';
import { Attestation } from '../../types';
import './AttestationBadge.css';

interface AttestationBadgeProps {
  attestation?: Attestation | null;
  status: 'active' | 'investigating' | 'resolved' | 'verified' | 'pending' | 'on_chain';
  compact?: boolean;
}

/**
 * AttestationBadge displays the verification status of a resolution.
 * Shows on-chain verification details with explorer link.
 */
const AttestationBadge: React.FC<AttestationBadgeProps> = ({ attestation, status, compact = false }) => {
  // Determine badge appearance based on status
  const getBadgeConfig = () => {
    if (attestation && attestation.verified) {
      return {
        className: 'badge-verified',
        icon: '🟢',
        label: 'On-Chain Verified',
        showDetails: true,
      };
    }
    
    switch (status) {
      case 'verified':
        return {
          className: 'badge-verified-pending',
          icon: '🔵',
          label: 'Verified (Off-Chain)',
          showDetails: false,
        };
      case 'resolved':
        return {
          className: 'badge-resolved',
          icon: '🟡',
          label: 'Resolved',
          showDetails: false,
        };
      case 'investigating':
        return {
          className: 'badge-investigating',
          icon: '🟠',
          label: 'Investigating',
          showDetails: false,
        };
      case 'on_chain':
        return {
          className: 'badge-verified',
          icon: '🟢',
          label: 'On-Chain',
          showDetails: true,
        };
      default:
        return {
          className: 'badge-active',
          icon: '🔴',
          label: 'Active',
          showDetails: false,
        };
    }
  };

  const config = getBadgeConfig();

  if (compact) {
    return (
      <span className={`attestation-badge compact ${config.className}`}>
        <span className="badge-icon">{config.icon}</span>
        <span className="badge-label">{config.label}</span>
      </span>
    );
  }

  return (
    <div className={`attestation-badge ${config.className}`}>
      <div className="badge-header">
        <span className="badge-icon">{config.icon}</span>
        <span className="badge-label">{config.label}</span>
      </div>
      
      {config.showDetails && attestation && (
        <div className="badge-details">
          <div className="detail-row">
            <span className="detail-label">Tx:</span>
            <a 
              href={attestation.explorer_url} 
              target="_blank" 
              rel="noopener noreferrer"
              className="tx-hash"
            >
              {attestation.transaction_hash.slice(0, 10)}...{attestation.transaction_hash.slice(-8)}
            </a>
          </div>
          
          <div className="detail-row">
            <span className="detail-label">Block:</span>
            <span className="detail-value">#{attestation.block_number.toLocaleString()}</span>
          </div>
          
          <div className="detail-row">
            <span className="detail-label">Time:</span>
            <span className="detail-value">
              {new Date(attestation.block_timestamp).toLocaleString()}
            </span>
          </div>
          
          <a 
            href={attestation.explorer_url} 
            target="_blank" 
            rel="noopener noreferrer"
            className="explorer-link"
          >
            View on Explorer →
          </a>
        </div>
      )}
    </div>
  );
};

export default AttestationBadge;
