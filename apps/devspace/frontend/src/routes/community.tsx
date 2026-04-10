import { createFileRoute } from "@tanstack/react-router";

import { fetchUsers } from "@/entities/user";
import { getPageSize } from "@/shared/lib/pageSize";
import { UserList } from "@/widgets/UserList";

export interface CommunitySearch {
  page?: number | undefined;
  search?: string | undefined;
  main_role?: string | undefined;
  skills?: string[] | undefined;
  limit?: number | undefined;
}

export const Route = createFileRoute("/community")({
  validateSearch: (search: Record<string, unknown>): CommunitySearch => {
    const rawSkills = search.skills;
    let skills: string[] | undefined;
    if (Array.isArray(rawSkills)) {
      skills = rawSkills.filter((s): s is string => typeof s === "string");
    } else if (typeof rawSkills === "string") {
      skills = [rawSkills];
    }
    return {
      page: Number(search.page) || 1,
      search: search.search as string | undefined,
      main_role: search.main_role as string | undefined,
      skills: skills && skills.length > 0 ? skills : undefined,
      limit: Number(search.limit) || undefined,
    };
  },
  loaderDeps: ({ search: { page, search, main_role, skills, limit } }) => ({
    page,
    search,
    main_role,
    skills,
    limit,
  }),
  loader: async ({ deps }) => {
    const limit = deps.limit ?? getPageSize();
    const { page, ...restDeps } = deps;
    const currentPage = page ?? 1;
    const start_at = (currentPage - 1) * limit;

    return fetchUsers({ ...restDeps, start_at, limit });
  },
  component: CommunityPage,
});

function CommunityPage(): React.JSX.Element {
  const data = Route.useLoaderData();

  return (
    <div style={{ padding: "24px", maxWidth: "1200px", margin: "0 auto" }}>
      <UserList users={data.items} totalPages={data.totalPages} total={data.total} />
    </div>
  );
}
