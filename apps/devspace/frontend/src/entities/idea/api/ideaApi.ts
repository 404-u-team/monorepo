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

export interface FetchIdeasParameters {
    start_at?: number | undefined;
    limit?: number | undefined;
    search?: string | undefined;
    author_id?: string | undefined;
}

export interface PaginatedIdeas {
    items: IIdea[];
    total: number;
    totalPages: number;
}

export async function fetchIdeas(parameters?: FetchIdeasParameters): Promise<PaginatedIdeas> {
    const response = await apiClient.get<IIdea[]>('/ideas', { params: parameters });

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

    return {
        items,
        total,
        totalPages,
    };
}
