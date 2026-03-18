import { createFileRoute } from "@tanstack/react-router";
import { ProfileForm } from "@/features/profile/ui/ProfileForm/ProfileForm";

import type { JSX } from "react";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/profile")({
  component: UsersMe,
});

function UsersMe(): JSX.Element | undefined {
  const [userId, setUserId] = useState<string>();
  useEffect(() => {
    let isMounted = true;
    const loadUser = async (): Promise<void> => {
      try {
        const response = await fetch("/users/me");
        if (!response.ok) {
          throw new Error(
            `Failed to fetch current user: ${String(response.status)}`,
          );
        }
        const data = (await response.json()) as { id?: string };
        if (isMounted && typeof data.id === "string") {
          setUserId(data.id);
        }
      } catch (error) {
        // eslint-disable-next-line no-console
        console.error("Error loading current user profile:", error);
      }
    };
    void loadUser();
    return (): void => {
      isMounted = false;
    };
  }, []);
  if (userId !== undefined) {
    return <ProfileForm id={userId} />;
  }
  return undefined;
}
