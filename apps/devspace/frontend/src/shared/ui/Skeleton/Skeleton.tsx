import { clsx } from "clsx";
import type { CSSProperties, JSX } from "react";

import styles from "./Skeleton.module.scss";

export interface SkeletonProps {
  width?: string | number | undefined;
  height?: string | number | undefined;
  borderRadius?: string | number | undefined;
  className?: string | undefined;
}

export function Skeleton({ width, height, borderRadius, className }: SkeletonProps): JSX.Element {
  const style: CSSProperties = {
    width,
    height,
    borderRadius,
  };

  return <div className={clsx(styles.skeleton, className)} style={style} />;
}
