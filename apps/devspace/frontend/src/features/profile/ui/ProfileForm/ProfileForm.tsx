import { Link } from "@tanstack/react-router";
import { clsx } from "clsx";
import { Camera, Check, Plus, User, X } from "lucide-react";
import type { JSX } from "react";
import { useCallback, useEffect, useRef, useState } from "react";

import {
  fetchProjects,
  ProjectCard,
  getMyRequests,
  acceptRequest,
  rejectRequest,
  type IRequest,
} from "@/entities/project";
import { fetchSkills } from "@/entities/skill";
import { apiClient } from "@/shared/api/client";
import { Button, Input, Skeleton, SkillSearch } from "@/shared/ui";
import type { SkillSearchOption } from "@/shared/ui";

// main_role comes as a full object from API now
import styles from "./ProfileForm.module.scss";

export type ProfileFormProps = Record<string, never>;

interface IMainRoleObject {
  id: string;
  name: string;
  parent_id?: string | null;
  color?: string;
}

interface IProfileData {
  id: string;
  email: string;
  nickname: string;
  avatar_url: string | undefined;
  bio: string | undefined;
  main_role: IMainRoleObject | null | undefined;
  skills: IProfileSkill[];
}

interface IProfileSkill {
  id: string;
  name: string;
}

interface IProjectItem {
  id: string;
}

