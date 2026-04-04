import { type JSX, useCallback } from 'react';
import { useNavigate, useSearch } from '@tanstack/react-router';
import { UserCard } from '@/entities/user';
import type { IUserResponse } from '@/entities/user';
import { fetchSkills } from '@/entities/skill';
import { DataListLayout, SkillSearch, SkillMultiSelect, type SkillSearchOption, type SkillMultiSelectOption } from '@/shared/ui';

interface SearchParameters {
    page?: number | undefined;
    search?: string | undefined;
    main_role?: string | undefined;
    skills?: string[] | undefined;
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

    const handleMainRoleChange = (option: SkillSearchOption | undefined): void => {
        void navigate({
            search: (previous: SearchParameters) => ({
                ...previous,
                main_role: option?.id,
                page: 1,
            }),
        });
    };

    const handleSkillsChange = (options: SkillMultiSelectOption[]): void => {
        void navigate({
            search: (previous: SearchParameters) => ({
                ...previous,
                skills: options.length > 0 ? options.map((o) => o.id) : undefined,
                page: 1,
            }),
        });
    };

    const loadRoles = useCallback(async (query: string): Promise<SkillSearchOption[]> => {
        const skills = await fetchSkills({ search: query || undefined, limit: 20 });
        return skills
            .filter((skill) => skill.parent_id === null)
            .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
    }, []);

    const loadSkills = useCallback(async (query: string): Promise<SkillMultiSelectOption[]> => {
        const skills = await fetchSkills({ search: query || undefined, limit: 30 });
        return skills
            .filter((skill) => skill.parent_id !== null)
            .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
    }, []);

    // Build current SkillSearchOption for main_role from URL (need id; name unknown — show empty or fetch)
    const mainRoleValue: SkillSearchOption | undefined = searchParameters.main_role
        ? { id: searchParameters.main_role, name: searchParameters.main_role, color: undefined }
        : undefined;

    // Build current SkillMultiSelectOption[] for skills from URL
    const skillsValue: SkillMultiSelectOption[] = (searchParameters.skills ?? []).map((id) => ({
        id,
        name: id,
        color: undefined,
    }));

    const filtersNode = (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
            <div>
                <SkillSearch
                    value={mainRoleValue}
                    onChange={handleMainRoleChange}
                    loadOptions={loadRoles}
                    placeholder="Основная роль..."
                    label="Роль"
                />
            </div>
            <div>
                <SkillMultiSelect
                    value={skillsValue}
                    onChange={handleSkillsChange}
                    loadOptions={loadSkills}
                    placeholder="Навыки..."
                    label="Навыки"
                    max={10}
                />
            </div>
        </div>
    );

    return (
        <DataListLayout
            title="Сообщество"
            subtitle="Найдите специалистов и единомышленников для совместной работы"
            searchValue={(searchParameters as Record<string, string>).search ?? ''}
            onSearchChange={handleSearch}
            filtersNode={filtersNode}
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
