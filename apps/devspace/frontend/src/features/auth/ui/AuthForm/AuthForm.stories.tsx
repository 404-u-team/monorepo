import type { Meta, StoryObj } from '@storybook/react';
import { AuthForm } from './AuthForm';

const meta = {
    title: 'features/auth/AuthForm',
    component: AuthForm,
    parameters: {
        layout: 'fullscreen',
    },
    decorators: [
        (Story) => (
            <div style={{ padding: '2rem 1rem', display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', backgroundColor: 'var(--bg--main)' }}>
                <Story />
            </div>
        ),
    ],
} satisfies Meta<typeof AuthForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
