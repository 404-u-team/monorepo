import type { JSX } from "react";
import { createFileRoute } from "@tanstack/react-router";
import UserCard from "@/entities/user/ui/UserCard";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index(): JSX.Element {
  return (
    <div className="p-2">
      <h3>Welcome! This is the home page.</h3>
      <div className="mt-4">
        <UserCard />
        <UserCard
          avatar_uri=""
          user_id="@msOur"
          mainRole="Frontend"
          description="Lorem ipsum text text text text text text text
          Lorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text textLorem ipsum text text text text text text text"
          skill_id={[
            "JavaScript",
            "Vue.js",
            "React",
            "MobX",
            "Nuxt",
            "JavaScript",
            "Vue.js",
            "React",
            "MobX",
            "Nuxt",
          ]}
          //onProfileButtonClick={handleProfileClick}
          //onInviteButtonClick={handleInviteClick}
        />
      </div>
    </div>
  );
}
