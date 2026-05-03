import { createBrowserRouter, Navigate } from "react-router";
import { isAuthenticated } from "@/lib/auth";
import { LoginPage } from "@/pages/login";
import { LoginRootPage } from "@/pages/login-root";
import { LoginIAMPage } from "@/pages/login-iam";
import { SetupPage } from "@/pages/setup";
import { AppShell } from "@/components/layout/app-shell";
import { ServiceRoute } from "@/pages/services/_router";

function RequireAuth({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export function DashboardPage() {
  return (
    <div className="content-area" style={{ padding: 24 }}>
      <h1 style={{ fontFamily: "var(--pixel-font)", fontSize: 11, color: "var(--accent)" }}>
        Dashboard
      </h1>
      <p style={{ color: "var(--text-muted)", marginTop: 12 }}>
        Select a service from the sidebar to browse resources.
      </p>
    </div>
  );
}

export const router = createBrowserRouter([
  {
    path: "/setup",
    element: <SetupPage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/login/root",
    element: <LoginRootPage />,
  },
  {
    path: "/login/iam",
    element: <LoginIAMPage />,
  },
  {
    path: "/",
    element: (
      <RequireAuth>
        <AppShell />
      </RequireAuth>
    ),
    children: [
      { index: true, element: <DashboardPage /> },
      { path: "services/:serviceId", element: <ServiceRoute /> },
    ],
  },
]);
