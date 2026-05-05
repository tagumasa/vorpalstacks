/**
 * TanStack Query hooks for admin_config RPC operations.
 * Provides ListConfig, UpdateConfig, ResetConfig, ShutdownServer,
 * ListServices, SetPortMode.
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
  ListServicesRequestSchema,
  SetPortModeRequestSchema,
  type ConfigEntry,
  type ServiceInfo,
} from "@/gen/admin_config_pb";
import { transport } from "@/lib/transport";

const client = createClient(AdminConfigService, transport);

const CONFIG_KEY = ["admin_config", "list"];
const SERVICES_KEY = ["admin_config", "services"];

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
      qc.invalidateQueries({ queryKey: SERVICES_KEY });
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

export function useServicesList() {
  return useQuery({
    queryKey: SERVICES_KEY,
    queryFn: async () => {
      const res = await client.listServices(
        create(ListServicesRequestSchema, {}),
      );
      return res.services;
    },
  });
}

export function useSetPortMode() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ serviceName, mode }: { serviceName: string; mode: string }) => {
      return client.setPortMode(
        create(SetPortModeRequestSchema, { serviceName, mode }),
      );
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: SERVICES_KEY });
      qc.invalidateQueries({ queryKey: CONFIG_KEY });
    },
  });
}

export type { ConfigEntry, ServiceInfo };
