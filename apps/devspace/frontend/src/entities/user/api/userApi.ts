import { apiClient } from "@/shared/api/client";
import type { IUserResponse } from "../model/IUserResponse";

export interface InviteToSlotParameters {
  project_id: string;
  slot_id: string;
  id: string;
}

export interface IInviteResponse {
  cover_letter: string;
  created_at: string;
  id: string;
  leader_id: string;
  project_id: string;
  slot_id: string;
  status: string;
  type: string;
  user_id: string;
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
  id,
}: InviteToSlotParameters): Promise<IInviteResponse> {
  return apiClient
    .post<IInviteResponse>(
      `/api/projects/${project_id}/slots/${slot_id}/invite`,
      {
        id: id,
      },
    )
    .then((response) => response.data);
}
