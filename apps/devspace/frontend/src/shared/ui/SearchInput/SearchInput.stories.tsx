import type { Meta, StoryObj } from '@storybook/react';
import { SearchInput } from './SearchInput';

const meta = {
    title: 'shared/SearchInput',
    component: SearchInput,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    argTypes: {
        onSearch: { action: 'onSearch' },
    },
} satisfies Meta<typeof SearchInput>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        placeholder: 'Поиск проектов...',
        onSearch: () => {},
    },
};

export const WithValue: Story = {
    args: {
        value: 'React',
        placeholder: 'Поиск проектов...',
        onSearch: () => {},
    },
};
