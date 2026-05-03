/**
 * CloudTrail service page. Lists trails via ListTrails RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { CloudTrailService, type TrailInfo } from "@/gen/cloudtrail_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<TrailInfo, any>[] = [
  {
    accessorKey: "name",
    header: "Trail Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "trailarn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "homeregion",
    header: "Home Region",
    size: 120,
  },
];

export function CloudTrailPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<TrailInfo | null>(null);

  const client = createClient(CloudTrailService, transport);
  const { queryKey } = useListKey("cloudtrail");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listTrails({}),
    refetchInterval: 30_000,
  });

  const items: TrailInfo[] = dropEmpty(data?.trails ?? [], "name");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📋</span>
          <h1>CloudTrail</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📋</span>
          <h1>CloudTrail</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📋</span>
        <h1>CloudTrail</h1>
        <span className="resource-count">{items.length} trails</span>
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
