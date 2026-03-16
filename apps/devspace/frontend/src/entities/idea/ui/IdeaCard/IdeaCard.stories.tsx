import { createElement } from 'react';
import { http, HttpResponse } from 'msw';
import type { Meta, StoryObj } from '@storybook/react';
import { StoreContext, type IRootStore } from '@/shared/lib/store';
import { IdeaCard } from './IdeaCard';

const MOCK_IDEA = {
    id: '1',
    title: 'AI Код-ревью',
    description: 'Автоматический инструмент использующий LLM для обнаружения проблем в коде',
    category: 'AI & ML',
    author_id: 'user-1',
    created_at: '2026-03-01T12:00:00Z',
    updated_at: '2026-03-10T15:30:00Z',
    views_count: 2600,
    favorites_count: 237,
};

const MOCK_AUTHOR = {
    id: 'user-1',
    nickname: 'TestUser1',
    avatar_uri: '',
    bio: 'Full-stack developer',
    main_role: 'Developer',
    skills: [],
};

/** Переопределяем глобальный стор — тут нужен авторизованный пользователь */
const authenticatedStore: IRootStore = {
    userStore: { isAuthenticated: true },
};

const meta = {
    title: 'entities/IdeaCard',
    component: IdeaCard,
    parameters: {
        layout: 'centered',
        msw: {
            handlers: [
                http.get('*/ideas/:ideaId', () => {
                    return HttpResponse.json(MOCK_IDEA);
                }),
                http.get('*/users/:userId', () => {
                    return HttpResponse.json(MOCK_AUTHOR);
                }),
                http.post('*/ideas/:ideaId/favorite', () => {
                    return HttpResponse.json({ is_favorite: true });
                }),
            ],
        },
    },
    tags: ['autodocs'],
    decorators: [
        (Story) => createElement(
            StoreContext.Provider,
            { value: authenticatedStore },
            createElement('div', { style: { width: '340px' } }, createElement(Story)),
        ),
    ],
} satisfies Meta<typeof IdeaCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        ideaId: '1',
    },
};

export const LongContent: Story = {
    args: {
        ideaId: '2',
    },
    parameters: {
        msw: {
            handlers: [
                http.get('*/ideas/:ideaId', () => {
                    return HttpResponse.json({
                        ...MOCK_IDEA,
                        id: '2',
                        title: 'Платформа для совместной разработки с интеграцией AI-ассистентов и автоматизированным код-ревью',
                        description: 'Разработка платформы которая объединяет возможности совместного редактирования кода, AI-ассистентов для генерации и рефакторинга, а также автоматизированного код-ревью с поддержкой множества языков программирования',
                        category: 'Веб-разработка',
                        favorites_count: 1500000,
                        views_count: 5200000,
                    });
                }),
                http.get('*/users/:userId', () => {
                    return HttpResponse.json(MOCK_AUTHOR);
                }),
            ],
        },
    },
};

export const Loading: Story = {
    args: {
        ideaId: '3',
    },
    parameters: {
        msw: {
            handlers: [
                http.get('*/ideas/:ideaId', async () => {
                    await new Promise(() => {});
                    return HttpResponse.json(MOCK_IDEA);
                }),
            ],
        },
    },
};
