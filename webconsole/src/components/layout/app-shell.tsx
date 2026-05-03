import { Header } from "./header";
import { Sidebar } from "./sidebar";
import { StatusBar } from "./statusbar";
import { Outlet } from "react-router";

/**
 * Root application shell layout. Uses CSS Grid to arrange:
 * - Header: full-width top bar (status, controls)
 * - Sidebar: left column with service list
 * - Main: right column content area (Outlet for nested routes)
 * - StatusBar: full-width bottom bar (telemetry)
 */
export function AppShell() {
  return (
    <div className="app">
      <Header />
      <Sidebar />
      <main className="main">
        <Outlet />
      </main>
      <StatusBar />
    </div>
  );
}
