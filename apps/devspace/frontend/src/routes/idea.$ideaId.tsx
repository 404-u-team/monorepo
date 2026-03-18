import { createFileRoute } from '@tanstack/react-router';
import { IdeaDetail } from '@/widgets/IdeaDetail';
import type { JSX } from 'react';

export const Route = createFileRoute('/idea/$ideaId')({
    component: IdeaDetailPage,
});

function IdeaDetailPage(): JSX.Element {
    return (
        <div style={{ padding: '0 24px' }}>
            <IdeaDetail />
        </div>
    );
}
