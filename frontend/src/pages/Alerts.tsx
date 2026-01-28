import React, { useState } from 'react';
import { 
  Bell, 
  AlertTriangle, 
  CheckCircle,
  Info,
  Clock,
  ExternalLink,
  Trash2,
  Check,
  Filter,
  Plus,
  XCircle
} from 'lucide-react';
import '../styles/Alerts.css';

interface Alert {
  id: number;
  type: 'critical' | 'warning' | 'info' | 'success';
  title: string;
  description: string;
  time: string;
  source: string;
  isRead: boolean;
}

const Alerts: React.FC = () => {
  const [activeTab, setActiveTab] = useState('all');
  
  const [alerts] = useState<Alert[]>([
    {
      id: 1,
      type: 'critical',
      title: 'High-Risk Rug Pull Detected',
      description: 'The project "QuickProfit Token" has been flagged as a potential rug pull. Liquidity has been removed from the pool within the last hour.',
      time: '5 minutes ago',
      source: 'Automated Scanner',
      isRead: false
    },
    {
      id: 2,
      type: 'critical',
      title: 'Phishing Website Identified',
      description: 'A new phishing website mimicking "Uniswap" has been detected at uni-swap-exchange.com. Multiple user complaints received.',
      time: '23 minutes ago',
      source: 'Community Report',
      isRead: false
    },
    {
      id: 3,
      type: 'warning',
      title: 'Suspicious Wallet Activity',
      description: 'Wallet 0x7a2...8f3 linked to previous scam operations has initiated large transfers to a new token contract.',
      time: '1 hour ago',
      source: 'Wallet Monitor',
      isRead: true
    },
    {
      id: 4,
      type: 'info',
      title: 'New Scam Pattern Identified',
      description: 'Our AI has identified a new type of honeypot contract pattern being used across multiple chains. Detection rules updated.',
      time: '3 hours ago',
      source: 'AI Analysis',
      isRead: true
    },
    {
      id: 5,
      type: 'success',
      title: 'Weekly Scan Completed',
      description: 'Automated weekly scan completed successfully. 156 new projects analyzed, 12 flagged for review.',
      time: '6 hours ago',
      source: 'System',
      isRead: true
    },
  ]);

  const tabs = [
    { id: 'all', label: 'All Alerts', count: alerts.length },
    { id: 'unread', label: 'Unread', count: alerts.filter(a => !a.isRead).length },
    { id: 'critical', label: 'Critical', count: alerts.filter(a => a.type === 'critical').length },
  ];

  const getAlertIcon = (type: string) => {
    switch (type) {
      case 'critical':
        return <AlertTriangle size={24} />;
      case 'warning':
        return <AlertTriangle size={24} />;
      case 'info':
        return <Info size={24} />;
      case 'success':
        return <CheckCircle size={24} />;
      default:
        return <Bell size={24} />;
    }
  };

  const filteredAlerts = alerts.filter(alert => {
    if (activeTab === 'all') return true;
    if (activeTab === 'unread') return !alert.isRead;
    if (activeTab === 'critical') return alert.type === 'critical';
    return true;
  });

  return (
    <div className="alerts-page">
      <div className="alerts-header">
        <div className="alerts-header-content">
          <h1>Alerts</h1>
          <p>Real-time notifications about detected threats and system events</p>
        </div>
        <div className="alerts-actions">
          <button className="btn btn-secondary">
            <Filter size={16} />
            Filters
          </button>
          <button className="btn btn-primary">
            <Plus size={16} />
            Create Alert Rule
          </button>
        </div>
      </div>

      <div className="alert-tabs">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            className={`alert-tab ${activeTab === tab.id ? 'active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            {tab.label}
            <span className="tab-count">{tab.count}</span>
          </button>
        ))}
      </div>

      <div className="alerts-list">
        {filteredAlerts.map((alert) => (
          <div 
            key={alert.id} 
            className={`alert-card ${!alert.isRead ? 'unread' : ''} ${alert.type === 'critical' ? 'critical' : ''}`}
          >
            <div className={`alert-icon ${alert.type}`}>
              {getAlertIcon(alert.type)}
            </div>
            
            <div className="alert-content">
              <h3 className="alert-title">{alert.title}</h3>
              <p className="alert-description">{alert.description}</p>
              <div className="alert-meta">
                <span><Clock size={14} /> {alert.time}</span>
                <span><Bell size={14} /> {alert.source}</span>
              </div>
            </div>

            <div className="alert-actions">
              <span className="alert-time">{alert.time}</span>
              <div className="alert-buttons">
                <button className="btn btn-ghost" title="Mark as read">
                  <Check size={16} />
                </button>
                <button className="btn btn-ghost" title="Dismiss">
                  <XCircle size={16} />
                </button>
                <button className="btn btn-secondary">
                  <ExternalLink size={16} />
                  View
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {filteredAlerts.length === 0 && (
        <div className="alerts-empty">
          <Bell size={64} />
          <h3>No alerts to display</h3>
          <p>You're all caught up! New alerts will appear here.</p>
        </div>
      )}
    </div>
  );
};

export default Alerts;
