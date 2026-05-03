import { useState, useEffect, type FormEvent } from "react";
import { useTranslation } from "react-i18next";
import { createClient } from "@connectrpc/connect";
import { transport } from "@/lib/transport";
import { setTokens } from "@/lib/auth";
import { AdminAuthService } from "@/gen/admin_auth_pb";
import { useNavigate } from "react-router";

const authClient = createClient(AdminAuthService, transport);

export function LoginRootPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    authClient.isRootInitialized({}).then((res) => {
      if (!res.isInitialized) {
        navigate("/setup", { replace: true });
      }
    }).catch(() => {});
  }, [navigate]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const res = await authClient.loginRoot({ password });
      setTokens(res.accessToken, res.refreshToken, res.idToken);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : t("login.loginFailed"));
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="login-page">
      <div className="login-box">
        <h1 className="login-title">VORPALSTACKS</h1>
        <p className="login-subtitle">{t("login.rootSubtitle")}</p>
        {error && <div className="login-error">{error}</div>}
        <form onSubmit={handleSubmit}>
          <label>
            {t("login.password")}
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              required
            />
          </label>
          <button type="submit" disabled={loading}>
            {loading ? t("login.signingIn") : t("login.signInAsRoot")}
          </button>
        </form>
        <button
          type="button"
          className="login-switch"
          onClick={() => navigate("/login/iam")}
        >
          {t("login.switchToIam")}
        </button>
      </div>
    </div>
  );
}
