import { type ReactNode, type JSX } from 'react';
import { clsx } from 'clsx';
import { SearchInput } from '../SearchInput/SearchInput';
import { Pagination } from '../Pagination/Pagination';
import styles from './DataListLayout.module.scss';

export interface DataListLayoutProps {
    title: string;
    subtitle?: string | undefined;
    
    // Search
    searchValue?: string | undefined;
    onSearchChange?: ((value: string) => void) | undefined;
    searchPlaceholder?: string | undefined;

    // Additional controls (e.g. inline dropdown next to search)
    controlsNode?: ReactNode | undefined;

    // Sidebar filters
    filtersNode?: ReactNode | undefined;

    // Grid content
    children: ReactNode;

    // Empty state
    isEmpty?: boolean | undefined;
    emptyMessage?: string | undefined;

    // Pagination
    currentPage?: number | undefined;
    totalPages?: number | undefined;
    onPageChange?: ((page: number) => void) | undefined;

    className?: string | undefined;
}

export function DataListLayout({
    title,
    subtitle,
    searchValue,
    onSearchChange,
    searchPlaceholder = 'Поиск...',
    controlsNode,
    filtersNode,
    children,
    isEmpty = false,
    emptyMessage = 'Ничего не найдено',
    currentPage,
    totalPages,
    onPageChange,
    className,
}: DataListLayoutProps): JSX.Element {
    const hasSearch = onSearchChange !== undefined;
    const hasPagination = totalPages !== undefined && totalPages > 1 && onPageChange !== undefined && currentPage !== undefined;

    return (
        <div className={clsx(styles.wrapper, className)}>
            <div className={styles.header}>
                <h1 className={styles.title}>{title}</h1>
                {subtitle !== undefined && <p className={styles.subtitle}>{subtitle}</p>}
            </div>

            {(hasSearch || controlsNode !== undefined) && (
                <div className={styles.controlsBlock}>
                    {hasSearch && (
                        <SearchInput
                            value={searchValue ?? ''}
                            onSearch={onSearchChange}
                            placeholder={searchPlaceholder}
                            className={styles.search ?? ''}
                        />
                    )}
                    {controlsNode !== undefined && (
                        <div className={styles.extraControls}>
                            {controlsNode}
                        </div>
                    )}
                </div>
            )}

            <div className={styles.mainArea}>
                {filtersNode !== undefined && (
                    <aside className={styles.sidebar}>
                        {filtersNode}
                    </aside>
                )}

                <div className={styles.contentArea}>
                    {!isEmpty ? (
                        <div className={styles.grid}>
                            {children}
                        </div>
                    ) : (
                        <div className={styles.empty}>
                            <p>{emptyMessage}</p>
                        </div>
                    )}

                    {hasPagination && (
                        <div className={styles.pagination}>
                            <Pagination
                                currentPage={currentPage}
                                totalPages={totalPages}
                                onPageChange={onPageChange}
                            />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
