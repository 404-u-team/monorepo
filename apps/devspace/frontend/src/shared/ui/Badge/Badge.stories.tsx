import type { Meta, StoryObj } from "@storybook/react";
import { Zap, Clock, Code, Layout, Brain } from "lucide-react";

import { Badge } from "./Badge";

const meta = {
  title: "shared/Badge",
  component: Badge,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: "Default Badge",
  },
};

export const Green: Story = {
  args: {
    children: "Success",
    color: "34D399",
  },
};

export const Blue: Story = {
  args: {
    children: "Info",
    color: "3B82F6",
  },
};

export const Red: Story = {
  args: {
    children: "Error",
    color: "EF4444",
  },
};

export const Orange: Story = {
  args: {
    children: "Warning",
    color: "F59E0B",
  },
};

export const WithIcon: Story = {
  args: {
    children: "AI & ML",
    icon: <Brain size={14} />,
  },
};

export const WithIconAndColor: Story = {
  args: {
    children: "Backend",
    color: "EF4444",
    icon: <Code size={14} />,
  },
};

export const TooltipLike: Story = {
  args: {
    children: "Design",
    color: "8B5CF6",
    icon: <Layout size={14} />,
  },
};

export const LongText: Story = {
  args: {
    children: "Fullstack Development",
    color: "F472B6",
    icon: <Zap size={14} />,
  },
};

export const Timed: Story = {
  args: {
    children: "2 hours left",
    color: "64748B",
    icon: <Clock size={14} />,
  },
};
