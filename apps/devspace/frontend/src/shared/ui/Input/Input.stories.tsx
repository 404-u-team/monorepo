import type { Meta, StoryObj } from '@storybook/react';
import { Input } from './Input';
import { Mail, Lock, Eye } from 'lucide-react';

const meta = {
    title: 'shared/Input',
    component: Input,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    decorators: [
        (Story) => (
            <div style={{ width: '300px' }}>
                <Story />
            </div>
        ),
    ],
} satisfies Meta<typeof Input>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        placeholder: 'Enter text...',
    },
};

export const Email: Story = {
    args: {
        placeholder: 'your.email@example.com',
        type: 'email',
        iconLeft: <Mail size={20} />,
    },
};

export const Password: Story = {
    args: {
        placeholder: 'Password',
        type: 'password',
        iconLeft: <Lock size={20} />,
        iconRight: (
            <button type="button" style={{ background: 'none', border: 'none', cursor: 'pointer', padding: 0, color: 'inherit', display: 'flex' }}>
                <Eye size={20} />
            </button>
        )
    },
};

export const WithError: Story = {
    args: {
        placeholder: 'Enter text...',
        error: 'Required field',
    },
};
