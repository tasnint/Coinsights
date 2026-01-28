import React from 'react';
import { NavLink, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Search, 
  Bell, 
  FileText, 
  Shield,
  TrendingUp,
  Settings,
  HelpCircle
} from 'lucide-react';
import '../styles/Sidebar.css';

const Sidebar: React.FC = () => {
  const location = useLocation();

  const mainNavItems = [
    { path: '/', icon: LayoutDashboard, label: 'Dashboard' },
    { path: '/search', icon: Search, label: 'Search & Analyze' },
    { path: '/alerts', icon: Bell, label: 'Alerts', badge: 3 },
    { path: '/reports', icon: FileText, label: 'Reports' },
  ];

  const secondaryNavItems = [
    { path: '/trends', icon: TrendingUp, label: 'Market Trends' },
    { path: '/protection', icon: Shield, label: 'Scam Protection' },
  ];

  const utilityNavItems = [
    { path: '/settings', icon: Settings, label: 'Settings' },
    { path: '/help', icon: HelpCircle, label: 'Help & Support' },
  ];

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <div className="sidebar-logo">
          <img src="/logo.png" alt="Coinsights Logo" />
          <h1>Coinsights</h1>
        </div>
      </div>

      <nav className="sidebar-nav">
        <div className="nav-section">
          <div className="nav-section-title">Main Menu</div>
          {mainNavItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            >
              <item.icon />
              <span>{item.label}</span>
              {item.badge && <span className="nav-badge">{item.badge}</span>}
            </NavLink>
          ))}
        </div>

        <div className="nav-section">
          <div className="nav-section-title">Insights</div>
          {secondaryNavItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            >
              <item.icon />
              <span>{item.label}</span>
            </NavLink>
          ))}
        </div>

        <div className="nav-section">
          <div className="nav-section-title">System</div>
          {utilityNavItems.map((item) => (
            <NavLink
              key={item.path}
              to={item.path}
              className={({ isActive }) => `nav-item ${isActive ? 'active' : ''}`}
            >
              <item.icon />
              <span>{item.label}</span>
            </NavLink>
          ))}
        </div>
      </nav>

      <div className="sidebar-footer">
        <div className="sidebar-user">
          <div className="user-avatar">A</div>
          <div className="user-info">
            <h4>Admin User</h4>
            <p>Analyst</p>
          </div>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;
