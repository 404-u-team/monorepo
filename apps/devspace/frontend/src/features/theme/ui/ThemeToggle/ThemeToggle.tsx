import { Moon, Sun } from "lucide-react";
import type { JSX } from "react";

import { useTheme } from "@/shared/lib/hooks/useTheme";

import styles from "./ThemeToggle.module.scss";

export function ThemeToggle(): JSX.Element {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      className={styles.button}
      onClick={toggleTheme}
      aria-label={theme === "light" ? "Включить тёмную тему" : "Включить светлую тему"}
    >
      {theme === "light" ? <Moon size={20} /> : <Sun size={20} />}
    </button>
  );
}
