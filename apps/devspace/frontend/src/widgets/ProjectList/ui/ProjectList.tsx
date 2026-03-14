import { type JSX } from 'react';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { ProjectCard, type IProject } from '@/entities/project';
import { Dropdown, DataListLayout } from '@/shared/ui';

interface SearchParameters {
    page?: number | undefined;
    search?: string | undefined;
    status?: 'open' | 'closed' | undefined;
}

export interface ProjectListProps {
    projects: IProject[];
    totalPages: number;
}

const statusOptions = [
    { label: 'Все статусы', value: '' },
    { label: 'Открытые', value: 'open' },
    { label: 'Закрытые', value: 'closed' },
];

export function ProjectList({ projects, totalPages }: ProjectListProps): JSX.Element {
    const searchParameters = useSearch({ strict: false });
    const navigate = useNavigate({ from: '/projects' });

    const handleSearch = (value: string): void => {
        void navigate({
            search: (previous: SearchParameters) => ({ ...previous, search: value || undefined, page: 1 }),
        });
    };

    const handleStatusChange = (value: string): void => {
        void navigate({
            search: (previous: SearchParameters) => ({ 
                ...previous, 
                status: value === '' ? undefined : (value as 'open' | 'closed'), 
                page: 1 
            }),
        });
    };

    const handlePageChange = (page: number): void => {
        window.scrollTo({ top: 0, behavior: 'smooth' });
        void navigate({
            search: (previous: SearchParameters) => ({ ...previous, page }),
        });
    };

    const StatusFilter = (
        <Dropdown
            options={statusOptions}
            value={(searchParameters as Record<string, string>).status ?? ''}
            onChange={handleStatusChange}
        />
    );

    return (
        <DataListLayout
            title="Проекты"
            subtitle="Найдите интересный проект и присоединяйтесь к командой работе"
            searchValue={(searchParameters as Record<string, string>).search ?? ''}
            onSearchChange={handleSearch}
            controlsNode={StatusFilter}
            isEmpty={projects.length === 0}
            emptyMessage="Проекты не найдены"
            currentPage={Number((searchParameters as Record<string, string>).page) || 1}
            totalPages={totalPages}
            onPageChange={handlePageChange}
        >
            {projects.map((project) => (
                <ProjectCard
                    key={project.id}
                    projectId={project.id}
                    to={`/projects/${project.id}`}
                />
            ))}
        </DataListLayout>
    );
}
