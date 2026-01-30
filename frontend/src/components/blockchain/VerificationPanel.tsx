import React, { useState } from 'react';
import { attestationsApi } from '../../services/api';
import { VerificationResponse } from '../../types';
import './VerificationPanel.css';

interface VerificationPanelProps {
  resolutionId?: string;
  evidenceHash?: string;
  initialVerification?: VerificationResponse;
}

/**
 * VerificationPanel allows users to verify attestations on-chain.
 * Provides both resolution ID and hash-based verification.
 */
const VerificationPanel: React.FC<VerificationPanelProps> = ({
  resolutionId,
  evidenceHash,
  initialVerification,
}) => {
  const [hashInput, setHashInput] = useState(evidenceHash || '');
  const [verification, setVerification] = useState<VerificationResponse | null>(
    initialVerification || null
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleVerify = async () => {
    setLoading(true);
    setError(null);

    try {
      const params = resolutionId 
        ? { resolution_id: resolutionId }
        : { evidence_hash: hashInput };
      
      const response = await attestationsApi.verify(params);
      setVerification(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Verification failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="verification-panel">
      <h3 className="panel-title">
        <span className="title-icon">🔍</span>
        On-Chain Verification
      </h3>

      {!resolutionId && (
        <div className="hash-input-section">
          <label htmlFor="evidence-hash">Evidence Hash</label>
          <div className="input-row">
            <input
              id="evidence-hash"
              type="text"
              value={hashInput}
              onChange={(e) => setHashInput(e.target.value)}
              placeholder="0x..."
              className="hash-input"
            />
            <button 
              onClick={handleVerify} 
              disabled={loading || !hashInput}
              className="verify-button"
            >
              {loading ? 'Verifying...' : 'Verify'}
            </button>
          </div>
          <p className="input-hint">
            Enter the Keccak256 hash of the resolution evidence
          </p>
        </div>
      )}

      {resolutionId && (
        <div className="resolution-verify-section">
          <p className="resolution-id">
            Resolution: <code>{resolutionId}</code>
          </p>
          <button 
            onClick={handleVerify} 
            disabled={loading}
            className="verify-button full-width"
          >
            {loading ? 'Verifying...' : 'Verify On-Chain'}
          </button>
        </div>
      )}

      {error && (
        <div className="verification-error">
          <span className="error-icon">❌</span>
          {error}
        </div>
      )}

      {verification && (
        <div className={`verification-result ${verification.verified ? 'success' : 'failure'}`}>
          <div className="result-header">
            <span className="result-icon">
              {verification.verified ? '✅' : '❌'}
            </span>
            <span className="result-title">
              {verification.verified ? 'Verified' : 'Not Verified'}
            </span>
          </div>

          <div className="result-details">
            <div className="detail-item">
              <span className="detail-label">On-Chain:</span>
              <span className={`detail-value ${verification.on_chain ? 'yes' : 'no'}`}>
                {verification.on_chain ? 'Yes' : 'No'}
              </span>
            </div>

            <div className="detail-item">
              <span className="detail-label">Hash Match:</span>
              <span className={`detail-value ${verification.hash_match ? 'yes' : 'no'}`}>
                {verification.hash_match ? 'Yes' : 'No'}
              </span>
            </div>

            <div className="detail-item">
              <span className="detail-label">Timestamp Valid:</span>
              <span className={`detail-value ${verification.timestamp_valid ? 'yes' : 'no'}`}>
                {verification.timestamp_valid ? 'Yes' : 'N/A'}
              </span>
            </div>
          </div>

          <p className="result-message">{verification.message}</p>

          {verification.attestation && (
            <div className="attestation-details">
              <h4>Attestation Details</h4>
              <table className="attestation-table">
                <tbody>
                  <tr>
                    <td>Transaction</td>
                    <td>
                      <a 
                        href={verification.attestation.explorer_url}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        {verification.attestation.transaction_hash.slice(0, 16)}...
                      </a>
                    </td>
                  </tr>
                  <tr>
                    <td>Block</td>
                    <td>#{verification.attestation.block_number.toLocaleString()}</td>
                  </tr>
                  <tr>
                    <td>Timestamp</td>
                    <td>{new Date(verification.attestation.block_timestamp).toLocaleString()}</td>
                  </tr>
                  <tr>
                    <td>Attestor</td>
                    <td className="monospace">{verification.attestation.attestor}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      <div className="panel-footer">
        <p className="footer-text">
          Attestations are recorded on Base (Coinbase L2) for tamper-proof verification.
        </p>
      </div>
    </div>
  );
};

export default VerificationPanel;
