import { forwardRef, type InputHTMLAttributes, type ReactNode } from 'react';
import { clsx } from 'clsx';
import styles from './Input.module.scss';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    iconLeft?: ReactNode;
    iconRight?: ReactNode;
    error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ className, iconLeft, iconRight, error, ...props }, reference) => {
        return (
            <div className={clsx(styles.wrapper, className)}>
                <div className={clsx(styles.inputContainer, error !== undefined && error !== '' && styles.hasError)}>
                    {iconLeft !== undefined && <span className={styles.iconLeft}>{iconLeft}</span>}
                    <input ref={reference} className={styles.input} {...props} />
                    {iconRight !== undefined && <span className={styles.iconRight}>{iconRight}</span>}
                </div>
                {error !== undefined && error !== '' && <span className={styles.error}>{error}</span>}
            </div>
        );
    }
);

Input.displayName = 'Input';
