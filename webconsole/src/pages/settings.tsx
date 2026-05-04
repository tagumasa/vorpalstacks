/**
 * Settings page with 4 tabs (Server/Features/Endpoints/Ports),
 * read-only info card, and danger zone.
 */
import { useState, useMemo, useCallback } from "react";
import { useTranslation } from "react-i18next";
import { useConfigList, useConfigUpdate, useConfigReset, useShutdownServer } from "@/hooks/use-config";

type TabKey = "server" | "features" | "endpoints" | "ports";

const TABS: { key: TabKey; labelKey: string }[] = [
  { key: "server", labelKey: "settings.tabServer" },
  { key: "features", labelKey: "settings.tabFeatures" },
  { key: "endpoints", labelKey: "settings.tabEndpoints" },
  { key: "ports", labelKey: "settings.tabPorts" },
];

const PORT_KEY_LABELS: Record<string, string> = {
  "ports.s3_website": "S3 Website",
  "ports.apigateway": "API Gateway",
  "ports.cognito_hosted": "Cognito Hosted UI",
  "ports.cloudfront": "CloudFront",
  "ports.lambda_url": "Lambda URL",
  "ports.appsync_events": "AppSync Events",
  "ports.neptune": "Neptune",
  "ports.route53_dns": "Route53 DNS",
  "ports.route53_healthcheck": "Route53 Health Check",
};

export function SettingsPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("server");
  const { data: entries, isLoading } = useConfigList();
  const updateMut = useConfigUpdate();
  const resetMut = useConfigReset();
  const shutdownMut = useShutdownServer();

  const readOnlyEntries = useMemo(
    () => (entries ?? []).filter((e) => e.readOnly),
    [entries],
  );

  const tabEntries = useMemo(() => {
    if (!entries) return [];
    const cats = categoryForTab(tab);
    return entries.filter((e) => !e.readOnly && cats.includes(e.category));
  }, [entries, tab]);

  const [dirty, setDirty] = useState<Record<string, string>>({});
  const [shutdownConfirm, setShutdownConfirm] = useState("");

  const handleDirty = useCallback((key: string, value: string) => {
    setDirty((prev) => ({ ...prev, [key]: value }));
  }, []);

  const handleSave = useCallback(() => {
    const updates = Object.entries(dirty);
    if (updates.length === 0) return;
    Promise.all(updates.map(([key, value]) => updateMut.mutateAsync({ key, value })))
      .then(() => setDirty({}));
  }, [dirty, updateMut]);

  const handleResetTab = useCallback(() => {
    tabEntries.forEach((e) => resetMut.mutate(e.key));
    setDirty({});
  }, [tabEntries, resetMut]);

  const handleShutdown = useCallback(() => {
    if (shutdownConfirm === "SHUTDOWN") {
      shutdownMut.mutate();
    }
  }, [shutdownConfirm, shutdownMut]);

  const applyNote = tab === "server" || tab === "ports"
    ? t("settings.restartNote")
    : undefined;

  const currentVal = (key: string, saved: string) => dirty[key] ?? saved;

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">⚙</span>
        <h1>{t("settings.title")}</h1>
      </div>

      <div className="settings-info-card">
        {readOnlyEntries.map((e) => (
          <span key={e.key} className="settings-info-item">
            <span className="settings-info-label">{labelForReadOnly(e.key, t)}</span>
            <span className="settings-info-value">{e.value || "—"}</span>
          </span>
        ))}
      </div>

      <div className="tab-bar">
        {TABS.map((tb) => (
          <button
            key={tb.key}
            className={`tab-btn ${tab === tb.key ? "active" : ""}`}
            onClick={() => { setTab(tb.key); setDirty({}); }}
          >
            {t(tb.labelKey)}
          </button>
        ))}
      </div>

      <div className="settings-tab-content">
        {tab === "ports" ? (
          <PortsTab
            entries={tabEntries}
            dirty={dirty}
            onDirty={handleDirty}
            currentVal={currentVal}
          />
        ) : (
          <ConfigForm
            entries={tabEntries}
            dirty={dirty}
            onDirty={handleDirty}
            currentVal={currentVal}
          />
        )}

        {applyNote && <div className="settings-note">⚠ {applyNote}</div>}

        <div className="settings-actions">
          <button
            className="btn btn-primary"
            disabled={Object.keys(dirty).length === 0 || updateMut.isPending}
            onClick={handleSave}
          >
            {t("settings.save")}
          </button>
          <button
            className="btn btn-secondary"
            disabled={resetMut.isPending}
            onClick={handleResetTab}
          >
            {t("settings.reset")}
          </button>
        </div>
      </div>

      <div className="settings-danger-zone">
        <h3>{t("settings.danger.shutdown")}</h3>
        <p>{t("settings.danger.shutdownDesc")}</p>
        <div className="settings-shutdown-row">
          <input
            className="settings-input"
            placeholder={t("settings.danger.shutdownConfirm")}
            value={shutdownConfirm}
            onChange={(e) => setShutdownConfirm(e.target.value)}
          />
          <button
            className="btn btn-danger"
            disabled={shutdownConfirm !== "SHUTDOWN" || shutdownMut.isPending}
            onClick={handleShutdown}
          >
            {t("settings.danger.shutdown")}
          </button>
        </div>
      </div>
    </div>
  );
}

