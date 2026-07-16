import { Header } from "./header";
import { Sidebar } from "./sidebar";
import { StatusBar } from "./statusbar";
import { Outlet } from "react-router";
import { useAppState } from "@/lib/app-state";
import { useIdleTimeout } from "@/lib/idle-timeout";

/**
 * Root application shell layout. Uses CSS Grid to arrange:
 * - Header: full-width top bar (status, controls)
 * - Sidebar: left column with service list (collapsible on mobile)
 * - Main: right column content area (Outlet for nested routes)
 * - StatusBar: full-width bottom bar (telemetry)
 */
export function AppShell() {
  const { sidebarCollapsed, setSidebarCollapsed } = useAppState();
  useIdleTimeout();

  return (
    <div className="app">
      <Header />
      {!sidebarCollapsed && (
        <div
          className="sidebar-overlay"
          onClick={() => setSidebarCollapsed(true)}
        />
      )}
      <div className={sidebarCollapsed ? "sidebar collapsed" : "sidebar"}>
        <Sidebar />
      </div>
      <main className="main">
        <Outlet />
      </main>
      <StatusBar />
    </div>
  );
}
