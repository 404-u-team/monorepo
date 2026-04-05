import { createFileRoute } from "@tanstack/react-router";

import { fetchUsers } from "@/entities/user";
import { UserList } from "@/widgets/UserList";

export interface CommunitySearch {
  page?: number | undefined;
  search?: string | undefined;
  main_role?: string | undefined;
  skills?: string[] | undefined;
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
    };
  },
  loaderDeps: ({ search: { page, search, main_role, skills } }) => ({
    page,
    search,
    main_role,
    skills,
  }),
  loader: async ({ deps }) => {
    const limit = 20;
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
      <UserList users={data.items} totalPages={data.totalPages} />
    </div>
  );
}
