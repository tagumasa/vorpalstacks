/**
 * Athena service page. Lists work groups via ListWorkGroups RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { AthenaService, type WorkGroupSummary } from "@/gen/athena_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<WorkGroupSummary, any>[] = [
  {
    accessorKey: "name",
    header: "Work Group",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "state",
    header: "State",
    size: 90,
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      return v || "\u2014";
    },
  },
  {
    accessorKey: "creationtime",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

export function AthenaPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<WorkGroupSummary | null>(null);

  const client = createClient(AthenaService, transport);
  const { queryKey } = useListKey("athena");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listWorkGroups({}),
    refetchInterval: 30_000,
  });

  const items: WorkGroupSummary[] = dropEmpty(data?.workgroups ?? [], "name");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔬</span>
          <h1>Athena</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔬</span>
          <h1>Athena</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔬</span>
        <h1>Athena</h1>
        <span className="resource-count">{items.length} work groups</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.name}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.name}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.name}</h2>
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
