import React, { useState } from 'react';
import { 
  Search as SearchIcon, 
  AlertTriangle, 
  Globe,
  Youtube,
  Twitter,
  Clock,
  ExternalLink,
  Bookmark,
  Calendar,
  Eye,
  Tag
} from 'lucide-react';
import '../styles/Search.css';

const Search: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [activeFilter, setActiveFilter] = useState('all');

  const filters = [
    { id: 'all', label: 'All Sources', icon: Globe },
    { id: 'web', label: 'Web', icon: Globe },
    { id: 'youtube', label: 'YouTube', icon: Youtube },
    { id: 'twitter', label: 'Twitter', icon: Twitter },
  ];

  const searchResults = [
    {
      id: 1,
      title: 'SafeMoon Token Analysis',
      description: 'Multiple red flags detected including locked liquidity issues, suspicious wallet activity, and promotional patterns consistent with pump-and-dump schemes.',
      riskLevel: 'high',
      riskScore: 87,
      source: 'Web Analysis',
      date: 'Jan 28, 2026',
      views: '12.4K',
      tags: ['Rug Pull', 'DeFi', 'BSC']
    },
    {
      id: 2,
      title: 'CryptoGains YouTube Channel',
      description: 'Channel promoting high-risk investments with unrealistic return promises. Multiple flagged videos featuring undisclosed partnerships.',
      riskLevel: 'medium',
      riskScore: 65,
      source: 'YouTube',
      date: 'Jan 27, 2026',
      views: '8.2K',
      tags: ['Influencer', 'Promotion', 'High Returns']
    },
    {
      id: 3,
      title: 'MetaVault Protocol Review',
      description: 'Smart contract analysis reveals potential backdoor functions. Team anonymity and lack of audit reports raise concerns.',
      riskLevel: 'high',
      riskScore: 92,
      source: 'Contract Analysis',
      date: 'Jan 26, 2026',
      views: '5.7K',
      tags: ['Smart Contract', 'Honeypot', 'Ethereum']
    },
  ];

  const recentSearches = [
    'Bitcoin ETF scam',
    'SafeMoon token',
    'Crypto Ponzi scheme',
    'NFT rug pull',
    'DeFi exploit'
  ];

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    console.log('Searching for:', searchQuery);
  };

  return (
    <div className="search-page">
      <div className="search-hero">
        <h1>Search & Analyze</h1>
        <p>
          Investigate cryptocurrencies, projects, and influencers for potential 
          scams and red flags across the web.
        </p>

        <form className="search-box" onSubmit={handleSearch}>
          <div className="search-input-container">
            <SearchIcon size={22} />
            <input
              type="text"
              placeholder="Enter a cryptocurrency name, wallet address, or project URL..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
            <button type="submit" className="btn btn-primary search-btn">
              Search
            </button>
          </div>
        </form>

        <div className="search-filters">
          {filters.map((filter) => (
            <button
              key={filter.id}
              className={`filter-chip ${activeFilter === filter.id ? 'active' : ''}`}
              onClick={() => setActiveFilter(filter.id)}
            >
              <filter.icon size={14} />
              {filter.label}
            </button>
          ))}
        </div>
      </div>

      <div className="search-results">
        <div className="results-header">
          <h2>Search Results</h2>
          <span className="results-count">Showing 3 results</span>
        </div>

        {searchResults.map((result) => (
          <div className="result-card" key={result.id}>
            <div className="result-header">
              <div className="result-title">
                <AlertTriangle 
                  size={24} 
                  color={result.riskLevel === 'high' ? '#ef4444' : '#f59e0b'} 
                />
                <h3>{result.title}</h3>
              </div>
              <div className="result-risk">
                <div className={`risk-score ${result.riskLevel}`}>
                  <AlertTriangle size={14} />
                  Risk: {result.riskScore}/100
                </div>
              </div>
            </div>

            <div className="result-body">
              <p className="result-description">{result.description}</p>
              <div className="result-tags">
                {result.tags.map((tag, index) => (
                  <span className="result-tag" key={index}>
                    <Tag size={10} /> {tag}
                  </span>
                ))}
              </div>
            </div>

            <div className="result-footer">
              <div className="result-meta">
                <span><Globe size={14} /> {result.source}</span>
                <span><Calendar size={14} /> {result.date}</span>
                <span><Eye size={14} /> {result.views} views</span>
              </div>
              <div className="result-actions">
                <button className="btn btn-ghost" title="Save">
                  <Bookmark size={16} />
                </button>
                <button className="btn btn-secondary">
                  <ExternalLink size={16} />
                  View Details
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="recent-searches">
        <h3>Recent Searches</h3>
        <div className="search-history">
          {recentSearches.map((search, index) => (
            <button 
              key={index} 
              className="search-history-item"
              onClick={() => setSearchQuery(search)}
            >
              <Clock size={14} />
              {search}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
};

export default Search;
