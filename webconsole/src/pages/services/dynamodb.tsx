/**
 * DynamoDB service page. Lists table names via ListTables RPC.
 */
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { DynamoDBService } from "@/gen/dynamodb_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { useState } from "react";

interface TableRow {
  name: string;
}

const columns: ColumnDef<TableRow, any>[] = [
  {
    accessorKey: "name",
    header: "Table Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
];

export function DynamoDBPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);

  const client = createClient(DynamoDBService, transport);
  const { queryKey } = useListKey("dynamodb");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listTables({}),
    refetchInterval: 30_000,
  });

  const items: TableRow[] = dropEmpty((data?.tablenames ?? []).map((n) => ({ name: n })), "name");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🗃️</span>
          <h1>DynamoDB</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🗃️</span>
          <h1>DynamoDB</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🗃️</span>
        <h1>DynamoDB</h1>
        <span className="resource-count">{items.length} tables</span>
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
