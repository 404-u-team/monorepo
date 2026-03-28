export interface ISkill {
    id: string;
    name: string;
    parent_id: string | null;
    children: ISkill[];
    color?: string | undefined;
}
