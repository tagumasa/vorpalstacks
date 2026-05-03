import { Navigate } from "react-router";

/**
 * Login redirect. /login always redirects to the root user login page.
 * From there the user can switch to the IAM user login if needed.
 */
export function LoginPage() {
  return <Navigate to="/login/root" replace />;
}
