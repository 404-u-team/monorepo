import type { JSX } from 'react';
import { clsx } from 'clsx';
import MDEditor from '@uiw/react-md-editor';
import '@uiw/react-markdown-preview/markdown.css';
import { useTheme } from '@/shared/lib/hooks/useTheme';
import styles from './MdRenderer.module.scss';

export interface MdRendererProps {
    source: string;
    className?: string | undefined;
}

export function MdRenderer({ source, className }: MdRendererProps): JSX.Element {
    const { theme } = useTheme();
    return (
        <div className={clsx(styles.wrapper, className)} data-color-mode={theme}>
            <div className="wmde-markdown-var" />
            <MDEditor.Markdown source={source} />
        </div>
    );
}
