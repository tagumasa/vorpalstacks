/**
 * KMS service page. Lists keys via ListKeys RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { KMSService, type KeyListEntry } from "@/gen/kms_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<KeyListEntry, any>[] = [
  {
    accessorKey: "keyid",
    header: "Key ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "keyarn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
];

export function KMSPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<KeyListEntry | null>(null);

  const client = createClient(KMSService, transport);
  const { queryKey } = useListKey("kms");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listKeys({}),
    refetchInterval: 30_000,
  });

  const items: KeyListEntry[] = dropEmpty(data?.keys ?? [], "keyid");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔑</span>
          <h1>KMS</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔑</span>
          <h1>KMS</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔑</span>
        <h1>KMS</h1>
        <span className="resource-count">{items.length} keys</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.keyid}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.keyid}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.keyid}</h2>
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
