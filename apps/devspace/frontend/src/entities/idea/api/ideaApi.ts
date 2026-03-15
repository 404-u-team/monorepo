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

    const limit = parameters?.limit ?? (items.length > 0 ? items.length : 1);
    const startAt = parameters?.start_at ?? 0;

    let total: number;
    let totalPages: number;

    if (Number.isFinite(parsedTotal)) {
        total = parsedTotal;
        totalPages = Math.ceil(total / limit);
    } else {
        // Fallback when the backend does not provide X-Total-Count:
        // use at least the number of items we've seen so far.
        total = startAt + items.length;
        totalPages = Math.max(1, Math.ceil(total / limit));
    }

    return {
        items,
        total,
        totalPages,
    };
}
