import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { SkillMultiSelect, type SkillMultiSelectOption } from './SkillMultiSelect';

const MOCK_SKILLS: SkillMultiSelectOption[] = [
    { id: '1', name: 'React', color: '61DAFB' },
    { id: '2', name: 'TypeScript', color: '3178C6' },
    { id: '3', name: 'MobX', color: 'EF6C00' },
    { id: '4', name: 'Node.js', color: '339933' },
    { id: '5', name: 'PostgreSQL', color: '336791' },
    { id: '6', name: 'Docker', color: '2496ED' },
    { id: '7', name: 'Kubernetes', color: '326CE5' },
    { id: '8', name: 'GraphQL', color: 'E10098' },
    { id: '9', name: 'Redis', color: 'DC382D' },
    { id: '10', name: 'AWS', color: 'FF9900' },
    { id: '11', name: 'Terraform', color: '7B42BC' },
    { id: '12', name: 'Python', color: '3776AB' },
];

const mockLoadOptions = async (query: string): Promise<SkillMultiSelectOption[]> => {
    await new Promise((resolve) => { setTimeout(resolve, 300); });
    if (query === '') return MOCK_SKILLS;
    const lower = query.toLowerCase();
    return MOCK_SKILLS.filter((skill) => skill.name.toLowerCase().includes(lower));
};

const meta = {
    title: 'shared/SkillMultiSelect',
    component: SkillMultiSelect,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    decorators: [
        (Story) => (
            <div style={{ width: 400 }}>
                <Story />
            </div>
        ),
    ],
} satisfies Meta<typeof SkillMultiSelect>;

export default meta;
type Story = StoryObj<typeof meta>;

/* ── Interactive stories ─────────────────────────────────────────────────── */

function InteractiveDefault(): React.JSX.Element {
    const [selected, setSelected] = useState<SkillMultiSelectOption[]>([]);
    return (
        <SkillMultiSelect
            value={selected}
            onChange={setSelected}
            loadOptions={mockLoadOptions}
            placeholder="Добавить навыки..."
            max={10}
        />
    );
}

export const Default: Story = {
    args: {
        value: [],
        onChange: () => {},
        loadOptions: mockLoadOptions,
        placeholder: 'Добавить навыки...',
        max: 10,
    },
    render: () => <InteractiveDefault />,
};

function InteractiveWithValues(): React.JSX.Element {
    const [selected, setSelected] = useState<SkillMultiSelectOption[]>(MOCK_SKILLS.slice(0, 3));
    return (
        <SkillMultiSelect
            value={selected}
            onChange={setSelected}
            loadOptions={mockLoadOptions}
            max={10}
        />
    );
}

export const WithValues: Story = {
    args: {
        value: MOCK_SKILLS.slice(0, 3),
        onChange: () => {},
        loadOptions: mockLoadOptions,
        max: 10,
    },
    render: () => <InteractiveWithValues />,
};

function InteractiveMaxReached(): React.JSX.Element {
    const [selected, setSelected] = useState<SkillMultiSelectOption[]>(MOCK_SKILLS.slice(0, 3));
    return (
        <SkillMultiSelect
            value={selected}
            onChange={setSelected}
            loadOptions={mockLoadOptions}
            max={3}
        />
    );
}

export const MaxReached: Story = {
    args: {
        value: MOCK_SKILLS.slice(0, 3),
        onChange: () => {},
        loadOptions: mockLoadOptions,
        max: 3,
    },
    render: () => <InteractiveMaxReached />,
};

export const Disabled: Story = {
    args: {
        value: MOCK_SKILLS.slice(0, 2),
        onChange: () => {},
        loadOptions: mockLoadOptions,
        disabled: true,
        max: 10,
    },
};
