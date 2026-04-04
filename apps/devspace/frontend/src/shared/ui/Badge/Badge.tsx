import { clsx } from "clsx";
import type { JSX, ReactNode } from "react";

import styles from "./Badge.module.scss";

export interface BadgeProps {
  children: ReactNode;
  color?: string | undefined;
  icon?: ReactNode | undefined;
  className?: string | undefined;
}

export function Badge({ children, color, icon, className }: BadgeProps): JSX.Element {
  const style =
    color !== undefined
      ? ({ "--badge--bg": `#${color}20`, "--badge--text": `#${color}` } as React.CSSProperties)
      : undefined;

  return (
    <span className={clsx(styles.badge, className)} style={style}>
      {icon !== undefined && <span className={styles.icon}>{icon}</span>}
      {children}
    </span>
  );
}
