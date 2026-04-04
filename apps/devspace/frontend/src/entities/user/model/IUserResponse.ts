import type { ISkill } from '@/entities/skill';

export const ZERO_UUID = '00000000-0000-0000-0000-000000000000';

/** Skill-like object returned inside main_role (SkillWithoutChildren on backend) */
export interface IMainRole {
    id: string;
    name: string;
    parent_id: string | null;
    icon?: string;
    color?: string;
}

export interface IUserResponse {
    id: string;
    nickname: string;
    avatar_url: string;
    bio: string;
    /** Backend returns a SkillCategory object (or null), not a UUID string. */
    main_role: IMainRole | null;
    skills: ISkill[];
}

/** Check whether the main_role field carries a useful value. */
export function isValidMainRole(value: IMainRole | null | undefined): value is IMainRole {
    return value !== null && value !== undefined && typeof value === 'object' && value.id !== ZERO_UUID;
}
