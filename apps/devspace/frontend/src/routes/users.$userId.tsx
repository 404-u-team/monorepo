import { createFileRoute, useParams } from "@tanstack/react-router";
import type { JSX } from "react";
import { useEffect, useState } from "react";

import { fetchUserById, type IUserResponse } from "@/entities/user";
import { UserPublicProfile, UserPublicProfileSkeleton } from "@/widgets/UserPublicProfile";

export const Route = createFileRoute("/users/$userId")({
  component: UserPublicProfilePage,
});

function UserPublicProfilePage(): JSX.Element {
  const routeParameters = useParams({ strict: false });
  const userId = (routeParameters as Record<string, string>).userId ?? "";
  const [user, setUser] = useState<IUserResponse | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    let cancelled = false;

    async function load(): Promise<void> {
      try {
        const data = await fetchUserById(userId);
        if (!cancelled) setUser(data);
      } catch {
        // user not found - keep undefined
      } finally {
        if (!cancelled) setIsLoading(false);
      }
    }

    void load();
    return (): void => {
      cancelled = true;
    };
  }, [userId]);

  return (
    <div style={{ padding: "40px 24px", maxWidth: "1200px", margin: "0 auto" }}>
      {isLoading && <UserPublicProfileSkeleton />}
      {!isLoading && user !== undefined && <UserPublicProfile user={user} />}
      {!isLoading && user === undefined && (
        <p style={{ color: "var(--text--secondary)" }}>Пользователь не найден</p>
      )}
    </div>
  );
}
