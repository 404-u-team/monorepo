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
    page?: number | undefined;
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
    
    // Attempt to read standard pagination headers, fallback to mocked total if absent
    const totalCountHeader = String(response.headers['x-total-count'] ?? '');
    const totalCount = totalCountHeader !== '' ? Number(totalCountHeader) : 100;
    const limit = parameters?.limit ?? 20;

    return {
        items: response.data,
        total: totalCount,
        totalPages: Math.ceil(totalCount / limit)
    };
}
