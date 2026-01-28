import axios from 'axios';

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

export default api;
