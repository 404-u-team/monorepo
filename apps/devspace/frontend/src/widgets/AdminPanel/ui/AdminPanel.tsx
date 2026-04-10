import { clsx } from "clsx";
import { Trash2, Plus } from "lucide-react";
import { observer } from "mobx-react-lite";
import { useCallback, useEffect, useState, type JSX } from "react";

import { fetchIdeas, deleteIdea } from "@/entities/idea";
import type { IIdea } from "@/entities/idea";
import { fetchProjects, deleteProject } from "@/entities/project";
import type { IProject } from "@/entities/project";
import { fetchSkills, createSkill, deleteSkill } from "@/entities/skill";
import type { ISkill } from "@/entities/skill";
import { useUserStore } from "@/entities/user";
import { Button, ConfirmModal, Input, Pagination, SearchInput } from "@/shared/ui";

import styles from "./AdminPanel.module.scss";

type Tab = "skills" | "ideas" | "projects";

interface SkillCategory extends ISkill {
  children: ISkill[];
}

const PAGE_LIMIT = 20;

const TAB_LABELS: Record<Tab, string> = {
  skills: "Навыки",
  ideas: "Идеи",
  projects: "Проекты",
};

export const AdminPanel = observer(function AdminPanel(): JSX.Element {
  const userStore = useUserStore();
  const [tab, setTab] = useState<Tab>("skills");

  if (!userStore.isAuthenticated || userStore.user?.nickname !== "admin") {
    return (
      <div className={styles.noAccess}>
        <p className={styles.noAccessTitle}>Нет доступа</p>
        <p>Эта страница доступна только администраторам.</p>
      </div>
    );
  }

  return (
    <div className={styles.wrapper}>
      <h1 className={styles.title}>Панель администратора</h1>

      <div className={styles.tabs}>
        {(["skills", "ideas", "projects"] as Tab[]).map((t) => (
          <button
            key={t}
            type="button"
            className={clsx(styles.tab, tab === t && styles.active)}
            onClick={() => {
              setTab(t);
            }}
          >
            {TAB_LABELS[t]}
          </button>
        ))}
      </div>

      {tab === "skills" && <SkillsTab />}
      {tab === "ideas" && <IdeasTab />}
      {tab === "projects" && <ProjectsTab />}
    </div>
  );
});

