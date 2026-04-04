import type { Meta, StoryObj } from "@storybook/react";

import { Skeleton } from "./Skeleton";

const meta = {
  title: "shared/Skeleton",
  component: Skeleton,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof Skeleton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    width: 200,
    height: 20,
  },
};

export const Circle: Story = {
  args: {
    width: 40,
    height: 40,
    borderRadius: "50%",
  },
};

export const Card: Story = {
  args: {
    width: 300,
    height: 180,
    borderRadius: 12,
  },
};
