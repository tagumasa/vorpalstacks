import { useEffect, useRef, useCallback } from "react";
import { clearTokens } from "./auth";

const STORAGE_KEY = "vs_idle_timeout_min";
const DEFAULT_MINUTES = 30;
const MIN_MINUTES = 1;
const MAX_MINUTES = 480;

let cachedMinutes: number | null = null;

function readFromStorage(): number {
  const raw = localStorage.getItem(STORAGE_KEY);
  if (!raw) return DEFAULT_MINUTES;
  const n = parseInt(raw, 10);
  if (isNaN(n) || n < MIN_MINUTES) return DEFAULT_MINUTES;
  if (n > MAX_MINUTES) return MAX_MINUTES;
  return n;
}

export function getIdleTimeoutMin(): number {
  if (cachedMinutes === null) {
    cachedMinutes = readFromStorage();
  }
  return cachedMinutes;
}

export function setIdleTimeoutMin(minutes: number): void {
  const clamped = Math.max(MIN_MINUTES, Math.min(MAX_MINUTES, Math.round(minutes)));
  localStorage.setItem(STORAGE_KEY, String(clamped));
  cachedMinutes = clamped;
}

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

    const onStorage = (e: StorageEvent) => {
      if (e.key === STORAGE_KEY) {
        cachedMinutes = null;
        resetTimer();
      }
    };
    window.addEventListener("storage", onStorage);

    return () => {
      if (timerRef.current) clearTimeout(timerRef.current);
      for (const ev of events) {
        window.removeEventListener(ev, handler);
      }
      window.removeEventListener("storage", onStorage);
    };
  }, [resetTimer]);
}
