import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import { useEffect } from "react";

import { useUserStore } from "@/entities/user";
import { ProfileForm } from "@/features/profile";

export const Route = createFileRoute("/profile")({
  component: ProfilePage,
});

function ProfilePage(): JSX.Element {
  const userStore = useUserStore();
  const navigate = useNavigate();

  useEffect(() => {
    if (!userStore.isAuthLoading && !userStore.isAuthenticated) {
      void navigate({ to: "/auth" });
    }
  }, [userStore.isAuthLoading, userStore.isAuthenticated, navigate]);

  return <ProfileForm />;
}
