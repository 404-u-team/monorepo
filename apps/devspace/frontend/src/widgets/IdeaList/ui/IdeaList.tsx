import { type JSX } from "react";
import { useNavigate, useSearch, Link } from "@tanstack/react-router";
import { observer } from "mobx-react-lite";
import { Plus } from "lucide-react";
import { IdeaCard, type IIdea } from "@/entities/idea";
import { DataListLayout, Button } from "@/shared/ui";
import { useStore } from "@/shared/lib/store";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
}

export interface IdeaListProps {
  ideas: IIdea[];
  totalPages: number;
}

export const IdeaList = observer(function IdeaList({ ideas, totalPages }: IdeaListProps): JSX.Element {
  const { userStore } = useStore();
  const searchParameters: { page?: number; search?: string } = useSearch({
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

  const createButton = userStore.isAuthenticated ? (
    <Link to="/idea/new">
      <Button>
        <Plus size={18} />
        Создать идею
      </Button>
    </Link>
  ) : undefined;

  return (
    <DataListLayout
      title="Идеи"
      subtitle="Найдите вдохновение или присоединяйтесь к реализации новой задумки"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={createButton}
      isEmpty={ideas.length === 0}
      emptyMessage="Идеи не найдены"
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