// ─────────────────────────────────────────────────────────────
// Skills Tab
// ─────────────────────────────────────────────────────────────
function SkillsTab(): JSX.Element {
  const [categories, setCategories] = useState<SkillCategory[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // New category form
  const [newCategoryName, setNewCategoryName] = useState("");
  const [isAddingCategory, setIsAddingCategory] = useState(false);

  // New skill form
  const [newSkillName, setNewSkillName] = useState("");
  const [newSkillParentId, setNewSkillParentId] = useState("");
  const [isAddingSkill, setIsAddingSkill] = useState(false);

  // Delete confirmation
  const [deletingId, setDeletingId] = useState<string | undefined>(undefined);
  const [deletingName, setDeletingName] = useState("");
  const [deletingIsCategory, setDeletingIsCategory] = useState(false);

  const loadSkills = useCallback(async (): Promise<void> => {
    setIsLoading(true);
    try {
      const all = await fetchSkills({ limit: 200 });
      const roots = all.filter((s) => s.parent_id === null);
      const built: SkillCategory[] = roots.map((root) => ({
        ...root,
        children: all.filter((s) => s.parent_id === root.id),
      }));
      setCategories(built);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadSkills();
  }, [loadSkills]);

  const handleAddCategory = async (): Promise<void> => {
    if (newCategoryName.trim() === "") return;
    setIsAddingCategory(true);
    try {
      await createSkill({ name: newCategoryName.trim() });
      setNewCategoryName("");
      await loadSkills();
    } finally {
      setIsAddingCategory(false);
    }
  };

  const handleAddSkill = async (): Promise<void> => {
    if (newSkillName.trim() === "" || newSkillParentId === "") return;
    setIsAddingSkill(true);
    try {
      await createSkill({ name: newSkillName.trim(), parent_id: newSkillParentId });
      setNewSkillName("");
      await loadSkills();
    } finally {
      setIsAddingSkill(false);
    }
  };

  const confirmDelete = (id: string, name: string, isCategory: boolean): void => {
    setDeletingId(id);
    setDeletingName(name);
    setDeletingIsCategory(isCategory);
  };

  const handleDelete = async (): Promise<void> => {
    if (deletingId === undefined) return;
    try {
      await deleteSkill(deletingId, deletingIsCategory);
      await loadSkills();
    } finally {
      setDeletingId(undefined);
    }
  };

  return (
    <div className={styles.section}>
      {/* Add category */}
      <div className={styles.addForm}>
        <p className={styles.addFormTitle}>Новая категория</p>
        <div className={styles.addFormField}>
          <label htmlFor="cat-name">Название</label>
          <Input
            id="cat-name"
            value={newCategoryName}
            onChange={(event) => {
              setNewCategoryName(event.target.value);
            }}
            placeholder="Frontend, Backend…"
          />
        </div>
        <Button
          onClick={() => {
            void handleAddCategory();
          }}
          disabled={isAddingCategory || newCategoryName.trim() === ""}
        >
          <Plus size={16} />
          Добавить
        </Button>
      </div>

      {/* Add skill to category */}
      <div className={styles.addForm}>
        <p className={styles.addFormTitle}>Новый навык</p>
        <div className={styles.addFormField}>
          <label htmlFor="skill-name">Название</label>
          <Input
            id="skill-name"
            value={newSkillName}
            onChange={(event) => {
              setNewSkillName(event.target.value);
            }}
            placeholder="React, Node.js…"
          />
        </div>
        <div className={styles.addFormField}>
          <label htmlFor="skill-parent">Категория</label>
          <select
            id="skill-parent"
            value={newSkillParentId}
            onChange={(event) => {
              setNewSkillParentId(event.target.value);
            }}
            style={{
              padding: "8px 12px",
              borderRadius: "8px",
              border: "1px solid var(--border--main)",
              backgroundColor: "var(--bg--main)",
              color: "var(--text--primary)",
              fontSize: "14px",
            }}
          >
            <option value="">Выберите категорию</option>
            {categories.map((c) => (
              <option key={c.id} value={c.id}>
                {c.name}
              </option>
            ))}
          </select>
        </div>
        <Button
          onClick={() => {
            void handleAddSkill();
          }}
          disabled={isAddingSkill || newSkillName.trim() === "" || newSkillParentId === ""}
        >
          <Plus size={16} />
          Добавить
        </Button>
      </div>

      {/* Tree */}
      {isLoading ? (
        <p style={{ color: "var(--text--secondary)" }}>Загрузка навыков…</p>
      ) : (
        <div className={styles.skillTree}>
          {categories.map((cat) => (
            <div key={cat.id} className={styles.categoryBlock}>
              <div className={styles.categoryRow}>
                <div className={styles.categoryInfo}>
                  {cat.color !== undefined && (
                    <span
                      className={styles.categoryDot}
                      style={{ backgroundColor: `#${cat.color}` }}
                    />
                  )}
                  <span className={styles.categoryName}>{cat.name}</span>
                </div>
                <div className={styles.rowActions}>
                  <Button
                    variant="outline"
                    onClick={() => {
                      confirmDelete(cat.id, cat.name, true);
                    }}
                    style={{ padding: "4px 8px" }}
                    aria-label="Удалить категорию"
                  >
                    <Trash2 size={14} />
                  </Button>
                </div>
              </div>

              {cat.children.length > 0 && (
                <div className={styles.childList}>
                  {cat.children.map((child) => (
                    <div key={child.id} className={styles.childRow}>
                      <span className={styles.childName}>{child.name}</span>
                      <div className={styles.rowActions}>
                        <Button
                          variant="outline"
                          onClick={() => {
                            confirmDelete(child.id, child.name, false);
                          }}
                          style={{ padding: "4px 8px" }}
                          aria-label="Удалить навык"
                        >
                          <Trash2 size={14} />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      <ConfirmModal
        isOpen={deletingId !== undefined}
        title={deletingIsCategory ? "Удалить категорию?" : "Удалить навык?"}
        description={
          deletingIsCategory
            ? `Категория «${deletingName}» и все вложенные навыки будут удалены. Это действие нельзя отменить.`
            : `Навык «${deletingName}» будет удалён. Это действие нельзя отменить.`
        }
        severity="danger"
        confirmLabel="Удалить"
        cancelLabel="Отмена"
        onConfirm={() => {
          void handleDelete();
        }}
        onCancel={() => {
          setDeletingId(undefined);
        }}
      />
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Ideas Tab
// ─────────────────────────────────────────────────────────────
function IdeasTab(): JSX.Element {
  const [ideas, setIdeas] = useState<IIdea[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [deletingId, setDeletingId] = useState<string | undefined>(undefined);
  const [isDeleting, setIsDeleting] = useState(false);

  const load = useCallback(async (p: number, q: string): Promise<void> => {
    setIsLoading(true);
    try {
      const result = await fetchIdeas({
        start_at: (p - 1) * PAGE_LIMIT,
        limit: PAGE_LIMIT,
        search: q || undefined,
      });
      setIdeas(result.items);
      setTotal(result.total);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(page, search);
  }, [load, page, search]);

  const handleSearch = (value: string): void => {
    setSearch(value);
    setPage(1);
  };

  const handleDelete = async (): Promise<void> => {
    if (deletingId === undefined) return;
    setIsDeleting(true);
    try {
      await deleteIdea(deletingId);
      await load(page, search);
    } finally {
      setIsDeleting(false);
      setDeletingId(undefined);
    }
  };

  const totalPages = Math.ceil(total / PAGE_LIMIT) || 1;

  return (
    <div className={styles.section}>
      <div className={styles.searchRow}>
        <SearchInput value={search} onSearch={handleSearch} placeholder="Поиск по идеям…" />
        <span style={{ fontSize: "13px", color: "var(--text--secondary)" }}>Всего: {total}</span>
      </div>

      <div className={styles.table}>
        <div className={clsx(styles.tableHeader, styles.ideasGrid)}>
          <span>Название</span>
          <span>Автор</span>
          <span>Дата</span>
          <span />
        </div>
        {isLoading && (
          <div style={{ padding: "24px", textAlign: "center", color: "var(--text--secondary)" }}>
            Загрузка…
          </div>
        )}
        {!isLoading && ideas.length === 0 && (
          <div style={{ padding: "24px", textAlign: "center", color: "var(--text--secondary)" }}>
            Ничего не найдено
          </div>
        )}
        {!isLoading &&
          ideas.length > 0 &&
          ideas.map((idea) => (
            <div key={idea.id} className={clsx(styles.tableRow, styles.ideasGrid)}>
              <span className={styles.cellTitle}>{idea.title}</span>
              <span className={styles.cellSecondary}>{idea.author_id.slice(0, 8)}…</span>
              <span className={styles.cellSecondary}>
                {new Date(idea.created_at).toLocaleDateString("ru")}
              </span>
              <div className={styles.cellActions}>
                <Button
                  variant="outline"
                  onClick={() => {
                    setDeletingId(idea.id);
                  }}
                  style={{ padding: "4px 8px" }}
                  aria-label="Удалить идею"
                >
                  <Trash2 size={14} />
                </Button>
              </div>
            </div>
          ))}
      </div>

      {totalPages > 1 && (
        <div className={styles.pagination}>
          <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
        </div>
      )}

      <ConfirmModal
        isOpen={deletingId !== undefined}
        title="Удалить идею?"
        description="Идея будет удалена безвозвратно."
        severity="danger"
        confirmLabel="Удалить"
        cancelLabel="Отмена"
        isLoading={isDeleting}
        onConfirm={() => {
          void handleDelete();
        }}
        onCancel={() => {
          setDeletingId(undefined);
        }}
      />
    </div>
  );
}

// ─────────────────────────────────────────────────────────────
// Projects Tab
// ─────────────────────────────────────────────────────────────
function ProjectsTab(): JSX.Element {
  const [projects, setProjects] = useState<IProject[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [deletingId, setDeletingId] = useState<string | undefined>(undefined);
  const [isDeleting, setIsDeleting] = useState(false);

  const load = useCallback(async (p: number, q: string): Promise<void> => {
    setIsLoading(true);
    try {
      const result = await fetchProjects({
        start_at: (p - 1) * PAGE_LIMIT,
        limit: PAGE_LIMIT,
        search: q || undefined,
      });
      setProjects(result.items);
      setTotal(result.total);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(page, search);
  }, [load, page, search]);

  const handleSearch = (value: string): void => {
    setSearch(value);
    setPage(1);
  };

  const handleDelete = async (): Promise<void> => {
    if (deletingId === undefined) return;
    setIsDeleting(true);
    try {
      await deleteProject(deletingId);
      await load(page, search);
    } finally {
      setIsDeleting(false);
      setDeletingId(undefined);
    }
  };

  const totalPages = Math.ceil(total / PAGE_LIMIT) || 1;

  return (
    <div className={styles.section}>
      <div className={styles.searchRow}>
        <SearchInput value={search} onSearch={handleSearch} placeholder="Поиск по проектам…" />
        <span style={{ fontSize: "13px", color: "var(--text--secondary)" }}>Всего: {total}</span>
      </div>

      <div className={styles.table}>
        <div className={clsx(styles.tableHeader, styles.projectsGrid)}>
          <span>Название</span>
          <span>Лидер</span>
          <span>Статус</span>
          <span />
        </div>
        {isLoading && (
          <div style={{ padding: "24px", textAlign: "center", color: "var(--text--secondary)" }}>
            Загрузка…
          </div>
        )}
        {!isLoading && projects.length === 0 && (
          <div style={{ padding: "24px", textAlign: "center", color: "var(--text--secondary)" }}>
            Ничего не найдено
          </div>
        )}
        {!isLoading &&
          projects.length > 0 &&
          projects.map((project) => (
            <div key={project.id} className={clsx(styles.tableRow, styles.projectsGrid)}>
              <span className={styles.cellTitle}>{project.title}</span>
              <span className={styles.cellSecondary}>{project.leader_id.slice(0, 8)}…</span>
              <span>
                <span
                  className={project.status === "open" ? styles.statusOpen : styles.statusClosed}
                >
                  {project.status === "open" ? "Открыт" : "Закрыт"}
                </span>
              </span>
              <div className={styles.cellActions}>
                <Button
                  variant="outline"
                  onClick={() => {
                    setDeletingId(project.id);
                  }}
                  style={{ padding: "4px 8px" }}
                  aria-label="Удалить проект"
                >
                  <Trash2 size={14} />
                </Button>
              </div>
            </div>
          ))}
      </div>

      {totalPages > 1 && (
        <div className={styles.pagination}>
          <Pagination currentPage={page} totalPages={totalPages} onPageChange={setPage} />
        </div>
      )}

      <ConfirmModal
        isOpen={deletingId !== undefined}
        title="Удалить проект?"
        description="Проект будет удалён безвозвратно."
        severity="danger"
        confirmLabel="Удалить"
        cancelLabel="Отмена"
        isLoading={isDeleting}
        onConfirm={() => {
          void handleDelete();
        }}
        onCancel={() => {
          setDeletingId(undefined);
        }}
      />
    </div>
  );
}
