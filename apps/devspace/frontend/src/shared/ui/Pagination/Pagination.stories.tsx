import type { Meta, StoryObj } from '@storybook/react';
import { Pagination } from './Pagination';

const meta = {
    title: 'shared/Pagination',
    component: Pagination,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    argTypes: {
        onPageChange: { action: 'onPageChange' },
    },
} satisfies Meta<typeof Pagination>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        currentPage: 1,
        totalPages: 10,
        onPageChange: () => {},
    },
};

export const MiddlePage: Story = {
    args: {
        currentPage: 5,
        totalPages: 10,
        onPageChange: () => {},
    },
};

export const LastPage: Story = {
    args: {
        currentPage: 10,
        totalPages: 10,
        onPageChange: () => {},
    },
};

export const ManyPages: Story = {
    args: {
        currentPage: 10,
        totalPages: 100,
        onPageChange: () => {},
    },
};
