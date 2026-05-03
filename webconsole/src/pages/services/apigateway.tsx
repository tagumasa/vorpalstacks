/**
 * API Gateway service page. Lists REST APIs via GetRestApis RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { APIGatewayService, type RestApi } from "@/gen/apigateway_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<RestApi, any>[] = [
  {
    accessorKey: "name",
    header: "API Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "id",
    header: "API ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
    size: 100,
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
    accessorKey: "createddate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

export function APIGatewayPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<RestApi | null>(null);

  const client = createClient(APIGatewayService, transport);
  const { queryKey } = useListKey("apigateway");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.getRestApis({}),
    refetchInterval: 30_000,
  });

  const items: RestApi[] = dropEmpty(data?.items ?? [], "id");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🌉</span>
          <h1>API Gateway</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🌉</span>
          <h1>API Gateway</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🌉</span>
        <h1>API Gateway</h1>
        <span className="resource-count">{items.length} APIs</span>
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
