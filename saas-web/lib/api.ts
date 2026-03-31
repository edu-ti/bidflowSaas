// lib/api.ts
import axios from 'axios';
import jwtDecode from 'jwt-decode';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
});

// Request interceptor to add JWT token
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

// Response interceptor to handle 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      if (typeof window !== 'undefined') {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
    }

    // MOCK FOR NETWORK ERROR TO KEEP PROTOTYPE FUNCTIONAL
    if (error.message === 'Network Error' || error.code === 'ERR_NETWORK') {
      console.warn('Network Error detected, returning mock data for:', error.config.url);
      
      const url = error.config.url || '';
      let mockData: any = {};
      
      if (url.includes('/radar')) {
        mockData = {
          editais: [
            { id: 'mock1', title: 'Edital Mockado (Backend OFF)', value: 150000, deadline: '2025-12-31T00:00:00Z', priority_label: 'HIGH', win_probability: 0.85, similarity: 0.92, why_relevant: 'Mock criado pelo interceptor porque o backend não está online.' },
            { id: 'mock2', title: 'Fornecimento de Software Georreferenciado', value: 80000, deadline: '2025-10-15T00:00:00Z', priority_label: 'MEDIUM', win_probability: 0.55, similarity: 0.68, why_relevant: 'Aderência moderada. Requer ajustes de portfólio.' }
          ]
        };
      } else if (url.includes('/ai/insights')) {
        mockData = { average_score: 0.82, success_rate: 0.65, average_confidence: 0.94, average_win_probability: 0.70 };
      } else if (url.includes('/monitor/events')) {
        mockData = [
          { id: '1', type: 'STATUS_CHANGE', severity: 'HIGH', message: 'Licitação aberta para propostas.', timestamp: new Date().toISOString(), metadata: { suggested_action: 'Revisar requisitos' } }
        ];
      } else if (url.includes('/notifications')) {
        mockData = [
          { id: 'mock1', title: 'Bem-vindo ao BidFlow', message: 'O sistema está rodando localmente sem backend (modo simulação).', read: false, timestamp: new Date().toISOString() }
        ];
      } else if (url.includes('/ai/history')) {
        mockData = [
          { id: 'mock1', edital_id: '12351', ai_score: 0.89, smart_score: 0.88, confidence: 0.95 }
        ];
      } else if (url.includes('/billing/checkout')) {
         mockData = { checkout_url: 'https://checkout.stripe.com/mock-bidflow-url' };
      }
      
      return Promise.resolve({ data: mockData, status: 200, statusText: 'OK', headers: {}, config: error.config });
    }

    return Promise.reject(error);
  }
);

export const getRadar = async (params?: Record<string, any>) => {
  const res = await api.get('/api/v1/radar', { params });
  return res.data;
};

export const getRadarDetail = async (id: string) => {
  const res = await api.get(`/api/v1/radar/${id}`);
  return res.data;
};

export const getAIInsights = async (id?: string) => {
  const res = await api.get('/api/v1/ai/insights', { params: { id } });
  return res.data;
};

export const getMonitorEvents = async () => {
  const res = await api.get('/api/v1/monitor/events');
  return res.data;
};

export const getNotifications = async () => {
  const res = await api.get('/api/v1/notifications');
  return res.data;
};

export const postFeedback = async (payload: any) => {
  const res = await api.post('/api/v1/ai/feedback', payload);
  return res.data;
};

export const postCheckout = async (plan: string) => {
  const res = await api.post('/api/v1/billing/checkout', { plan });
  return res.data;
};

export const followRadar = async (id: string) => {
  const res = await api.post(`/api/v1/radar/follow/${id}`);
  return res.data;
};

export const readNotification = async (id: string) => {
  const res = await api.post(`/api/v1/notifications/read/${id}`);
  return res.data;
};

export const getAIHistory = async () => {
  const res = await api.get('/api/v1/ai/history');
  return res.data;
};

export default api;
