/**
 * STS service page. Shows caller identity via GetCallerIdentity RPC.
 */
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { createClient } from "@connectrpc/connect";
import { STSService } from "@/gen/sts_pb";
import { transport } from "@/lib/transport";
import { JsonViewer } from "@/components/shared/json-viewer";

export function STSPage() {
  const { t } = useTranslation();
  const client = createClient(STSService, transport);

  const { data, isLoading, error } = useQuery({
    queryKey: ["sts", "caller-identity"],
    queryFn: () => client.getCallerIdentity({}),
  });

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🎫</span>
          <h1>STS</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🎫</span>
          <h1>STS</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🎫</span>
        <h1>STS — Caller Identity</h1>
      </div>

      <div style={{ padding: 16 }}>
        <div style={{ marginBottom: 16 }}>
          <div className="detail-field">
            <span className="detail-label">Account</span>
            <span className="cell-mono">{data?.account ?? "\u2014"}</span>
          </div>
          <div className="detail-field">
            <span className="detail-label">User ID</span>
            <span className="cell-mono">{data?.userid ?? "\u2014"}</span>
          </div>
          <div className="detail-field">
            <span className="detail-label">ARN</span>
            <span className="cell-mono">{data?.arn ?? "\u2014"}</span>
          </div>
        </div>
        <JsonViewer data={data ?? {}} />
      </div>
    </div>
  );
}
