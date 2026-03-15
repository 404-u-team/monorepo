import type { JSX } from "react";
import styles from "./UserCard.module.scss";
import { useEffect, useRef, useState, useMemo } from "react";
import InviteButton from "@/entities/user/ui/UserCard/InviteButton";
import { apiClient } from "@/shared/api/client";
import { Button } from "@/shared/ui";

interface UserCardProps {
  user_id: string;
  project_id?: string;
  slot_id?: string;
  inviteButton?: React.ReactNode;
}

interface UserData {
  avatar_uri: string;
  main_role: string;
  nickname: string;
  bio: string;
  skills: string[];
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
  }, [userData?.skills, needsScroll]);

  if (loading) {
    return (
      <div className={styles.userCard}>
        <div>Загрузка...</div>
      </div>
    );
  }

  if (error !== undefined) {
    return (
      <div className={styles.userCard}>
        <div>{error}</div>
      </div>
    );
  }

  if (userData === undefined) {
    return (
      <div className={styles.userCard}>
        <div>Данные не найдены</div>
      </div>
    );
  }

  const { avatar_uri, main_role, bio, skills, nickname } = userData;

  return (
    <div className={styles.userCard}>
      <div className={styles.cardHeader}>
        <div className={styles.avatarUri}>
          <img src={avatar_uri} alt="avatar" />
        </div>
        <div className={styles.textInfo}>
          <div className={styles.userId}>{nickname}</div>
          <div className={styles.mainRole}>{main_role}</div>
        </div>
      </div>
      <div className={styles.bio}>{bio}</div>

      <div
        ref={skillBoxReference}
        className={`${styles.skillBox ?? ""} ${needsScroll ? (styles.skillBoxWithMask ?? "") : ""}`}
      >
        <div
          ref={scrollReference}
          className={`${styles.scroll ?? ""} ${needsScroll ? (styles.scrollAnimated ?? "") : ""}`}
        >
          {needsScroll
            ? [...skills, ...skills].map((uuid, index) => (
                <div
                  key={`${uuid}-${String(index)}`}
                  className={styles.skillName}
                >
                  {uuid}
                </div>
              ))
            : [...skills].map((uuid, index) => (
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
        <Button variant="outline" className={styles.profileButton}>
          Профиль
        </Button>
        {project_id !== undefined && slot_id !== undefined && (
          <InviteButton project_id={project_id} slot_id={slot_id} />
        )}
      </div>
    </div>
  );
}

export default UserCard;
