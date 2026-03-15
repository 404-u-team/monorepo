import { type JSX } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import { clsx } from 'clsx';
import { Button } from '../Button/Button';
import styles from './Pagination.module.scss';

export interface PaginationProps {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
    className?: string;
}

export function Pagination({ currentPage, totalPages, onPageChange, className }: PaginationProps): JSX.Element | undefined {
    if (totalPages <= 1) return undefined;

    const getVisiblePages = (current: number, total: number): (number | string)[] => {
        if (total <= 7) {
            return Array.from({ length: total }, (_, index) => index + 1);
        }

        if (current <= 4) {
            return [1, 2, 3, 4, 5, '...', total];
        }

        if (current >= total - 3) {
            return [1, '...', total - 4, total - 3, total - 2, total - 1, total];
        }

        return [1, '...', current - 1, current, current + 1, '...', total];
    };

    const pages = getVisiblePages(currentPage, totalPages);

    return (
        <div className={clsx(styles.pagination, className)}>
            <Button
                variant="outline"
                disabled={currentPage <= 1}
                onClick={() => { onPageChange(currentPage - 1); }}
                className={styles.navButton}
                aria-label="Previous Page"
            >
                <ChevronLeft size={16} />
            </Button>

            <div className={styles.pages}>
                {pages.map((page, index) => (
                    page === '...' ? (
                        <span key={`ellipsis-${String(index)}`} className={styles.ellipsis}>...</span>
                    ) : (
                        <Button
                            key={page}
                            variant={page === currentPage ? 'primary' : 'clear'}
                            onClick={() => { onPageChange(page as number); }}
                            className={clsx(styles.pageButton, page === currentPage && styles.active)}
                        >
                            {page}
                        </Button>
                    )
                ))}
            </div>

            <Button
                variant="outline"
                disabled={currentPage >= totalPages}
                onClick={() => { onPageChange(currentPage + 1); }}
                className={styles.navButton}
                aria-label="Next Page"
            >
                <ChevronRight size={16} />
            </Button>
        </div>
    );
}
