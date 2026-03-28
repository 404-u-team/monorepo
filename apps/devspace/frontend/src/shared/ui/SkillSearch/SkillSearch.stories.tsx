import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { SkillSearch, type SkillSearchOption } from './SkillSearch';

const MOCK_SKILLS: SkillSearchOption[] = [
    { id: '1', name: 'Frontend', color: '3B82F6' },
    { id: '2', name: 'Backend', color: 'EF4444' },
    { id: '3', name: 'DevOps', color: '7C3AED' },
    { id: '4', name: 'Mobile', color: '10B981' },
    { id: '5', name: 'Data Science', color: 'F59E0B' },
    { id: '6', name: 'Design', color: 'EC4899' },
    { id: '7', name: 'QA', color: '06B6D4' },
    { id: '8', name: 'Machine Learning', color: '8B5CF6' },
    { id: '9', name: 'Security', color: 'EF4444' },
    { id: '10', name: 'Blockchain', color: '64748B' },
];

const mockLoadOptions = async (query: string): Promise<SkillSearchOption[]> => {
    await new Promise((resolve) => { setTimeout(resolve, 300); });
    if (query === '') return MOCK_SKILLS;
    const lower = query.toLowerCase();
    return MOCK_SKILLS.filter((skill) => skill.name.toLowerCase().includes(lower));
};

const meta = {
    title: 'shared/SkillSearch',
    component: SkillSearch,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    decorators: [
        (Story) => (
            <div style={{ width: 360 }}>
                <Story />
            </div>
        ),
    ],
} satisfies Meta<typeof SkillSearch>;

export default meta;
type Story = StoryObj<typeof meta>;

/* ── Interactive stories ─────────────────────────────────────────────────── */

function InteractiveDefault(): React.JSX.Element {
    const [selected, setSelected] = useState<SkillSearchOption | undefined>(undefined);
    return (
        <SkillSearch
            value={selected}
            onChange={setSelected}
            loadOptions={mockLoadOptions}
            placeholder="Выберите основной навык..."
        />
    );
}

export const Default: Story = {
    args: {
        value: undefined,
        onChange: () => {},
        loadOptions: mockLoadOptions,
        placeholder: 'Выберите основной навык...',
    },
    render: () => <InteractiveDefault />,
};

function InteractiveWithValue(): React.JSX.Element {
    const [selected, setSelected] = useState<SkillSearchOption | undefined>(MOCK_SKILLS[0]);
    return (
        <SkillSearch
            value={selected}
            onChange={setSelected}
            loadOptions={mockLoadOptions}
            placeholder="Выберите основной навык..."
        />
    );
}

export const WithValue: Story = {
    args: {
        value: MOCK_SKILLS[0],
        onChange: () => {},
        loadOptions: mockLoadOptions,
    },
    render: () => <InteractiveWithValue />,
};

export const Disabled: Story = {
    args: {
        value: MOCK_SKILLS[2],
        onChange: () => {},
        loadOptions: mockLoadOptions,
        disabled: true,
    },
};
