import preview from "@/../.storybook/preview";
import UserCard from "./UserCard";

const meta = preview.meta({
  component: UserCard,
});

export const Default = meta.story({
  args: {
    user_id: "1\n",
  },
});
