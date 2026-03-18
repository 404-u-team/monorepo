import { useState } from "react";
import type { JSX } from "react";
import { Button } from "@/shared/ui/Button/Button";
import styles from "@/entities/user/ui/UserCard/UserCard.module.scss";
import { inviteUserToSlot } from "@/entities/user/api/userApi";

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
  const isMissingRequiredIds =
    project_id === undefined || slot_id === undefined || user_id === undefined;
  const handleInvite = async (): Promise<void> => {
    if (isMissingRequiredIds) {
      return;
    }

    setIsLoading(true);
    try {
      if (onInvite) {
        await onInvite(user_id);
      } else {
        await inviteUserToSlot({
          project_id,
          slot_id,
          user_id: user_id,
        });
      }
    } catch (error) {
      console.error("Ошибка приглашения:", error);
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
      disabled={isLoading || isMissingRequiredIds}
    >
      {isLoading ? "..." : "Пригласить"}
    </Button>
  );
};

export default InviteButton;
