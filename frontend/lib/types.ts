export type Role = 'employee' | 'admin';
export type RequestStatus = 'pending' | 'approved' | 'rejected';
export type LearningStatus = 'active' | 'completed';

export interface User {
  id: string;
  name: string;
  role: Role;
  email: string;
  department?: string;
  jobTitle?: string;
  telegram?: string;
}

export interface TrainingRequest {
  id: string;
  userId: string;
  topic: string;
  description: string;
  status: RequestStatus;
  createdAt: string;
  updatedAt: string;
}

export interface LearningProcess {
  id: string;
  requestId: string;
  userId: string;
  mentorId: string;
  mentorName: string;
  mentorEmail: string;
  mentorTg?: string;
  topic: string;
  status: LearningStatus;
  startDate: string;
  endDate?: string;
  plan: LearningPlanItem[];
  feedback?: LearningFeedback;
  notes?: string;
}

export interface LearningPlanItem {
  id: string;
  text: string;
  completed: boolean;
}

export interface LearningFeedback {
  rating: number;
  comment: string;
}

export interface Mentor {
  id: string;
  name: string;
  jobTitle: string;
  experience: string;
  workload: number; // 0-5
  email: string;
  telegram?: string;
  avatar?: string;
}

