import preview from '@/../.storybook/preview';
import UserCard from './UserCard';

const meta = preview.meta({
    component: UserCard,
});

export const Default = meta.story({
    args: {
        avatar_uri: 'https://placehold.co/150x150',
        user_id: 'john_doe',
        mainRole: 'Frontend Developer',
        description: 'Люблю React и TypeScript',
        skill_id: ['react', 'typescript', 'scss'],
    },
});