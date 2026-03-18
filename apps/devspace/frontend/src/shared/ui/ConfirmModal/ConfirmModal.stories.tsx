import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { ConfirmModal } from './ConfirmModal';
import { Button } from '../Button/Button';

const meta = {
    title: 'shared/ConfirmModal',
    component: ConfirmModal,
    parameters: {
        layout: 'centered',
    },
    tags: ['autodocs'],
    argTypes: {
        severity: {
            control: 'select',
            options: ['info', 'warning', 'danger'],
        },
        isLoading: {
            control: 'boolean',
        },
        onConfirm: { action: 'onConfirm' },
        onCancel: { action: 'onCancel' },
    },
} satisfies Meta<typeof ConfirmModal>;

export default meta;
type Story = StoryObj<typeof meta>;

/* ── Статические истории ─────────────────────────────────────────────────── */

export const Info: Story = {
    args: {
        isOpen: true,
        title: 'Подтвердите действие',
        description: 'Это информационное сообщение. Вы уверены, что хотите продолжить?',
        severity: 'info',
        confirmLabel: 'Продолжить',
        cancelLabel: 'Отмена',
        onConfirm: () => {},
        onCancel: () => {},
    },
};

export const Warning: Story = {
    args: {
        isOpen: true,
        title: 'Будьте осторожны',
        description: 'Это действие может повлиять на другие элементы и его нельзя будет легко отменить.',
        severity: 'warning',
        confirmLabel: 'Всё равно продолжить',
        cancelLabel: 'Отмена',
        onConfirm: () => {},
        onCancel: () => {},
    },
};

export const Danger: Story = {
    args: {
        isOpen: true,
        title: 'Удалить проект?',
        description: 'Это действие необратимо. Все данные проекта, включая слоты и заявки, будут удалены навсегда.',
        severity: 'danger',
        confirmLabel: 'Удалить',
        cancelLabel: 'Отмена',
        onConfirm: () => {},
        onCancel: () => {},
    },
};

export const DangerDeleteIdea: Story = {
    args: {
        isOpen: true,
        title: 'Удалить идею?',
        description: 'Вы уверены, что хотите удалить эту идею? Действие нельзя отменить.',
        severity: 'danger',
        confirmLabel: 'Удалить',
        cancelLabel: 'Отмена',
        onConfirm: () => {},
        onCancel: () => {},
    },
};

export const DangerDeleteSlot: Story = {
    args: {
        isOpen: true,
        title: 'Удалить слот?',
        description: 'Слот и все связанные заявки будут удалены.',
        severity: 'warning',
        confirmLabel: 'Удалить слот',
        cancelLabel: 'Отмена',
        onConfirm: () => {},
        onCancel: () => {},
    },
};

export const Loading: Story = {
    args: {
        isOpen: true,
        title: 'Удалить проект?',
        description: 'Это действие необратимо.',
        severity: 'danger',
        confirmLabel: 'Удалить',
        cancelLabel: 'Отмена',
        isLoading: true,
        onConfirm: () => {},
        onCancel: () => {},
    },
};

/* ── Интерактивные истории с кнопкой-триггером ───────────────────────────── */

function InteractiveDemo({ severity }: { severity: 'info' | 'warning' | 'danger' }): React.JSX.Element {
    const [isOpen, setIsOpen] = useState(false);

    const labelMap = {
        info: 'Открыть Info-диалог',
        warning: 'Открыть Warning-диалог',
        danger: 'Открыть Danger-диалог',
    };

    const titleMap = {
        info: 'Подтвердите действие',
        warning: 'Опасное изменение',
        danger: 'Удалить безвозвратно?',
    };

    const descMap = {
        info: 'Это просто подтверждение. Всё под контролем.',
        warning: 'Это действие сложно отменить. Убедитесь, что вы точно хотите продолжить.',
        danger: 'Все данные будут безвозвратно удалены. Восстановление невозможно.',
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '1rem' }}>
            <Button onClick={() => { setIsOpen(true); }}>{labelMap[severity]}</Button>
            <ConfirmModal
                isOpen={isOpen}
                title={titleMap[severity]}
                description={descMap[severity]}
                severity={severity}
                confirmLabel="Подтвердить"
                cancelLabel="Отмена"
                onConfirm={() => { setIsOpen(false); }}
                onCancel={() => { setIsOpen(false); }}
            />
        </div>
    );
}

export const InteractiveInfo: Story = {
    args: {
        isOpen: false,
        title: 'Подтвердите действие',
        onConfirm: () => {},
        onCancel: () => {},
    },
    render: () => <InteractiveDemo severity="info" />,
};

export const InteractiveWarning: Story = {
    args: {
        isOpen: false,
        title: 'Опасное изменение',
        onConfirm: () => {},
        onCancel: () => {},
    },
    render: () => <InteractiveDemo severity="warning" />,
};

export const InteractiveDanger: Story = {
    args: {
        isOpen: false,
        title: 'Удалить безвозвратно?',
        onConfirm: () => {},
        onCancel: () => {},
    },
    render: () => <InteractiveDemo severity="danger" />,
};
