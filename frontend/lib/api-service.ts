import axios from 'axios';
import { User, TrainingRequest, LearningProcess, Mentor, LearningPlanItem } from './types';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to attach the JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Add a response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized error (e.g., clear token and redirect to login)
      localStorage.removeItem('token');
      // TODO: Handle redirect to login
      // if (typeof window !== 'undefined' && !window.location.pathname.includes('/login')) {
      // }
    }
    return Promise.reject(error);
  }
);

export interface LoginResponse {
  token: string;
  user: User;
}

export const apiService = {
  // Auth
  register: async (data: any): Promise<User> => {
    const response = await api.post<User>('/auth/register', data);
    return response.data;
  },

  login: async (data: { email: string; password: string }): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me');
    return response.data;
  },

  // Requests
  createRequest: async (data: { topic: string; description: string }): Promise<TrainingRequest> => {
    const response = await api.post<TrainingRequest>('/requests', data);
    return response.data;
  },

  getMyRequests: async (): Promise<TrainingRequest[]> => {
    const response = await api.get<{ requests: TrainingRequest[] }>('/requests/my');
    return response.data.requests || [];
  },

  getAllRequests: async (): Promise<TrainingRequest[]> => {
    const response = await api.get<{ requests: TrainingRequest[] }>('/requests');
    return response.data.requests || [];
  },

  assignMentor: async (requestId: string, mentorId: string): Promise<LearningProcess> => {
    const response = await api.post<LearningProcess>(`/requests/${requestId}/assign`, { mentorId });
    return response.data;
  },

  // Learning Processes
  getMyLearnings: async (): Promise<LearningProcess[]> => {
    const response = await api.get<{ learnings: LearningProcess[] }>('/learnings');
    return response.data.learnings || [];
  },

  getLearning: async (id: string): Promise<LearningProcess> => {
    const response = await api.get<LearningProcess>(`/learnings/${id}`);
    return response.data;
  },

  getLearningProgress: async (id: string): Promise<any> => {
    const response = await api.get(`/learnings/${id}/progress`);
    return response.data;
  },

  addPlanItem: async (id: string, text: string): Promise<void> => {
    await api.post(`/learnings/${id}/plan`, { text });
  },

  updatePlan: async (id: string, plan: LearningPlanItem[]): Promise<void> => {
    await api.put(`/learnings/${id}/plan`, { plan });
  },

  updatePlanItem: async (id: string, itemId: string, data: { text?: string; completed?: boolean }): Promise<void> => {
    await api.put(`/learnings/${id}/plan/${itemId}`, data);
  },

  togglePlanItem: async (id: string, itemId: string): Promise<void> => {
    await api.patch(`/learnings/${id}/plan/${itemId}/toggle`);
  },

  removePlanItem: async (id: string, itemId: string): Promise<void> => {
    await api.delete(`/learnings/${id}/plan/${itemId}`);
  },

  updateNotes: async (id: string, notes: string): Promise<void> => {
    await api.put(`/learnings/${id}/notes`, { notes });
  },

  completeLearning: async (id: string, feedback: { rating: number; comment: string }): Promise<void> => {
    await api.post(`/learnings/${id}/complete`, feedback);
  },

  // Mentors
  getMentors: async (): Promise<Mentor[]> => {
    const response = await api.get<Mentor[]>('/mentors');
    return response.data;
  },

  getMentor: async (id: string): Promise<Mentor> => {
    const response = await api.get<Mentor>(`/mentors/${id}`);
    return response.data;
  },

  createMentor: async (data: any): Promise<Mentor> => {
    const response = await api.post<Mentor>('/mentors', data);
    return response.data;
  },
};
