import React, { useState, useEffect } from 'react';
import '../styles/Dashboard.css';

interface Issue {
  id: string;
  title: string;
  source: string;
  date: string;
  description: string;
}

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState({
    videosAnalyzed: 0,
    queriesMade: 0,
    googleSearches: 0,
    uniqueIssues: 0
  });

  const [issues, setIssues] = useState<Issue[]>([]);

  useEffect(() => {
    // TODO: Fetch real data from backend
    // For now, using placeholder data
    setStats({
      videosAnalyzed: 0,
      queriesMade: 0,
      googleSearches: 0,
      uniqueIssues: 0
    });

    setIssues([]);
  }, []);

  return (
    <div className="dashboard">
      <div className="dashboard-content">
        <div className="stats-section">
          <div className="stats-grid">
            <div className="stat-card">
              <h3 className="stat-card-title">Videos Analyzed</h3>
              <p className="stat-card-value">{stats.videosAnalyzed}</p>
            </div>

            <div className="stat-card">
              <h3 className="stat-card-title">Queries Made</h3>
              <p className="stat-card-value">{stats.queriesMade}</p>
            </div>

            <div className="stat-card">
              <h3 className="stat-card-title">Google searches</h3>
              <p className="stat-card-value">{stats.googleSearches}</p>
            </div>

            <div className="stat-card">
              <h3 className="stat-card-title">Unique Issues Identified</h3>
              <p className="stat-card-value">{stats.uniqueIssues}</p>
            </div>
          </div>
        </div>

        <div className="issues-section">
          <h2 className="issues-title">Issues:</h2>
          <div className="issues-list">
            {issues.length === 0 ? (
              <p className="no-issues">No issues detected yet</p>
            ) : (
              issues.map((issue) => (
                <div key={issue.id} className="issue-item">
                  <h4 className="issue-item-title">{issue.title}</h4>
                  <p className="issue-item-source">{issue.source}</p>
                  <p className="issue-item-description">{issue.description}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
