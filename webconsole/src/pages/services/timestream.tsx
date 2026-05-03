/**
 * Timestream service page. Lists databases and tables.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { TimestreamWriteService, type Database } from "@/gen/timestreamwrite_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<Database, any>[] = [
  {
    accessorKey: "databasename",
    header: "Database",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "tablecount",
    header: "Tables",
    size: 80,
  },
  {
    accessorKey: "arn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
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

export function TimestreamPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<Database | null>(null);

  const client = createClient(TimestreamWriteService, transport);
  const { queryKey } = useListKey("timestream");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listDatabases({}),
    refetchInterval: 30_000,
  });

  const items: Database[] = dropEmpty(data?.databases ?? [], "databasename");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">⏳</span>
          <h1>Timestream</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">⏳</span>
          <h1>Timestream</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">⏳</span>
        <h1>Timestream</h1>
        <span className="resource-count">{items.length} databases</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.databasename}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.databasename}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.databasename}</h2>
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
