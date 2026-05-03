import { useState } from "react";
import { useTranslation } from "react-i18next";
import { clearTokens } from "@/lib/auth";
import { useTheme } from "@/lib/theme";
import { useAppState, REGIONS } from "@/lib/app-state";
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
  const [online] = useState(true);
  const { cycleTheme, label: themeLabel } = useTheme();
  const { region, setRegion } = useAppState();

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
        VORPALSTACKS<span className="sub">Inspector {__APP_VERSION__}</span>
      </div>
      <div className="header-controls">
        <span className={`status-dot${online ? "" : " offline"}`} />
        <span className="header-btn" style={{ border: "none", padding: 0, cursor: "default", fontSize: 7 }}>
          {online ? t("header.online") : t("header.offline")}
        </span>
        <div className="header-separator" />
        <button className="header-btn" onClick={cycleTheme} title={t("header.theme")}>
          {"🎨"} {themeLabel}
        </button>
        <Dropdown
          trigger={
            <button className="header-btn" title={t("header.region")}>
              {"🌐"} {region}
            </button>
          }
        >
          {REGIONS.map((r) => (
            <div
              key={r}
              className={`dropdown-item${r === region ? " active" : ""}`}
              onClick={() => setRegion(r)}
            >
              {r}
            </div>
          ))}
        </Dropdown>
        <button className="header-btn" title={t("header.settings")}>
          {"⚙"}
        </button>
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
        <button className="header-btn" title={t("header.logout")} onClick={handleLogout}>
          {"🚪"}
        </button>
      </div>
    </header>
  );
}
