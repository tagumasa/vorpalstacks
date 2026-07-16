/**
 * Settings page with 4 tabs:
 *   Tab 1 (Server)       — CORS + TLS + Dynamic Port Range
 *   Tab 2 (Listener Ports) — Fixed server ports + Bind Mode
 *   Tab 3 (Port Mappings)  — Per-service mode + per-resource port assignments
 *   Tab 4 (Services)       — Enable/disable (grouped)
 */
import { useState, useMemo, useCallback } from "react";
import { useTranslation } from "react-i18next";
import {
  useConfigList,
  useConfigUpdate,
  useConfigReset,
  useServicesList,
  type ServiceInfo,
} from "@/hooks/use-config";
import { getIdleTimeoutMin, setIdleTimeoutMin } from "@/lib/idle-timeout";

type TabKey = "server" | "ports" | "mappings" | "services";

const TABS: { key: TabKey; labelKey: string }[] = [
  { key: "server", labelKey: "settings.tabServer" },
  { key: "ports", labelKey: "settings.tabPorts" },
  { key: "mappings", labelKey: "settings.tabMappings" },
  { key: "services", labelKey: "settings.tabServices" },
];

const SERVER_TAB_KEYS = [
  "http.cors_allowed_origins",
  "http.cors_allowed_methods",
  "http.cors_allowed_headers",
  "http.cors_expose_headers",
  "server.tls_enabled",
  "server.tls_cert_path",
  "server.tls_key_path",
  "ports.dynamic_range_start",
  "ports.dynamic_range_end",
];

const PORTS_TAB_KEYS = [
  "server.port",
  "server.grpc_web_port",
  "server.tls_port",
  "ports.route53_dns",
  "endpoints.base_url",
];

const FQDN_SERVICE_NAMES = [
  "s3_website", "apigateway", "cognito_hosted", "cloudfront",
  "lambda_url", "appsync_events",
];

const SERVICE_DISPLAY_GROUPS: { keys: string[]; displayKey: string }[] = [
  { keys: ["acm"], displayKey: "settings.services.svc_acm" },
  { keys: ["apigateway"], displayKey: "settings.services.svc_apigateway" },
  { keys: ["appsync"], displayKey: "settings.services.svc_appsync" },
  { keys: ["athena"], displayKey: "settings.services.svc_athena" },
  { keys: ["cloudfront"], displayKey: "settings.services.svc_cloudfront" },
  { keys: ["cloudtrail"], displayKey: "settings.services.svc_cloudtrail" },
  { keys: ["cloudwatch"], displayKey: "settings.services.svc_cloudwatch" },
  { keys: ["logs"], displayKey: "settings.services.svc_logs" },
  { keys: ["cognito"], displayKey: "settings.services.svc_cognito" },
  { keys: ["cognito_identity"], displayKey: "settings.services.svc_cognito_identity" },
  { keys: ["dynamodb"], displayKey: "settings.services.svc_dynamodb" },
  { keys: ["events"], displayKey: "settings.services.svc_events" },
  { keys: ["iam"], displayKey: "settings.services.svc_iam" },
  { keys: ["kinesis"], displayKey: "settings.services.svc_kinesis" },
  { keys: ["kms"], displayKey: "settings.services.svc_kms" },
  { keys: ["lambda"], displayKey: "settings.services.svc_lambda" },
  { keys: ["neptune", "neptune_data", "neptune_graph"], displayKey: "settings.services.svc_neptune" },
  { keys: ["route53"], displayKey: "settings.services.svc_route53" },
  { keys: ["s3"], displayKey: "settings.services.svc_s3" },
  { keys: ["scheduler"], displayKey: "settings.services.svc_scheduler" },
  { keys: ["secretsmanager"], displayKey: "settings.services.svc_secretsmanager" },
  { keys: ["sesv2"], displayKey: "settings.services.svc_sesv2" },
  { keys: ["sns"], displayKey: "settings.services.svc_sns" },
  { keys: ["sqs"], displayKey: "settings.services.svc_sqs" },
  { keys: ["ssm"], displayKey: "settings.services.svc_ssm" },
  { keys: ["stepfunctions"], displayKey: "settings.services.svc_stepfunctions" },
  { keys: ["sts"], displayKey: "settings.services.svc_sts" },
  { keys: ["timestream_query", "timestream_write"], displayKey: "settings.services.svc_timestream" },
  { keys: ["wafv2"], displayKey: "settings.services.svc_wafv2" },
];

