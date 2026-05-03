/**
 * Theme management hook. Provides three retro colour themes:
 * - IceCavern (default): cool blue palette
 * - DungeonTorch: warm amber palette
 * - Zelda: forest green palette
 * Persists selection to localStorage and applies via data-theme attribute.
 */
import { useState, useEffect } from "react";

export type Theme = "icecavern" | "dungeontorch" | "zelda";

const THEME_LABELS: Record<Theme, string> = {
  icecavern: "IceCavern",
  dungeontorch: "DungeonTorch",
  zelda: "Zelda",
};

const STORAGE_KEY = "vs-theme";

/** Hook that returns the current theme, setter, cycle function, and display label. */
export function useTheme() {
  const [theme, setThemeState] = useState<Theme>(() => {
    const stored = localStorage.getItem(STORAGE_KEY) as Theme | null;
    return stored ?? "icecavern";
  });

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
    localStorage.setItem(STORAGE_KEY, theme);
  }, [theme]);

  function setTheme(t: Theme) {
    setThemeState(t);
  }

  /** Advance to the next theme in the cycle. */
  function cycleTheme() {
    const themes: Theme[] = ["icecavern", "dungeontorch", "zelda"];
    const idx = themes.indexOf(theme);
    setTheme(themes[(idx + 1) % themes.length]!);
  }

  return { theme, setTheme, cycleTheme, label: THEME_LABELS[theme] } as const;
}
