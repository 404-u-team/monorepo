import type { JSX } from 'react';
import { clsx } from 'clsx';
import { Skeleton } from '@/shared/ui';
import styles from './IdeaCardSkeleton.module.scss';

export interface IdeaCardSkeletonProps {
    className?: string | undefined;
}

export function IdeaCardSkeleton({ className }: IdeaCardSkeletonProps): JSX.Element {
    return (
        <div className={clsx(styles.card, className)}>
            <Skeleton className={styles.image} />

            <div className={styles.content}>
                <Skeleton width="80%" height={20} />
                <Skeleton width="100%" height={14} />
                <Skeleton width="60%" height={14} />
            </div>

            <div className={styles.footer}>
                <div className={styles.author}>
                    <Skeleton width={24} height={24} borderRadius="50%" />
                    <Skeleton width={80} height={14} />
                </div>
                <div className={styles.stats}>
                    <Skeleton width={50} height={14} />
                    <Skeleton width={50} height={14} />
                </div>
            </div>
        </div>
    );
}
