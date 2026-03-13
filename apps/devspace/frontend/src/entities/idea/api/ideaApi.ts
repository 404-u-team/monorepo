import { apiClient } from '@/shared/api/client';
import type { IIdea } from '../model/IIdea';

export async function fetchIdeaById(ideaId: string): Promise<IIdea> {
    const response = await apiClient.get<IIdea>(`/ideas/${ideaId}`);
    return response.data;
}

export async function toggleIdeaFavorite(ideaId: string): Promise<{ is_favorite: boolean }> {
    const response = await apiClient.post<{ is_favorite: boolean }>(`/ideas/${ideaId}/favorite`);
    return response.data;
}