const BIND_MODE_OPTIONS = [
  { value: "localhost", labelKey: "settings.bindMode.localhost" },
  { value: "all", labelKey: "settings.bindMode.all" },
  { value: "interface", labelKey: "settings.bindMode.interface" },
];

export function SettingsPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("server");
  const { data: entries, isLoading } = useConfigList();
  const updateMut = useConfigUpdate();
  const resetMut = useConfigReset();
  const { data: services } = useServicesList();

  const readOnlyEntries = useMemo(
    () => (entries ?? []).filter((e) => e.readOnly),
    [entries],
  );

  const tabEntries = useMemo(() => {
    if (!entries) return [];
    switch (tab) {
      case "server":
        return entries.filter((e) => !e.readOnly && SERVER_TAB_KEYS.includes(e.key));
      case "ports":
        return entries.filter((e) =>
          !e.readOnly && (PORTS_TAB_KEYS.includes(e.key) || e.key === "server.bind_mode" || e.key === "server.bind_interface"),
        );
      case "mappings":
        return entries.filter((e) => {
          if (e.readOnly) return false;
          if (e.key.endsWith(".mode") && e.key.startsWith("ports.")) return false;
          if (FQDN_SERVICE_NAMES.some((s) => e.key === `ports.${s}`)) return true;
          if (e.key.startsWith("ports.") && e.key !== "ports.route53_dns" && e.key !== "ports.dynamic_range_start" && e.key !== "ports.dynamic_range_end") {
            return e.key.split(".").length >= 3;
          }
          return false;
        });
      case "services":
        return [];
    }
  }, [entries, tab]);

  const [dirty, setDirty] = useState<Record<string, string>>({});
  const [saveMsg, setSaveMsg] = useState("");

  const handleDirty = useCallback((key: string, value: string) => {
    setDirty((prev) => ({ ...prev, [key]: value }));
  }, []);

  const handleSave = useCallback(() => {
    const updates = Object.entries(dirty);
    if (updates.length === 0) return;
    Promise.allSettled(updates.map(([key, value]) => updateMut.mutateAsync({ key, value })))
      .then((results) => {
        const failed = results.filter((r) => r.status === "rejected");
        if (failed.length === 0) {
          setDirty({});
          setSaveMsg(t("settings.saved"));
        } else {
          const succeeded = results.filter((r) => r.status === "fulfilled").length;
          setDirty((prev) => {
            const next = { ...prev };
            updates.forEach(([key], i) => {
              const r = results[i];
              if (r && r.status === "fulfilled") delete next[key];
            });
            return next;
          });
          setSaveMsg(`${t("settings.saveFailed")} (${succeeded}/${updates.length})`);
        }
        setTimeout(() => setSaveMsg(""), 3000);
      });
  }, [dirty, updateMut, t]);

  const handleResetTab = useCallback(() => {
    let keys = tabEntries.map((e) => e.key);
    if (tab === "ports") {
      keys = [...keys, "server.bind_mode", "server.bind_interface"].filter(
        (k, i, arr) => arr.indexOf(k) === i,
      );
    }
    Promise.allSettled(keys.map((k) => resetMut.mutateAsync(k))).then(() => {
      setDirty({});
    });
  }, [tabEntries, tab, resetMut]);

  const applyNote = tab === "services"
    ? t("settings.immediateNote")
    : tab === "ports"
      ? t("settings.restartNote")
      : t("settings.saveNote");

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
        <div className="page-header-row">
          <span className="page-icon">⚙</span>
          <h1>{t("settings.title")}</h1>
          <span className="resource-count">
            {readOnlyEntries.map((e) => `${labelForReadOnly(e.key, t)}: ${e.value || "—"}`).join(" · ")}
          </span>
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
      </div>

      <div className="settings-tab-content">
        <div className="settings-note">{applyNote}</div>

        {tab === "server" && (
          <ServerTab
            allEntries={entries ?? []}
            dirty={dirty}
            onDirty={handleDirty}
            currentVal={currentVal}
          />
        )}

        {tab === "ports" && (
          <PortsTab
            allEntries={entries ?? []}
            dirty={dirty}
            onDirty={handleDirty}
            currentVal={currentVal}
          />
        )}

        {tab === "mappings" && (
          <MappingsTab
            services={services ?? []}
            entries={entries ?? []}
            dirty={dirty}
            onDirty={handleDirty}
            currentVal={currentVal}
          />
        )}

        {tab === "services" && (
          <ServicesTab services={services ?? []} />
        )}

        {(tab === "server" || tab === "ports" || tab === "mappings") && (
          <div className="settings-actions">
            <button
              className="btn btn-primary"
              disabled={Object.keys(dirty).length === 0 || updateMut.isPending}
              onClick={handleSave}
            >
              {updateMut.isPending ? "..." : t("settings.save")}
            </button>
            <button
              className="btn btn-secondary"
              disabled={resetMut.isPending}
              onClick={handleResetTab}
            >
              {t("settings.reset")}
            </button>
            {saveMsg && <span style={{ fontSize: 11, color: "var(--accent-green)" }}>{saveMsg}</span>}
          </div>
        )}
      </div>
    </div>
  );
}