function ConfigForm({
  entries,
  dirty,
  onDirty,
  currentVal,
}: {
  entries: import("@/gen/admin_config_pb").ConfigEntry[];
  dirty: Record<string, string>;
  onDirty: (key: string, value: string) => void;
  currentVal: (key: string, saved: string) => string;
}) {
  return (
    <div className="settings-form">
      {entries.map((e) => (
        <div key={e.key} className="settings-entry">
          <label className="settings-label">
            {e.description}
            {dirty[e.key] !== undefined && <span className="settings-dirty-dot" />}
          </label>
          {e.type === "BOOL" ? (
            <label className="settings-toggle">
              <input
                type="checkbox"
                checked={currentVal(e.key, e.value) === "true"}
                onChange={(ev) => onDirty(e.key, ev.target.checked ? "true" : "false")}
              />
              <span className="settings-toggle-slider" />
            </label>
          ) : e.type === "PORT" || e.type === "INT" ? (
            <input
              type="number"
              className="settings-input"
              value={currentVal(e.key, e.value)}
              onChange={(ev) => onDirty(e.key, ev.target.value)}
            />
          ) : (
            <input
              type="text"
              className="settings-input"
              value={currentVal(e.key, e.value)}
              onChange={(ev) => onDirty(e.key, ev.target.value)}
            />
          )}
        </div>
      ))}
    </div>
  );
}

function PortsTab({
  entries,
  dirty,
  onDirty,
  currentVal,
}: {
  entries: import("@/gen/admin_config_pb").ConfigEntry[];
  dirty: Record<string, string>;
  onDirty: (key: string, value: string) => void;
  currentVal: (key: string, saved: string) => string;
}) {
  const { t } = useTranslation();

  const defaultPorts = entries.filter((e) => {
    const parts = e.key.split(".");
    return parts.length === 2;
  });

  const resourcePorts = entries.filter((e) => {
    const parts = e.key.split(".");
    return parts.length >= 3;
  });

  return (
    <div className="settings-form">
      <h4 className="settings-section-title">{t("settings.ports.defaultPorts")}</h4>
      {defaultPorts.map((e) => (
        <div key={e.key} className="settings-entry settings-entry-row">
          <span className="settings-port-label">
            {PORT_KEY_LABELS[e.key] ?? e.key}
          </span>
          <input
            type="number"
            className="settings-input settings-port-input"
            value={currentVal(e.key, e.value)}
            onChange={(ev) => onDirty(e.key, ev.target.value)}
          />
          {dirty[e.key] !== undefined && <span className="settings-dirty-dot" />}
        </div>
      ))}

      {resourcePorts.length > 0 && (
        <>
          <h4 className="settings-section-title">{t("settings.ports.resourceMappings")}</h4>
          {resourcePorts.map((e) => {
            const parts = e.key.split(".");
            const svc = parts.slice(1, -1).join(".");
            const resId = parts[parts.length - 1];
            return (
              <div key={e.key} className="settings-entry settings-entry-row">
                <span className="settings-port-label">
                  {PORT_KEY_LABELS[`ports.${svc}`] ?? svc}
                </span>
                <span className="settings-resource-id">{resId}</span>
                <span className="settings-port-val">{e.value}</span>
              </div>
            );
          })}
        </>
      )}
    </div>
  );
}

function categoryForTab(tab: TabKey): string[] {
  switch (tab) {
    case "server": return ["server"];
    case "features": return ["features"];
    case "endpoints": return ["endpoints", "http"];
    case "ports": return ["ports"];
  }
}

function labelForReadOnly(key: string, t: (k: string) => string): string {
  switch (key) {
    case "aws.account_id": return t("settings.info.account");
    case "aws.region": return t("settings.info.region");
    case "storage.data_path": return t("settings.info.dataPath");
    case "storage.metadata_path": return t("settings.info.metadataPath");
    default: return key;
  }
}
