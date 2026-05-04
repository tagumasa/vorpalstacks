/**
 * Combined frontend + server telemetry hook for StatusBar.
 * Polls GetServerMetrics every 30 seconds for MEM, GOR, uptime.
 * Reads frontend MetricsStore for REQ, ERR, LAT.
 */
import { useState, useEffect } from "react";
import { create } from "@bufbuild/protobuf";
import { createClient } from "@connectrpc/connect";
import {
  AdminConfigService,
  GetServerMetricsRequestSchema,
  type GetServerMetricsResponse,
} from "@/gen/admin_config_pb";
import { transport } from "@/lib/transport";
import { getFrontendMetrics } from "@/lib/metrics-store";
import { formatBytes, formatUptime } from "@/lib/format";

const client = createClient(AdminConfigService, transport);

const POLL_INTERVAL = 30_000;

export interface TelemetryMetrics {
  requests: number;
  errors: number;
  avgLatency: number;
  memory: string;
  goroutines: number;
  uptime: string;
  serverMetrics: GetServerMetricsResponse | null;
}

export function useTelemetry(): TelemetryMetrics {
  const [serverMetrics, setServerMetrics] = useState<GetServerMetricsResponse | null>(null);
  const [tick, setTick] = useState(0);

  useEffect(() => {
    const poll = async () => {
      try {
        const res = await client.getServerMetrics(
          create(GetServerMetricsRequestSchema, {}),
        );
        setServerMetrics(res);
      } catch {
        // Server unreachable — keep stale data.
      }
    };

    poll();
    const timer = setInterval(poll, POLL_INTERVAL);
    return () => clearInterval(timer);
  }, []);

  useEffect(() => {
    const timer = setInterval(() => setTick((t) => t + 1), 5_000);
    return () => clearInterval(timer);
  }, []);

  void tick;

  const fe = getFrontendMetrics();
  const mem = serverMetrics ? formatBytes(serverMetrics.processMemorySysBytes) : "--";
  const gor = serverMetrics ? Number(serverMetrics.goroutineCount) : 0;
  const up = serverMetrics ? formatUptime(serverMetrics.uptimeSeconds) : "--";

  return {
    requests: fe.requests,
    errors: fe.errors,
    avgLatency: fe.avgLatency,
    memory: mem,
    goroutines: gor,
    uptime: up,
    serverMetrics,
  };
}
