import { type ButtonHTMLAttributes, type JSX } from 'react';
import { clsx } from 'clsx';
import styles from './Button.module.scss';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'primary' | 'outline' | 'clear';
    fullWidth?: boolean;
}

export function Button({
    className,
    variant = 'primary',
    fullWidth,
    children,
    ...props
}: ButtonProps): JSX.Element {
    return (
        <button
            className={clsx(
                styles.button,
                styles[variant as keyof typeof styles],
                fullWidth === true ? styles.fullWidth : undefined,
                className
            )}
            {...props}
        >
            {children}
        </button>
    );
}