function ServerTab({
  allEntries,
  dirty,
  onDirty,
  currentVal,
}: {
  allEntries: import("@/gen/admin_config_pb").ConfigEntry[];
  dirty: Record<string, string>;
  onDirty: (key: string, value: string) => void;
  currentVal: (key: string, saved: string) => string;
}) {
  const { t } = useTranslation();

  const find = (key: string) => allEntries.find((e) => e.key === key);
  const tlsEnabled = currentVal("server.tls_enabled", find("server.tls_enabled")?.value ?? "false") === "true";
  const [idleMin, setIdleMin] = useState(getIdleTimeoutMin);

  const handleIdleChange = (val: string) => {
    const n = parseInt(val, 10);
    if (!isNaN(n) && n >= 1 && n <= 480) {
      setIdleMin(n);
      setIdleTimeoutMin(n);
    }
  };

  const renderEntry = (key: string) => {
    const e = find(key);
    if (!e) return null;
    return (
      <div key={e.key} className="settings-entry">
        <label className="settings-label">
          {t(`settings.entries.${e.key}` as const, e.description)}
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
    );
  };

  return (
    <div className="settings-form">
      <h4 className="settings-section-title">{t("settings.sections.session")}</h4>
      <div className="settings-entry">
        <label className="settings-label">{t("settings.idleTimeout")}</label>
        <input
          type="number"
          className="settings-input"
          min={1}
          max={480}
          value={idleMin}
          onChange={(ev) => handleIdleChange(ev.target.value)}
        />
      </div>

      <h4 className="settings-section-title">{t("settings.sections.cors")}</h4>
      {renderEntry("http.cors_allowed_origins")}
      {renderEntry("http.cors_allowed_methods")}
      {renderEntry("http.cors_allowed_headers")}
      {renderEntry("http.cors_expose_headers")}

      <h4 className="settings-section-title">{t("settings.sections.tls")}</h4>
      {renderEntry("server.tls_enabled")}
      {tlsEnabled && renderEntry("server.tls_cert_path")}
      {tlsEnabled && renderEntry("server.tls_key_path")}

      <h4 className="settings-section-title">{t("settings.sections.dynamicRange")}</h4>
      {renderEntry("ports.dynamic_range_start")}
      {renderEntry("ports.dynamic_range_end")}
    </div>
  );
}

