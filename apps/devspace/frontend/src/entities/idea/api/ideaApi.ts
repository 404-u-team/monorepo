import { apiClient } from '@/shared/api/client';
import type { IIdea } from '../model/IIdea';

interface RawIdea {
    id: string;
    title: string;
    description: string;
    content?: string;
    category?: string;
    author_id: string;
    created_at: string;
    updated_at: string;
    views_count?: number;
    favorites_count?: number;
}

function mapIdea(data: RawIdea): IIdea {
    return {
        id: data.id,
        title: data.title,
        description: data.description,
        content: data.content,
        category: data.category,
        author_id: data.author_id,
        created_at: data.created_at,
        updated_at: data.updated_at,
        views_count: data.views_count ?? 0,
        favorites_count: data.favorites_count ?? 0,
    };
}

export async function fetchIdeaById(ideaId: string): Promise<IIdea> {
    const response = await apiClient.get<RawIdea>(`/ideas/${ideaId}`);
    return mapIdea(response.data);
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
    hasMore: boolean;
}

export async function fetchIdeas(parameters?: FetchIdeasParameters): Promise<PaginatedIdeas> {
    const limit = parameters?.limit ?? 20;
    const startAt = parameters?.start_at ?? 0;
    const response = await apiClient.get<RawIdea[]>('/ideas', { params: parameters });

    const items = response.data.map((item) => mapIdea(item));

    const hasMore = items.length >= limit;
    const currentPage = Math.floor(startAt / limit) + 1;
    const totalPages = hasMore ? currentPage + 1 : currentPage;
    const total = startAt + items.length + (hasMore ? 1 : 0);

    return {
        items,
        total,
        totalPages,
        hasMore,
    };
}

export async function createIdea(data: {
    title: string;
    description?: string;
    content?: string;
    category?: string;
}): Promise<IIdea> {
    const response = await apiClient.post<RawIdea>('/ideas', data);
    return mapIdea(response.data);
}

export async function updateIdea(ideaId: string, data: {
    title?: string;
    description?: string;
    content?: string;
    category?: string;
}): Promise<IIdea> {
    const response = await apiClient.put<RawIdea>(`/ideas/${ideaId}`, data);
    return mapIdea(response.data);
}

export async function deleteIdea(ideaId: string): Promise<void> {
    await apiClient.delete(`/ideas/${ideaId}`);
}

// Note: There's no dedicated endpoint for creating a project from an idea.
// Use createProject from projectApi with the idea_id parameter instead.
