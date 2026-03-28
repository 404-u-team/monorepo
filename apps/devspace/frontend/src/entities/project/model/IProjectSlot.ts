export interface IProjectSlotSkill {
    id: string;
    name: string;
    color?: string | undefined;
    icon?: string | undefined;
}

export interface IProjectSlot {
    id: string;
    project_id: string;
    /** Primary (first-level) skills — can be multiple */
    primary_skills: IProjectSlotSkill[];
    /** Secondary (second-level) skill objects */
    secondary_skills?: IProjectSlotSkill[] | undefined;
    title: string;
    description: string;
    status: 'open' | 'closed';
    user_id: string | null;
    created_at: string;
}
