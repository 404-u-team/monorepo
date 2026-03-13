import { createElement } from 'react';
import { http, HttpResponse } from 'msw';
import type { Meta, StoryObj } from '@storybook/react';
import { ProjectCard } from './ProjectCard';

const MOCK_PROJECT = {
    id: 'proj-1',
    title: 'AI Код-ревью',
    description: 'Создаем мощный AI инструмент для автономного ревью кода',
    leader_id: 'user-1',
    status: 'open',
    idea_id: null,
    created_at: '2026-01-15T10:00:00Z',
    updated_at: '2026-03-10T15:00:00Z',
    slots: [
        {
            id: 'slot-1',
            project_id: 'proj-1',
            skill_category_id: 'skill-react',
            skill: { id: 'skill-react', name: 'React', color: '3B82F6', icon: 'react' },
            title: 'Frontend Developer',
            description: '',
            status: 'closed',
            user_id: 'user-2',
            created_at: '2026-01-20T10:00:00Z',
        },
        {
            id: 'slot-2',
            project_id: 'proj-1',
            skill_category_id: 'skill-mobx',
            skill: { id: 'skill-mobx', name: 'MobX', color: 'EF6C00', icon: 'mobx' },
            title: 'State Management',
            description: '',
            status: 'closed',
            user_id: 'user-3',
            created_at: '2026-01-20T10:00:00Z',
        },
        {
            id: 'slot-3',
            project_id: 'proj-1',
            skill_category_id: 'skill-mysql',
            skill: { id: 'skill-mysql', name: 'MySQL', color: '00758F', icon: 'mysql' },
            title: 'Database Engineer',
            description: '',
            status: 'closed',
            user_id: 'user-4',
            created_at: '2026-01-20T10:00:00Z',
        },
        {
            id: 'slot-4',
            project_id: 'proj-1',
            skill_category_id: 'skill-fastapi',
            skill: { id: 'skill-fastapi', name: 'FastAPI', color: '009688', icon: 'fastapi' },
            title: 'Backend Developer',
            description: '',
            status: 'closed',
            user_id: 'user-5',
            created_at: '2026-01-20T10:00:00Z',
        },
        {
            id: 'slot-5',
            project_id: 'proj-1',
            skill_category_id: 'skill-devops',
            skill: { id: 'skill-devops', name: 'DevOps', color: '7C3AED', icon: 'devops' },
            title: 'DevOps Engineer',
            description: '',
            status: 'open',
            user_id: null,
            created_at: '2026-02-01T10:00:00Z',
        },
        {
            id: 'slot-6',
            project_id: 'proj-1',
            skill_category_id: 'skill-devops',
            skill: { id: 'skill-devops', name: 'DevOps', color: '7C3AED', icon: 'devops' },
            title: 'DevOps Engineer #2',
            description: '',
            status: 'open',
            user_id: null,
            created_at: '2026-02-01T10:00:00Z',
        },
    ],
};

const MOCK_USERS: Record<string, object> = {
    'user-2': { id: 'user-2', nickname: 'Alice', avatar_uri: '', bio: '', main_role: 'Dev', skills: [] },
    'user-3': { id: 'user-3', nickname: 'Bob', avatar_uri: '', bio: '', main_role: 'Dev', skills: [] },
    'user-4': { id: 'user-4', nickname: 'Charlie', avatar_uri: '', bio: '', main_role: 'DBA', skills: [] },
    'user-5': { id: 'user-5', nickname: 'Diana', avatar_uri: '', bio: '', main_role: 'Backend', skills: [] },
};

const meta = {
    title: 'entities/ProjectCard',
    component: ProjectCard,
    parameters: {
        layout: 'centered',
        msw: {
            handlers: [
                http.get('*/projects/:projectId', () => {
                    return HttpResponse.json(MOCK_PROJECT);
                }),
                http.get('*/users/:userId', ({ params }) => {
                    const user = MOCK_USERS[params.userId as string];
                    if (user !== undefined) return HttpResponse.json(user);
                    return new HttpResponse(null, { status: 404 });
                }),
            ],
        },
    },
    tags: ['autodocs'],
    decorators: [
        (Story) => createElement('div', { style: { width: '420px' } }, createElement(Story)),
    ],
} satisfies Meta<typeof ProjectCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        projectId: 'proj-1',
    },
};

