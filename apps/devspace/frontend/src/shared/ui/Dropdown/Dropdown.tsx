import { type SelectHTMLAttributes, type JSX, forwardRef } from 'react';
import { clsx } from 'clsx';
import { ChevronDown } from 'lucide-react';
import styles from './Dropdown.module.scss';

export interface DropdownOption {
    label: string;
    value: string;
}

export interface DropdownProps extends Omit<SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
    options: DropdownOption[];
    value?: string;
    onChange: (value: string) => void;
    placeholder?: string;
}

export const Dropdown = forwardRef<HTMLSelectElement, DropdownProps>(
    ({ className, options, value, onChange, placeholder, disabled, ...props }, reference): JSX.Element => {
        return (
            <div className={clsx(styles.wrapper, disabled === true && styles.disabled, className)}>
                <select
                    ref={reference}
                    className={styles.dropdown}
                    value={value ?? ''}
                    disabled={disabled}
                    onChange={(event) => { onChange(event.target.value); }}
                    {...props}
                >
                    {placeholder !== undefined && (
                        <option value="" disabled hidden>
                            {placeholder}
                        </option>
                    )}
                    {options.map((option) => (
                        <option key={option.value} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
                <div className={styles.iconWrapper}>
                    <ChevronDown size={16} />
                </div>
            </div>
        );
    }
);

Dropdown.displayName = 'Dropdown';
