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

// ============================================
// BLOCKCHAIN / ATTESTATION TYPES
// ============================================

// Issue represents a tracked issue
export interface Issue {
  id: string;
  exchange: string;
  category: string;
  title: string;
  description: string;
  first_detected: string;
  last_updated: string;
  complaint_count: number;
  severity: 'critical' | 'high' | 'medium' | 'low';
  status: 'active' | 'investigating' | 'resolved' | 'verified';
  resolution?: Resolution;
  attestation?: Attestation;
}

// Resolution represents a resolved issue with evidence
export interface Resolution {
  id: string;
  exchange: string;
  issue_category: string;
  summary: string;
  evidence: ResolutionEvidence;
  confidence: number;
  resolution_window: number;
  status: 'pending' | 'verified' | 'on_chain';
  created_at: string;
  verified_at?: string;
  attestation?: Attestation;
}

// ResolutionEvidence contains the data that gets hashed for on-chain attestation
export interface ResolutionEvidence {
  complaints_before: number;
  complaints_after: number;
  percentage_decrease: number;
  sentiment_shift: number;
  sample_complaints: string[];
  data_sources: string[];
  measurement_start: string;
  measurement_end: string;
  analysis_methodology: string;
}

// Attestation represents an on-chain verification record
export interface Attestation {
  id: number;
  transaction_hash: string;
  block_number: number;
  block_timestamp: string;
  chain_id: number;
  contract_address: string;
  evidence_hash: string;
  previous_hash?: string;
  attestor: string;
  explorer_url: string;
  verified: boolean;
}

// Chain configuration
export interface ChainConfig {
  name: string;
  chain_id: number;
  rpc_url: string;
  explorer_url: string;
  contract_address: string;
  is_testnet: boolean;
}

// Verification response
export interface VerificationResponse {
  verified: boolean;
  on_chain: boolean;
  attestation?: Attestation;
  hash_match: boolean;
  timestamp_valid: boolean;
  message: string;
}

// Blockchain stats
export interface BlockchainStats {
  total_issues: number;
  total_resolutions: number;
  issues_by_status: Record<string, number>;
  attestation_count: number;
  on_chain_attestation_count?: number;
}

