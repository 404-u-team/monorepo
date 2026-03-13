import { createFileRoute } from '@tanstack/react-router';
import { fetchProjects, type PaginatedProjects } from '@/entities/project';
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
        return fetchProjects({
            ...deps,
            limit: 20,
        });
    },
    component: ProjectsPage,
});

function ProjectsPage(): React.JSX.Element {
    const data = Route.useLoaderData() as PaginatedProjects;

    return (
        <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
            <ProjectList
                projects={data.items}
                totalPages={data.totalPages}
            />
        </div>
    );
}