function PortsTab({
  allEntries,
  dirty,
  onDirty,
  currentVal,
}: {
  allEntries: import("@/gen/admin_config_pb").ConfigEntry[];
  dirty: Record<string, string>;
  onDirty: (key: string, value: string) => void;
  currentVal: (key: string, saved: string) => string;
}) {
  const { t } = useTranslation();
  const find = (key: string) => allEntries.find((e) => e.key === key);
  const bindMode = currentVal("server.bind_mode", find("server.bind_mode")?.value ?? "localhost");
  const tlsEnabled = currentVal("server.tls_enabled", find("server.tls_enabled")?.value ?? "false") === "true";

  const portEntry = (key: string) => {
    const e = find(key);
    if (!e) return null;
    if (key === "server.tls_port" && !tlsEnabled) return null;
    return (
      <div key={key} className="settings-entry">
        <label className="settings-label">
          {t(`settings.entries.${key}` as const, e.description)}
          {dirty[key] !== undefined && <span className="settings-dirty-dot" />}
        </label>
        <input
          type="number"
          className="settings-input"
          value={currentVal(key, e.value)}
          onChange={(ev) => onDirty(key, ev.target.value)}
        />
      </div>
    );
  };

  return (
    <div className="settings-form">
      <h4 className="settings-section-title">{t("settings.sections.listenerPorts")}</h4>
      {portEntry("server.port")}
      {portEntry("server.grpc_web_port")}
      {portEntry("server.tls_port")}

      <h4 className="settings-section-title" style={{ marginTop: "1rem" }}>
        {t("settings.sections.endpoints")}
      </h4>
      {portEntry("endpoints.base_url")}

      <h4 className="settings-section-title" style={{ marginTop: "1rem" }}>
        {t("settings.sections.bindMode")}
      </h4>
      <div className="settings-entry">
        <label className="settings-label">
          {t("settings.entries.server.bind_mode")}
          {dirty["server.bind_mode"] !== undefined && <span className="settings-dirty-dot" />}
        </label>
        <select
          className="settings-select"
          value={bindMode}
          onChange={(ev) => onDirty("server.bind_mode", ev.target.value)}
        >
          {BIND_MODE_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>{t(o.labelKey)}</option>
          ))}
        </select>
      </div>

      {bindMode === "interface" && (
        <div className="settings-entry">
          <label className="settings-label">
            {t("settings.entries.server.bind_interface")}
            {dirty["server.bind_interface"] !== undefined && <span className="settings-dirty-dot" />}
          </label>
          <input
            type="text"
            className="settings-input"
            placeholder="192.168.1.10"
            value={currentVal("server.bind_interface", find("server.bind_interface")?.value ?? "")}
            onChange={(ev) => onDirty("server.bind_interface", ev.target.value)}
          />
        </div>
      )}
    </div>
  );
}

