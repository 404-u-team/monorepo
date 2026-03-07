import type { JSX } from "react";
import styles from "./UserCard.module.scss";
import { useEffect, useRef, useState } from "react";
import InviteButton from "@/entities/user/ui/UserCard/InviteButton";
// import { apiClient } from "../api-client";
// import type {
//   PrivateUserProfile,
//   PublicUserProfile,
//   UUID,
// } from "@/types/api.types";

interface UserCardProps {
  user_id: string;
  //description?: string;
  //skill_id?: string[];
  //project_id?: string;
  //slot_id?: string;
  inviteButton?: React.ReactNode;
}

export function UserCard({
  //avatar_uri,
  user_id,
  //mainRole,
  //description,
  //skill_id = [],
  //project_id,
  //slot_id,
}: UserCardProps): JSX.Element {
  const skillBoxReference = useRef<HTMLDivElement>(null);
  const scrollReference = useRef<HTMLDivElement>(null);
  const [needsScroll, setNeedsScroll] = useState(false);
  const project_id = "2";
  const slot_id = "sdad";
  const avatar_uri = "https://placehold.co/150x150";
  const mainRole = "Frontend Developer";
  const description =
    "Люблю React и TypeScript Lorem ipsum, dolor sit amet consectetur adipisicing elit. Minima consequatur cum doloribus asperiores exercitationem ipsum incidunt odit provident quia, delectus sequi ullam perspiciatis facilis eveniet cumque deleniti at laboriosam. Impedit!";
  const skill_id = ["react", "css"];

  // export const usersApi = {
  //   getMyProfile: (token: string) =>
  //     apiClient.request<PrivateUserProfile>("/users/me", {
  //       method: "GET",
  //       token,
  //     }),

  //   updateMyProfile: (
  //     token: string,
  //     data: {
  //       nickname?: string;
  //       avatar_uri?: string;
  //       bio?: string;
  //     },
  //   ) =>
  //     apiClient.request<PrivateUserProfile>("/users/me", {
  //       method: "PUT",
  //       token,
  //       body: JSON.stringify(data),
  //     }),

  //   addSkill: (token: string, skill_id: UUID) =>
  //     apiClient.request<void>("/users/me/skills", {
  //       method: "POST",
  //       token,
  //       body: JSON.stringify({ skill_id }),
  //     }),

  //   removeSkill: (token: string, skill_id: UUID) =>
  //     apiClient.request<void>(`/users/me/skills/${skill_id}`, {
  //       method: "DELETE",
  //       token,
  //     }),

  //   getUserProfile: (user_id: UUID) =>
  //     apiClient.request<PublicUserProfile>(`/users/${user_id}`, {
  //       method: "GET",
  //     }),
  // };

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

      <div
        ref={skillBoxReference}
        className={`${styles.skillBox ?? ""} ${needsScroll ? (styles.skillBoxWithMask ?? "") : ""}`}
      >
        <div
          ref={scrollReference}
          className={`${styles.Scroll ?? ""} ${needsScroll ? (styles.ScrollAnimated ?? "") : ""}`}
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
