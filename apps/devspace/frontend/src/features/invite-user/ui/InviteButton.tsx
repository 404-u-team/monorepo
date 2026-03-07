import { useState } from "react";
import styles from "@/entities/user/ui/UserCard/UserCard.module.scss";

interface InviteButtonProps {
  project_id?: string;
  slot_id?: string;
  user_id: string;
  onInvite?: (userId: string) => Promise<void>;
}

const InviteButton = ({
  project_id,
  slot_id,
  user_id,
  onInvite,
}: InviteButtonProps) => {
  const [isLoading, setIsLoading] = useState(false);

  const handleInvite = async () => {
    if (!project_id || !slot_id) return;

    setIsLoading(true);
    try {
      if (onInvite) {
        await onInvite(user_id);
      } else {
        const token = localStorage.getItem("authToken");
        const response = await fetch(
          `/api/projects/${project_id}/slots/${slot_id}/invite`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({
              user_id: user_id,
            }),
          },
        );

        if (!response.ok) throw new Error("Ошибка отправки");
      }
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return project_id && slot_id ? (
    <button
      className={styles.inviteButton} // Используем стили из UserCard
      onClick={handleInvite}
      disabled={isLoading}
    >
      {isLoading ? "..." : "Пригласить"}
    </button>
  ) : null;
};

export default InviteButton;
