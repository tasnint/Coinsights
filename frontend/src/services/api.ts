import axios from 'axios';
import { Issue, Resolution, Attestation, VerificationResponse, BlockchainStats, ChainConfig, ResolutionEvidence } from '../types';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Complaints/Scams API
export const complaintsApi = {
  getAll: () => api.get('/complaints'),
  getById: (id: string) => api.get(`/complaints/${id}`),
  search: (query: string) => api.get(`/complaints/search?q=${encodeURIComponent(query)}`),
  getRecent: (limit: number = 10) => api.get(`/complaints/recent?limit=${limit}`),
};

// Search API
export const searchApi = {
  search: (query: string, filters?: { source?: string }) => 
    api.get('/search', { params: { q: query, ...filters } }),
  getHistory: () => api.get('/search/history'),
};

// Alerts API
export const alertsApi = {
  getAll: () => api.get('/alerts'),
  markAsRead: (id: number) => api.patch(`/alerts/${id}/read`),
  dismiss: (id: number) => api.delete(`/alerts/${id}`),
  getUnreadCount: () => api.get('/alerts/unread/count'),
};

// Reports API
export const reportsApi = {
  getAll: () => api.get('/reports'),
  getById: (id: number) => api.get(`/reports/${id}`),
  generate: (type: string, params?: object) => api.post('/reports/generate', { type, ...params }),
  download: (id: number) => api.get(`/reports/${id}/download`, { responseType: 'blob' }),
};

// Stats API
export const statsApi = {
  getDashboard: () => api.get('/stats/dashboard'),
  getScamTrends: () => api.get('/stats/trends'),
};

// Scraper API (trigger manual scrapes)
export const scraperApi = {
  triggerGoogleSearch: (query: string) => api.post('/scraper/google', { query }),
  triggerYouTubeSearch: (query: string) => api.post('/scraper/youtube', { query }),
  triggerAiAnalysis: (query: string) => api.post('/scraper/ai', { query }),
  getLatestResults: () => api.get('/scraper/results/latest'),
};

// ============================================
// BLOCKCHAIN / ATTESTATION API
// ============================================

// Issues API
export const issuesApi = {
  getAll: (status?: string) => 
    api.get<{ issues: Issue[]; count: number }>('/issues', { params: { status } }),
  getById: (id: string) => 
    api.get<Issue>(`/issues/${id}`),
  create: (issue: Partial<Issue>) => 
    api.post<Issue>('/issues', issue),
  update: (id: string, update: Partial<Issue>) => 
    api.patch<Issue>(`/issues/${id}`, update),
};

// Resolutions API
export const resolutionsApi = {
  getAll: (status?: string) => 
    api.get<{ resolutions: Resolution[]; count: number }>('/resolutions', { params: { status } }),
  getById: (id: string) => 
    api.get<Resolution>(`/resolutions/${id}`),
  create: (issueId: string, summary: string, evidence: ResolutionEvidence) => 
    api.post<Resolution>('/resolutions', { issue_id: issueId, summary, evidence }),
  getAttestation: (id: string) => 
    api.get<Attestation>(`/resolutions/${id}/attestation`),
};

// Attestations API
export const attestationsApi = {
  create: (resolutionId: string) => 
    api.post<{ success: boolean; attestation: Attestation }>('/attestations', { resolution_id: resolutionId }),
  verify: (params: { evidence_hash?: string; resolution_id?: string }) => 
    api.post<VerificationResponse>('/attestations/verify', params),
};

// Blockchain Info API
export const blockchainApi = {
  getInfo: () => 
    api.get<{ chain: ChainConfig; wallet_address: string; supported_chains: Record<string, ChainConfig> }>('/blockchain/info'),
  getStats: () => 
    api.get<BlockchainStats>('/blockchain/stats'),
  hashEvidence: (evidence: ResolutionEvidence) => 
    api.post<{ hash: string }>('/blockchain/hash', evidence),
  // Demo endpoint for testing the full workflow
  runDemoWorkflow: () => 
    api.post<{
      workflow_complete: boolean;
      issue: Issue;
      resolution: Resolution;
      evidence_hash: string;
      next_steps: Record<string, string>;
    }>('/demo/full-workflow'),
};

export default api;
