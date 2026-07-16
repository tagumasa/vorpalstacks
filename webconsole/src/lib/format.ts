/**
 * Formatting utilities for byte sizes and durations.
 */

export function formatBytes(bytes: number | bigint): string {
  const b = typeof bytes === "bigint" ? Number(bytes) : bytes;
  if (b < 1024) return `${b}B`;
  if (b < 1024 * 1024) return `${(b / 1024).toFixed(0)}KB`;
  if (b < 1024 * 1024 * 1024) return `${(b / (1024 * 1024)).toFixed(0)}MB`;
  return `${(b / (1024 * 1024 * 1024)).toFixed(1)}GB`;
}

export function formatUptime(seconds: number | bigint): string {
  const s = typeof seconds === "bigint" ? Number(seconds) : seconds;
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  if (h > 0) return `${h}h${m}m`;
  if (m > 0) return `${m}m`;
  return `${Math.floor(s)}s`;
}
