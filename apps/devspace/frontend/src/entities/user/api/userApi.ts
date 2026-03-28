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
}

export async function fetchUsers(parameters?: FetchUsersParameters): Promise<PaginatedUsers> {
  const response = await apiClient.get<IUserResponse[]>('/users', { params: parameters });
  const items = response.data;
  const totalCountHeader = response.headers['x-total-count'] as string | undefined;
  const parsedTotal = typeof totalCountHeader === 'string' ? Number(totalCountHeader) : Number.NaN;

  let total: number;
  let totalPages: number;

  if (Number.isFinite(parsedTotal)) {
    const limit = parameters?.limit ?? (items.length > 0 ? items.length : 1);
    total = parsedTotal;
    totalPages = Math.ceil(total / limit);
  } else {
    total = items.length;
    totalPages = 1;
  }

  return { items, total, totalPages };
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
