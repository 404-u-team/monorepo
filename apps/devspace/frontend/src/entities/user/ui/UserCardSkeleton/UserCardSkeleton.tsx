import type { JSX } from "react";
import { clsx } from "clsx";
import { Skeleton } from "@/shared/ui";
import styles from "./UserCardSkeleton.module.scss";

export interface UserCardSkeletonProps {
  className?: string | undefined;
}

export function UserCardSkeleton({
  className,
}: UserCardSkeletonProps): JSX.Element {
  return (
    <div className={clsx(styles.card, className)}>
      <div className={styles.header}>
        <div className={styles.avatarWrapper}>
          <Skeleton className={styles.avatar} />
        </div>
        <div className={styles.textInfo}>
          <Skeleton width={120} height={20} />
          <Skeleton width={80} height={16} />
        </div>
      </div>

      <div className={styles.bio}>
        <Skeleton width="100%" height={14} />
        <Skeleton width="100%" height={14} />
        <Skeleton width="90%" height={14} />
      </div>

      <div className={styles.skills}>
        <Skeleton width={60} height={22} borderRadius={100} />
        <Skeleton width={50} height={22} borderRadius={100} />
        <Skeleton width={55} height={22} borderRadius={100} />
        <Skeleton width={45} height={22} borderRadius={100} />
      </div>

      <div className={styles.actions}>
        <Skeleton width={100} height={32} borderRadius={8} />
        <Skeleton width={100} height={32} borderRadius={8} />
      </div>
    </div>
  );
}
