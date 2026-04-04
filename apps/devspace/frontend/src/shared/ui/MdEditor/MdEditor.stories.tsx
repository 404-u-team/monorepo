import type { Meta, StoryObj } from "@storybook/react";

import { MdEditor } from "./MdEditor";

const meta: Meta<typeof MdEditor> = {
  title: "shared/ui/MdEditor",
  component: MdEditor,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
  },
  argTypes: {
    onChange: { action: "changed" },
  },
};

export default meta;
type Story = StoryObj<typeof MdEditor>;

export const Default: Story = {
  args: {
    value: "",
    placeholder: "Начните вводить текст...",
  },
};

export const WithContent: Story = {
  args: {
    value: "# Заголовок\n\nЭто текст в редакторе.",
    height: 400,
  },
};

export const Disabled: Story = {
  args: {
    value: "Редактирование запрещено.",
    disabled: true,
  },
};

export const Placeholder: Story = {
  args: {
    value: "",
    placeholder: "Введите подробное описание вашей идеи...",
    height: 200,
  },
};
