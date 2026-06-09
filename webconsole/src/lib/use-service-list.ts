/**
 * Shared hooks for service list pages. Provides region-aware query keys,
 * pagination support, and data helpers.
 */
import { useState, useCallback } from "react";
import { useQuery, useQueryClient, useIsFetching } from "@tanstack/react-query";
import { useAppState } from "@/lib/app-state";

/** Default refetch interval for all service list queries (30 seconds). */
export const REFETCH_INTERVAL = 30_000;

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

/**
 * Configuration for usePaginatedList.
 * @typeparam T - The item type (e.g. LogGroupSummary, StreamSummary)
 * @typeparam R - The RPC response type
 */
export interface PaginatedListConfig<T, R> {
  /** Unique query key prefix for this list. */
  queryKeyBase: unknown[];
  /** The RPC call that accepts a next-token and returns a response. */
  fetchPage: (nextToken: string) => Promise<R>;
  /** Extract the items array from the response. */
  getItems: (response: R) => T[];
  /** Extract the next-token from the response (empty string if none). */
  getNextToken: (response: R) => string;
}

/**
 * Paginated list hook that automatically fetches all pages.
 * Tracks accumulated items across pages and provides a loadMore callback
 * when more pages are available.
 *
 * Usage:
 * ```ts
 * const { items, hasMore, loadMore, isLoading, isFetchingMore } = usePaginatedList({
 *   queryKeyBase: queryKey,
 *   fetchPage: (token) => client.listLogGroups({ nexttoken: token || undefined }),
 *   getItems: (r) => r.loggroups ?? [],
 *   getNextToken: (r) => r.nexttoken ?? "",
 * });
 * ```
 */
export function usePaginatedList<T, R>(config: PaginatedListConfig<T, R>) {
  const { queryKeyBase, fetchPage, getItems, getNextToken } = config;

  const [continuationToken, setContinuationToken] = useState("");
  const [accumulatedItems, setAccumulatedItems] = useState<T[]>([]);

  const queryKey = [...queryKeyBase, continuationToken];

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => fetchPage(continuationToken),
    refetchInterval: REFETCH_INTERVAL,
  });

  const currentPageItems = data ? getItems(data) : [];
  const nextToken = data ? getNextToken(data) : "";

  const items = continuationToken
    ? [...accumulatedItems, ...currentPageItems]
    : currentPageItems;

  const hasMore = nextToken !== "";

  const isFetchingMore = useIsFetching({ queryKey }) > 0;

  const loadMore = useCallback(() => {
    setAccumulatedItems(items);
    setContinuationToken(nextToken);
  }, [items, nextToken]);

  const queryClient = useQueryClient();
  const invalidate = useCallback(() => {
    setContinuationToken("");
    setAccumulatedItems([]);
    queryClient.invalidateQueries({ queryKey: queryKeyBase });
  }, [queryClient, queryKeyBase]);

  return {
    items,
    hasMore,
    loadMore,
    isLoading,
    isFetchingMore,
    error,
    invalidate,
    data: data as R,
  };
}
