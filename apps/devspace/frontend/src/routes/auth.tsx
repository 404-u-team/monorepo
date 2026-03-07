import { createFileRoute } from '@tanstack/react-router';
import { AuthForm } from '@/features/auth';

export const Route = createFileRoute('/auth')({
    component: AuthPage,
});

function AuthPage(): React.JSX.Element {
    return (
        <div style={{
            display: 'flex',
            minHeight: '100vh',
            padding: '2rem 1rem',
            backgroundColor: 'var(--bg--main)',
            alignItems: 'center',
            justifyContent: 'center'
        }}>
            <AuthForm />
        </div>
    );
}
