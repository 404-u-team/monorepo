import type { JSX } from "react";
import styles from "./UserCard.module.scss";
import { useEffect, useRef, useState } from "react";

interface UserCardProps {
  avatar_uri?: string;
  user_id: string;
  mainRole?: string;
  description?: string;
  skill_id?: string[];
  onInviteButtonClick?: () => void;
}

export function UserCard({
  avatar_uri,
  user_id,
  mainRole,
  description,
  skill_id = [],
  //onInviteButtonClick
}: UserCardProps): JSX.Element {
  const skillBoxRef = useRef<HTMLDivElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [needsScroll, setNeedsScroll] = useState(false);

  useEffect(() => {
    const checkOverflow = () => {
      const box = skillBoxRef.current;
      const inner = scrollRef.current;
      if (!box || !inner) return;
      const contentWidth = inner.scrollWidth / 2;
      const containerWidth = box.offsetWidth;

      setNeedsScroll(contentWidth > containerWidth);
    };

    checkOverflow();

    const ro = new ResizeObserver(checkOverflow);
    if (skillBoxRef.current) ro.observe(skillBoxRef.current);

    return () => ro.disconnect();
  }, [skill_id]);

  return (
    <div className={styles.userCard}>
      <div className={styles.cardHeader}>
        <div className={styles.avatar_uri}>
          <img src={avatar_uri} alt="avatar" />
        </div>
        <div className={styles.textInfo}>
          <div className={styles.user_id}>{user_id}</div>
          <div className={styles.mainRole}>{mainRole}</div>
        </div>
      </div>
      <div className={styles.description}>{description}</div>

      <div ref={skillBoxRef} className={styles.skillBox}>
        <div
          ref={scrollRef}
          className={`${styles.Scroll} ${needsScroll ? styles.ScrollAnimated : ""}`}
        >
          {[...skill_id, ...skill_id].map((uuid, index) => (
            <div key={`${uuid}-${index}`} className={styles.skillName}>
              {uuid}
            </div>
          ))}
        </div>
      </div>

      <div className={styles.profileButtonsBox}>
        <button className={styles.profileButton}>Профиль</button>
        <button className={styles.inviteButton}>Пригласить</button>
      </div>
    </div>
  );
}

export default UserCard;
