export type UUID = string;

export interface User {
  id: UUID;
  email: string;
  nickname: string;
  avatar_uri: string;
  bio: string;
  created_at: string;
}

export interface PrivateUserProfile extends User {
  skills: Skill[];
}

export interface PublicUserProfile {
  id: UUID;
  nickname: string;
  avatar_uri: string;
  bio: string;
  skills: Skill[];
}

export interface Skill {
  id: UUID;
  name: string;
  parent_id: UUID | null;
  children?: Skill[];
}

export interface Idea {
  id: UUID;
  title: string;
  description: string;
  author_id: UUID;
  created_at: string;
  updated_at: string;
}

export interface IdeaDetails extends Idea {
  views_count: number;
  favorites_count: number;
}

export interface Project {
  id: UUID;
  title: string;
  description: string;
  leader_id: UUID;
  status: "open" | "closed";
  idea_id: UUID | null;
  created_at: string;
  updated_at: string;
}

export interface ProjectDetails extends Project {
  slots: ProjectSlot[];
}

export interface ProjectSlot {
  id: UUID;
  project_id: UUID;
  skill_category_id: UUID;
  title: string;
  description: string;
  status: "open" | "closed";
  user_id: UUID | null;
  created_at: string;
}

export type RequestType = "apply" | "invite";
export type RequestStatus = "pending" | "accepted" | "rejected";

export interface Request {
  id: UUID;
  type: RequestType;
  status: RequestStatus;
  project_id: UUID;
  slot_id: UUID;
  user_id: UUID;
  leader_id: UUID;
  cover_letter?: string;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
}

export interface PaginationParameters {
  page?: number;
  limit?: number;
}

export interface ApiError {
  statusCode: number;
  message: string;
  error?: string;
}
