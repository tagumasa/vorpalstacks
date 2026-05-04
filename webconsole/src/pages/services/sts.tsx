/**
 * STS service page. Shows caller identity via GetCallerIdentity RPC.
 */
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import { STSService } from "@/gen/sts_pb";
import {
  ServicePageLayout,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";
import { useListKey } from "@/lib/use-service-list";

/** STS page displaying the current caller identity (account, user ID, ARN). */
export function STSPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(STSService);
  const { queryKey } = useListKey("sts-caller-identity");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.getCallerIdentity({}),
  });

  return (
    <ServicePageLayout
      icon="🎫"
      title={t("services.sts.title")}
      isLoading={isLoading}
      error={error}
    >
      <div style={{ padding: 16 }}>
        <div style={{ marginBottom: 16 }}>
          <div className="detail-field">
            <span className="detail-label">{t("services.sts.accountLabel")}</span>
            <span className="cell-mono">{data?.account ?? "\u2014"}</span>
          </div>
          <div className="detail-field">
            <span className="detail-label">{t("services.sts.userIdLabel")}</span>
            <span className="cell-mono">{data?.userid ?? "\u2014"}</span>
          </div>
          <div className="detail-field">
            <span className="detail-label">{t("services.sts.arnLabel")}</span>
            <span className="cell-mono">{data?.arn ?? "\u2014"}</span>
          </div>
        </div>
        <JsonViewer data={data ?? {}} />
      </div>
    </ServicePageLayout>
  );
}
