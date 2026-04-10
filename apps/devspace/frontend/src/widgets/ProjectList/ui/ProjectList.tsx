import { useNavigate, useSearch, Link } from "@tanstack/react-router";
import { Plus } from "lucide-react";
import { observer } from "mobx-react-lite";
import { type JSX, useCallback } from "react";

import { fetchSkills } from "@/entities/skill";
import { ProjectCard } from "@/entities/project";
import { useStore } from "@/shared/lib/store";
import { setPageSize, PAGE_SIZE_OPTIONS, type PageSize } from "@/shared/lib/pageSize";
import {
  Dropdown,
  DataListLayout,
  Button,
  SkillMultiSelect,
  type SkillMultiSelectOption,
} from "@/shared/ui";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
  status?: "open" | "closed" | undefined;
  open_slots?: boolean | undefined;
  slots_skills?: string[] | undefined;
  min_people?: number | undefined;
  max_people?: number | undefined;
  limit?: number | undefined;
}

export interface ProjectListProps {
  projects: { id: string }[];
  totalPages: number;
  total: number;
}

const statusOptions = [
  { label: "Все статусы", value: "" },
  { label: "Открытые", value: "open" },
  { label: "Закрытые", value: "closed" },
];

const pageSizeOptions = PAGE_SIZE_OPTIONS.map((n) => ({ label: `${String(n)} / стр.`, value: String(n) }));

export const ProjectList = observer(function ProjectList({
  projects,
  totalPages,
  total,
}: ProjectListProps): JSX.Element {
  const { userStore } = useStore();
  const searchParameters: SearchParameters = useSearch({ strict: false });
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

  const handleOpenSlotsToggle = (): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        open_slots: previous.open_slots === true ? undefined : true,
        page: 1,
      }),
    });
  };

  const handleSlotsSkillsChange = (options: SkillMultiSelectOption[]): void => {
    void navigate({
      search: (previous: SearchParameters) => ({
        ...previous,
        slots_skills: options.length > 0 ? options.map((o) => o.id) : undefined,
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

  const handlePageSizeChange = (value: string): void => {
    const size = Number(value) as PageSize;
    setPageSize(size);
    void navigate({
      search: (previous: SearchParameters) => ({ ...previous, limit: size, page: 1 }),
    });
  };

  const loadSlotsSkills = useCallback(async (query: string): Promise<SkillMultiSelectOption[]> => {
    const skills = await fetchSkills({ search: query || undefined, limit: 30 });
    return skills
      .filter((skill) => skill.parent_id !== null)
      .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
  }, []);

  const slotsSkillsValue: SkillMultiSelectOption[] = (searchParameters.slots_skills ?? []).map(
    (id) => ({ id, name: id, color: undefined }),
  );

  const currentLimit = searchParameters.limit ?? 20;
  const isOpenSlotsActive = searchParameters.open_slots === true;

  const filtersNode = (
    <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
      <div>
        <Dropdown
          options={statusOptions}
          value={(searchParameters as Record<string, string>).status ?? ""}
          onChange={handleStatusChange}
        />
      </div>
      <div>
        <Button
          variant={isOpenSlotsActive ? "primary" : "outline"}
          onClick={handleOpenSlotsToggle}
          style={{ width: "100%" }}
        >
          Только со слотами
        </Button>
      </div>
      <div>
        <SkillMultiSelect
          value={slotsSkillsValue}
          onChange={handleSlotsSkillsChange}
          loadOptions={loadSlotsSkills}
          placeholder="Навыки в слотах..."
          label="Навыки в слотах"
          max={5}
        />
      </div>
    </div>
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
      title="Проекты"
      subtitle="Найдите интересный проект и присоединяйтесь к командной работе"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={controlsNode}
      filtersNode={filtersNode}
      isEmpty={projects.length === 0}
      emptyMessage="Проекты не найдены"
      currentPage={Number((searchParameters as Record<string, string>).page) || 1}
      totalPages={totalPages}
      total={total}
      onPageChange={handlePageChange}
    >
      {projects.map((project) => (
        <ProjectCard key={project.id} projectId={project.id} to={`/project/${project.id}`} fromRoute="/projects" />
      ))}
    </DataListLayout>
  );
});
