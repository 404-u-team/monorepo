import type { JSX } from "react";
import styles from "./UserCard.module.scss";
import { useEffect, useRef, useState } from "react";
import InviteButton from "@/entities/user/ui/UserCard/InviteButton";
import { apiClient } from "@/shared/api/client";

interface UserCardProps {
  user_id: string;
  project_id?: string;
  slot_id?: string;
  inviteButton?: React.ReactNode;
}
interface UserData {
  avatar_uri: string;
  mainRole: string;
  description: string;
  skill_id: string[];
}

export function UserCard({
  user_id,
  project_id,
  slot_id,
}: UserCardProps): JSX.Element {
  const skillBoxReference = useRef<HTMLDivElement>(null);
  const scrollReference = useRef<HTMLDivElement>(null);
  const [needsScroll, setNeedsScroll] = useState(false);
  const [userData, setUserData] = useState<UserData | undefined>(undefined);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | undefined>(undefined);

  useEffect(() => {
    const fetchUserData = async (): Promise<void> => {
      try {
        setLoading(true);
        const response = await apiClient.get<UserData>(`/users/${user_id}`);
        setUserData(response.data);
        setError(undefined);
      } catch (error_) {
        setError("Ошибка загрузки данных пользователя");
        console.error("Error fetching user data:", error_);
      } finally {
        setLoading(false);
      }
    };

    void fetchUserData();
  }, [user_id]);

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
  }, [userData?.skill_id, needsScroll]);

  if (loading) {
    return (
      <div className={styles.userCard}>
        <div className={styles.loading}>Загрузка...</div>
      </div>
    );
  }

  if (error !== undefined) {
    return (
      <div className={styles.userCard}>
        <div className={styles.error}>{error}</div>
      </div>
    );
  }

  if (userData === undefined) {
    return (
      <div className={styles.userCard}>
        <div className={styles.error}>Данные не найдены</div>
      </div>
    );
  }

  const { avatar_uri, mainRole, description, skill_id } = userData;

  return (
    <div className={styles.userCard}>
      <div className={styles.cardHeader}>
        <div className={styles.avatarUri}>
          <img src={avatar_uri} alt="avatar" />
        </div>
        <div className={styles.textInfo}>
          <div className={styles.userId}>{user_id}</div>
          <div className={styles.mainRole}>{mainRole}</div>
        </div>
      </div>
      <div className={styles.description}>{description}</div>

      <div
        ref={skillBoxReference}
        className={`${styles.skillBox ?? ""} ${needsScroll ? (styles.skillBoxWithMask ?? "") : ""}`}
      >
        <div
          ref={scrollReference}
          className={`${styles.scroll ?? ""} ${needsScroll ? (styles.scrollAnimated ?? "") : ""}`}
        >
          {needsScroll
            ? [...skill_id, ...skill_id].map((uuid, index) => (
                <div
                  key={`${uuid}-${String(index)}`}
                  className={styles.skillName}
                >
                  {uuid}
                </div>
              ))
            : [...skill_id].map((uuid, index) => (
                <div
                  key={`${uuid}-${String(index)}`}
                  className={styles.skillName}
                >
                  {uuid}
                </div>
              ))}
        </div>
      </div>

      <div className={styles.profileButtonsBox}>
        <button className={styles.profileButton}>Профиль</button>
        <InviteButton project_id={project_id} slot_id={slot_id} />
      </div>
    </div>
  );
}

export default UserCard;
