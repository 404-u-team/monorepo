import { type JSX } from "react";
import { useNavigate, useSearch, Link } from "@tanstack/react-router";
import { observer } from "mobx-react-lite";
import { Plus, Star } from "lucide-react";
import { IdeaCard, type IIdea } from "@/entities/idea";
import { DataListLayout, Button } from "@/shared/ui";
import { useStore } from "@/shared/lib/store";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
  favorites?: boolean | undefined;
}

export interface IdeaListProps {
  ideas: IIdea[];
  totalPages: number;
}

export const IdeaList = observer(function IdeaList({ ideas, totalPages }: IdeaListProps): JSX.Element {
  const { userStore } = useStore();
  const searchParameters: SearchParameters = useSearch({
    strict: false,
  });
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
        favorites: previous.favorites ? undefined : true,
        page: 1,
      }),
    });
  };

  const isFavoritesActive = searchParameters.favorites === true;

  const favoritesButton = userStore.isAuthenticated ? (
    <Button
      variant={isFavoritesActive ? "primary" : "outline"}
      onClick={handleToggleFavorites}
    >
      <Star size={16} fill={isFavoritesActive ? "currentColor" : "none"} />
      Избранное
    </Button>
  ) : undefined;

  const createButton = userStore.isAuthenticated ? (
    <Link to="/idea/new">
      <Button>
        <Plus size={18} />
        Создать идею
      </Button>
    </Link>
  ) : undefined;

  const controls = (
    <>
      {favoritesButton}
      {createButton}
    </>
  );

  return (
    <DataListLayout
      title="Идеи"
      subtitle="Найдите вдохновение или присоединяйтесь к реализации новой задумки"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={userStore.isAuthenticated ? controls : createButton}
      isEmpty={ideas.length === 0}
      emptyMessage={isFavoritesActive ? "В избранном пока нет идей" : "Идеи не найдены"}
      currentPage={
        Number((searchParameters as Record<string, string>).page) || 1
      }
      totalPages={totalPages}
      onPageChange={handlePageChange}
    >
      {ideas.map((idea) => (
        <IdeaCard key={idea.id} ideaId={idea.id} />
      ))}
    </DataListLayout>
  );
});
