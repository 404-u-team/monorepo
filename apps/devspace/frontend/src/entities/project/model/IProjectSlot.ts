export interface IProjectSlotSkill {
    id: string;
    name: string;
    color?: string | undefined;
    icon?: string | undefined;
}

export interface IProjectSlot {
    id: string;
    project_id: string;
    skill_category_id: string;
    skill?: IProjectSlotSkill | undefined;
    /** Secondary (second-level) skill IDs */
    secondary_skills_ids?: string[] | undefined;
    /** Resolved secondary skill objects (may be populated by server) */
    secondary_skills?: IProjectSlotSkill[] | undefined;
    title: string;
    description: string;
    status: 'open' | 'closed';
    user_id: string | null;
    created_at: string;
}
