import type { JSX, ReactNode } from 'react';
import { clsx } from 'clsx';
import styles from './Badge.module.scss';

export interface BadgeProps {
    children: ReactNode;
    className?: string | undefined;
}

export function Badge({ children, className }: BadgeProps): JSX.Element {
    return (
        <span className={clsx(styles.badge, className)}>
            {children}
        </span>
    );
}
