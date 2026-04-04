import { clsx } from "clsx";
import type { JSX } from "react";

import { Skeleton } from "@/shared/ui";

import styles from "./UserCardSkeleton.module.scss";

export interface UserCardSkeletonProps {
  className?: string | undefined;
}

export function UserCardSkeleton({ className }: UserCardSkeletonProps): JSX.Element {
  return (
    <div className={clsx(styles.card, className)}>
      <div className={styles.header}>
        <Skeleton width={56} height={56} borderRadius={28} />
        <div className={styles.textInfo}>
          <Skeleton width={130} height={18} />
          <Skeleton width={90} height={14} />
        </div>
      </div>

      <div className={styles.bio}>
        <Skeleton width="100%" height={13} />
        <Skeleton width="75%" height={13} />
      </div>

      <div className={styles.skills}>
        <Skeleton width={70} height={24} borderRadius={20} />
        <Skeleton width={60} height={24} borderRadius={20} />
        <Skeleton width={80} height={24} borderRadius={20} />
      </div>

      <div className={styles.actions}>
        <Skeleton width={100} height={32} borderRadius={8} />
      </div>
    </div>
  );
}
