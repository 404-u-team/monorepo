import type { JSX } from "react";
import styles from "./UserCard.module.scss";

interface UserCardProps {
  avatar_uri?: string;
  user_id?: string;
  mainRole?: string;
  description?: string;
  skill_id?: string[];
  onProfileButtonClick?: () => void;
  onInviteButtonClick?: () => void;
}
export function UserCard({
  avatar_uri,
  user_id,
  mainRole,
  description,
  skill_id = [],
  //onProfileButtonClick,
  //onInviteButtonClick
}: UserCardProps): JSX.Element {
  return (
    <div className={styles.userCard}>
      <div className={styles.cardHeader}>
        <div className={styles.avatar_uri}>
          {avatar_uri}
          <img src={avatar_uri} alt="avatar" />
        </div>
        <div className={styles.textInfo}>
          <div className={styles.user_id}>{user_id}</div>
          <div className={styles.mainRole}>{mainRole}</div>
        </div>
      </div>
      <div className={styles.description}>{description}</div>
      <div className={styles.skillBox}>
        {skill_id.map((uuid) => (
          <div key={uuid} className={styles.skillName}>
            {uuid}
          </div>
        ))}
      </div>
      <div className={styles.profileButtonsBox}>
        <button className={styles.profileButton}>Профиль</button>
        <button className={styles.inviteButton}>Пригласить</button>
      </div>
    </div>
  );
}
export default UserCard;
