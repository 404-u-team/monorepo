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
    title: string;
    description: string;
    status: 'open' | 'closed';
    user_id: string | null;
    created_at: string;
}
