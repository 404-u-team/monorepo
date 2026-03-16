import { apiClient } from '@/shared/api/client';
import type { IUserResponse } from '../model/IUserResponse';

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
