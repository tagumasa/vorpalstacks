/**
 * TanStack Query hooks for admin_config RPC operations.
 * Provides ListConfig, UpdateConfig, ResetConfig, and ShutdownServer.
 */
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { create } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import {
  AdminConfigService,
  ListConfigRequestSchema,
  UpdateConfigRequestSchema,
  ResetConfigRequestSchema,
  ShutdownServerRequestSchema,
  type ConfigEntry,
} from "@/gen/admin_config_pb";
import { transport } from "@/lib/transport";

const client = createClient(AdminConfigService, transport);

const CONFIG_KEY = ["admin_config", "list"];

export function useConfigList(category?: string) {
  return useQuery({
    queryKey: [...CONFIG_KEY, category ?? "all"],
    queryFn: async () => {
      const res = await client.listConfig(
        create(ListConfigRequestSchema, { category: category ?? "" }),
      );
      return res.entries;
    },
  });
}

export function useConfigUpdate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ key, value }: { key: string; value: string }) => {
      return client.updateConfig(create(UpdateConfigRequestSchema, { key, value }));
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: CONFIG_KEY });
    },
  });
}

export function useConfigReset() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (key: string) => {
      return client.resetConfig(create(ResetConfigRequestSchema, { key }));
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: CONFIG_KEY });
    },
  });
}

export function useShutdownServer() {
  return useMutation({
    mutationFn: async () => {
      return client.shutdownServer(create(ShutdownServerRequestSchema, {}));
    },
  });
}

export type { ConfigEntry };
