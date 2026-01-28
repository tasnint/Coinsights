import React from 'react';
import { 
  FileText, 
  BarChart3, 
  PieChart, 
  TrendingUp,
  Download,
  Plus,
  Calendar,
  CheckCircle,
  Clock,
  Loader,
  Eye,
  Trash2,
  Share2
} from 'lucide-react';
import '../styles/Reports.css';

const Reports: React.FC = () => {
  const reportTypes = [
    {
      icon: BarChart3,
      title: 'Scam Summary Report',
      description: 'Comprehensive overview of all detected scams, trends, and statistics.',
      color: 'blue'
    },
    {
      icon: PieChart,
      title: 'Risk Analysis Report',
      description: 'Detailed breakdown of risk categories and their distribution.',
      color: 'purple'
    },
    {
      icon: TrendingUp,
      title: 'Trend Report',
      description: 'Analysis of emerging scam patterns and market trends.',
      color: 'green'
    },
    {
      icon: FileText,
      title: 'Custom Report',
      description: 'Create a custom report with your selected parameters.',
      color: 'orange'
    },
  ];

  const recentReports = [
    {
      id: 1,
      name: 'Weekly Scam Summary - Jan 21-27',
      type: 'Scam Summary',
      date: 'Jan 27, 2026',
      status: 'completed',
      size: '2.4 MB'
    },
    {
      id: 2,
      name: 'Monthly Risk Analysis - January',
      type: 'Risk Analysis',
      date: 'Jan 26, 2026',
      status: 'completed',
      size: '5.1 MB'
    },
    {
      id: 3,
      name: 'DeFi Scam Trends Q4 2025',
      type: 'Trend Report',
      date: 'Jan 25, 2026',
      status: 'completed',
      size: '8.7 MB'
    },
    {
      id: 4,
      name: 'Custom Report - Ethereum Analysis',
      type: 'Custom Report',
      date: 'Jan 28, 2026',
      status: 'generating',
      size: '-'
    },
    {
      id: 5,
      name: 'YouTube Influencer Audit',
      type: 'Custom Report',
      date: 'Jan 24, 2026',
      status: 'completed',
      size: '3.2 MB'
    },
  ];

  const stats = [
    { value: '47', label: 'Reports Generated' },
    { value: '12', label: 'This Month' },
    { value: '156', label: 'Scams Documented' },
    { value: '98%', label: 'Accuracy Rate' },
  ];

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle size={14} />;
      case 'pending':
        return <Clock size={14} />;
      case 'generating':
        return <Loader size={14} className="spinning" />;
      default:
        return null;
    }
  };

  return (
    <div className="reports-page">
      <div className="reports-header">
        <div className="reports-header-content">
          <h1>Reports</h1>
          <p>Generate and manage scam detection reports and analytics</p>
        </div>
        <div className="reports-actions">
          <button className="btn btn-secondary">
            <Calendar size={16} />
            Schedule Report
          </button>
          <button className="btn btn-primary">
            <Plus size={16} />
            New Report
          </button>
        </div>
      </div>

      <div className="report-stats">
        {stats.map((stat, index) => (
          <div className="report-stat" key={index}>
            <div className="report-stat-value">{stat.value}</div>
            <div className="report-stat-label">{stat.label}</div>
          </div>
        ))}
      </div>

      <div className="report-types">
        {reportTypes.map((type, index) => (
          <div className="report-type-card" key={index}>
            <div className={`report-type-icon ${type.color}`}>
              <type.icon size={28} />
            </div>
            <h3>{type.title}</h3>
            <p>{type.description}</p>
          </div>
        ))}
      </div>

      <div className="reports-table-card">
        <div className="reports-table-header">
          <h2>Recent Reports</h2>
          <button className="btn btn-ghost">View All</button>
        </div>
        <table className="reports-table">
          <thead>
            <tr>
              <th>Report Name</th>
              <th>Type</th>
              <th>Generated</th>
              <th>Status</th>
              <th>Size</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {recentReports.map((report) => (
              <tr key={report.id}>
                <td>
                  <div className="report-name">
                    <FileText size={18} />
                    <span>{report.name}</span>
                  </div>
                </td>
                <td>{report.type}</td>
                <td>{report.date}</td>
                <td>
                  <span className={`report-status ${report.status}`}>
                    {getStatusIcon(report.status)}
                    {report.status.charAt(0).toUpperCase() + report.status.slice(1)}
                  </span>
                </td>
                <td>{report.size}</td>
                <td>
                  <div className="report-actions">
                    <button className="report-action-btn" title="View">
                      <Eye size={14} />
                    </button>
                    <button className="report-action-btn" title="Download">
                      <Download size={14} />
                    </button>
                    <button className="report-action-btn" title="Share">
                      <Share2 size={14} />
                    </button>
                    <button className="report-action-btn" title="Delete">
                      <Trash2 size={14} />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default Reports;
