import { User, TrainingRequest, LearningProcess, Mentor } from './types';

const DELAY = 800;

// Mock Data
const MOCK_USERS: User[] = [
  {
    id: '1',
    name: 'Иван Иванов',
    role: 'employee',
    email: 'user@skillbridge.com',
    department: 'IT',
    jobTitle: 'Frontend Developer',
    telegram: '@ivan_dev'
  },
  {
    id: '2',
    name: 'Анна Петрова',
    role: 'admin',
    email: 'admin@skillbridge.com',
    department: 'HR',
    jobTitle: 'Learning Manager',
    telegram: '@anna_hr'
  }
];

const MOCK_MENTORS: Mentor[] = [
  {
    id: 'm1',
    name: 'Сергей Сергеев',
    jobTitle: 'Senior Backend Dev',
    experience: '8 лет',
    workload: 3,
    email: 'sergey@skillbridge.com',
    telegram: '@sergey_backend'
  },
  {
    id: 'm2',
    name: 'Елена Сидорова',
    jobTitle: 'Team Lead Design',
    experience: '6 лет',
    workload: 5,
    email: 'elena@skillbridge.com',
    telegram: '@elena_design'
  },
  {
    id: 'm3',
    name: 'Дмитрий Козлов',
    jobTitle: 'Product Manager',
    experience: '5 лет',
    workload: 1,
    email: 'dmitry@skillbridge.com',
    telegram: '@dmitry_pm'
  }
];

let requests: TrainingRequest[] = [
  {
    id: 'r1',
    userId: '1',
    user: MOCK_USERS[0],
    topic: 'Advanced React Patterns',
    description: 'Хочу углубить знания в оптимизации рендеринга и сложных хуках.',
    status: 'pending',
    createdAt: new Date(Date.now() - 86400000 * 2).toISOString(), // 2 days ago
    updatedAt: new Date(Date.now() - 86400000 * 2).toISOString()
  }
];

let learnings: LearningProcess[] = [];

// Mock Service
export const mockService = {
  login: async (email: string): Promise<User | null> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return MOCK_USERS.find(u => u.email === email) || null;
  },

  register: async (data: Omit<User, 'id' | 'role'>): Promise<User> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    const newUser: User = {
      ...data,
      id: Math.random().toString(36).substr(2, 9),
      role: 'employee'
    };
    MOCK_USERS.push(newUser);
    return newUser;
  },

  getRequests: async (): Promise<TrainingRequest[]> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return [...requests];
  },

  getUserRequests: async (userId: string): Promise<TrainingRequest[]> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return requests.filter(r => r.userId === userId);
  },

  createRequest: async (userId: string, data: Pick<TrainingRequest, 'topic' | 'description'>): Promise<TrainingRequest> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    const user = MOCK_USERS.find(u => u.id === userId);
    if (!user) throw new Error('User not found');

    const newRequest: TrainingRequest = {
      id: Math.random().toString(36).substr(2, 9),
      userId,
      user,
      ...data,
      status: 'pending',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };
    requests.push(newRequest);
    return newRequest;
  },

  getMentors: async (): Promise<Mentor[]> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return [...MOCK_MENTORS];
  },

  assignMentor: async (requestId: string, mentorId: string): Promise<LearningProcess> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    const request = requests.find(r => r.id === requestId);
    if (!request) throw new Error('Request not found');
    
    const mentor = MOCK_MENTORS.find(m => m.id === mentorId);
    if (!mentor) throw new Error('Mentor not found');

    // Update request status
    request.status = 'approved';
    request.updatedAt = new Date().toISOString();

    // Create learning process
    const learning: LearningProcess = {
      id: Math.random().toString(36).substr(2, 9),
      userId: request.userId,
      mentor,
      request,
      status: 'active',
      startDate: new Date().toISOString(),
      plan: [],
    };
    learnings.push(learning);
    return learning;
  },

  getUserLearnings: async (userId: string): Promise<LearningProcess[]> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return learnings.filter(l => l.userId === userId);
  },

  getLearning: async (id: string): Promise<LearningProcess | undefined> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    return learnings.find(l => l.id === id);
  },

  updateLearningPlan: async (id: string, plan: any[]): Promise<void> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    const learning = learnings.find(l => l.id === id);
    if (learning) {
      learning.plan = plan;
    }
  },

  completeLearning: async (id: string, feedback: { rating: number, comment: string }): Promise<void> => {
    await new Promise(resolve => setTimeout(resolve, DELAY));
    const learning = learnings.find(l => l.id === id);
    if (learning) {
      learning.status = 'completed';
      learning.endDate = new Date().toISOString();
      learning.feedback = feedback;
    }
  }
};

