export interface IProject {
  id: string;
  title: string;
  description: string;
  content?: string | null | undefined;
  leader_id: string;
  is_leader?: boolean | undefined;
  is_favorite?: boolean | undefined;
  status: "open" | "closed";
  idea_id: string | null;
  created_at: string;
  updated_at: string;
}
