export interface IProject {
    id: string;
    title: string;
    description: string;
    leader_id: string;
    status: 'open' | 'closed';
    idea_id: string | null;
    created_at: string;
    updated_at: string;
}
