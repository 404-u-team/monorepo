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

interface RawUser {
    id?: string;
    ID?: string;
    nickname?: string;
    Nickname?: string;
    avatar_uri?: string;
    AvatarUrl?: string;
    AvatarUri?: string;
    bio?: string;
    Bio?: string;
    main_role?: string;
    MainRole?: string;
    skills?: string[];
    Skills?: string[];
}

function mapUser(data: RawUser): IUserResponse {
    return {
        id: data.id ?? data.ID ?? '',
        nickname: data.nickname ?? data.Nickname ?? '',
        avatar_uri: data.avatar_uri ?? data.AvatarUrl ?? data.AvatarUri ?? '',
        bio: data.bio ?? data.Bio ?? '',
        main_role: data.main_role ?? data.MainRole ?? '',
        skills: data.skills ?? data.Skills ?? [],
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
