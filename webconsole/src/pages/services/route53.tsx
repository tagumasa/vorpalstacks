/**
 * Route 53 service page. Lists hosted zones via ListHostedZones RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { Route53Service, type HostedZone } from "@/gen/route53_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<HostedZone, any>[] = [
  {
    accessorKey: "name",
    header: "Zone Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "id",
    header: "Zone ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "resourcerecordsetcount",
    header: "Records",
    size: 80,
  },
];

export function Route53Page() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<HostedZone | null>(null);

  const client = createClient(Route53Service, transport);
  const { queryKey } = useListKey("route53");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listHostedZones({});
      return dropEmpty(resp.hostedzones ?? [], "id");
    },
    refetchInterval: 30_000,
  });

  const items: HostedZone[] = data ?? [];

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🌐</span>
          <h1>Route 53</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🌐</span>
          <h1>Route 53</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🌐</span>
        <h1>Route 53</h1>
        <span className="resource-count">{items.length} hosted zones</span>
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
