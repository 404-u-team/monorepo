import type { Meta, StoryObj } from '@storybook/react';
import { DataListLayout } from './DataListLayout';
import { Badge } from '../../ui';

const meta = {
    title: 'shared/DataListLayout',
    component: DataListLayout,
    parameters: {
        layout: 'padded',
    },
    tags: ['autodocs'],
    argTypes: {
        onSearchChange: { action: 'onSearchChange' },
        onPageChange: { action: 'onPageChange' },
    },
} satisfies Meta<typeof DataListLayout>;

export default meta;
type Story = StoryObj<typeof meta>;

// Вспомогательный компонент для карточек-заглушек
const MockCard = ({ title, desc }: { title: string, desc: string }) => (
    <div style={{ border: '1px solid var(--border--main)', borderRadius: '12px', padding: '16px', background: 'var(--bg--main)' }}>
        <h3 style={{ margin: '0 0 8px 0', fontSize: '18px' }}>{title}</h3>
        <p style={{ margin: '0 0 16px 0', color: 'var(--text--secondary)', fontSize: '14px' }}>{desc}</p>
        <Badge>Mock Category</Badge>
    </div>
);

const MOCK_CARDS = Array.from({ length: 9 }).map((_, i) => (
    <MockCard key={i} title={`Element ${i + 1}`} desc={`Brief description for element number ${i + 1}`} />
));

// С имитацией Dropdown в controlsNode
const MockDropdown = () => (
    <select style={{ padding: '8px 16px', borderRadius: '8px', border: '1px solid var(--border--main)', background: 'var(--bg--main)', cursor: 'pointer' }}>
        <option>All items</option>
        <option>Open</option>
        <option>Closed</option>
    </select>
);

// С имитацией фильтра в боковой панели
const MockSidebarFilters = () => (
    <div style={{ padding: '24px', background: 'var(--bg--surface)', borderRadius: '12px' }}>
        <h4 style={{ margin: '0 0 16px', fontSize: '18px' }}>Filters</h4>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
            <label style={{ display: 'flex', gap: '8px', alignItems: 'center', cursor: 'pointer' }}>
                <input type="checkbox" /> Category A
            </label>
            <label style={{ display: 'flex', gap: '8px', alignItems: 'center', cursor: 'pointer' }}>
                <input type="checkbox" /> Category B
            </label>
            <label style={{ display: 'flex', gap: '8px', alignItems: 'center', cursor: 'pointer' }}>
                <input type="checkbox" /> Category C
            </label>
        </div>
    </div>
);

export const Default: Story = {
    args: {
        title: 'List Title',
        subtitle: 'This is a subtitle describing what this list contains.',
        searchValue: '',
        currentPage: 1,
        totalPages: 5,
        children: MOCK_CARDS,
    },
};

export const WithControlsNextToSearch: Story = {
    args: {
        ...Default.args,
        controlsNode: <MockDropdown />,
    },
};

export const WithSidebarFilters: Story = {
    args: {
        ...Default.args,
        filtersNode: <MockSidebarFilters />,
    },
};

export const EmptyState: Story = {
    args: {
        title: 'Empty List',
        searchValue: 'blablabla',
        isEmpty: true,
        emptyMessage: 'No items matching your search criteria were found.',
        children: [],
    },
};
