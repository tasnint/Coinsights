// Complaint/Scam Types
export interface Complaint {
  id: string;
  title: string;
  description: string;
  source: string;
  url: string;
  riskLevel: 'high' | 'medium' | 'low';
  riskScore: number;
  scamType: string;
  detectedAt: string;
  status: 'active' | 'resolved' | 'investigating';
  tags: string[];
  affectedAmount?: string;
}

// Alert Types
export interface Alert {
  id: number;
  type: 'critical' | 'warning' | 'info' | 'success';
  title: string;
  description: string;
  time: string;
  source: string;
  isRead: boolean;
}

// Report Types
export interface Report {
  id: number;
  name: string;
  type: string;
  date: string;
  status: 'completed' | 'pending' | 'generating';
  size: string;
}

// Search Result Types
export interface SearchResult {
  id: number;
  title: string;
  description: string;
  riskLevel: 'high' | 'medium' | 'low';
  riskScore: number;
  source: string;
  date: string;
  views: string;
  tags: string[];
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  error?: string;
}

// Stats Types
export interface DashboardStats {
  scamsDetected: number;
  fundsProtected: string;
  activeMonitors: number;
  detectionRate: number;
}

// Activity Feed Types
export interface Activity {
  id: number;
  type: 'scan' | 'alert' | 'report';
  text: string;
  time: string;
}
