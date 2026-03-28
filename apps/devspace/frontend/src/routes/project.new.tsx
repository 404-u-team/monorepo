import { createFileRoute } from '@tanstack/react-router';
import { CreateProjectForm } from '@/features/project/create';
import type { JSX } from 'react';

export interface ProjectNewSearch {
    title?: string | undefined;
    description?: string | undefined;
    idea_id?: string | undefined;
}

export const Route = createFileRoute('/project/new')({
    validateSearch: (search: Record<string, unknown>): ProjectNewSearch => {
        return {
            title: typeof search.title === 'string' ? search.title : undefined,
            description: typeof search.description === 'string' ? search.description : undefined,
            idea_id: typeof search.idea_id === 'string' ? search.idea_id : undefined,
        };
    },
    component: CreateProjectPage,
});

function CreateProjectPage(): JSX.Element {
    const { title, description, idea_id } = Route.useSearch();

    return (
        <div style={{ padding: '0 24px' }}>
            <CreateProjectForm
                initialTitle={title}
                initialDescription={description}
                initialIdeaId={idea_id}
            />
        </div>
    );
}
