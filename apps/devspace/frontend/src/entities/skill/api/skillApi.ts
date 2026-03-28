import { apiClient } from '@/shared/api/client';
import type { ISkill } from '../model/ISkill';

export interface FetchSkillsParameters {
    parent_id?: string | undefined;
    search?: string | undefined;
    page?: number | undefined;
    limit?: number | undefined;
}

export async function fetchSkills(parameters?: FetchSkillsParameters): Promise<ISkill[]> {
    const response = await apiClient.get<ISkill[]>('/skills', { params: parameters });
    return response.data;
}

export async function fetchSkillById(skillId: string): Promise<ISkill> {
    const response = await apiClient.get<ISkill>(`/skills/${skillId}`);
    return response.data;
}
