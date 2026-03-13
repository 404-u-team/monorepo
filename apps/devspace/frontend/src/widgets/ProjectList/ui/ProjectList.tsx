import { type JSX } from 'react';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { ProjectCard, type IProject } from '@/entities/project';
import { Pagination, SearchInput, Dropdown } from '@/shared/ui';
import styles from './ProjectList.module.scss';
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

    return (
        <div className={styles.wrapper}>
            <div className={styles.header}>
                <h1 className={styles.title}>Проекты</h1>
                <p className={styles.subtitle}>Найдите интересный проект и присоединяйтесь к командой работе</p>
            </div>

            <div className={styles.controls}>
                <SearchInput
                    value={(searchParameters as Record<string, string>).search ?? ''}
                    onSearch={handleSearch}
                    placeholder="Название или описание..."
                    className={styles.search ?? ''}
                />
                <Dropdown
                    options={statusOptions}
                    value={(searchParameters as Record<string, string>).status ?? ''}
                    onChange={handleStatusChange}
                    className={styles.filter ?? ''}
                />
            </div>

            <div className={styles.grid}>
                {projects.map((project) => (
                    <ProjectCard
                        key={project.id}
                        projectId={project.id}
                        to={`/projects/${project.id}`}
                    />
                ))}
            </div>

            {projects.length === 0 && (
                <div className={styles.empty}>
                    <p>Проекты не найдены</p>
                </div>
            )}

            {totalPages > 1 && (
                <div className={styles.pagination}>
                    <Pagination
                        currentPage={Number((searchParameters as Record<string, string>).page) || 1}
                        totalPages={totalPages}
                        onPageChange={handlePageChange}
                    />
                </div>
            )}
        </div>
    );
}
