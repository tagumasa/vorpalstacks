import { useState, useEffect } from "react";

/** Telemetry metrics displayed in the status bar. */
interface Metrics {
  /** Total RPC requests made in this session. */
  requests: number;
  /** Total RPC errors encountered. */
  errors: number;
  /** Average response latency in milliseconds. */
  latency: string;
  /** Server memory usage (placeholder until telemetry RPC exists). */
  memory: string;
  /** Server uptime (placeholder until telemetry RPC exists). */
  uptime: string;
}

/**
 * Status bar component at the bottom of the app shell.
 * Shows session metrics: REQ count, ERR count, LAT, MEM, and uptime.
 * Metrics are placeholders until a telemetry RPC is implemented.
 */
export function StatusBar() {
  const [metrics] = useState<Metrics>({
    requests: 0,
    errors: 0,
    latency: "0ms",
    memory: "--",
    uptime: formatUptime(startTime),
  });

  /** Refresh uptime display every 60 seconds. */
  const [, setTick] = useState(0);
  useEffect(() => {
    const timer = setInterval(() => setTick((t) => t + 1), 60_000);
    return () => clearInterval(timer);
  }, []);

  return (
    <footer className="statusbar">
      <div className="sb-item">
        REQ <span className="sb-val">{metrics.requests.toLocaleString()}</span>
      </div>
      <div className="sb-item">
        ERR <span className="sb-val green">{metrics.errors}</span>
      </div>
      <div className="sb-item">
        LAT <span className="sb-val cyan">{metrics.latency}</span>
      </div>
      <div className="sb-item">
        MEM <span className="sb-val">{metrics.memory}</span>
      </div>
      <div className="sb-item">
        UPTIME <span className="sb-val yellow">{metrics.uptime}</span>
      </div>
      <div className="sb-item">
        VER <span className="sb-val">{__APP_VERSION__}</span>
      </div>
    </footer>
  );
}

/** Session start time for uptime calculation. */
const startTime = Date.now();

/** Formats milliseconds into a human-readable uptime string (e.g. "2h15m"). */
function formatUptime(ms: number): string {
  const elapsed = Date.now() - ms;
  const seconds = Math.floor(elapsed / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  if (hours > 0) return `${hours}h${minutes % 60}m`;
  if (minutes > 0) return `${minutes}m`;
  return `${seconds}s`;
}
