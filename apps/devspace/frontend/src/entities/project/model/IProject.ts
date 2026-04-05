export interface IProject {
  id: string;
  title: string;
  description: string;
  content: string | null;
  leader_id: string;
  status: "open" | "closed";
  idea_id: string | null;
  created_at: string;
  updated_at: string;
}
