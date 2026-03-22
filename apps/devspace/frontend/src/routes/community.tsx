import { createFileRoute } from '@tanstack/react-router';
import { fetchUsers } from '@/entities/user';
import { UserList } from '@/widgets/UserList';

export interface CommunitySearch {
    page?: number | undefined;
    search?: string | undefined;
}

export const Route = createFileRoute('/community')({
    validateSearch: (search: Record<string, unknown>): CommunitySearch => {
        return {
            page: Number(search.page) || 1,
            search: search.search as string | undefined,
        };
    },
    loaderDeps: ({ search: { page, search } }) => ({ page, search }),
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
        <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
            <UserList users={data.items} totalPages={data.totalPages} />
        </div>
    );
}
