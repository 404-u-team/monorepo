import type { Meta, StoryObj } from '@storybook/react';
import { MdRenderer } from './MdRenderer';

const meta: Meta<typeof MdRenderer> = {
    title: 'shared/ui/MdRenderer',
    component: MdRenderer,
    tags: ['autodocs'],
    parameters: {
        layout: 'centered',
    },
};

export default meta;
type Story = StoryObj<typeof MdRenderer>;

const sampleMarkdown = `
# Привет, мир!

Это пример **Markdown** контента.

### Список возможностей:
1. Жирный текст
2. *Курсив*
3. [Ссылка](https://google.com)
4. \`Код\`

> Цитата дня: Чем глубже в лес тем шкебеде доп ес ес.

\`\`\`typescript
const greeting = "Hello, world!";
console.log(greeting);
\`\`\`
`;

export const Default: Story = {
    args: {
        source: sampleMarkdown,
    },
};

export const LongContent: Story = {
    args: {
        source: sampleMarkdown.repeat(5),
    },
};

export const Empty: Story = {
    args: {
        source: '',
    },
};
