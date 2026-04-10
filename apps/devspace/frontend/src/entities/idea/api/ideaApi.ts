import { apiClient } from "@/shared/api/client";

import type { IIdea } from "../model/IIdea";

interface RawIdea {
  id: string;
  title: string;
  description: string;
  content?: string;
  category?: string;
  author_id: string;
  is_author?: boolean;
  is_favorite?: boolean;
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
    is_author: data.is_author,
    is_favorite: data.is_favorite,
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
  is_favorite?: boolean | undefined;
  views?: "asc" | "desc" | undefined;
  favorites?: "asc" | "desc" | undefined;
}

export interface PaginatedIdeas {
  items: IIdea[];
  total: number;
  totalPages: number;
  hasMore: boolean;
}

interface GetIdeasResponse {
  total: number;
  ideas: RawIdea[];
}

export async function fetchIdeas(parameters?: FetchIdeasParameters): Promise<PaginatedIdeas> {
  const limit = parameters?.limit ?? 20;
  const startAt = parameters?.start_at ?? 0;
  const response = await apiClient.get<GetIdeasResponse>("/ideas", { params: parameters });

  const { total, ideas } = response.data;
  const items = ideas.map((item) => mapIdea(item));

  const totalPages = Math.ceil(total / limit) || 1;
  const hasMore = startAt + items.length < total;

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
  const response = await apiClient.post<RawIdea>("/ideas", data);
  return mapIdea(response.data);
}

export async function updateIdea(
  ideaId: string,
  data: {
    title?: string;
    description?: string;
    content?: string;
    category?: string;
  },
): Promise<IIdea> {
  const response = await apiClient.put<RawIdea>(`/ideas/${ideaId}`, data);
  return mapIdea(response.data);
}

export async function deleteIdea(ideaId: string): Promise<void> {
  await apiClient.delete(`/ideas/${ideaId}`);
}

// Note: There's no dedicated endpoint for creating a project from an idea.
// Use createProject from projectApi with the idea_id parameter instead.
