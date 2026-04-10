const PAGE_SIZE_KEY = "devspace_page_size";

export const PAGE_SIZE_OPTIONS = [10, 20, 50] as const;
export type PageSize = (typeof PAGE_SIZE_OPTIONS)[number];
export const DEFAULT_PAGE_SIZE: PageSize = 20;

export function getPageSize(): PageSize {
  try {
    const stored = localStorage.getItem(PAGE_SIZE_KEY);
    if (stored !== null) {
      const value = Number(stored) as PageSize;
      if ((PAGE_SIZE_OPTIONS as readonly number[]).includes(value)) return value;
    }
  } catch {
    // localStorage unavailable
  }
  return DEFAULT_PAGE_SIZE;
}

export function setPageSize(size: PageSize): void {
  try {
    localStorage.setItem(PAGE_SIZE_KEY, String(size));
  } catch {
    // localStorage unavailable
  }
}
