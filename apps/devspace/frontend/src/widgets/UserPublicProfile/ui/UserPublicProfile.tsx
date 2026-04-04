import { type JSX } from "react";

import type { IUserResponse } from "@/entities/user";
import { isValidMainRole } from "@/entities/user/model/IUserResponse";
import { Badge, Skeleton, UserAvatar } from "@/shared/ui";

import styles from "./UserPublicProfile.module.scss";

export interface UserPublicProfileProps {
  user: IUserResponse;
}

export function UserPublicProfile({ user }: UserPublicProfileProps): JSX.Element {
  const mainRoleName = isValidMainRole(user.main_role) ? user.main_role.name : undefined;

  return (
    <div className={styles.wrapper}>
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
