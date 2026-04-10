import { useNavigate, useSearch, Link } from "@tanstack/react-router";
import { Plus, Star } from "lucide-react";
import { observer } from "mobx-react-lite";
import { type JSX } from "react";

import { IdeaCard } from "@/entities/idea";
import { useStore } from "@/shared/lib/store";
import { setPageSize, PAGE_SIZE_OPTIONS, type PageSize } from "@/shared/lib/pageSize";
import { DataListLayout, Button, Dropdown } from "@/shared/ui";

type SortOrder = "asc" | "desc";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
  is_favorite?: boolean | undefined;
  views?: SortOrder | undefined;
  favorites?: SortOrder | undefined;
  limit?: number | undefined;
}

export interface IdeaListProps {
  ideas: { id: string }[];
  totalPages: number;
  total: number;
}

const pageSizeOptions = PAGE_SIZE_OPTIONS.map((n) => ({
  label: `${String(n)} / стр.`,
  value: String(n),
}));

const sortViewsOptions = [
  { label: "Просмотры: —", value: "" },
  { label: "Просмотры: ↑", value: "asc" },
  { label: "Просмотры: ↓", value: "desc" },
];

const sortFavoritesOptions = [
  { label: "Избранное: —", value: "" },
  { label: "Избранное: ↑", value: "asc" },
  { label: "Избранное: ↓", value: "desc" },
];

export const IdeaList = observer(function IdeaList({
  ideas,
  totalPages,
  total,
}: IdeaListProps): JSX.Element {
  const { userStore } = useStore();
  const searchParameters: SearchParameters = useSearch({ strict: false });
  const navigate = useNavigate({ from: "/ideas" });

  const handleSearch = (value: string): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        search: value || undefined,
        page: 1,
      }),
    });
  };

  const handlePageChange = (page: number): void => {
    window.scrollTo({ top: 0, behavior: "smooth" });
    void navigate({
      search: (previous: SearchParameters) => ({ ...previous, page }),
    });
  };

  const handleToggleFavorites = (): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        is_favorite: previous.is_favorite === true ? undefined : true,
        page: 1,
      }),
    });
  };

  const handleSortViews = (value: string): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        views: value === "" ? undefined : (value as SortOrder),
        page: 1,
      }),
    });
  };

  const handleSortFavorites = (value: string): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        favorites: value === "" ? undefined : (value as SortOrder),
        page: 1,
      }),
    });
  };

  const handlePageSizeChange = (value: string): void => {
    const size = Number(value) as PageSize;
    setPageSize(size);
    void navigate({
      search: (previous: SearchParameters) => ({ ...previous, limit: size, page: 1 }),
    });
  };

  const isFavoritesActive = searchParameters.is_favorite === true;
  const currentLimit = searchParameters.limit ?? 20;

  const filtersNode = (
    <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
      <Dropdown
        options={sortViewsOptions}
        value={searchParameters.views ?? ""}
        onChange={handleSortViews}
      />
      <Dropdown
        options={sortFavoritesOptions}
        value={searchParameters.favorites ?? ""}
        onChange={handleSortFavorites}
      />
      {userStore.isAuthenticated && (
        <Button
          variant={isFavoritesActive ? "primary" : "outline"}
          onClick={handleToggleFavorites}
          style={{ width: "100%" }}
        >
          <Star size={16} fill={isFavoritesActive ? "currentColor" : "none"} />
          Только избранные
        </Button>
      )}
    </div>
  );

  const createButton = userStore.isAuthenticated ? (
    <Link to="/idea/new">
      <Button>
        <Plus size={18} />
        Создать идею
      </Button>
    </Link>
  ) : undefined;

  const controlsNode = (
    <div style={{ display: "flex", gap: "0.75rem", alignItems: "center" }}>
      <Dropdown
        options={pageSizeOptions}
        value={String(currentLimit)}
        onChange={handlePageSizeChange}
      />
      {createButton}
    </div>
  );

  return (
    <DataListLayout
      title="Идеи"
      subtitle="Найдите вдохновение или присоединяйтесь к реализации новой задумки"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={controlsNode}
      filtersNode={filtersNode}
      isEmpty={ideas.length === 0}
      emptyMessage={isFavoritesActive ? "В избранном пока нет идей" : "Идеи не найдены"}
      currentPage={Number((searchParameters as Record<string, string>).page) || 1}
      totalPages={totalPages}
      total={total}
      onPageChange={handlePageChange}
    >
      {ideas.map((idea) => (
        <IdeaCard key={idea.id} ideaId={idea.id} fromRoute="/ideas" />
      ))}
    </DataListLayout>
  );
});
