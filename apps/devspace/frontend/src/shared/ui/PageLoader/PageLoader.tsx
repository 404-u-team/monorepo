import { useRouterState } from "@tanstack/react-router";
import { clsx } from "clsx";
import { useEffect, useReducer, useRef, type JSX } from "react";

import styles from "./PageLoader.module.scss";

type Phase = "idle" | "loading" | "completing";

interface BarState {
  phase: Phase;
  width: number;
}

type BarAction =
  | { type: "loading_start" }
  | { type: "advance"; width: number }
  | { type: "complete" }
  | { type: "reset" };

function barReducer(_state: BarState, action: BarAction): BarState {
  switch (action.type) {
    case "loading_start":
      return { phase: "loading", width: 0 };
    case "advance":
      return { ..._state, width: action.width };
    case "complete":
      return { phase: "completing", width: 100 };
    case "reset":
      return { phase: "idle", width: 0 };
  }
}

export function PageLoader(): JSX.Element | undefined {
  const isLoading = useRouterState({ select: (s) => s.isLoading });
  const [{ phase, width }, dispatch] = useReducer(barReducer, { phase: "idle", width: 0 });

  const timersReference = useRef<ReturnType<typeof setTimeout>[]>([]);

  const clearTimers = (): void => {
    for (const t of timersReference.current) clearTimeout(t);
    timersReference.current = [];
  };

  const schedule = (callback: () => void, delay: number): void => {
    timersReference.current.push(setTimeout(callback, delay));
  };

  useEffect(() => {
    if (isLoading) {
      clearTimers();
      dispatch({ type: "loading_start" });
      schedule(() => {
        dispatch({ type: "advance", width: 20 });
      }, 30);
      schedule(() => {
        dispatch({ type: "advance", width: 50 });
      }, 300);
      schedule(() => {
        dispatch({ type: "advance", width: 72 });
      }, 900);
      schedule(() => {
        dispatch({ type: "advance", width: 85 });
      }, 2000);
    } else if (phase === "loading") {
      clearTimers();
      dispatch({ type: "complete" });
      schedule(() => {
        dispatch({ type: "reset" });
      }, 450);
    }

    return clearTimers;
  }, [isLoading, phase]);

  if (phase === "idle") return undefined;

  return (
    <div
      className={clsx(styles.bar, phase === "completing" && styles.completing)}
      style={{ width: `${String(width)}%` }}
      role="progressbar"
      aria-valuenow={width}
      aria-valuemin={0}
      aria-valuemax={100}
    />
  );
}
