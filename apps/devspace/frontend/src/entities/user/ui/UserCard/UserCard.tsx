import { Link } from "@tanstack/react-router";
import { clsx } from "clsx";
import { useEffect, useRef, useState, type JSX } from "react";

import { UserAvatar } from "@/shared/ui";

import { fetchUserById } from "../../api/userApi";
// main_role is now returned as an object from the API, no need to fetch separately
import type { IUserResponse } from "../../model/IUserResponse";
import { isValidMainRole } from "../../model/IUserResponse";
import { UserCardSkeleton } from "../UserCardSkeleton/UserCardSkeleton";
import InviteButton from "./InviteButton";

import styles from "./UserCard.module.scss";

export interface UserCardProps {
  id: string;
  to?: string | undefined;
  className?: string | undefined;
  project_id?: string | undefined;
  slot_id?: string | undefined;
  onInvite?: (userId: string) => Promise<void>;
}

export function UserCard({
  id,
  to,
  className,
  project_id,
  slot_id,
  onInvite,
}: UserCardProps): JSX.Element {
  const [user, setUser] = useState<IUserResponse | undefined>(undefined);
  const [mainRoleName, setMainRoleName] = useState<string | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);
  const skillBoxReference = useRef<HTMLDivElement>(null);
  const scrollReference = useRef<HTMLDivElement>(null);
  const [needsScroll, setNeedsScroll] = useState(false);

  useEffect(() => {
    let cancelled = false;

    async function load(): Promise<void> {
      try {
        const data = await fetchUserById(id);
        if (cancelled) return;
        setUser(data);

        if (isValidMainRole(data.main_role)) {
          setMainRoleName(data.main_role.name);
        }
      } catch {
        // handled by empty state
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    void load();
    return (): void => {
      cancelled = true;
    };
  }, [id]);

  useEffect(() => {
    const checkOverflow = (): void => {
      const box = skillBoxReference.current;
      const inner = scrollReference.current;
      if (box === null || inner === null) return;
      const contentWidth = inner.scrollWidth / (needsScroll ? 2 : 1);
      setNeedsScroll(contentWidth > box.offsetWidth);
    };

    checkOverflow();

    const observer = new ResizeObserver(checkOverflow);
    if (skillBoxReference.current !== null) observer.observe(skillBoxReference.current);
    return (): void => {
      observer.disconnect();
    };
  }, [user?.skills, needsScroll]);

  if (isLoading || user === undefined) {
    return <UserCardSkeleton className={className} />;
  }

  const Wrapper = to !== undefined ? Link : "article";
  const wrapperProps = to !== undefined ? { to } : {};

  const skills = user.skills;

  return (
    <Wrapper
      {...wrapperProps}
      className={clsx(styles.card, to !== undefined && styles.link, className)}
    >
      <div className={styles.header}>
        <UserAvatar
          avatarUrl={user.avatar_url}
          nickname={user.nickname}
          size={44}
          className={styles.avatarWrapper}
        />

        <div className={styles.textInfo}>
          <span className={styles.nickname}>{user.nickname}</span>
          {mainRoleName !== undefined && <span className={styles.mainRole}>{mainRoleName}</span>}
        </div>
      </div>

      {user.bio !== "" && <p className={styles.bio}>{user.bio}</p>}

      {skills.length > 0 && (
        <div
          ref={skillBoxReference}
          className={clsx(styles.skillBox, needsScroll && styles.skillBoxWithMask)}
        >
          <div
            ref={scrollReference}
            className={clsx(styles.scroll, needsScroll && styles.scrollAnimated)}
          >
            {(needsScroll ? [...skills, ...skills] : skills).map((skill, index) => {
              const skillName = typeof skill === "string" ? skill : skill.name;
              const uniqueKey = `${typeof skill === "string" ? skill : skill.id}-${String(index)}`;
              return (
                <span key={uniqueKey} className={styles.skillName}>
                  {skillName}
                </span>
              );
            })}
          </div>
        </div>
      )}

      <div className={styles.profileButtonsBox}>
        {to === undefined && (
          <Link to="/users/$userId" params={{ userId: id }} className={styles.profileButton}>
            Профиль
          </Link>
        )}
        {project_id !== undefined && slot_id !== undefined && (
          <InviteButton
            project_id={project_id}
            slot_id={slot_id}
            user_id={id}
            {...(onInvite !== undefined ? { onInvite } : {})}
          />
        )}
      </div>
    </Wrapper>
  );
}

export default UserCard;
