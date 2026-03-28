import { useEffect, type JSX, type ReactNode } from 'react';
import { createPortal } from 'react-dom';
import { X, AlertTriangle, AlertCircle, Info } from 'lucide-react';
import { clsx } from 'clsx';
import { Button } from '../Button/Button';
import styles from './ConfirmModal.module.scss';

export type ConfirmModalSeverity = 'info' | 'warning' | 'danger';

export interface ConfirmModalProps {
    isOpen: boolean;
    title: string;
    description?: string | undefined;
    confirmLabel?: string | undefined;
    cancelLabel?: string | undefined;
    severity?: ConfirmModalSeverity | undefined;
    isLoading?: boolean | undefined;
    children?: ReactNode | undefined;
    onConfirm: () => void;
    onCancel: () => void;
}

const severityIcon: Record<ConfirmModalSeverity, JSX.Element> = {
    info: <Info size={22} />,
    warning: <AlertTriangle size={22} />,
    danger: <AlertCircle size={22} />,
};

const iconClass: Record<ConfirmModalSeverity, string> = {
    info: 'iconInfo',
    warning: 'iconWarning',
    danger: 'iconDanger',
};

const confirmClass: Record<ConfirmModalSeverity, string> = {
    info: '',
    warning: 'confirmWarning',
    danger: 'confirmDanger',
};

export function ConfirmModal({
    isOpen,
    title,
    description,
    confirmLabel = 'Подтвердить',
    cancelLabel = 'Отмена',
    severity = 'info',
    isLoading = false,
    children,
    onConfirm,
    onCancel,
}: ConfirmModalProps): JSX.Element | null {
    useEffect(() => {
        if (!isOpen) return undefined;

        const handleKeyDown = (event: KeyboardEvent): void => {
            if (event.key === 'Escape') {
                onCancel();
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        document.body.style.overflow = 'hidden';

        return (): void => {
            document.removeEventListener('keydown', handleKeyDown);
            document.body.style.overflow = '';
        };
    }, [isOpen, onCancel]);

    // no component if modal is closed
    // eslint-disable-next-line unicorn/no-null
    if (!isOpen) return null;

    return createPortal(
        <div
            className={styles.overlay}
            role="dialog"
            aria-modal="true"
            aria-labelledby="confirm-modal-title"
            onClick={(event) => {
                if (event.target === event.currentTarget) {
                    onCancel();
                }
            }}
        >
            <div className={clsx(styles.modal, styles[severity])}>
                <button
                    type="button"
                    className={styles.closeButton}
                    onClick={onCancel}
                    aria-label="Закрыть"
                >
                    <X size={18} />
                </button>

                <div className={clsx(styles.iconWrapper, styles[iconClass[severity]])}>
                    {severityIcon[severity]}
                </div>

                <div className={styles.body}>
                    <h2 id="confirm-modal-title" className={styles.title}>
                        {title}
                    </h2>
                    {description !== undefined && description !== '' && (
                        <p className={styles.description}>{description}</p>
                    )}
                    {children !== undefined && (
                        <div className={styles.content}>{children}</div>
                    )}
                </div>

                <div className={styles.actions}>
                    <Button
                        variant="outline"
                        onClick={onCancel}
                        disabled={isLoading}
                    >
                        {cancelLabel}
                    </Button>
                    <Button
                        className={clsx(styles.confirmButton, confirmClass[severity] !== '' && styles[confirmClass[severity]])}
                        onClick={onConfirm}
                        disabled={isLoading}
                    >
                        {isLoading ? 'Загрузка...' : confirmLabel}
                    </Button>
                </div>
            </div>
        </div>,
        document.body,
    );
}
