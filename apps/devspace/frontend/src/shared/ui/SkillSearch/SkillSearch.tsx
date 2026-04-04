import { clsx } from "clsx";
import { Search, X, ChevronDown } from "lucide-react";
import { useState, useEffect, useRef, useCallback, type JSX } from "react";

import { useDebounce } from "@/shared/lib/hooks/useDebounce";

import styles from "./SkillSearch.module.scss";

export interface SkillSearchOption {
  id: string;
  name: string;
  color?: string | undefined;
}

export interface SkillSearchProps {
  /** Currently selected skill */
  value: SkillSearchOption | undefined;
  /** Called when user picks a skill or clears selection */
  onChange: (skill: SkillSearchOption | undefined) => void;
  /** Async loader that returns matching skills for the given query */
  loadOptions: (query: string) => Promise<SkillSearchOption[]>;
  /** Placeholder while nothing is selected */
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

export function SkillSearch({
  value,
  onChange,
  loadOptions,
  placeholder = "Найти навык...",
  delay = 300,
  disabled = false,
  className,
  label,
}: SkillSearchProps): JSX.Element {
  const [query, setQuery] = useState("");
  const [isOpen, setIsOpen] = useState(false);
  const [options, setOptions] = useState<SkillSearchOption[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [highlightedIndex, setHighlightedIndex] = useState(-1);
  const debouncedQuery = useDebounce(query, delay);

  const containerReference = useRef<HTMLDivElement>(null);
  const inputReference = useRef<HTMLInputElement>(null);
  const listReference = useRef<HTMLUListElement>(null);

  // Load options when debounced query changes
  useEffect(() => {
    if (!isOpen) return;

    let cancelled = false;

    void loadOptions(debouncedQuery)
      .then((result) => {
        if (!cancelled) {
          setOptions(result);
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
  }, [debouncedQuery, isOpen, loadOptions]);

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

  const handleOpen = useCallback((): void => {
    if (disabled) return;
    setIsOpen(true);
    setQuery("");
    setTimeout(() => {
      inputReference.current?.focus();
    }, 0);
  }, [disabled]);

  const handleSelect = useCallback(
    (option: SkillSearchOption): void => {
      onChange(option);
      setIsOpen(false);
      setQuery("");
    },
    [onChange],
  );

  const handleClear = useCallback(
    (event: React.MouseEvent): void => {
      event.stopPropagation();
      onChange(undefined);
      setQuery("");
    },
    [onChange],
  );

  const handleKeyDown = (event: React.KeyboardEvent): void => {
    if (event.key === "Escape") {
      setIsOpen(false);
      return;
    }
    if (!isOpen) {
      if (event.key === "Enter" || event.key === "ArrowDown") {
        handleOpen();
      }
      return;
    }
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
            aria-selected={value?.id === option.id}
            className={clsx(
              styles.option,
              index === highlightedIndex && styles.highlighted,
              value?.id === option.id && styles.selected,
            )}
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
      {/* Trigger / selected value */}
      {!isOpen ? (
        <button
          type="button"
          className={styles.trigger}
          onClick={handleOpen}
          disabled={disabled}
          aria-haspopup="listbox"
        >
          {value !== undefined ? (
            <span className={styles.selectedValue}>
              {value.color !== undefined && (
                <span className={styles.colorDot} style={{ backgroundColor: `#${value.color}` }} />
              )}
              <span className={styles.selectedName}>{value.name}</span>
              <button
                type="button"
                className={styles.clearButton}
                onClick={handleClear}
                aria-label="Очистить"
              >
                <X size={14} />
              </button>
            </span>
          ) : (
            <span className={styles.placeholder}>{placeholder}</span>
          )}
          <ChevronDown size={16} className={styles.chevron} />
        </button>
      ) : (
        <div className={styles.inputWrapper}>
          <Search size={16} className={styles.searchIcon} />
          <input
            ref={inputReference}
            type="text"
            className={styles.input}
            value={query}
            onChange={(event) => {
              setQuery(event.target.value);
              setIsLoading(true);
            }}
            placeholder={placeholder}
            role="combobox"
            aria-expanded={isOpen}
            aria-autocomplete="list"
          />
        </div>
      )}

      {/* Dropdown list */}
      {isOpen && (
        <ul ref={listReference} className={styles.dropdown} role="listbox">
          {renderDropdownContent()}
        </ul>
      )}
    </div>
  );
}
