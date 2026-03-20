import { createFileRoute } from '@tanstack/react-router';
import { ProjectDetail } from '@/widgets/ProjectDetail';
import type { JSX } from 'react';

export const Route = createFileRoute('/project/$projectId')({
    component: ProjectDetailPage,
});

function ProjectDetailPage(): JSX.Element {
    return (
        <div style={{ padding: '0 24px' }}>
            <ProjectDetail />
        </div>
    );
}
