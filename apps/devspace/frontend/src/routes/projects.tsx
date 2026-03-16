import { createFileRoute } from '@tanstack/react-router';
import { fetchProjects } from '@/entities/project';
import { ProjectList } from '@/widgets/ProjectList';

export interface ProjectSearch {
    page?: number | undefined;
    search?: string | undefined;
    status?: 'open' | 'closed' | undefined;
}

export const Route = createFileRoute('/projects')({
    validateSearch: (search: Record<string, unknown>): ProjectSearch => {
        return {
            page: Number(search.page) || 1,
            search: search.search as string | undefined,
            status: (search.status === 'open' || search.status === 'closed') ? search.status : undefined,
        };
    },
    loaderDeps: ({ search: { page, search, status } }) => ({
        page,
        search,
        status,
    }),
    loader: async ({ deps }) => {
        const limit = 20;
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
        <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
            <ProjectList
                projects={data.items}
                totalPages={data.totalPages}
            />
        </div>
    );
}
