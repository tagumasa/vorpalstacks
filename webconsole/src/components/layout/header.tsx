import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import { clearTokens } from "@/lib/auth";
import { useTheme } from "@/lib/theme";
import { useAppState, REGIONS, REGION_LABELS } from "@/lib/app-state";
import { useNavigate } from "react-router";
import { Dropdown } from "@/components/shared/dropdown";

const LANGUAGES = [
  { code: "en", label: "EN" },
  { code: "ja", label: "JA" },
  { code: "zh", label: "ZH" },
] as const;

export function Header() {
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();
  const [online, setOnline] = useState(navigator.onLine);
  useEffect(() => {
    const goOnline = () => setOnline(true);
    const goOffline = () => setOnline(false);
    window.addEventListener("online", goOnline);
    window.addEventListener("offline", goOffline);
    return () => {
      window.removeEventListener("online", goOnline);
      window.removeEventListener("offline", goOffline);
    };
  }, []);
  const { cycleTheme, label: themeLabel } = useTheme();
  const { region, setRegion, sidebarCollapsed, setSidebarCollapsed } = useAppState();

  function handleLogout() {
    clearTokens();
    navigate("/login");
  }

  function setLanguage(code: string) {
    i18n.changeLanguage(code);
  }

  const currentLang = LANGUAGES.find((l) => i18n.language.startsWith(l.code)) ?? LANGUAGES[0];

  return (
    <header className="header">
      <div className="header-logo">
        <button className="header-hamburger icon-hamburger" onClick={() => setSidebarCollapsed(!sidebarCollapsed)} />
        VORPALSTACKS<span className="sub">Inspector {__APP_VERSION__}</span>
      </div>
      <div className="header-controls">
        <span className={`status-dot${online ? "" : " offline"}`} />
        <span className="header-btn" style={{ border: "none", padding: 0, cursor: "default", fontSize: 7 }}>
          {online ? t("header.online") : t("header.offline")}
        </span>
        <div className="header-separator" />
        <button className="header-btn icon-theme" onClick={cycleTheme} title={t("header.theme")}>
          {themeLabel}
        </button>
        <Dropdown
          trigger={
            <button className="header-btn icon-region" title={t("header.region")}>
              {REGION_LABELS[region] ?? region}
            </button>
          }
        >
          {REGIONS.map((r) => (
            <div
              key={r}
              className={`dropdown-item${r === region ? " active" : ""}`}
              onClick={() => setRegion(r)}
            >
              {REGION_LABELS[r] ?? r}
            </div>
          ))}
        </Dropdown>
        <button className="header-btn icon-settings" title={t("header.settings")} onClick={() => navigate("/settings")} />
        <Dropdown
          trigger={
            <button className="header-btn" title={t("header.language")}>
              {currentLang.label}
            </button>
          }
        >
          {LANGUAGES.map((l) => (
            <div
              key={l.code}
              className={`dropdown-item${l.code === currentLang.code ? " active" : ""}`}
              onClick={() => setLanguage(l.code)}
            >
              {l.label}
            </div>
          ))}
        </Dropdown>
        <div className="header-separator" />
        <button className="header-btn icon-logout" title={t("header.logout")} onClick={handleLogout} />
      </div>
    </header>
  );
}
