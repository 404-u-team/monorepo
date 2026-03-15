import { apiClient } from '@/shared/api/client';
import type { IProject } from '../model/IProject';
import type { IProjectSlot } from '../model/IProjectSlot';

export interface IProjectDetailResponse extends IProject {
    slots: IProjectSlot[];
}

export async function fetchProjectById(projectId: string): Promise<IProjectDetailResponse> {
    const response = await apiClient.get<IProjectDetailResponse>(`/projects/${projectId}`);
    return response.data;
}

export async function fetchProjectSlots(projectId: string): Promise<IProjectSlot[]> {
    const response = await apiClient.get<IProjectSlot[]>(`/projects/${projectId}/slots`);
    return response.data;
}

export interface FetchProjectsParameters {
    start_at?: number | undefined;
    limit?: number | undefined;
    search?: string | undefined;
    status?: 'open' | 'closed' | undefined;
    leader_id?: string | undefined;
}

export interface PaginatedProjects {
    items: IProject[];
    total: number;
    totalPages: number;
}

export async function fetchProjects(parameters?: FetchProjectsParameters): Promise<PaginatedProjects> {
    const response = await apiClient.get<IProject[]>('/projects', { params: parameters });

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
        // No valid total from the server; fall back to the number of items we actually received.
        total = items.length;
        totalPages = 1;
    }

    return {
        items,
        total,
        totalPages
    };
}
