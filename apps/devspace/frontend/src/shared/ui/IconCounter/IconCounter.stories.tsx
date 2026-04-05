import type { Meta, StoryObj } from "@storybook/react";
import { Heart, Eye } from "lucide-react";

import { IconCounter } from "./IconCounter";

const meta = {
  title: "shared/IconCounter",
  component: IconCounter,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    active: { control: "boolean" },
  },
} satisfies Meta<typeof IconCounter>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Favorites: Story = {
  args: {
    icon: <Heart size={16} />,
    count: 237,
  },
};

export const FavoritesActive: Story = {
  args: {
    icon: <Heart size={16} fill="currentColor" />,
    count: 238,
    active: true,
    onClick: () => {},
  },
};

export const Views: Story = {
  args: {
    icon: <Eye size={16} />,
    count: 2600,
  },
};

export const LargeCount: Story = {
  args: {
    icon: <Eye size={16} />,
    count: 1_500_000,
  },
};

export const Clickable: Story = {
  args: {
    icon: <Heart size={16} />,
    count: 42,
    onClick: () => {},
  },
};
