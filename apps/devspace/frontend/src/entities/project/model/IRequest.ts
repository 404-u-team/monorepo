export interface IRequest {
  id: string;
  type: "apply" | "invite";
  status: "pending" | "accepted" | "rejected";
  project_id: string;
  slot_id: string;
  user_id: string;
  leader_id: string;
  cover_letter: string;
  created_at: string;
}
