/**
 * Sidebar navigation rendering the static service catalog by category.
 * All entries come from SERVICE_CATALOG — no backend query required.
 * Each service page handles its own error state if the admin handler
 * is unavailable.
 */
import { useMemo, useState } from "react";
import { useNavigate, useLocation } from "react-router";
import { useTranslation } from "react-i18next";
import { SERVICE_CATALOG, SERVICE_CATEGORIES } from "@/lib/service-catalog";
import { useShutdownServer } from "@/hooks/use-config";
import { Modal } from "@/components/shared/modal";

export function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();
  const shutdownMut = useShutdownServer();
  const [showDanger, setShowDanger] = useState(false);
  const [showShutdownModal, setShowShutdownModal] = useState(false);
  const [shutdownConfirm, setShutdownConfirm] = useState("");

  function getActiveService(): string {
    const match = location.pathname.match(/^\/services\/([^/]+)/);
    return match?.[1] ?? "";
  }

  const activeService = getActiveService();

  const catalogByCategory = useMemo(() => {
    const map = new Map<string, typeof SERVICE_CATALOG>();
    for (const cat of SERVICE_CATEGORIES) {
      const entries = SERVICE_CATALOG
        .filter((s) => s.category === cat)
        .sort((a, b) => a.displayName.localeCompare(b.displayName));
      map.set(cat, entries);
    }
    return map;
  }, []);

  return (
    <>
      <div className="sidebar-title">Services</div>
      {[...catalogByCategory.entries()].map(([category, entries]) => (
        <div key={category}>
          <div className="sidebar-category">{category}</div>
          {entries.map((s) => (
            <div
              key={s.id}
              className={`sidebar-item${activeService === s.id ? " active" : ""}`}
              onClick={() => navigate(`/services/${s.id}`)}
            >
              <span className="icon">{s.icon}</span>
              {s.displayName}
            </div>
          ))}
        </div>
      ))}
      <div className="sidebar-spacer" />
      <div className="sidebar-danger">
        <button
          className="sidebar-danger-toggle"
          onClick={() => setShowDanger(!showDanger)}
        >
          {t("settings.danger.toggle")}
        </button>
        {showDanger && (
          <div className="sidebar-danger-content">
            <button
              className="btn btn-danger btn-sm"
              onClick={() => setShowShutdownModal(true)}
            >
              {t("settings.danger.shutdown")}
            </button>
          </div>
        )}
      </div>

      <Modal
        open={showShutdownModal}
        onClose={() => { setShowShutdownModal(false); setShutdownConfirm(""); }}
      >
        <h2>{t("settings.danger.shutdown")}</h2>
        <p>{t("settings.danger.shutdownDesc")}</p>
        <label>
          {t("settings.danger.shutdownConfirm")}
          <input
            className="modal-input"
            value={shutdownConfirm}
            onChange={(e) => setShutdownConfirm(e.target.value)}
            autoFocus
          />
        </label>
        <div className="modal-actions">
          <button
            className="btn btn-secondary"
            onClick={() => { setShowShutdownModal(false); setShutdownConfirm(""); }}
          >
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-danger"
            disabled={shutdownConfirm !== "SHUTDOWN" || shutdownMut.isPending}
            onClick={() => shutdownMut.mutate()}
          >
            {shutdownMut.isPending ? "..." : t("settings.danger.shutdown")}
          </button>
        </div>
      </Modal>
    </>
  );
}
