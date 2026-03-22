import { type JSX } from 'react';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { UserCard } from '@/entities/user';
import type { IUserResponse } from '@/entities/user';
import { DataListLayout } from '@/shared/ui';

interface SearchParameters {
    page?: number | undefined;
    search?: string | undefined;
}

export interface UserListProps {
    users: IUserResponse[];
    totalPages: number;
}

export function UserList({ users, totalPages }: UserListProps): JSX.Element {
    const searchParameters: SearchParameters = useSearch({ strict: false });
    const navigate = useNavigate({ from: '/community' });

    const handleSearch = (value: string): void => {
        void navigate({
            search: (previous: SearchParameters) => ({
                ...previous,
                search: value || undefined,
                page: 1,
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
        <DataListLayout
            title="Сообщество"
            subtitle="Найдите специалистов и единомышленников для совместной работы"
            searchValue={(searchParameters as Record<string, string>).search ?? ''}
            onSearchChange={handleSearch}
            isEmpty={users.length === 0}
            emptyMessage="Пользователи не найдены"
            currentPage={Number((searchParameters as Record<string, string>).page) || 1}
            totalPages={totalPages}
            onPageChange={handlePageChange}
        >
            {users.map((user) => (
                <UserCard
                    key={user.id}
                    id={user.id}
                    to={`/users/${user.id}`}
                />
            ))}
        </DataListLayout>
    );
}
