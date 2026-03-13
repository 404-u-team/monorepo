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
