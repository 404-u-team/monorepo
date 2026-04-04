import { apiClient } from "@/shared/api/client";

import type { IProject } from "../model/IProject";
import type { IProjectSlot } from "../model/IProjectSlot";
import type { IRequest } from "../model/IRequest";

export interface IProjectDetailResponse extends IProject {
  slots?: IProjectSlot[];
}

export async function fetchProjectById(projectId: string): Promise<IProjectDetailResponse> {
  const response = await apiClient.get<IProjectDetailResponse>(`/projects/${projectId}`);
  return response.data;
}

export async function fetchProjectSlots(projectId: string): Promise<IProjectSlot[]> {
  const response = await apiClient.get<IProjectSlot[]>(`/projects/${projectId}/slots`);
  return response.data;
}

export interface FetchProjectsParameters {
  start_at?: number | undefined;
  limit?: number | undefined;
  search?: string | undefined;
  status?: "open" | "closed" | undefined;
  leader_id?: string | undefined;
}

export interface PaginatedProjects {
  items: IProject[];
  total: number;
  totalPages: number;
  hasMore: boolean;
}

export async function createProject(data: {
  title: string;
  description?: string;
  content?: string;
  idea_id?: string;
}): Promise<IProject> {
  const response = await apiClient.post<IProject>("/projects", data);
  return response.data;
}

export async function updateProject(
  projectId: string,
  data: {
    title: string;
    description: string;
    content?: string | undefined;
    status: "open" | "closed";
  },
): Promise<IProject> {
  const response = await apiClient.put<IProject>(`/projects/${projectId}`, data);
  return response.data;
}

export async function deleteProject(projectId: string): Promise<void> {
  await apiClient.delete(`/projects/${projectId}`);
}

export async function createProjectSlot(
  projectId: string,
  data: {
    primary_skills_id: string[];
    title: string;
    description?: string;
    status?: "open" | "closed";
    secondary_skills_id?: string[];
  },
): Promise<IProjectSlot> {
  const response = await apiClient.post<IProjectSlot>(`/projects/${projectId}/slots`, data);
  return response.data;
}

export async function deleteProjectSlot(projectId: string, slotId: string): Promise<void> {
  await apiClient.delete(`/projects/${projectId}/slots/${slotId}`);
}

export async function updateProjectSlot(
  projectId: string,
  slotId: string,
  data: {
    primary_skills_id?: string[];
    secondary_skills_id?: string[];
    title?: string;
    description?: string;
    status?: "open" | "closed";
  },
): Promise<IProjectSlot> {
  const response = await apiClient.put<IProjectSlot>(
    `/projects/${projectId}/slots/${slotId}`,
    data,
  );
  return response.data;
}

export async function applyToSlot(
  projectId: string,
  slotId: string,
  data?: {
    cover_letter?: string;
  },
): Promise<IRequest> {
  const response = await apiClient.post<IRequest>(
    `/projects/${projectId}/slots/${slotId}/apply`,
    data ?? {},
  );
  return response.data;
}

export async function getProjectRequests(
  projectId: string,
  parameters?: {
    slot_id?: string;
    status?: "pending" | "accepted" | "rejected";
  },
): Promise<IRequest[]> {
  const response = await apiClient.get<IRequest[]>(`/projects/${projectId}/requests`, {
    params: parameters,
  });
  return response.data;
}

export async function acceptRequest(requestId: string): Promise<IRequest> {
  const response = await apiClient.put<IRequest>(`/requests/${requestId}/accept`);
  return response.data;
}

export async function rejectRequest(requestId: string): Promise<IRequest> {
  const response = await apiClient.put<IRequest>(`/requests/${requestId}/reject`);
  return response.data;
}

export async function getMyRequests(parameters?: {
  type?: "apply" | "invite";
  status?: "pending" | "accepted" | "rejected";
}): Promise<IRequest[]> {
  const response = await apiClient.get<IRequest[]>("/users/me/requests", { params: parameters });
  return response.data;
}

export async function fetchProjects(
  parameters?: FetchProjectsParameters,
): Promise<PaginatedProjects> {
  const limit = parameters?.limit ?? 20;
  const startAt = parameters?.start_at ?? 0;
  const response = await apiClient.get<IProject[]>("/projects", { params: parameters });

  const items = response.data;

  // Backend doesn't provide x-total-count header.
  // Detect "hasMore" heuristically: if we received exactly `limit` items, there's probably a next page.
  const hasMore = items.length >= limit;
  const currentPage = Math.floor(startAt / limit) + 1;
  // totalPages is at least currentPage; if there are more items, assume at least one more page.
  const totalPages = hasMore ? currentPage + 1 : currentPage;
  const total = startAt + items.length + (hasMore ? 1 : 0);

  return {
    items,
    total,
    totalPages,
    hasMore,
  };
}
