import { apiClient } from "@/shared/api/client";

import type { ISkill } from "../model/ISkill";

export interface FetchSkillsParameters {
  parent_id?: string | undefined;
  search?: string | undefined;
  start_at?: number | undefined;
  limit?: number | undefined;
}

export async function fetchSkills(parameters?: FetchSkillsParameters): Promise<ISkill[]> {
  const response = await apiClient.get<ISkill[]>("/skills", { params: parameters });
  return response.data;
}

export async function fetchSkillById(skillId: string): Promise<ISkill> {
  const response = await apiClient.get<ISkill>(`/skills/${skillId}`);
  return response.data;
}

export async function createSkill(data: {
  name: string;
  parent_id?: string | undefined;
}): Promise<ISkill> {
  const response = await apiClient.post<ISkill>("/skills", data);
  return response.data;
}

export async function deleteSkill(skillId: string, cascade = false): Promise<void> {
  await apiClient.delete(`/skills/${skillId}`, { params: cascade ? { cascade: true } : undefined });
}
