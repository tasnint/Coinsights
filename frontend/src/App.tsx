import React, { useState, useEffect } from 'react';

interface Issue {
  id: string;
  exchange: string;
  category: string;
  title: string;
  description: string;
  first_detected: string;
  severity: string;
  status: string;
  count: number;
  examples: string[];
}

interface Resolution {
  id: string;
  exchange: string;
  issue_category: string;
  summary: string;
  status: string;
  created_at: string;
}

interface Stats {
  videos_analyzed: number;
  comments_analyzed: number;
  issues_found: number;
  categories: number;
}

function App() {
  const [issues, setIssues] = useState<Issue[]>([]);
  const [resolutions, setResolutions] = useState<Resolution[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const API_BASE = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [issuesRes, resolutionsRes, statsRes] = await Promise.all([
          fetch(`${API_BASE}/issues`),
          fetch(`${API_BASE}/resolutions`),
          fetch(`${API_BASE}/stats`)
        ]);

        if (issuesRes.ok) {
          const issuesData = await issuesRes.json();
          setIssues(issuesData.issues || []);
        }

        if (resolutionsRes.ok) {
          const resolutionsData = await resolutionsRes.json();
          setResolutions(resolutionsData.resolutions || []);
        }

        if (statsRes.ok) {
          const statsData = await statsRes.json();
          setStats(statsData);
        }

        setError(null);
      } catch (err) {
        setError('Backend API not running. Start the API server with: cd backend/cmd/api && go run main.go');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [API_BASE]);

  // Sort issues by count (most complaints first)
  const sortedIssues = [...issues].sort((a, b) => b.count - a.count);

  return (
    <div style={{ fontFamily: 'monospace', padding: '20px', maxWidth: '1000px' }}>
      <h1>Coinsights - Blockchain Issue Tracker</h1>
      <p>
        A backend-focused project demonstrating interest in <strong>Coinbase</strong> and <strong>blockchain technologies</strong>.
        This system scrapes cryptocurrency exchange complaints, analyzes them with AI, and creates on-chain attestations for verified resolutions.
      </p>

      {loading && <p>Loading data from backend...</p>}
      
      {error && (
        <div style={{ background: '#ffe6e6', padding: '10px', marginBottom: '20px' }}>
          <strong>Note:</strong> {error}
        </div>
      )}

      {stats && (
        <>
          <hr />
          <h2>Scraping Stats</h2>
          <ul>
            <li>YouTube Videos Analyzed: <strong>{stats.videos_analyzed}</strong></li>
            <li>Comments Processed: <strong>{stats.comments_analyzed}</strong></li>
            <li>Total Issues Found: <strong>{stats.issues_found}</strong></li>
            <li>Categories Identified: <strong>{stats.categories}</strong></li>
          </ul>
        </>
      )}

      <hr />

      <h2>Issues ({sortedIssues.length} categories)</h2>
      {sortedIssues.length === 0 ? (
        <p>No issues loaded. Make sure the backend API is running.</p>
      ) : (
        <ul>
          {sortedIssues.map((issue) => (
            <li key={issue.id} style={{ marginBottom: '20px' }}>
              <strong>{issue.title}</strong>
              {issue.status === 'resolved' || issue.status === 'verified' ? (
                <span style={{ color: 'green' }}> [RESOLVED]</span>
              ) : (
                <span style={{ color: 'orange' }}> [{issue.status.toUpperCase()}]</span>
              )}
              <span style={{ color: '#666' }}> - {issue.count} complaints</span>
              <br />
              <small>
                Exchange: {issue.exchange} | Category: {issue.category} | Severity: <span style={{ 
                  color: issue.severity === 'high' ? 'red' : issue.severity === 'medium' ? 'orange' : 'gray' 
                }}>{issue.severity}</span>
              </small>
              <br />
              <small>{issue.description}</small>
            </li>
          ))}
        </ul>
      )}

      <hr />

      <h2>Resolutions ({resolutions.length} verified)</h2>
      {resolutions.length === 0 ? (
        <p>No resolutions recorded yet. Resolutions are created when issues are verified as fixed.</p>
      ) : (
        <ul>
          {resolutions.map((resolution) => (
            <li key={resolution.id} style={{ marginBottom: '15px' }}>
              <strong>{resolution.exchange} - {resolution.issue_category}</strong>
              {resolution.status === 'on_chain' ? (
                <span style={{ color: 'green' }}> [ON-CHAIN VERIFIED]</span>
              ) : resolution.status === 'verified' ? (
                <span style={{ color: 'blue' }}> [VERIFIED]</span>
              ) : (
                <span style={{ color: 'gray' }}> [PENDING]</span>
              )}
              <br />
              <small>{resolution.summary}</small>
            </li>
          ))}
        </ul>
      )}

      <hr />

      <h2>Tech Stack</h2>
      <ul>
        <li><strong>Backend:</strong> Go (Golang) - REST API, YouTube/Gemini scrapers, blockchain integration</li>
        <li><strong>Smart Contracts:</strong> Solidity - ResolutionAttestation.sol on Base Sepolia</li>
        <li><strong>Blockchain:</strong> Ethereum/Base (Coinbase L2) - On-chain attestations for verified resolutions</li>
        <li><strong>Data Sources:</strong> YouTube Data API v3, Google Search, Gemini AI with search grounding</li>
      </ul>

      <hr />

      <h2>Blockchain Concepts Used</h2>
      <ul>
        <li>Smart Contracts (Solidity)</li>
        <li>On-Chain Attestations</li>
        <li>Keccak256 Evidence Hashing</li>
        <li>Chain-of-Custody (linked hashes)</li>
        <li>Transaction Signing (ECDSA)</li>
        <li>Events/Logs for indexing</li>
        <li>Layer 2 (Base/Coinbase L2)</li>
        <li>Append-Only Audit Trail</li>
      </ul>
    </div>
  );
}

export default App;
