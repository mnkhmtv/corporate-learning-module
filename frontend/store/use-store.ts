import { create } from 'zustand';
import { User, TrainingRequest, LearningProcess, Mentor } from '@/lib/types';
import { mockService } from '@/lib/mock-service';

interface AppState {
  user: User | null;
  isLoading: boolean;
  requests: TrainingRequest[];
  learnings: LearningProcess[];
  mentors: Mentor[];
  
  // Actions
  login: (email: string) => Promise<boolean>;
  logout: () => void;
  register: (data: any) => Promise<void>;
  
  fetchUserData: () => Promise<void>;
  fetchAllRequests: () => Promise<void>; // Admin only
  createRequest: (data: { topic: string; description: string }) => Promise<void>;
  fetchMentors: () => Promise<void>;
  assignMentor: (requestId: string, mentorId: string) => Promise<void>;
  
  fetchLearning: (id: string) => Promise<LearningProcess | undefined>;
  updateLearningPlan: (id: string, plan: any[]) => Promise<void>;
  completeLearning: (id: string, feedback: { rating: number, comment: string }) => Promise<void>;
}

export const useStore = create<AppState>((set, get) => ({
  user: null,
  isLoading: false,
  requests: [],
  learnings: [],
  mentors: [],

  login: async (email: string) => {
    set({ isLoading: true });
    try {
      const user = await mockService.login(email);
      if (user) {
        set({ user });
        return true;
      }
      return false;
    } finally {
      set({ isLoading: false });
    }
  },

  logout: () => {
    set({ user: null, requests: [], learnings: [] });
  },

  register: async (data) => {
    set({ isLoading: true });
    try {
      const user = await mockService.register(data);
      set({ user });
    } finally {
      set({ isLoading: false });
    }
  },

  fetchUserData: async () => {
    const { user } = get();
    if (!user) return;
    
    set({ isLoading: true });
    try {
      const [requests, learnings] = await Promise.all([
        mockService.getUserRequests(user.id),
        mockService.getUserLearnings(user.id)
      ]);
      set({ requests, learnings });
    } finally {
      set({ isLoading: false });
    }
  },

  fetchAllRequests: async () => {
    set({ isLoading: true });
    try {
      const requests = await mockService.getRequests();
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
      await mockService.createRequest(user.id, data);
      await get().fetchUserData();
    } finally {
      set({ isLoading: false });
    }
  },

  fetchMentors: async () => {
    const mentors = await mockService.getMentors();
    set({ mentors });
  },

  assignMentor: async (requestId, mentorId) => {
    set({ isLoading: true });
    try {
      await mockService.assignMentor(requestId, mentorId);
      await get().fetchAllRequests(); // Refresh list
    } finally {
      set({ isLoading: false });
    }
  },

  fetchLearning: async (id) => {
    set({ isLoading: true });
    try {
      return await mockService.getLearning(id);
    } finally {
      set({ isLoading: false });
    }
  },

  updateLearningPlan: async (id, plan) => {
    await mockService.updateLearningPlan(id, plan);
    const learning = await mockService.getLearning(id);
    // update local state if needed
  },

  completeLearning: async (id, feedback) => {
    set({ isLoading: true });
    try {
      await mockService.completeLearning(id, feedback);
      await get().fetchUserData();
    } finally {
      set({ isLoading: false });
    }
  }
}));

