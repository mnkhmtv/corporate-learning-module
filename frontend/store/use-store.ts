import { create } from 'zustand';
import { User, TrainingRequest, LearningProcess, Mentor, LearningPlanItem } from '@/lib/types';
import { apiService } from '@/lib/api-service';

interface AppState {
  user: User | null;
  isLoading: boolean;
  isAuthLoading: boolean;
  requests: TrainingRequest[];
  learnings: LearningProcess[];
  mentors: Mentor[];
  
  // Actions
  login: (data: { email: string; password: string }) => Promise<boolean>;
  logout: () => void;
  register: (data: any) => Promise<void>;
  checkAuth: () => Promise<void>;
  
  fetchUserData: () => Promise<void>;
  fetchAllRequests: () => Promise<void>; // Admin only
  createRequest: (data: { topic: string; description: string }) => Promise<void>;
  fetchMentors: () => Promise<void>;
  assignMentor: (requestId: string, mentorId: string) => Promise<void>;
  
  fetchLearning: (id: string) => Promise<LearningProcess | undefined>;
  updateLearningPlan: (id: string, plan: LearningPlanItem[]) => Promise<void>;
  updateNotes: (id: string, notes: string) => Promise<void>;
  completeLearning: (id: string, feedback: { rating: number, comment: string }) => Promise<void>;
}

export const useStore = create<AppState>((set, get) => ({
  user: null,
  isLoading: false,
  isAuthLoading: true,
  requests: [],
  learnings: [],
  mentors: [],

  login: async (data) => {
    set({ isLoading: true });
    try {
      const { token, user } = await apiService.login(data);
      localStorage.setItem('token', token);
      set({ user });
      // Fetch initial data after login
      if (user.role === 'admin') {
        // For now, component logic handles what to fetch
      }
      return true;
    } catch (error) {
      console.error('Login failed:', error);
      return false;
    } finally {
      set({ isLoading: false });
    }
  },

  logout: () => {
    localStorage.removeItem('token');
    set({ user: null, requests: [], learnings: [] });
  },

  register: async (data) => {
    set({ isLoading: true });
    try {
      await apiService.register(data);
      // After register, usually we want to login automatically or redirect to login
      // The UI handles redirection probably.
    } finally {
      set({ isLoading: false });
    }
  },

  checkAuth: async () => {
    const token = localStorage.getItem('token');
    if (!token) {
      set({ isAuthLoading: false });
      return;
    }

    // No need to set isLoading here, we have a dedicated flag
    try {
      const user = await apiService.getMe();
      set({ user });
    } catch (error) {
      console.error('Auth check failed:', error);
      localStorage.removeItem('token');
      set({ user: null });
    } finally {
      set({ isAuthLoading: false });
    }
  },

  fetchUserData: async () => {
    const { user } = get();
    if (!user) return;
    
    set({ isLoading: true });
    try {
      const [requests, learnings] = await Promise.all([
        apiService.getMyRequests(),
        apiService.getMyLearnings()
      ]);
      set({ requests, learnings });
    } finally {
      set({ isLoading: false });
    }
  },

  fetchAllRequests: async () => {
    set({ isLoading: true });
    try {
      const requests = await apiService.getAllRequests();
      set({ requests });
    } finally {
      set({ isLoading: false });
    }
  },

  createRequest: async (data) => {
    const { user } = get();
    if (!user) return;
    set({ isLoading: true });
    try {
      await apiService.createRequest(data);
      await get().fetchUserData();
    } finally {
      set({ isLoading: false });
    }
  },

  fetchMentors: async () => {
    const mentors = await apiService.getMentors();
    set({ mentors });
  },

  assignMentor: async (requestId, mentorId) => {
    set({ isLoading: true });
    try {
      await apiService.assignMentor(requestId, mentorId);
      await get().fetchAllRequests(); // Refresh list
    } finally {
      set({ isLoading: false });
    }
  },

  fetchLearning: async (id) => {
    set({ isLoading: true });
    try {
      return await apiService.getLearning(id);
    } finally {
      set({ isLoading: false });
    }
  },

  updateLearningPlan: async (id, plan) => {
    await apiService.updatePlan(id, plan);
    // TODO: Update local state optimistically if needed.
  },

  updateNotes: async (id, notes) => {
    await apiService.updateNotes(id, notes);
    // Note: We don't need to update local state here as the component
    // that calls this will manage its own state.
  },

  completeLearning: async (id, feedback) => {
    set({ isLoading: true });
    try {
      await apiService.completeLearning(id, feedback);
      await get().fetchUserData();
    } finally {
      set({ isLoading: false });
    }
  }
}));
