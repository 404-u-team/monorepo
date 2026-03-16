import { useEffect, useRef, useState, type JSX } from "react";
import { clsx } from "clsx";
import { Button } from "@/shared/ui";
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
  onInvite?: (id: string) => Promise<void>;
}

export function UserCard({
  id,
  to,
  className,
  project_id,
  slot_id,
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
        // handled by future error state
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
      if (!box || !inner) return;
      const contentWidth = inner.scrollWidth / (needsScroll ? 2 : 1);
      const containerWidth = box.offsetWidth;

      setNeedsScroll(contentWidth > containerWidth);
    };

    checkOverflow();

    const ro = new ResizeObserver(checkOverflow);
    if (skillBoxReference.current) ro.observe(skillBoxReference.current);

    return (): void => {
      ro.disconnect();
    };
  }, [user?.skills, needsScroll]);

  if (isLoading || user === undefined) {
    return <UserCardSkeleton className={className} />;
  }

  const Wrapper = to !== undefined ? "a" : "article";
  const wrapperProps = to !== undefined ? { href: to } : {};

  return (
    <Wrapper
      {...wrapperProps}
      className={clsx(styles.card, to !== undefined && styles.link, className)}
    >
      <div className={styles.header}>
        <div className={styles.avatarWrapper}>
          <img
            className={styles.avatar}
            src={user.avatar_uri}
            alt={user.nickname}
            onError={(event) => {
              (event.target as HTMLImageElement).style.display = "none";
            }}
          />
        </div>

        <div className={styles.textInfo}>
          <div className={styles.nickname}>{user.nickname}</div>
          <span className={styles.mainRole}>{user.main_role}</span>
        </div>
      </div>

      <div className={styles.bio}>{user.bio}</div>

      <div
        ref={skillBoxReference}
        className={clsx(
          styles.skillBox,
          needsScroll && styles.skillBoxWithMask,
        )}
      >
        <div
          ref={scrollReference}
          className={clsx(styles.scroll, needsScroll && styles.scrollAnimated)}
        >
          {(needsScroll ? [...user.skills, ...user.skills] : user.skills).map(
            (skill, index) => {
              const skillName = typeof skill === "string" ? skill : skill.name;
              const uniqueKey = `${typeof skill === "string" ? skill : JSON.stringify(skill)}-${String(index)}`;

              return (
                <div key={uniqueKey} className={styles.skillName}>
                  {skillName}
                </div>
              );
            },
          )}
        </div>
      </div>

      <div className={styles.profileButtonsBox}>
        <Button variant="outline" className={styles.profileButton}>
          Профиль
        </Button>
        {project_id !== undefined && slot_id !== undefined && (
          <InviteButton project_id={project_id} slot_id={slot_id} />
        )}
      </div>
    </Wrapper>
  );
}

export default UserCard;
