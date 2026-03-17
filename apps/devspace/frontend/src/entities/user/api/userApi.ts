import { apiClient } from '@/shared/api/client';
import type { IUserResponse } from '../model/IUserResponse';

const userCache = new Map<string, Promise<IUserResponse>>();

function mapUser(data: any): IUserResponse {
    return {
        id: data.id || data.ID,
        nickname: data.nickname || data.Nickname,
        avatar_uri: data.avatar_uri || data.AvatarUrl || data.AvatarUri || '',
        bio: data.bio || data.Bio || '',
        main_role: data.main_role || data.MainRole || '',
        skills: data.skills || data.Skills || [],
    };
}

export function fetchUserById(userId: string): Promise<IUserResponse> {
    const cached = userCache.get(userId);
    if (cached !== undefined) {
        return cached;
    }

    const request = apiClient
        .get<any>(`/users/${userId}`)
        .then((response) => mapUser(response.data))
        .catch((error: unknown) => {
            userCache.delete(userId);
            throw error;
        });

    userCache.set(userId, request);
    return request;
}
