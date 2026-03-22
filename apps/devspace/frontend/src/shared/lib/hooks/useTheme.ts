import { useState, useEffect } from 'react';

export type Theme = 'light' | 'dark';

function getInitialTheme(): Theme {
    const stored = localStorage.getItem('theme');
    if (stored === 'dark' || stored === 'light') return stored;
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

export function useTheme(): { theme: Theme; toggleTheme: () => void } {
    const [theme, setTheme] = useState<Theme>(getInitialTheme);

    useEffect(() => {
        document.documentElement.classList.toggle('dark', theme === 'dark');
        localStorage.setItem('theme', theme);
    }, [theme]);

    const toggleTheme = (): void => {
        setTheme((previous) => (previous === 'light' ? 'dark' : 'light'));
    };

    return { theme, toggleTheme };
}
