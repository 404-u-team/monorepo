import { apiClient } from "@/shared/api/client";

import type { IUserResponse } from "../model/IUserResponse";

export interface FetchUsersParameters {
  search?: string | undefined;
  start_at?: number | undefined;
  limit?: number | undefined;
  main_role?: string | undefined;
  skills?: string[] | undefined;
}

export interface PaginatedUsers {
  items: IUserResponse[];
  total: number;
  totalPages: number;
  hasMore: boolean;
}

interface GetUsersResponse {
  total: number;
  profiles: IUserResponse[];
}

export async function fetchUsers(parameters?: FetchUsersParameters): Promise<PaginatedUsers> {
  const limit = parameters?.limit ?? 20;
  const startAt = parameters?.start_at ?? 0;

  const { skills, ...rest } = parameters ?? {};
  const apiParameters: Record<string, unknown> = { ...rest };
  if (skills && skills.length > 0) {
    apiParameters.skills = skills;
  }

  const response = await apiClient.get<GetUsersResponse>("/users", {
    params: apiParameters,
    paramsSerializer: { indexes: false },
  });
  const { total, profiles } = response.data;

  const totalPages = Math.ceil(total / limit) || 1;
  const hasMore = startAt + profiles.length < total;

  return { items: profiles, total, totalPages, hasMore };
}

export interface InviteToSlotParameters {
  project_id: string;
  slot_id: string;
  user_id: string;
}

export interface IInviteResponse {
  id: string;
  slot_id: string;
  user_id: string;
  project_id: string;
  leader_id: string;
  cover_letter: string;
  created_at: string;
  status: string;
  type: string;
}

const userCache = new Map<string, Promise<IUserResponse>>();

export function fetchUserById(userId: string): Promise<IUserResponse> {
  const cached = userCache.get(userId);
  if (cached !== undefined) {
    return cached;
  }
  const request = apiClient
    .get<IUserResponse>(`/users/${userId}`)
    .then((response) => response.data)
    .catch((error: unknown) => {
      userCache.delete(userId);
      throw error;
    });

  userCache.set(userId, request);
  return request;
}

export function inviteUserToSlot({
  project_id,
  slot_id,
  user_id,
}: InviteToSlotParameters): Promise<IInviteResponse> {
  return apiClient
    .post<IInviteResponse>(`/projects/${project_id}/slots/${slot_id}/invite`, {
      user_id,
    })
    .then((response) => response.data);
}
