import React, { useState, useEffect } from 'react';
import { Issue, BlockchainStats, ChainConfig } from '../types';
import { issuesApi, blockchainApi } from '../services/api';
import IssueCard from '../components/blockchain/IssueCard';
import VerificationPanel from '../components/blockchain/VerificationPanel';
import '../styles/Verification.css';

/**
 * Verification page displays all tracked issues and their on-chain verification status.
 * Provides tools for recording attestations and verifying them.
 */
const Verification: React.FC = () => {
  const [issues, setIssues] = useState<Issue[]>([]);
  const [stats, setStats] = useState<BlockchainStats | null>(null);
  const [chainInfo, setChainInfo] = useState<{
    chain: ChainConfig;
    wallet_address: string;
  } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'issues' | 'verify'>('issues');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [runningDemo, setRunningDemo] = useState(false);

  useEffect(() => {
    loadData();
  }, [statusFilter]);

  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      const [issuesRes, statsRes, chainRes] = await Promise.all([
        issuesApi.getAll(statusFilter || undefined),
        blockchainApi.getStats(),
        blockchainApi.getInfo().catch(() => null), // May fail if blockchain not configured
      ]);

      setIssues(issuesRes.data.issues || []);
      setStats(statsRes.data);
      if (chainRes) {
        setChainInfo({
          chain: chainRes.data.chain,
          wallet_address: chainRes.data.wallet_address,
        });
      }
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleIssueUpdate = (updatedIssue: Issue) => {
    setIssues(prev => 
      prev.map(issue => 
        issue.id === updatedIssue.id ? updatedIssue : issue
      )
    );
    // Refresh stats
    blockchainApi.getStats().then(res => setStats(res.data));
  };

  const runDemoWorkflow = async () => {
    setRunningDemo(true);
    try {
      await blockchainApi.runDemoWorkflow();
      await loadData(); // Refresh to show new issue
    } catch (err: any) {
      setError(err.response?.data?.error || 'Demo workflow failed');
    } finally {
      setRunningDemo(false);
    }
  };

  return (
    <div className="verification-page">
      {/* Header */}
      <div className="page-header">
        <div className="header-content">
          <h1 className="page-title">
            <span className="title-icon">⛓️</span>
            On-Chain Verification
          </h1>
          <p className="page-subtitle">
            Tamper-proof attestations of resolved issues on Base (Coinbase L2)
          </p>
        </div>

        {chainInfo && (
          <div className="chain-info">
            <div className="chain-badge">
              <span className="chain-icon">🔗</span>
              <span className="chain-name">{chainInfo.chain.name}</span>
              {chainInfo.chain.is_testnet && (
                <span className="testnet-badge">Testnet</span>
              )}
            </div>
            <div className="wallet-info">
              <span className="wallet-label">Attestor:</span>
              <code className="wallet-address">
                {chainInfo.wallet_address.slice(0, 6)}...{chainInfo.wallet_address.slice(-4)}
              </code>
            </div>
          </div>
        )}
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <span className="stat-icon">📋</span>
            <div className="stat-content">
              <span className="stat-value">{stats.total_issues}</span>
              <span className="stat-label">Total Issues</span>
            </div>
          </div>
          <div className="stat-card">
            <span className="stat-icon">✅</span>
            <div className="stat-content">
              <span className="stat-value">{stats.total_resolutions}</span>
              <span className="stat-label">Resolutions</span>
            </div>
          </div>
          <div className="stat-card highlight">
            <span className="stat-icon">⛓️</span>
            <div className="stat-content">
              <span className="stat-value">{stats.attestation_count}</span>
              <span className="stat-label">On-Chain Attestations</span>
            </div>
          </div>
          <div className="stat-card">
            <span className="stat-icon">🔴</span>
            <div className="stat-content">
              <span className="stat-value">{stats.issues_by_status?.active || 0}</span>
              <span className="stat-label">Active Issues</span>
            </div>
          </div>
        </div>
      )}

      {/* Tabs */}
      <div className="tabs">
        <button 
          className={`tab ${activeTab === 'issues' ? 'active' : ''}`}
          onClick={() => setActiveTab('issues')}
        >
          📋 Issues & Resolutions
        </button>
        <button 
          className={`tab ${activeTab === 'verify' ? 'active' : ''}`}
          onClick={() => setActiveTab('verify')}
        >
          🔍 Verify Hash
        </button>
      </div>

      {/* Content */}
      <div className="page-content">
        {activeTab === 'issues' && (
          <>
            {/* Filters */}
            <div className="filters-row">
              <div className="filter-group">
                <label>Status:</label>
                <select 
                  value={statusFilter} 
                  onChange={(e) => setStatusFilter(e.target.value)}
                  className="filter-select"
                >
                  <option value="">All</option>
                  <option value="active">Active</option>
                  <option value="investigating">Investigating</option>
                  <option value="resolved">Resolved</option>
                  <option value="verified">Verified (On-Chain)</option>
                </select>
              </div>

              <button 
                className="demo-button"
                onClick={runDemoWorkflow}
                disabled={runningDemo}
              >
                {runningDemo ? '⏳ Running...' : '🧪 Run Demo Workflow'}
              </button>
            </div>

            {/* Issues List */}
            {loading ? (
              <div className="loading">Loading issues...</div>
            ) : error ? (
              <div className="error-message">{error}</div>
            ) : issues.length === 0 ? (
              <div className="empty-state">
                <span className="empty-icon">📭</span>
                <h3>No Issues Found</h3>
                <p>Click "Run Demo Workflow" to create a sample issue with resolution.</p>
              </div>
            ) : (
              <div className="issues-grid">
                {issues.map(issue => (
                  <IssueCard 
                    key={issue.id} 
                    issue={issue}
                    onAttestationCreated={handleIssueUpdate}
                  />
                ))}
              </div>
            )}
          </>
        )}

        {activeTab === 'verify' && (
          <div className="verify-section">
            <VerificationPanel />
            
            <div className="verify-info">
              <h3>🔐 How Verification Works</h3>
              <ol>
                <li>
                  <strong>Evidence Hashing:</strong> Resolution evidence (complaint counts, 
                  sentiment data, sources) is hashed using Keccak256
                </li>
                <li>
                  <strong>On-Chain Commit:</strong> The hash is stored on Base (Coinbase L2) 
                  via our smart contract
                </li>
                <li>
                  <strong>Independent Verification:</strong> Anyone can verify the hash 
                  exists on-chain without trusting our backend
                </li>
              </ol>
              
              <div className="code-example">
                <code>
                  {`// Evidence gets hashed like this:
const evidence = {
  complaints_before: 150,
  complaints_after: 22,
  percentage_decrease: 0.85,
  ...
};
const hash = keccak256(JSON.stringify(evidence));
// → 0x93fa2c...b81e`}
                </code>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Verification;
