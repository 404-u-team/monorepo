import type { Meta, StoryObj } from "@storybook/react";

import { Dropdown } from "./Dropdown";

const meta = {
  title: "shared/Dropdown",
  component: Dropdown,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {
    onChange: { action: "onChange" },
  },
} satisfies Meta<typeof Dropdown>;

export default meta;
type Story = StoryObj<typeof meta>;

const mockOptions = [
  { label: "Все", value: "" },
  { label: "Открытые", value: "open" },
  { label: "Закрытые", value: "closed" },
];

export const Default: Story = {
  args: {
    options: mockOptions,
    onChange: () => {},
    placeholder: "Выберите статус...",
  },
};

export const WithValue: Story = {
  args: {
    options: mockOptions,
    value: "open",
    onChange: () => {},
  },
};

export const Disabled: Story = {
  args: {
    options: mockOptions,
    disabled: true,
    placeholder: "Выберите статус...",
    onChange: () => {},
  },
};
