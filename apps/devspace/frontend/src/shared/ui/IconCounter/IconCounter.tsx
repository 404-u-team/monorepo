import { clsx } from "clsx";
import { type JSX, type ReactNode, type MouseEvent } from "react";

import styles from "./IconCounter.module.scss";

export interface IconCounterProps {
  icon: ReactNode;
  count: number;
  active?: boolean | undefined;
  onClick?: ((event: MouseEvent<HTMLButtonElement>) => void) | undefined;
  className?: string | undefined;
}

function formatCount(count: number): string {
  if (count >= 1_000_000) {
    return `${(count / 1_000_000).toFixed(1).replace(/\.0$/, "")}M`;
  }
  if (count >= 1_000) {
    return `${(count / 1_000).toFixed(1).replace(/\.0$/, "")}k`;
  }
  return String(count);
}

export function IconCounter({
  icon,
  count,
  active,
  onClick,
  className,
}: IconCounterProps): JSX.Element {
  const content = (
    <>
      <span className={styles.icon}>{icon}</span>
      <span className={styles.count}>{formatCount(count)}</span>
    </>
  );

  if (onClick !== undefined) {
    return (
      <button
        type="button"
        className={clsx(
          styles.wrapper,
          styles.clickable,
          active === true && styles.active,
          className,
        )}
        onClick={onClick}
      >
        {content}
      </button>
    );
  }

  return (
    <span className={clsx(styles.wrapper, active === true && styles.active, className)}>
      {content}
    </span>
  );
}
