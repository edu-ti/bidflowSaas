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
