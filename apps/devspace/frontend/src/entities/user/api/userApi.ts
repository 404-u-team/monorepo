import { apiClient } from "@/shared/api/client";
import type { IUserResponse, ISkill } from "../model/IUserResponse";

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
        .get<RawUser>(`/users/${userId}`)
        .then((response) => mapUser(response.data))
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
