export interface IIdea {
  id: string;
  title: string;
  description: string;
  content?: string | undefined;
  category?: string | undefined;
  author_id: string;
  is_author?: boolean | undefined;
  is_favorite?: boolean | undefined;
  created_at: string;
  updated_at: string;
  views_count: number;
  favorites_count: number;
}
