import { createFileRoute } from '@tanstack/react-router';
import { fetchIdeas } from '@/entities/idea';
import { IdeaList } from '@/widgets/IdeaList';
import type { JSX } from 'react';

export interface IdeaSearch {
    page?: number | undefined;
    search?: string | undefined;
}

export const Route = createFileRoute('/ideas')({
    validateSearch: (search: Record<string, unknown>): IdeaSearch => {
        return {
            page: Number(search.page) || 1,
            search: search.search as string | undefined,
        };
    },
    loaderDeps: ({ search: { page, search } }) => ({
        page,
        search,
    }),
    loader: async ({ deps }) => {
        const limit = 20;
        const currentPage = deps.page ?? 1;
        const start_at = (currentPage - 1) * limit;

        return fetchIdeas({
            search: deps.search,
            start_at,
            limit,
        });
    },
    component: IdeasPage,
});

function IdeasPage(): JSX.Element {
    const data = Route.useLoaderData();

    return (
        <div style={{ padding: '24px', maxWidth: '1200px', margin: '0 auto' }}>
            <IdeaList
                ideas={data.items}
                totalPages={data.totalPages}
            />
        </div>
    );
}
