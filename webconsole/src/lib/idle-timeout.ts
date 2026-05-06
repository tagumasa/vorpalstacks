import { useEffect, useRef, useCallback } from "react";
import { clearTokens } from "./auth";

const STORAGE_KEY = "vs_idle_timeout_min";
const DEFAULT_MINUTES = 30;
const MIN_MINUTES = 1;
const MAX_MINUTES = 480;

/** Read idle timeout in minutes from localStorage. */
export function getIdleTimeoutMin(): number {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return DEFAULT_MINUTES;
  const n = parseInt(raw, 10);
  if (isNaN(n) || n < MIN_MINUTES) return DEFAULT_MINUTES;
  if (n > MAX_MINUTES) return MAX_MINUTES;
  return n;
}

/** Persist idle timeout in minutes to localStorage. */
export function setIdleTimeoutMin(minutes: number): void {
  const clamped = Math.max(MIN_MINUTES, Math.min(MAX_MINUTES, Math.round(minutes)));
  localStorage.setItem(STORAGE_KEY, String(clamped));
}

/**
 * Hook that monitors user activity and auto-logouts after
 * the configured idle timeout. Installs listeners for
 * mousemove, keydown, scroll, touchstart, and pointerdown.
 * Resets the timer on any activity.
 */
export function useIdleTimeout() {
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const resetTimer = useCallback(() => {
    if (timerRef.current) clearTimeout(timerRef.current);
    const minutes = getIdleTimeoutMin();
    timerRef.current = setTimeout(() => {
      clearTokens();
      const base = document.querySelector("base")?.getAttribute("href") ?? "/webconsole/";
      window.location.href = base + "login";
    }, minutes * 60 * 1000);
  }, []);

  useEffect(() => {
    const events = ["mousemove", "keydown", "scroll", "touchstart", "pointerdown"] as const;
    const handler = () => resetTimer();
    for (const ev of events) {
      window.addEventListener(ev, handler, { passive: true });
    }
    resetTimer();

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
      for (const ev of events) {
        window.removeEventListener(ev, handler);
      }
    };
  }, [resetTimer]);
}
