import { useNavigate, useSearch } from "@tanstack/react-router";
import { type JSX, useCallback } from "react";

import { fetchSkills } from "@/entities/skill";
import { UserCard } from "@/entities/user";
import type { IUserResponse } from "@/entities/user";
import { setPageSize, PAGE_SIZE_OPTIONS, type PageSize } from "@/shared/lib/pageSize";
import {
  DataListLayout,
  Dropdown,
  SkillSearch,
  SkillMultiSelect,
  type SkillSearchOption,
  type SkillMultiSelectOption,
} from "@/shared/ui";

interface SearchParameters {
  page?: number | undefined;
  search?: string | undefined;
  main_role?: string | undefined;
  skills?: string[] | undefined;
  limit?: number | undefined;
}

export interface UserListProps {
  users: IUserResponse[];
  totalPages: number;
  total: number;
}

const pageSizeOptions = PAGE_SIZE_OPTIONS.map((n) => ({
  label: `${String(n)} / стр.`,
  value: String(n),
}));

export function UserList({ users, totalPages, total }: UserListProps): JSX.Element {
  const searchParameters: SearchParameters = useSearch({ strict: false });
  const navigate = useNavigate({ from: "/community" });

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
    window.scrollTo({ top: 0, behavior: "smooth" });
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

  const handlePageSizeChange = (value: string): void => {
    const size = Number(value) as PageSize;
    setPageSize(size);
    void navigate({
      search: (previous: SearchParameters) => ({ ...previous, limit: size, page: 1 }),
    });
  };

  // Root-level skills (main roles) — no parent_id means root skills per API
  const loadRoles = useCallback(async (query: string): Promise<SkillSearchOption[]> => {
    const skills = await fetchSkills({ search: query || undefined, limit: 20 });
    return skills
      .filter((skill) => skill.parent_id === null)
      .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
  }, []);

  // Child skills for the skills filter
  const loadSkills = useCallback(async (query: string): Promise<SkillMultiSelectOption[]> => {
    const skills = await fetchSkills({ search: query || undefined, limit: 30 });
    return skills
      .filter((skill) => skill.parent_id !== null)
      .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
  }, []);

  const mainRoleValue: SkillSearchOption | undefined =
    searchParameters.main_role !== undefined && searchParameters.main_role !== ""
      ? { id: searchParameters.main_role, name: searchParameters.main_role, color: undefined }
      : undefined;

  const skillsValue: SkillMultiSelectOption[] = (searchParameters.skills ?? []).map((id) => ({
    id,
    name: id,
    color: undefined,
  }));

  const currentLimit = searchParameters.limit ?? 20;

  const filtersNode = (
    <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
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

  const controlsNode = (
    <Dropdown
      options={pageSizeOptions}
      value={String(currentLimit)}
      onChange={handlePageSizeChange}
    />
  );

  return (
    <DataListLayout
      title="Сообщество"
      subtitle="Найдите специалистов и единомышленников для совместной работы"
      searchValue={(searchParameters as Record<string, string>).search ?? ""}
      onSearchChange={handleSearch}
      controlsNode={controlsNode}
      filtersNode={filtersNode}
      isEmpty={users.length === 0}
      emptyMessage="Пользователи не найдены"
      currentPage={Number((searchParameters as Record<string, string>).page) || 1}
      totalPages={totalPages}
      total={total}
      onPageChange={handlePageChange}
    >
      {users.map((user) => (
        <UserCard key={user.id} id={user.id} to={`/users/${user.id}`} fromRoute="/community" />
      ))}
    </DataListLayout>
  );
}
