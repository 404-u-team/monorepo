import type { ISkill } from '@/entities/skill';

export const ZERO_UUID = '00000000-0000-0000-0000-000000000000';

export interface IUserResponse {
    id: string;
    nickname: string;
    avatar_uri: string;
    bio: string;
    main_role: string | null;
    skills: ISkill[];
}

export function isValidMainRole(value: string | null | undefined): value is string {
    return value !== null && value !== undefined && value !== '' && value !== ZERO_UUID;
}