export function ProfileForm(_props: ProfileFormProps): JSX.Element {
  const [userData, setUserData] = useState<IProfileData | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState<string | undefined>(undefined);

  const [nickname, setNickname] = useState("");
  const [bio, setBio] = useState("");
  const [mainRole, setMainRole] = useState<SkillSearchOption | undefined>(undefined);
  const [initialNickname, setInitialNickname] = useState("");
  const [initialBio, setInitialBio] = useState("");
  const [initialMainRole, setInitialMainRole] = useState("");
  const [initialMainRoleOption, setInitialMainRoleOption] = useState<SkillSearchOption | undefined>(
    undefined,
  );
  const [isSaving, setIsSaving] = useState(false);
  const [saveSuccess, setSaveSuccess] = useState(false);
  const [saveError, setSaveError] = useState<string | undefined>(undefined);

  const [skills, setSkills] = useState<IProfileSkill[]>([]);
  const [skillSearchValue, setSkillSearchValue] = useState<SkillSearchOption | undefined>(
    undefined,
  );
  const [isAddingSkill, setIsAddingSkill] = useState(false);
  const [removingSkillId, setRemovingSkillId] = useState<string | undefined>(undefined);

  const [projects, setProjects] = useState<IProjectItem[]>([]);
  const [isProjectsLoading, setIsProjectsLoading] = useState(false);

  const [myRequests, setMyRequests] = useState<IRequest[]>([]);
  const [isRequestsLoading, setIsRequestsLoading] = useState(false);
  const [requestActionId, setRequestActionId] = useState<string | undefined>(undefined);

  const successTimerReference = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    return (): void => {
      if (successTimerReference.current !== null) {
        clearTimeout(successTimerReference.current);
      }
    };
  }, []);

  useEffect(() => {
    let cancelled = false;

    async function load(): Promise<void> {
      try {
        setIsLoading(true);
        setLoadError(undefined);

        const response = await apiClient.get<IProfileData>("/users/me");
        const data = response.data;

        if (!cancelled) {
          setUserData(data);
          setNickname(data.nickname);
          setInitialNickname(data.nickname);
          setBio(data.bio ?? "");
          setInitialBio(data.bio ?? "");
          setSkills(data.skills);

          const roleObject = data.main_role;
          setInitialMainRole(roleObject?.id ?? "");
          if (roleObject?.id !== undefined) {
            const roleOption: SkillSearchOption = {
              id: roleObject.id,
              name: roleObject.name,
              color: roleObject.color,
            };
            setMainRole(roleOption);
            setInitialMainRoleOption(roleOption);
          }

          void loadProjects(data.id);
          void loadRequests();
        }
      } catch {
        if (!cancelled) setLoadError("Не удалось загрузить данные профиля");
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    void load();
    return (): void => {
      cancelled = true;
    };
  }, []);

  const loadProjects = async (userId: string): Promise<void> => {
    try {
      setIsProjectsLoading(true);
      const result = await fetchProjects({ leader_id: userId, limit: 20 });
      setProjects(result.items.map((project) => ({ id: project.id })));
    } catch {
      // no-op — projects section will be empty
    } finally {
      setIsProjectsLoading(false);
    }
  };

  const loadRequests = async (): Promise<void> => {
    try {
      setIsRequestsLoading(true);
      const data = await getMyRequests();
      setMyRequests(data);
    } catch {
      // no-op — section will remain empty
    } finally {
      setIsRequestsLoading(false);
    }
  };

  const handleAcceptInvite = async (requestId: string): Promise<void> => {
    setRequestActionId(requestId);
    try {
      const updated = await acceptRequest(requestId);
      setMyRequests((previous) => previous.map((r) => (r.id === requestId ? updated : r)));
    } catch {
      // no-op
    } finally {
      setRequestActionId(undefined);
    }
  };

  const handleRejectInvite = async (requestId: string): Promise<void> => {
    setRequestActionId(requestId);
    try {
      const updated = await rejectRequest(requestId);
      setMyRequests((previous) => previous.map((r) => (r.id === requestId ? updated : r)));
    } catch {
      // no-op
    } finally {
      setRequestActionId(undefined);
    }
  };

  const handleSave = async (): Promise<void> => {
    if (!hasChanges) return;

    try {
      setIsSaving(true);
      setSaveError(undefined);
      setSaveSuccess(false);

      const currentMainRoleId = mainRole?.id;
      await apiClient.put("/users/me", {
        nickname: nickname.trim(),
        bio,
        main_role: currentMainRoleId,
      });

      const saved = nickname.trim();
      setInitialNickname(saved);
      setNickname(saved);
      setInitialBio(bio);
      setInitialMainRole(currentMainRoleId ?? "");
      setSaveSuccess(true);

      if (successTimerReference.current !== null) clearTimeout(successTimerReference.current);
      successTimerReference.current = setTimeout(() => {
        setSaveSuccess(false);
      }, 3000);
    } catch {
      setSaveError("Не удалось сохранить изменения");
    } finally {
      setIsSaving(false);
    }
  };

  const handleCancel = (): void => {
    setNickname(initialNickname);
    setBio(initialBio);
    setMainRole(initialMainRoleOption);
    setSaveError(undefined);
  };

  const loadMainRoleOptions = useCallback(async (query: string): Promise<SkillSearchOption[]> => {
    const results = await fetchSkills({ search: query !== "" ? query : undefined, limit: 20 });
    return results
      .filter((skill) => skill.parent_id === null)
      .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
  }, []);

  const loadSkillOptions = useCallback(
    async (query: string): Promise<SkillSearchOption[]> => {
      const results = await fetchSkills({ search: query !== "" ? query : undefined, limit: 30 });
      const existingIds = new Set(skills.map((skill) => skill.id));
      return results
        .filter((skill) => !existingIds.has(skill.id))
        .map((skill) => ({ id: skill.id, name: skill.name, color: skill.color }));
    },
    [skills],
  );

  const handleAddSkill = async (selected: SkillSearchOption | undefined): Promise<void> => {
    if (selected === undefined) return;
    if (skills.some((skill) => skill.id === selected.id)) return;

    try {
      setIsAddingSkill(true);
      await apiClient.post("/users/me/skills", { skill_id: selected.id });
      setSkills((previous) => [...previous, { id: selected.id, name: selected.name }]);
      setSkillSearchValue(undefined);
    } catch {
      // skill might already exist or invalid
    } finally {
      setIsAddingSkill(false);
    }
  };

  const handleRemoveSkill = async (skillId: string): Promise<void> => {
    try {
      setRemovingSkillId(skillId);
      await apiClient.delete(`/users/me/skills/${skillId}`);
      setSkills((previous) => previous.filter((skill) => skill.id !== skillId));
    } catch {
      // no-op
    } finally {
      setRemovingSkillId(undefined);
    }
  };

  const hasChanges =
    nickname.trim() !== initialNickname ||
    bio !== initialBio ||
    (mainRole?.id ?? "") !== initialMainRole;

  function renderProjects(): JSX.Element {
    if (isProjectsLoading) {
      return (
        <div className={styles.projectsGrid}>
          <Skeleton height={180} borderRadius={12} />
          <Skeleton height={180} borderRadius={12} />
        </div>
      );
    }
    if (projects.length > 0) {
      return (
        <div className={styles.projectsGrid}>
          {projects.map((project) => (
            <ProjectCard key={project.id} projectId={project.id} to={`/project/${project.id}`} />
          ))}
        </div>
      );
    }
    return (
      <div className={styles.emptyProjects}>
        <p className={styles.emptyHint}>У вас пока нет проектов</p>
        <Link to="/project/new">
          <Button variant="outline">Создать первый проект</Button>
        </Link>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className={styles.layout}>
        <aside className={styles.sidebar}>
          <div className={styles.avatarCard}>
            <Skeleton width={96} height={96} borderRadius={48} />
            <Skeleton width={140} height={20} />
            <Skeleton width={100} height={16} />
          </div>
          <div className={styles.sideNav}>
            <Skeleton width="100%" height={40} borderRadius={8} />
            <Skeleton width="100%" height={40} borderRadius={8} />
          </div>
        </aside>
        <main className={styles.main}>
          <div className={styles.section}>
            <Skeleton width={220} height={24} />
            <div className={styles.formRow}>
              <Skeleton width="100%" height={44} borderRadius={8} />
              <Skeleton width="100%" height={44} borderRadius={8} />
            </div>
          </div>
          <div className={styles.section}>
            <Skeleton width={120} height={24} />
            <Skeleton width="100%" height={120} borderRadius={8} />
          </div>
        </main>
      </div>
    );
  }

  if (loadError !== undefined) {
    return (
      <div className={styles.errorPage}>
        <p className={styles.errorText}>{loadError}</p>
        <Button
          onClick={() => {
            window.location.reload();
          }}
        >
          Попробовать снова
        </Button>
      </div>
    );
  }

  return (
    <div className={styles.layout}>
      {/* Sidebar */}
      <aside className={styles.sidebar}>
        <div className={styles.avatarCard}>
          <div className={styles.avatarContainer}>
            {userData?.avatar_url !== undefined && userData.avatar_url !== "" ? (
              <img
                src={userData.avatar_url}
                alt={userData.nickname}
                className={styles.avatarImage}
                onError={(event) => {
                  (event.target as HTMLImageElement).style.display = "none";
                }}
              />
            ) : (
              <div className={styles.avatarPlaceholder}>
                <User size={36} />
              </div>
            )}
            <button className={styles.avatarEditButton} aria-label="Сменить фото">
              <Camera size={14} />
            </button>
          </div>
          <div className={styles.userSummary}>
            <span className={styles.summaryNickname}>{userData?.nickname}</span>
            {mainRole !== undefined && <span className={styles.summaryRole}>{mainRole.name}</span>}
            <span className={styles.summaryEmail}>{userData?.email}</span>
          </div>
        </div>

        <nav className={styles.sideNav}>
          <button className={clsx(styles.navItem, styles.navItemActive)}>Профиль</button>
          <button className={styles.navItem} disabled>
            Уведомления
          </button>
        </nav>
      </aside>

      {/* Main content */}
      <main className={styles.main}>
        {/* Basic Info */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Основная информация</h2>
          <div className={styles.formRow}>
            <div className={styles.field}>
              <label className={styles.label} htmlFor="profile-nickname">
                Никнейм
              </label>
              <Input
                id="profile-nickname"
                placeholder="Например: john_doe"
                value={nickname}
                onChange={(event) => {
                  setNickname(event.target.value);
                }}
              />
            </div>
            <div className={styles.field}>
              <label className={styles.label} htmlFor="profile-email">
                E-mail
              </label>
              <Input id="profile-email" type="email" value={userData?.email ?? ""} disabled />
            </div>
          </div>
          <div className={styles.field}>
            <label className={styles.label}>Основная роль</label>
            <SkillSearch
              value={mainRole}
              onChange={setMainRole}
              loadOptions={loadMainRoleOptions}
              placeholder="Выберите основную роль..."
            />
          </div>
        </section>

        {/* Bio */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>О себе</h2>
          <textarea
            className={styles.bioTextarea}
            value={bio}
            onChange={(event) => {
              setBio(event.target.value);
            }}
            placeholder="Расскажите о себе, своём опыте и интересах..."
            rows={4}
          />
        </section>

        {/* Save bar */}
        {(hasChanges || saveSuccess || saveError !== undefined) && (
          <div className={styles.saveBar}>
            {saveSuccess && (
              <span className={styles.successText}>
                <Check size={16} />
                Изменения сохранены
              </span>
            )}
            {saveError !== undefined && <span className={styles.errorText}>{saveError}</span>}
            <div className={styles.saveButtons}>
              <Button variant="outline" onClick={handleCancel} disabled={isSaving || !hasChanges}>
                Отмена
              </Button>
              <Button
                onClick={() => {
                  void handleSave();
                }}
                disabled={isSaving || !hasChanges}
              >
                {isSaving ? "Сохранение..." : "Сохранить"}
              </Button>
            </div>
          </div>
        )}

        {/* Skills */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Навыки</h2>

          {skills.length > 0 ? (
            <div className={styles.skillsList}>
              {skills.map((skill) => (
                <span key={skill.id} className={styles.skillTag}>
                  {skill.name}
                  <button
                    className={styles.skillRemove}
                    onClick={() => {
                      void handleRemoveSkill(skill.id);
                    }}
                    disabled={removingSkillId === skill.id}
                    aria-label={`Удалить навык ${skill.name}`}
                  >
                    <X size={12} />
                  </button>
                </span>
              ))}
            </div>
          ) : (
            <p className={styles.emptyHint}>Навыки не добавлены</p>
          )}

          <div className={styles.skillAdd}>
            <SkillSearch
              value={skillSearchValue}
              onChange={(selected) => {
                void handleAddSkill(selected);
              }}
              loadOptions={loadSkillOptions}
              placeholder="Найти и добавить навык..."
              disabled={isAddingSkill}
            />
            {skillSearchValue !== undefined && (
              <Button
                onClick={() => {
                  void handleAddSkill(skillSearchValue);
                }}
                disabled={isAddingSkill}
              >
                <Plus size={16} />
                Добавить
              </Button>
            )}
          </div>
        </section>

        {/* My Requests & Invites */}
        <section className={styles.section}>
          <h2 className={styles.sectionTitle}>Мои заявки и приглашения</h2>
          {((): JSX.Element => {
            if (isRequestsLoading) {
              return <Skeleton height={80} borderRadius={8} />;
            }
            if (myRequests.length === 0) {
              return <p className={styles.emptyHint}>Заявок и приглашений нет</p>;
            }
            return (
              <div className={styles.requestsList}>
                {myRequests.map((request) => {
                  let statusClass = styles.requestStatusRejected;
                  let statusLabel = "Отклонено";

                  if (request.status === "pending") {
                    statusClass = styles.requestStatusPending;
                    statusLabel = "Ожидает";
                  } else if (request.status === "accepted") {
                    statusClass = styles.requestStatusAccepted;
                    statusLabel = "Принято";
                  }

                  return (
                    <div key={request.id} className={styles.requestItem}>
                      <div className={styles.requestMeta}>
                        <span className={styles.requestType}>
                          {request.type === "apply" ? "Отклик" : "Приглашение"}
                        </span>
                        <span className={statusClass}>{statusLabel}</span>
                      </div>
                      {request.cover_letter !== "" && (
                        <p className={styles.requestCoverLetter}>«{request.cover_letter}»</p>
                      )}
                      {request.type === "invite" && request.status === "pending" && (
                        <div className={styles.requestActions}>
                          <Button
                            onClick={() => {
                              void handleAcceptInvite(request.id);
                            }}
                            disabled={requestActionId === request.id}
                          >
                            Принять
                          </Button>
                          <Button
                            variant="outline"
                            onClick={() => {
                              void handleRejectInvite(request.id);
                            }}
                            disabled={requestActionId === request.id}
                          >
                            Отклонить
                          </Button>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            );
          })()}
        </section>

        {/* My Projects */}
        <section className={styles.section}>
          <div className={styles.sectionHeader}>
            <h2 className={styles.sectionTitle}>Мои проекты</h2>
            <Link to="/project/new" className={styles.createLink}>
              <Plus size={14} />
              Создать проект
            </Link>
          </div>

          {renderProjects()}
        </section>
      </main>
    </div>
  );
}
