import { QueryClient } from "@tanstack/react-query";

/**
 * Shared TanStack Query client instance.
 * Configured with conservative defaults suitable for an admin console
 * that talks to a local Connect-RPC server.
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      /** Data is considered fresh for 30 seconds before refetching. */
      staleTime: 30_000,
      /** Only retry once on failure to avoid hammering the server. */
      retry: 1,
      /** Do not refetch when the browser window regains focus. */
      refetchOnWindowFocus: false,
    },
  },
});
