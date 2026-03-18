import { apiClient } from "@/shared/api/client";
import type { IUserResponse } from "../model/IUserResponse";

export interface InviteToSlotParameters {
  project_id: string;
  slot_id: string;
  user_id: string;
}

export interface IInviteResponse {
  ID: string;
  SlotID: string;
  UserID: string;
  ProjectID: string;
  LeaderID: string;
  CoverLetter: string;
  CreatedAt: string;
  Status: string;
  Type: string;
}

const userCache = new Map<string, Promise<IUserResponse>>();

interface RawUser {
    id: string;
    nickname: string;
    avatar_uri: string;
    bio?: string;
    main_role?: string;
    skills?: ISkill[];
}

function mapUser(data: RawUser): IUserResponse {
    return {
        id: data.id,
        nickname: data.nickname,
        avatar_uri: data.avatar_uri,
        bio: data.bio ?? '',
        main_role: data.main_role ?? '',
        skills: data.skills ?? [],
    };
}

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
