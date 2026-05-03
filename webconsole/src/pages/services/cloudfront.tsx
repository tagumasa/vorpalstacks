import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { CloudFrontService, type DistributionSummary } from "@/gen/cloudfront_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<DistributionSummary, any>[] = [
  {
    accessorKey: "id",
    header: "Distribution ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "domainname",
    header: "Domain",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "enabled",
    header: "Enabled",
    cell: ({ getValue }) => {
      const v = getValue() as boolean;
      return v ? <span className="badge badge-green">Yes</span> : <span className="badge">No</span>;
    },
    size: 80,
  },
  {
    accessorKey: "comment",
    header: "Comment",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      return v || "\u2014";
    },
  },
];

function CloudFrontDetail({ item }: { item: DistributionSummary }) {
  const origins = item.origins?.items ?? [];
  const aliases = item.aliases?.items ?? [];
  const cacheBehaviors = item.cachebehaviors?.items ?? [];
  const defBehavior = item.defaultcachebehavior;

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">ID</span><span className="cell-mono">{item.id}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Domain</span><span className="cell-mono">{item.domainname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Status</span><span>{item.status || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Enabled</span>{item.enabled ? <span className="badge badge-green">Yes</span> : <span className="badge badge-red">No</span>}</div>
        <div className="detail-field"><span className="detail-label">IPv6</span>{item.isipv6enabled ? "Yes" : "No"}</div>
        <div className="detail-field"><span className="detail-label">Staging</span>{item.staging ? "Yes" : "No"}</div>
        <div className="detail-field"><span className="detail-label">Comment</span><span>{item.comment || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">ETag</span><span className="cell-mono">{item.etag || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Last Modified</span><span>{item.lastmodifiedtime ? new Date(item.lastmodifiedtime).toLocaleString() : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Web ACL ID</span><span className="cell-mono">{item.webaclid || "\u2014"}</span></div>
      </section>

      {aliases.length > 0 && (
        <section className="detail-section">
          <h3>Aliases ({aliases.length})</h3>
          {aliases.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono">{a}</span></div>
          ))}
        </section>
      )}

      {origins.length > 0 && (
        <section className="detail-section">
          <h3>Origins ({origins.length})</h3>
          {origins.map((o) => (
            <div key={o.id} className="detail-field">
              <span className="detail-label">{o.id}</span>
              <span className="cell-mono">{o.domainname || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      {defBehavior && (
        <section className="detail-section">
          <h3>Default Cache Behaviour</h3>
          <div className="detail-field"><span className="detail-label">Target Origin</span><span className="cell-mono">{defBehavior.targetoriginid || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">Viewer Protocol</span><span>{String(defBehavior.viewerprotocolpolicy) || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">Allowed Methods</span><span>{defBehavior.allowedmethods?.items?.join(", ") || "\u2014"}</span></div>
        </section>
      )}

      {cacheBehaviors.length > 0 && (
        <section className="detail-section">
          <h3>Cache Behaviours ({cacheBehaviors.length})</h3>
          {cacheBehaviors.map((cb) => (
            <div key={cb.pathpattern} className="detail-field">
              <span className="detail-label">{cb.pathpattern}</span>
              <span className="cell-mono">{cb.targetoriginid || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      {item.viewercertificate && (
        <section className="detail-section">
          <h3>Viewer Certificate</h3>
          <div className="detail-field"><span className="detail-label">Certificate Source</span><span>{(item.viewercertificate as any).certificatesource || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">Minimum Protocol</span><span>{(item.viewercertificate as any).minimumprotocolversion || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">SSL Support</span><span>{(item.viewercertificate as any).sslsupportmethod || "\u2014"}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

export function CloudFrontPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<DistributionSummary | null>(null);

  const client = createClient(CloudFrontService, transport);
  const { queryKey } = useListKey("cloudfront");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listDistributions({});
      return dropEmpty(resp.distributionlist?.items ?? [], "id");
    },
    refetchInterval: 30_000,
  });

  const items: DistributionSummary[] = data ?? [];

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">☁️</span>
          <h1>CloudFront</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">☁️</span>
          <h1>CloudFront</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">☁️</span>
        <h1>CloudFront</h1>
        <span className="resource-count">{items.length} distributions</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.id}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.id}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.id}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <CloudFrontDetail item={selectedItem} />
          </div>
        )}
      </div>
    </div>
  );
}
