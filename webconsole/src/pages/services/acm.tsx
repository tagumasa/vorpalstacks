import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { ACMService, type CertificateSummary } from "@/gen/acm_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<CertificateSummary, any>[] = [
  {
    accessorKey: "domainname",
    header: "Domain",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "ISSUED" ? "badge-green" : v === "PENDING_VALIDATION" ? "badge-yellow" : "";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 160,
  },
  {
    accessorKey: "type",
    header: "Type",
    size: 100,
  },
  {
    accessorKey: "notafter",
    header: "Expires",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try {
        return new Date(v).toLocaleDateString();
      } catch {
        return v;
      }
    },
  },
];

function fmtDate(v: string | undefined): string {
  if (!v) return "\u2014";
  try { return new Date(v).toLocaleString(); } catch { return v; }
}

function ACMDetail({ item }: { item: CertificateSummary }) {
  const sans = item.subjectalternativenamesummaries ?? [];
  const keyUsages = item.keyusages ?? [];
  const extKeyUsages = item.extendedkeyusages ?? [];

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Domain</span><span className="cell-mono">{item.domainname}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.certificatearn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Status</span><span className="badge">{String(item.status)}</span></div>
        <div className="detail-field"><span className="detail-label">Type</span><span>{String(item.type)}</span></div>
        <div className="detail-field"><span className="detail-label">Key Algorithm</span><span>{String(item.keyalgorithm)}</span></div>
        <div className="detail-field"><span className="detail-label">In Use</span>{item.inuse ? <span className="badge badge-green">Yes</span> : <span className="badge">No</span>}</div>
        <div className="detail-field"><span className="detail-label">Managed By</span><span>{String(item.managedby)}</span></div>
        <div className="detail-field"><span className="detail-label">Renewal</span><span>{String(item.renewaleligibility)}</span></div>
      </section>

      <section className="detail-section">
        <h3>Validity</h3>
        <div className="detail-field"><span className="detail-label">Created</span><span>{fmtDate(item.createdat)}</span></div>
        <div className="detail-field"><span className="detail-label">Issued</span><span>{fmtDate(item.issuedat)}</span></div>
        <div className="detail-field"><span className="detail-label">Not Before</span><span>{fmtDate(item.notbefore)}</span></div>
        <div className="detail-field"><span className="detail-label">Not After</span><span>{fmtDate(item.notafter)}</span></div>
        {item.importedat && <div className="detail-field"><span className="detail-label">Imported</span><span>{fmtDate(item.importedat)}</span></div>}
        {item.revokedat && <div className="detail-field"><span className="detail-label">Revoked</span><span>{fmtDate(item.revokedat)}</span></div>}
      </section>

      {sans.length > 0 && (
        <section className="detail-section">
          <h3>Subject Alternative Names ({sans.length})</h3>
          {sans.map((s, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono">{s}</span></div>
          ))}
        </section>
      )}

      {(keyUsages.length > 0 || extKeyUsages.length > 0) && (
        <section className="detail-section">
          <h3>Key Usage</h3>
          {keyUsages.length > 0 && <div className="detail-field"><span className="detail-label">Usages</span><span>{keyUsages.join(", ")}</span></div>}
          {extKeyUsages.length > 0 && <div className="detail-field"><span className="detail-label">Extended</span><span>{extKeyUsages.join(", ")}</span></div>}
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

export function ACMPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<CertificateSummary | null>(null);

  const client = createClient(ACMService, transport);
  const { queryKey } = useListKey("acm");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listCertificates({});
      return dropEmpty(resp.certificatesummarylist ?? [], "certificatearn");
    },
    refetchInterval: 30_000,
  });

  const items: CertificateSummary[] = data ?? [];

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔐</span>
          <h1>ACM</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔐</span>
          <h1>ACM</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔐</span>
        <h1>ACM</h1>
        <span className="resource-count">{items.length} certificates</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.certificatearn}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.certificatearn}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.domainname}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <ACMDetail item={selectedItem} />
          </div>
        )}
      </div>
    </div>
  );
}
