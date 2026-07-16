/**
 * Status bar component at the bottom of the app shell.
 * Shows session metrics: REQ, ERR, LAT, MEM, GOR, UPTIME, VER.
 * Metrics are sourced from the useTelemetry() hook which polls
 * GetServerMetrics and reads the frontend interceptor store.
 */
import { useTelemetry } from "@/hooks/use-telemetry";

export function StatusBar() {
  const m = useTelemetry();

  return (
    <footer className="statusbar">
      <div className="sb-item">
        REQ <span className="sb-val">{m.requests.toLocaleString()}</span>
      </div>
      <div className="sb-item">
        ERR <span className="sb-val green">{m.errors}</span>
      </div>
      <div className="sb-item">
        LAT <span className="sb-val cyan">{m.avgLatency}ms</span>
      </div>
      <div className="sb-item">
        MEM <span className="sb-val">{m.memory}</span>
      </div>
      <div className="sb-item">
        GOR <span className="sb-val">{m.goroutines}</span>
      </div>
      <div className="sb-item">
        UPTIME <span className="sb-val yellow">{m.uptime}</span>
      </div>
      <div className="sb-item">
        VER <span className="sb-val">{__APP_VERSION__}</span>
      </div>
    </footer>
  );
}