export const LongContent: Story = {
    args: {
        projectId: 'proj-long',
    },
    parameters: {
        msw: {
            handlers: [
                http.get('*/projects/:projectId', () => {
                    return HttpResponse.json({
                        ...MOCK_PROJECT,
                        id: 'proj-long',
                        title: 'Глобальная Платформа для Масштабируемого и Распределенного Машинного Обучения с Интеграцией Нейросетей Вавилона Кода',
                        description: 'Этот проект представляет собой революционную попытку объединить все известные алгоритмы в единую мета-модель. Мы хотим создать экосистему, где сотни микросервисов, написанных на десятке разных языков программирования, динамически обмениваются весами моделей в режиме реального времени на базе распределенного графа знаний и квантовых вычислений.',
                        slots: [
                            ...MOCK_PROJECT.slots,
                            {
                                id: 'slot-7',
                                project_id: 'proj-long',
                                skill_category_id: 'skill-k8s',
                                skill: { id: 'skill-k8s', name: 'Kubernetes', color: '326CE5', icon: 'kubernetes' },
                                title: 'Infrastructure Lead',
                                description: '',
                                status: 'open',
                                user_id: null,
                                created_at: '2026-02-10T10:00:00Z',
                            },
                            {
                                id: 'slot-8',
                                project_id: 'proj-long',
                                skill_category_id: 'skill-rust',
                                skill: { id: 'skill-rust', name: 'Rust', color: '000000', icon: 'rust' },
                                title: 'Performance Engineer',
                                description: '',
                                status: 'open',
                                user_id: null,
                                created_at: '2026-02-11T10:00:00Z',
                            },
                            {
                                id: 'slot-9',
                                project_id: 'proj-long',
                                skill_category_id: 'skill-go',
                                skill: { id: 'skill-go', name: 'Golang', color: '00ADD8', icon: 'go' },
                                title: 'Microservices Architect',
                                description: '',
                                status: 'closed',
                                user_id: 'user-6',
                                created_at: '2026-02-12T10:00:00Z',
                            },
                            {
                                id: 'slot-10',
                                project_id: 'proj-long',
                                skill_category_id: 'skill-graphql',
                                skill: { id: 'skill-graphql', name: 'GraphQL', color: 'E10098', icon: 'graphql' },
                                title: 'API Designer',
                                description: '',
                                status: 'closed',
                                user_id: 'user-7',
                                created_at: '2026-02-13T10:00:00Z',
                            },
                            {
                                id: 'slot-11',
                                project_id: 'proj-long',
                                skill_category_id: 'skill-aws',
                                skill: { id: 'skill-aws', name: 'AWS', color: 'FF9900', icon: 'aws' },
                                title: 'Cloud Engineer',
                                description: '',
                                status: 'closed',
                                user_id: 'user-8',
                                created_at: '2026-02-14T10:00:00Z',
                            },
                        ],
                    });
                }),
                http.get('*/users/:userId', ({ params }) => {
                    const longUsers: Record<string, object> = {
                        ...MOCK_USERS,
                        'user-6': { id: 'user-6', nickname: 'Edward', avatar_uri: '', bio: '', main_role: 'Go Dev', skills: [] },
                        'user-7': { id: 'user-7', nickname: 'Fiona', avatar_uri: '', bio: '', main_role: 'API Dev', skills: [] },
                        'user-8': { id: 'user-8', nickname: 'George', avatar_uri: '', bio: '', main_role: 'Cloud', skills: [] },
                    };
                    const user = longUsers[params.userId as string];
                    if (user !== undefined) return HttpResponse.json(user);
                    return new HttpResponse(null, { status: 404 });
                }),
            ],
        },
    },
};

export const AsLink: Story = {
    args: {
        projectId: 'proj-1',
        to: '/projects/proj-1',
    },
};

export const Closed: Story = {
    args: {
        projectId: 'proj-closed',
    },
    parameters: {
        msw: {
            handlers: [
                http.get('*/projects/:projectId', () => {
                    return HttpResponse.json({
                        ...MOCK_PROJECT,
                        id: 'proj-closed',
                        status: 'closed',
                        slots: MOCK_PROJECT.slots.map((s) => ({ ...s, status: 'closed' as const })),
                    });
                }),
                http.get('*/users/:userId', ({ params }) => {
                    const user = MOCK_USERS[params.userId as string];
                    if (user !== undefined) return HttpResponse.json(user);
                    return new HttpResponse(null, { status: 404 });
                }),
            ],
        },
    },
};

export const Loading: Story = {
    args: {
        projectId: 'proj-loading',
    },
    parameters: {
        msw: {
            handlers: [
                http.get('*/projects/:projectId', async () => {
                    await new Promise(() => {});
                    return HttpResponse.json(MOCK_PROJECT);
                }),
            ],
        },
    },
};
