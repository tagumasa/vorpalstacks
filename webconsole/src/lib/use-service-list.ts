/**
 * Shared hook for service list pages. Provides region-aware query key so
 * switching region triggers a refetch, plus an invalidation helper.
 */
import { useQueryClient } from "@tanstack/react-query";
import { useAppState } from "@/lib/app-state";

/**
 * Returns [queryKey, region]. The queryKey embeds region so React Query
 * refetches when the user switches region.
 */
export function useListKey(prefix: string) {
  const { region } = useAppState();
  return { queryKey: [prefix, "list", region], region };
}

/** Invalidates the given query key, triggering a refetch. */
export function useListInvalidator() {
  const queryClient = useQueryClient();
  return (queryKey: unknown[]) =>
    queryClient.invalidateQueries({ queryKey });
}

/**
 * Filters out empty proto objects from a list. Backend may return extra
 * entries with all-default fields (empty strings, zero numbers) that
 * render as ghost rows. An item is considered empty if the given key
 * field is an empty string (or undefined/null).
 */
export function dropEmpty<T>(items: T[], keyField: keyof T): T[] {
  return items.filter((item) => {
    const val = item[keyField];
    return val !== "" && val !== undefined && val !== null;
  });
}
