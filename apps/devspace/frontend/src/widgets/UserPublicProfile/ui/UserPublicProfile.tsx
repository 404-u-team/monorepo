import { Link, useLocation, useRouter } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { useEffect, useState, type JSX } from "react";

import { fetchIdeas } from "@/entities/idea";
import type { IIdea } from "@/entities/idea";
import { fetchProjects } from "@/entities/project";
import type { IProject } from "@/entities/project";
import type { IUserResponse } from "@/entities/user";
import { isValidMainRole } from "@/entities/user/model/IUserResponse";
import { Badge, Skeleton, UserAvatar } from "@/shared/ui";

import styles from "./UserPublicProfile.module.scss";

export interface UserPublicProfileProps {
  user: IUserResponse;
}

export function UserPublicProfile({ user }: UserPublicProfileProps): JSX.Element {
  const mainRoleName = isValidMainRole(user.main_role) ? user.main_role.name : undefined;
  const location = useLocation();
  const router = useRouter();
  const backTo = (location.state as { backTo?: string } | null)?.backTo;

  const [ideas, setIdeas] = useState<IIdea[]>([]);
  const [projects, setProjects] = useState<IProject[]>([]);

  useEffect(() => {
    let cancelled = false;

    void fetchIdeas({ author_id: user.id, limit: 6 }).then((result) => {
      if (!cancelled) setIdeas(result.items);
    });

    void fetchProjects({ leader_id: user.id, limit: 6 }).then((result) => {
      if (!cancelled) setProjects(result.items);
    });

    return (): void => {
      cancelled = true;
    };
  }, [user.id]);

  return (
    <div className={styles.wrapper}>
      {backTo !== undefined && (
        <button
          type="button"
          className={styles.backLink}
          onClick={() => {
            router.history.back();
          }}
        >
          <ArrowLeft size={16} />
          Назад к списку
        </button>
      )}

      <div className={styles.header}>
        <UserAvatar
          avatarUrl={user.avatar_url}
          nickname={user.nickname}
          size={120}
          className={styles.avatarWrapper}
        />

        <div className={styles.info}>
          <h1 className={styles.nickname}>{user.nickname}</h1>
          {mainRoleName !== undefined ? (
            <p className={styles.mainRole}>{mainRoleName}</p>
          ) : (
            <p className={styles.emptyHint}>Основная роль не указана</p>
          )}
        </div>
      </div>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>О себе</h2>
        {user.bio !== "" ? (
          <p className={styles.bio}>{user.bio}</p>
        ) : (
          <p className={styles.emptyHint}>Это поле пока не заполнено</p>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Навыки</h2>
        {user.skills.length > 0 ? (
          <div className={styles.skills}>
            {user.skills.map((skill) => (
              <Badge key={skill.id}>{skill.name}</Badge>
            ))}
          </div>
        ) : (
          <p className={styles.emptyHint}>Навыки пока не добавлены</p>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Идеи</h2>
        {ideas.length > 0 ? (
          <div className={styles.itemList}>
            {ideas.map((idea) => (
              <Link key={idea.id} to="/idea/$ideaId" params={{ ideaId: idea.id }} className={styles.itemCard}>
                <span className={styles.itemTitle}>{idea.title}</span>
                {idea.category !== undefined && idea.category !== "" && (
                  <Badge className={styles.itemBadge}>{idea.category}</Badge>
                )}
                <p className={styles.itemDescription}>{idea.description}</p>
              </Link>
            ))}
          </div>
        ) : (
          <p className={styles.emptyHint}>Идей пока нет</p>
        )}
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionTitle}>Проекты</h2>
        {projects.length > 0 ? (
          <div className={styles.itemList}>
            {projects.map((project) => (
              <Link key={project.id} to="/project/$projectId" params={{ projectId: project.id }} className={styles.itemCard}>
                <div className={styles.itemCardHeader}>
                  <span className={styles.itemTitle}>{project.title}</span>
                  <Badge className={styles[project.status]}>
                    {project.status === "open" ? "Открыт" : "Закрыт"}
                  </Badge>
                </div>
                {project.description !== "" && (
                  <p className={styles.itemDescription}>{project.description}</p>
                )}
              </Link>
            ))}
          </div>
        ) : (
          <p className={styles.emptyHint}>Проектов пока нет</p>
        )}
      </section>
    </div>
  );
}

export function UserPublicProfileSkeleton(): JSX.Element {
  return (
    <div className={styles.wrapper}>
      <div className={styles.header}>
        <Skeleton width={120} height={120} borderRadius={60} />
        <div className={styles.info}>
          <Skeleton width={200} height={32} />
          <Skeleton width={140} height={20} />
        </div>
      </div>

      <div className={styles.section}>
        <Skeleton width={80} height={24} />
        <div style={{ marginTop: "12px", display: "flex", flexDirection: "column", gap: "8px" }}>
          <Skeleton width="100%" height={16} />
          <Skeleton width="90%" height={16} />
          <Skeleton width="70%" height={16} />
        </div>
      </div>

      <div className={styles.section}>
        <Skeleton width={80} height={24} />
        <div style={{ marginTop: "12px", display: "flex", gap: "8px", flexWrap: "wrap" }}>
          <Skeleton width={80} height={28} borderRadius={100} />
          <Skeleton width={70} height={28} borderRadius={100} />
          <Skeleton width={90} height={28} borderRadius={100} />
          <Skeleton width={65} height={28} borderRadius={100} />
        </div>
      </div>
    </div>
  );
}
