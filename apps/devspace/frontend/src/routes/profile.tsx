import { createFileRoute } from "@tanstack/react-router";
import { apiClient } from "@/shared/api/client";
import { ProfileForm } from "@/features/profile";

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
        const response = await apiClient.get("/users/me");

        if (response.data !== undefined) {
          throw new Error("No data received from server");
        }

        const data = response.data as { id?: string };

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
    return <ProfileForm />;
  }
  return undefined;
}
