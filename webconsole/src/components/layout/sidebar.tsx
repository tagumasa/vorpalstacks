/**
 * Sidebar navigation rendering the static service catalog by category.
 * All entries come from SERVICE_CATALOG — no backend query required.
 * Each service page handles its own error state if the admin handler
 * is unavailable.
 */
import { useMemo } from "react";
import { useNavigate, useLocation } from "react-router";
import { SERVICE_CATALOG, SERVICE_CATEGORIES } from "@/lib/service-catalog";

export function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

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
    <nav className="sidebar">
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
    </nav>
  );
}
