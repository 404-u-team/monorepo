import { apiClient } from "@/shared/api/client";
import type { IUserResponse } from "../model/IUserResponse";

export interface FetchUsersParameters {
  search?: string | undefined;
  start_at?: number | undefined;
  limit?: number | undefined;
}

export interface PaginatedUsers {
  items: IUserResponse[];
  total: number;
  totalPages: number;
  hasMore: boolean;
}

export async function fetchUsers(parameters?: FetchUsersParameters): Promise<PaginatedUsers> {
  const limit = parameters?.limit ?? 20;
  const startAt = parameters?.start_at ?? 0;

  // Backend expects "username" query parameter, not "search"
  const { search, ...rest } = parameters ?? {};
  const apiParams = { ...rest, username: search };

  const response = await apiClient.get<IUserResponse[]>('/users', { params: apiParams });
  const items = response.data;

  const hasMore = items.length >= limit;
  const currentPage = Math.floor(startAt / limit) + 1;
  const totalPages = hasMore ? currentPage + 1 : currentPage;
  const total = startAt + items.length + (hasMore ? 1 : 0);

  return { items, total, totalPages, hasMore };
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
