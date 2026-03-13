import { useState } from "react";
import type { JSX } from "react";
import { Button } from "@/shared/ui/Button/Button";
import styles from "@/entities/user/ui/UserCard/UserCard.module.scss";

interface InviteButtonProps {
  project_id?: string | undefined;
  slot_id?: string | undefined;
  user_id?: string;
  onInvite?: (userId: string) => Promise<void>;
}

const InviteButton = ({
  project_id,
  slot_id,
  user_id,
  onInvite,
}: InviteButtonProps): JSX.Element | undefined => {
  const [isLoading, setIsLoading] = useState(false);

  const handleInvite = async (): Promise<void> => {
    if (project_id === undefined || slot_id === undefined) return;

    setIsLoading(true);
    try {
      if (onInvite) {
        await onInvite(String(user_id));
      } else {
        const response = await fetch(
          `/api/projects/${project_id}/slots/${slot_id}/invite`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
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

  return (
    <Button
      variant="primary"
      className={styles.inviteButton}
      onClick={() => {
        handleInvite().catch(console.error);
      }}
      disabled={isLoading}
    >
      {isLoading ? "..." : "Пригласить"}
    </Button>
  );
};

export default InviteButton;
