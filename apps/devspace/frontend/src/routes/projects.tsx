import { createFileRoute } from "@tanstack/react-router";

import { fetchProjects } from "@/entities/project";
import { getPageSize } from "@/shared/lib/pageSize";
import { ProjectList } from "@/widgets/ProjectList";

export interface ProjectSearch {
  page?: number | undefined;
  search?: string | undefined;
  status?: "open" | "closed" | undefined;
  open_slots?: boolean | undefined;
  slots_skills?: string[] | undefined;
  min_people?: number | undefined;
  max_people?: number | undefined;
  limit?: number | undefined;
}

export const Route = createFileRoute("/projects")({
  validateSearch: (search: Record<string, unknown>): ProjectSearch => {
    const rawSkills = search.slots_skills;
    let slots_skills: string[] | undefined;
    if (Array.isArray(rawSkills)) {
      slots_skills = rawSkills.filter((s): s is string => typeof s === "string");
    } else if (typeof rawSkills === "string") {
      slots_skills = [rawSkills];
    }

    const minPeople = Number(search.min_people);
    const maxPeople = Number(search.max_people);

    return {
      page: Number(search.page) || 1,
      search: search.search as string | undefined,
      status: search.status === "open" || search.status === "closed" ? search.status : undefined,
      open_slots: search.open_slots === true || search.open_slots === "true" ? true : undefined,
      slots_skills: slots_skills && slots_skills.length > 0 ? slots_skills : undefined,
      min_people: minPeople > 0 ? minPeople : undefined,
      max_people: maxPeople > 0 ? maxPeople : undefined,
      limit: Number(search.limit) || undefined,
    };
  },
  loaderDeps: ({ search: { page, search, status, open_slots, slots_skills, min_people, max_people, limit } }) => ({
    page,
    search,
    status,
    open_slots,
    slots_skills,
    min_people,
    max_people,
    limit,
  }),
  loader: async ({ deps }) => {
    const limit = deps.limit ?? getPageSize();
    const { page, ...restDeps } = deps;
    const currentPage = page ?? 1;
    const start_at = (currentPage - 1) * limit;

    return fetchProjects({
      ...restDeps,
      start_at,
      limit,
    });
  },
  component: ProjectsPage,
});

function ProjectsPage(): React.JSX.Element {
  const data = Route.useLoaderData();

  return (
    <div style={{ padding: "24px", maxWidth: "1200px", margin: "0 auto" }}>
      <ProjectList projects={data.items} totalPages={data.totalPages} total={data.total} />
    </div>
  );
}
