import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import { fetchIdeas } from "@/entities/idea";
import { getPageSize } from "@/shared/lib/pageSize";
import { IdeaList } from "@/widgets/IdeaList";

export interface IdeaSearch {
  page?: number | undefined;
  search?: string | undefined;
  is_favorite?: boolean | undefined;
  views?: "asc" | "desc" | undefined;
  favorites?: "asc" | "desc" | undefined;
  limit?: number | undefined;
}

export const Route = createFileRoute("/ideas")({
  validateSearch: (search: Record<string, unknown>): IdeaSearch => {
    const views =
      search.views === "asc" || search.views === "desc" ? search.views : undefined;
    const favorites =
      search.favorites === "asc" || search.favorites === "desc" ? search.favorites : undefined;
    return {
      page: Number(search.page) || 1,
      search: search.search as string | undefined,
      is_favorite: search.is_favorite === true || search.is_favorite === "true" ? true : undefined,
      views,
      favorites,
      limit: Number(search.limit) || undefined,
    };
  },
  loaderDeps: ({ search: { page, search, is_favorite, views, favorites, limit } }) => ({
    page,
    search,
    is_favorite,
    views,
    favorites,
    limit,
  }),
  loader: async ({ deps }) => {
    const limit = deps.limit ?? getPageSize();
    const currentPage = deps.page ?? 1;
    const start_at = (currentPage - 1) * limit;

    return fetchIdeas({
      search: deps.search,
      start_at,
      limit,
      is_favorite: deps.is_favorite,
      views: deps.views,
      favorites: deps.favorites,
    });
  },
  component: IdeasPage,
});

function IdeasPage(): JSX.Element {
  const data = Route.useLoaderData();

  return (
    <div style={{ padding: "24px", maxWidth: "1200px", margin: "0 auto" }}>
      <IdeaList ideas={data.items} totalPages={data.totalPages} total={data.total} />
    </div>
  );
}
