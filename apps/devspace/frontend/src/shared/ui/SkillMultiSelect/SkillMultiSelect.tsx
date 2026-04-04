import { clsx } from "clsx";
import { Search, X } from "lucide-react";
import { useState, useEffect, useRef, useCallback, useMemo, type JSX } from "react";

import { useDebounce } from "@/shared/lib/hooks/useDebounce";

import styles from "./SkillMultiSelect.module.scss";

export interface SkillMultiSelectOption {
  id: string;
  name: string;
  color?: string | undefined;
}

export interface SkillMultiSelectProps {
  /** Currently selected skills */
  value: SkillMultiSelectOption[];
  /** Called when selection changes */
  onChange: (skills: SkillMultiSelectOption[]) => void;
  /** Async loader that returns matching skills for the given query */
  loadOptions: (query: string) => Promise<SkillMultiSelectOption[]>;
  /** Maximum number of selectable items */
  max?: number | undefined;
  /** Placeholder while input is empty */
  placeholder?: string | undefined;
  /** Debounce delay in ms */
  delay?: number | undefined;
  /** Whether the component is disabled */
  disabled?: boolean | undefined;
  /** Custom CSS class */
  className?: string | undefined;
  /** Label for accessibility */
  label?: string | undefined;
}

export function SkillMultiSelect({
  value,
  onChange,
  loadOptions,
  max = 10,
  placeholder = "Добавить навык...",
  delay = 300,
  disabled = false,
  className,
  label,
}: SkillMultiSelectProps): JSX.Element {
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [options, setOptions] = useState<SkillMultiSelectOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(-1);
  const debouncedQuery = useDebounce(query, delay);

  const containerReference = useRef<HTMLDivElement>(null);
  const inputReference = useRef<HTMLInputElement>(null);
  const listReference = useRef<HTMLUListElement>(null);

  const selectedIds = useMemo(() => new Set(value.map((skill) => skill.id)), [value]);
  const isAtMax = value.length >= max;

  // Load options when debounced query changes
  useEffect(() => {
    if (!isOpen) return;

    let cancelled = false;

    void loadOptions(debouncedQuery)
      .then((result) => {
        if (!cancelled) {
          // Filter out already selected items
          setOptions(result.filter((option) => !selectedIds.has(option.id)));
          setIsLoading(false);
          setHighlightedIndex(-1);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setOptions([]);
          setIsLoading(false);
        }
      });

    return (): void => {
      cancelled = true;
    };
  }, [debouncedQuery, isOpen, selectedIds, loadOptions]);

  // Close on outside click
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent): void => {
      if (
        containerReference.current !== null &&
        !containerReference.current.contains(event.target as Node)
      ) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return (): void => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  const handleFocus = useCallback((): void => {
    if (disabled || isAtMax) return;
    setIsOpen(true);
  }, [disabled, isAtMax]);

  const handleSelect = useCallback(
    (option: SkillMultiSelectOption): void => {
      if (isAtMax) return;
      onChange([...value, option]);
      setQuery("");
      inputReference.current?.focus();
    },
    [onChange, value, isAtMax],
  );

  const handleRemove = useCallback(
    (skillId: string): void => {
      onChange(value.filter((skill) => skill.id !== skillId));
    },
    [onChange, value],
  );

  const handleKeyDown = (event: React.KeyboardEvent): void => {
    if (event.key === "Escape") {
      setIsOpen(false);
      return;
    }
    if (event.key === "Backspace" && query === "" && value.length > 0) {
      const lastItem = value[value.length - 1];
      if (lastItem !== undefined) {
        handleRemove(lastItem.id);
      }
      return;
    }
    if (!isOpen) return;

    if (event.key === "ArrowDown") {
      event.preventDefault();
      setHighlightedIndex((previous) => Math.min(previous + 1, options.length - 1));
    } else if (event.key === "ArrowUp") {
      event.preventDefault();
      setHighlightedIndex((previous) => Math.max(previous - 1, 0));
    } else if (
      event.key === "Enter" &&
      highlightedIndex >= 0 &&
      highlightedIndex < options.length
    ) {
      event.preventDefault();
      const selected = options[highlightedIndex];
      if (selected !== undefined) {
        handleSelect(selected);
      }
    }
  };

  // Scroll highlighted item into view
  useEffect(() => {
    if (highlightedIndex < 0 || listReference.current === null) return;
    const items = listReference.current.querySelectorAll('[role="option"]');
    items[highlightedIndex]?.scrollIntoView({ block: "nearest" });
  }, [highlightedIndex]);

  const renderDropdownContent = (): JSX.Element => {
    if (isLoading) {
      return <li className={styles.statusItem}>Загрузка...</li>;
    }
    if (options.length === 0) {
      return (
        <li className={styles.statusItem}>
          {query !== "" ? "Ничего не найдено" : "Начните вводить название"}
        </li>
      );
    }
    return (
      <>
        {options.map((option, index) => (
          <li
            key={option.id}
            role="option"
            aria-selected={false}
            className={clsx(styles.option, index === highlightedIndex && styles.highlighted)}
            onMouseEnter={() => {
              setHighlightedIndex(index);
            }}
            onClick={() => {
              handleSelect(option);
            }}
          >
            {option.color !== undefined && (
              <span className={styles.optionDot} style={{ backgroundColor: `#${option.color}` }} />
            )}
            <span className={styles.optionName}>{option.name}</span>
          </li>
        ))}
      </>
    );
  };

  return (
    <div
      ref={containerReference}
      className={clsx(styles.container, disabled && styles.disabled, className)}
      onKeyDown={handleKeyDown}
    >
      {label !== undefined && <span className={styles.label}>{label}</span>}

      <div
        className={clsx(styles.inputArea, isOpen && styles.focused)}
        onClick={() => {
          inputReference.current?.focus();
        }}
      >
        {/* Selected tags */}
        {value.map((skill) => (
          <span key={skill.id} className={styles.tag}>
            {skill.color !== undefined && (
              <span className={styles.tagDot} style={{ backgroundColor: `#${skill.color}` }} />
            )}
            <span className={styles.tagName}>{skill.name}</span>
            <button
              type="button"
              className={styles.tagRemove}
              onClick={(event) => {
                event.stopPropagation();
                handleRemove(skill.id);
              }}
              aria-label={`Убрать ${skill.name}`}
              disabled={disabled}
            >
              <X size={12} />
            </button>
          </span>
        ))}

        {/* Input */}
        {!isAtMax && (
          <div className={styles.inputWrapper}>
            <Search size={14} className={styles.searchIcon} />
            <input
              ref={inputReference}
              type="text"
              className={styles.input}
              value={query}
              onChange={(event) => {
                setQuery(event.target.value);
                setIsLoading(true);
              }}
              onFocus={handleFocus}
              placeholder={value.length === 0 ? placeholder : ""}
              disabled={disabled}
              role="combobox"
              aria-expanded={isOpen}
              aria-autocomplete="list"
            />
          </div>
        )}

        {isAtMax && <span className={styles.maxHint}>Макс. {max}</span>}
      </div>

      {/* Dropdown list */}
      {isOpen && (
        <ul ref={listReference} className={styles.dropdown} role="listbox">
          {renderDropdownContent()}
        </ul>
      )}
    </div>
  );
}
