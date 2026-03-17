import type { JSX } from 'react';
import { clsx } from 'clsx';
import MDEditor from '@uiw/react-md-editor';
import '@uiw/react-markdown-preview/markdown.css';
import styles from './MdRenderer.module.scss';

export interface MdRendererProps {
    source: string;
    className?: string | undefined;
}

export function MdRenderer({ source, className }: MdRendererProps): JSX.Element {
    return (
        <div className={clsx(styles.wrapper, className)} data-color-mode="light">
            {/* wmde-markdown-var allows CSS variable inheritance from parent context */}
            <div className="wmde-markdown-var" />
            <MDEditor.Markdown source={source} />
        </div>
    );
}
