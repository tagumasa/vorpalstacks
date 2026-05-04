/**
 * Singleton store for frontend-tracked RPC metrics (REQ, ERR, LAT).
 * Updated by the telemetry interceptor in transport.ts.
 */
let requests = 0;
let errors = 0;
let totalLatency = 0;

export function recordRequest() { requests++; }
export function recordError() { errors++; }
export function recordLatency(ms: number) { totalLatency += ms; }

export function getFrontendMetrics() {
  return {
    requests,
    errors,
    avgLatency: requests > 0 ? Math.round(totalLatency / requests) : 0,
  };
}
