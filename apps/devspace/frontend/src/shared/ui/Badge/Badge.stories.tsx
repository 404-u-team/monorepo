import type { Meta, StoryObj } from '@storybook/react';
import { Badge } from './Badge';

const meta = {
    title: 'shared/Badge',
    component: Badge,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
} satisfies Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        children: 'AI & ML',
    },
};

export const LongText: Story = {
    args: {
        children: 'Веб-разработка',
    },
};
