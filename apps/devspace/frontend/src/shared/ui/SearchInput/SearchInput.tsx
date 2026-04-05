import { clsx } from "clsx";
import { Search } from "lucide-react";
import { useState, useEffect, type JSX, type ChangeEvent } from "react";

import { useDebounce } from "@/shared/lib/hooks/useDebounce";

import { Input } from "../Input/Input";

import styles from "./SearchInput.module.scss";

export interface SearchInputProps {
  value?: string;
  onSearch: (value: string) => void;
  placeholder?: string;
  className?: string;
  delay?: number;
}

export function SearchInput({
  value = "",
  onSearch,
  placeholder = "Поиск...",
  className,
  delay = 500,
}: SearchInputProps): JSX.Element {
  const [localValue, setLocalValue] = useState(value);
  const debouncedValue = useDebounce(localValue, delay);

  // Sync input when prop changes externally
  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  useEffect(() => {
    // Prevent calling onSearch immediately on mount if value hasn't changed
    if (debouncedValue !== value) {
      onSearch(debouncedValue);
    }
  }, [debouncedValue, value, onSearch]);

  const handleChange = (event: ChangeEvent<HTMLInputElement>): void => {
    setLocalValue(event.target.value);
  };

  return (
    <Input
      value={localValue}
      onChange={handleChange}
      placeholder={placeholder}
      className={clsx(styles.searchInput, className)}
      iconLeft={<Search size={18} className={styles.icon} />}
      type="search"
    />
  );
}
