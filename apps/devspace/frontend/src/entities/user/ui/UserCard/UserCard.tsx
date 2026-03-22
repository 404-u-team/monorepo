import { useEffect, useRef, useState, type JSX } from "react";
import { clsx } from "clsx";
import { User } from "lucide-react";
import { fetchUserById } from "../../api/userApi";
import type { IUserResponse } from "../../model/IUserResponse";
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
    return (): void => { observer.disconnect(); };
  }, [user?.skills, needsScroll]);

  if (isLoading || user === undefined) {
    return <UserCardSkeleton className={className} />;
  }

  const Wrapper = to !== undefined ? "a" : "article";
  const wrapperProps = to !== undefined ? { href: to } : {};

  const skills = user.skills;

  return (
    <Wrapper
      {...wrapperProps}
      className={clsx(styles.card, to !== undefined && styles.link, className)}
    >
      <div className={styles.header}>
        <div className={styles.avatarWrapper}>
          {user.avatar_uri ? (
            <img
              className={styles.avatar}
              src={user.avatar_uri}
              alt={user.nickname}
              onError={(event) => {
                (event.target as HTMLImageElement).style.display = "none";
              }}
            />
          ) : (
            <div className={styles.avatarPlaceholder}>
              <User size={22} />
            </div>
          )}
        </div>

        <div className={styles.textInfo}>
          <span className={styles.nickname}>{user.nickname}</span>
          {user.main_role !== '' && (
            <span className={styles.mainRole}>{user.main_role}</span>
          )}
        </div>
      </div>

      {user.bio !== '' && (
        <p className={styles.bio}>{user.bio}</p>
      )}

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
          <a href={`/users/${id}`} className={styles.profileButton}>
            Профиль
          </a>
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