function MappingsTab({
  services,
  entries,
  dirty,
  onDirty,
  currentVal,
}: {
  services: ServiceInfo[];
  entries: import("@/gen/admin_config_pb").ConfigEntry[];
  dirty: Record<string, string>;
  onDirty: (key: string, value: string) => void;
  currentVal: (key: string, saved: string) => string;
}) {
  const { t } = useTranslation();

  const FQDN_SVC_MAP: Record<string, string> = {
    s3_website: "settings.mappings.svc_s3_website",
    apigateway: "settings.mappings.svc_apigateway",
    cognito_hosted: "settings.mappings.svc_cognito_hosted",
    cloudfront: "settings.mappings.svc_cloudfront",
    lambda_url: "settings.mappings.svc_lambda_url",
    appsync_events: "settings.mappings.svc_appsync_events",
  };

  const portServices = useMemo(() => {
    return Object.entries(FQDN_SVC_MAP).map(([svcName, labelKey]) => {
      const svc = services.find((s) => s.name === svcName);
      const modeKey = `ports.${svcName}.mode`;
      const savedMode = dirty[modeKey] ?? svc?.portMode ?? "fqdn";
      return {
        name: svcName,
        labelKey,
        portMode: savedMode,
        modeKey,
        fixedPortEntry: entries.find((e) => e.key === `ports.${svcName}`),
        defaultPort: svc?.defaultPort ?? 0,
      };
    });
  }, [services, entries, dirty]);

  const resourcePorts = useMemo(
    () => entries.filter((e) => {
      if (e.readOnly || e.key.endsWith(".mode")) return false;
      const parts = e.key.split(".");
      return parts.length >= 3 && parts[0] === "ports";
    }),
    [entries],
  );

  return (
    <div className="settings-form">
      <h4 className="settings-section-title">{t("settings.mappings.modeTable")}</h4>
      <table className="settings-table">
        <thead>
          <tr>
            <th>{t("settings.mappings.service")}</th>
            <th>{t("settings.mappings.mode")}</th>
            <th>{t("settings.mappings.port")}</th>
          </tr>
        </thead>
        <tbody>
          {portServices.map((s) => (
            <tr key={s.name}>
              <td>{t(s.labelKey)}</td>
              <td>
                <label className="settings-toggle" style={{ cursor: "pointer" }}>
                  <input
                    type="checkbox"
                    checked={s.portMode === "fqdn"}
                    onChange={(e) => onDirty(s.modeKey, e.target.checked ? "fqdn" : "individual")}
                  />
                  <span className="settings-toggle-slider" />
                </label>
                {dirty[s.modeKey] !== undefined && (
                  <span className="settings-dirty-dot" />
                )}
              </td>
              <td>
                {s.portMode === "fqdn" ? (
                  <input
                    type="number"
                    className="settings-input settings-port-input"
                    value={currentVal(
                      s.fixedPortEntry?.key ?? `ports.${s.name}`,
                      s.fixedPortEntry?.value ?? (s.defaultPort > 0 ? String(s.defaultPort) : ""),
                    )}
                    onChange={(ev) => onDirty(s.fixedPortEntry?.key ?? `ports.${s.name}`, ev.target.value)}
                  />
                ) : (
                  <span className="settings-port-val">dynamic</span>
                )}
                {dirty[s.fixedPortEntry?.key ?? `ports.${s.name}`] !== undefined && (
                  <span className="settings-dirty-dot" />
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {resourcePorts.length > 0 && (
        <>
          <h4 className="settings-section-title" style={{ marginTop: "1rem" }}>
            {t("settings.mappings.resourcePorts")}
          </h4>
          <table className="settings-table">
            <thead>
              <tr>
                <th>{t("settings.mappings.service")}</th>
                <th>{t("settings.mappings.resourceId")}</th>
                <th>{t("settings.mappings.port")}</th>
              </tr>
            </thead>
            <tbody>
              {resourcePorts.map((e) => {
                const parts = e.key.split(".");
                const svcKey = parts.slice(0, 2).join(".");
                const resId = parts.slice(2).join(".");
                return (
                  <tr key={e.key}>
                    <td>{svcKey}</td>
                    <td>{resId}</td>
                    <td>
                      <input
                        type="number"
                        className="settings-input settings-port-input"
                        value={currentVal(e.key, e.value)}
                        onChange={(ev) => onDirty(e.key, ev.target.value)}
                      />
                      {dirty[e.key] !== undefined && <span className="settings-dirty-dot" />}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </>
      )}

      {resourcePorts.length === 0 && (
        <p className="settings-empty">{t("settings.mappings.noMappings")}</p>
      )}
    </div>
  );
}

function ServicesTab({ services }: { services: ServiceInfo[] }) {
  const { t } = useTranslation();
  const updateMut = useConfigUpdate();

  const serviceMap = useMemo(() => {
    const m = new Map<string, ServiceInfo>();
    for (const s of services) {
      m.set(s.name, s);
    }
    return m;
  }, [services]);

  const isGroupEnabled = useCallback(
    (keys: string[]) => keys.some((k) => serviceMap.get(k)?.enabled ?? false),
    [serviceMap],
  );

  const handleToggle = useCallback(
    async (keys: string[], checked: boolean) => {
      const results = await Promise.allSettled(
        keys.map((k) =>
          updateMut.mutateAsync({ key: `services.${k}.enabled`, value: checked ? "true" : "false" }),
        ),
      );
      const failed = results.filter((r) => r.status === "rejected").length;
      if (failed > 0) {
        console.error(`Service toggle: ${failed}/${keys.length} updates failed`);
      }
    },
    [updateMut],
  );

  return (
    <div className="settings-form">
      <h4 className="settings-section-title">{t("settings.services.title")}</h4>
      <table className="settings-table">
        <thead>
          <tr>
            <th style={{ width: "80px" }}>{t("settings.services.enabled")}</th>
            <th>{t("settings.services.service")}</th>
          </tr>
        </thead>
        <tbody>
          {SERVICE_DISPLAY_GROUPS.map((g) => {
            const enabled = isGroupEnabled(g.keys);
            return (
              <tr key={g.keys[0]} onClick={() => handleToggle(g.keys, !enabled)} style={{ cursor: "pointer" }}>
                <td>
                  <span className={`badge ${enabled ? "badge-green" : "badge-red"}`}>
                    {enabled ? "ON" : "OFF"}
                  </span>
                </td>
                <td>{t(g.displayKey)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
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
