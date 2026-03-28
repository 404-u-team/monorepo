import type { JSX } from 'react';
import { User } from 'lucide-react';
import { clsx } from 'clsx';
import styles from './UserAvatar.module.scss';

export interface UserAvatarProps {
    avatarUrl?: string | undefined;
    nickname?: string | undefined;
    size?: number | undefined;
    className?: string | undefined;
}

export function UserAvatar({
    avatarUrl,
    nickname,
    size = 40,
    className,
}: UserAvatarProps): JSX.Element {
    const iconSize = Math.round(size * 0.5);

    return (
        <div
            className={clsx(styles.wrapper, className)}
            style={{ width: size, height: size, minWidth: size, borderRadius: '50%' }}
        >
            {avatarUrl !== undefined && avatarUrl !== '' ? (
                <img
                    src={avatarUrl}
                    alt={nickname ?? ''}
                    className={styles.image}
                    onError={(event) => {
                        (event.currentTarget as HTMLImageElement).style.display = 'none';
                        const placeholder = event.currentTarget.nextElementSibling as HTMLElement | null;
                        if (placeholder !== null) placeholder.style.display = 'flex';
                    }}
                />
            ) : undefined}
            <div
                className={styles.placeholder}
                style={avatarUrl !== undefined && avatarUrl !== '' ? { display: 'none' } : {}}
            >
                <User size={iconSize} />
            </div>
        </div>
    );
}
