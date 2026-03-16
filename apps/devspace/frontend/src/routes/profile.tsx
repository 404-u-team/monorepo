import { createFileRoute } from "@tanstack/react-router";
import { ProfileForm } from "@/features/profile/ui/ProfileForm/ProfileForm";

import type { JSX } from "react";

export const Route = createFileRoute("/profile")({
  component: UsersMe,
});

function UsersMe(): JSX.Element {
  return <ProfileForm id="1" />;
}
