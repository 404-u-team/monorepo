import { apiClient } from '@/shared/api/client';
import type { IIdea } from '../model/IIdea';

interface RawIdea {
    id?: string;
    ID?: string;
    title?: string;
    Title?: string;
    description?: string;
    Description?: string;
    content?: string;
    Content?: string;
    category?: string;
    Category?: string;
    author_id?: string;
    AuthorID?: string;
    created_at?: string;
    CreatedAt?: string;
    updated_at?: string;
    UpdatedAt?: string;
    views_count?: number;
    ViewsCount?: number;
    favorites_count?: number;
    FavoritesCount?: number;
}

function mapIdea(data: RawIdea): IIdea {
    return {
        id: data.id ?? data.ID ?? '',
        title: data.title ?? data.Title ?? '',
        description: data.description ?? data.Description ?? '',
        content: data.content ?? data.Content,
        category: data.category ?? data.Category,
        author_id: data.author_id ?? data.AuthorID ?? '',
        created_at: data.created_at ?? data.CreatedAt ?? '',
        updated_at: data.updated_at ?? data.UpdatedAt ?? '',
        views_count: data.views_count ?? data.ViewsCount ?? 0,
        favorites_count: data.favorites_count ?? data.FavoritesCount ?? 0,
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
}

export async function fetchIdeas(parameters?: FetchIdeasParameters): Promise<PaginatedIdeas> {
    const response = await apiClient.get<RawIdea[]>('/ideas', { params: parameters });

    const items = response.data.map((item) => mapIdea(item));
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
        total = startAt + items.length;
        totalPages = Math.max(1, Math.ceil(total / limit));
    }

    return {
        items,
        total,
        totalPages,
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

export async function createProjectFromIdea(ideaId: string): Promise<void> {
    await apiClient.post(`/ideas/${ideaId}/project`);
}
