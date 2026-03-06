import preview from "@/../.storybook/preview";
import UserCard from "./UserCard";

const meta = preview.meta({
  component: UserCard,
});

export const Default = meta.story({
  args: {
    avatar_uri: "https://placehold.co/150x150",
    user_id: "john_doe",
    mainRole: "Frontend Developer",
    description:
      "Люблю React и TypeScript Lorem ipsum, dolor sit amet consectetur adipisicing elit. Minima consequatur cum doloribus asperiores exercitationem ipsum incidunt odit provident quia, delectus sequi ullam perspiciatis facilis eveniet cumque deleniti at laboriosam. Impedit!",
    skill_id: [
      "react",
      "css",
      "html",
      "react",
      "css",
      "html",
      "react",
      "css",
      "html",
    ],
  },
});
