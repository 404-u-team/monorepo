import { useNavigate, useSearch, Link } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { observer } from "mobx-react-lite";
import { type JSX } from "react";

import { ProjectCard, type IProject } from "@/entities/project";
import { useStore } from "@/shared/lib/store";
import { Dropdown, DataListLayout, Button } from "@/shared/ui";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
  status?: "open" | "closed" | undefined;
}

export interface ProjectListProps {
  projects: IProject[];
  totalPages: number;
}

const statusOptions = [
  { label: "Все статусы", value: "" },
  { label: "Открытые", value: "open" },
  { label: "Закрытые", value: "closed" },
];

export const ProjectList = observer(function ProjectList({
  projects,
  totalPages,
}: ProjectListProps): JSX.Element {
  const { userStore } = useStore();
  const searchParameters: SearchParameters = useSearch({
    strict: false,
  });
  const navigate = useNavigate({ from: "/projects" });

  const handleSearch = (value: string): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        search: value || undefined,
        page: 1,
      }),
    });
  };

  const handleStatusChange = (value: string): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        status: value === "" ? undefined : (value as "open" | "closed"),
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

  const StatusFilter = (
    <Dropdown
      options={statusOptions}
      value={(searchParameters as Record<string, string>).status ?? ""}
      onChange={handleStatusChange}
    />
  );

  const createButton = userStore.isAuthenticated ? (
    <Link to="/project/new">
      <Button>
        <Plus size={18} />
        Создать проект
      </Button>
    </Link>
  ) : undefined;

  const controlsNode = (
    <div style={{ display: "flex", gap: "0.75rem", alignItems: "center" }}>
      {StatusFilter}
      {createButton}
    </div>
  );

  return (
    <DataListLayout
      title="Проекты"
      subtitle="Найдите интересный проект и присоединяйтесь к командной работе"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={controlsNode}
      isEmpty={projects.length === 0}
      emptyMessage="Проекты не найдены"
      currentPage={Number((searchParameters as Record<string, string>).page) || 1}
      totalPages={totalPages}
      onPageChange={handlePageChange}
    >
      {projects.map((project) => (
        <ProjectCard key={project.id} projectId={project.id} to={`/project/${project.id}`} />
      ))}
    </DataListLayout>
  );
});
