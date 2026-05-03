/**
 * SES service page. Lists email identities via ListEmailIdentities RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { SESv2Service, type IdentityInfo } from "@/gen/sesv2_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<IdentityInfo, any>[] = [
  {
    accessorKey: "identityname",
    header: "Identity",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "identitytype",
    header: "Type",
    cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>,
    size: 90,
  },
  {
    accessorKey: "verificationstatus",
    header: "Verification",
    cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>,
    size: 110,
  },
  {
    accessorKey: "sendingenabled",
    header: "Sending",
    cell: ({ getValue }) => {
      const v = getValue() as boolean;
      return v ? <span className="badge badge-green">Enabled</span> : <span className="badge">Disabled</span>;
    },
    size: 90,
  },
];

export function SESPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<IdentityInfo | null>(null);

  const client = createClient(SESv2Service, transport);
  const { queryKey } = useListKey("ses");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listEmailIdentities({}),
    refetchInterval: 30_000,
  });

  const items: IdentityInfo[] = dropEmpty(data?.emailidentities ?? [], "identityname");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">✉️</span>
          <h1>SES</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">✉️</span>
          <h1>SES</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">✉️</span>
        <h1>SES</h1>
        <span className="resource-count">{items.length} identities</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.identityname}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.identityname}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.identityname}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <JsonViewer data={selectedItem} />
          </div>
        )}
      </div>
    </div>
  );
}
