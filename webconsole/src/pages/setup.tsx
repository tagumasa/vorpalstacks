import { useState, useEffect, type FormEvent } from "react";
import { useTranslation } from "react-i18next";
import { createClient } from "@connectrpc/connect";
import { transport } from "@/lib/transport";
import { setTokens } from "@/lib/auth";
import { AdminAuthService } from "@/gen/admin_auth_pb";
import { useNavigate } from "react-router";

const authClient = createClient(AdminAuthService, transport);

function downloadCredentialsCsv(creds: { accessKeyId: string; secretAccessKey: string }) {
  const csv = `Access key ID,Secret access key\n${creds.accessKeyId},${creds.secretAccessKey}\n`;
  const blob = new Blob([csv], { type: "text/csv" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "root-access-keys.csv";
  a.click();
  URL.revokeObjectURL(url);
}

export function SetupPage() {
  const navigate = useNavigate();
  const { t } = useTranslation();
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [credentials, setCredentials] = useState<{
    accessKeyId: string;
    secretAccessKey: string;
  } | null>(null);

  useEffect(() => {
    authClient.isRootInitialized({}).then((res) => {
      if (res.isInitialized) {
        navigate("/login", { replace: true });
      }
    }).catch(() => {});
  }, [navigate]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError("");
    if (password !== confirmPassword) {
      setError(t("setup.passwordMismatch"));
      return;
    }
    if (password.length < 8) {
      setError(t("setup.passwordTooShort"));
      return;
    }
    setLoading(true);
    try {
      const res = await authClient.initialSetup({ password });
      setTokens(res.accessToken, res.refreshToken, res.idToken);
      setCredentials({
        accessKeyId: res.accessKeyId,
        secretAccessKey: res.secretAccessKey,
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : t("setup.setupFailed"));
    } finally {
      setLoading(false);
    }
  }

  if (credentials) {
    return (
      <div className="login-page">
        <div className="login-box">
          <h1 className="login-title">VORPALSTACKS</h1>
          <p className="login-subtitle">{t("setup.keysCreated")}</p>
          <div className="setup-credential-warning">
            {t("setup.saveWarning")}
          </div>
          <div className="setup-credential-block">
            <label className="setup-cred-label">{t("setup.accessKeyId")}</label>
            <code className="setup-cred-value">{credentials.accessKeyId}</code>
          </div>
          <div className="setup-credential-block">
            <label className="setup-cred-label">{t("setup.secretAccessKey")}</label>
            <code className="setup-cred-value">{credentials.secretAccessKey}</code>
          </div>
          <div className="setup-actions">
            <button type="button" onClick={() => navigate("/")}>
              {t("setup.continueToDashboard")}
            </button>
            <button
              type="button"
              className="setup-download-btn"
              onClick={() => downloadCredentialsCsv(credentials)}
            >
              {t("setup.downloadCsv")}
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="login-page">
      <div className="login-box">
        <h1 className="login-title">VORPALSTACKS</h1>
        <p className="login-subtitle">{t("setup.subtitle")}</p>
        {error && <div className="login-error">{error}</div>}
        <form onSubmit={handleSubmit}>
          <label>
            {t("setup.password")}
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="new-password"
              required
            />
          </label>
          <label>
            {t("setup.confirmPassword")}
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              autoComplete="new-password"
              required
            />
          </label>
          <button type="submit" disabled={loading}>
            {loading ? t("common.creating") : t("setup.createRootUser")}
          </button>
        </form>
      </div>
    </div>
  );
}
