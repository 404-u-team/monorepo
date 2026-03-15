import type { JSX } from 'react';
import { clsx } from 'clsx';
import { Skeleton } from '@/shared/ui';
import styles from './ProjectCardSkeleton.module.scss';

export interface ProjectCardSkeletonProps {
    className?: string | undefined;
}

export function ProjectCardSkeleton({ className }: ProjectCardSkeletonProps): JSX.Element {
    return (
        <div className={clsx(styles.card, className)}>
            <Skeleton className={styles.image} />

            <div className={styles.body}>
                <div className={styles.header}>
                    <Skeleton width="75%" height={20} />
                    <Skeleton width={20} height={20} borderRadius="50%" />
                </div>

                <Skeleton width="100%" height={14} />
                <Skeleton width="60%" height={14} />

                <div className={styles.skills}>
                    <Skeleton width={70} height={22} borderRadius={100} />
                    <Skeleton width={55} height={22} borderRadius={100} />
                    <Skeleton width={65} height={22} borderRadius={100} />
                </div>

                <div className={styles.slots}>
                    <Skeleton width={100} height={12} />
                    <div className={styles.slotBadges}>
                        <Skeleton width={60} height={22} borderRadius={100} />
                        <Skeleton width={60} height={22} borderRadius={100} />
                    </div>
                </div>

                <div className={styles.members}>
                    <Skeleton width={28} height={28} borderRadius="50%" />
                    <Skeleton width={28} height={28} borderRadius="50%" />
                    <Skeleton width={28} height={28} borderRadius="50%" />
                </div>
            </div>
        </div>
    );
}
