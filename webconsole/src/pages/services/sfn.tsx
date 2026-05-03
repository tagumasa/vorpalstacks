/**
 * Step Functions service page. Lists state machines via ListStateMachines RPC.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { SFNService, type StateMachineListItem } from "@/gen/sfn_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<StateMachineListItem, any>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "statemachinearn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "type",
    header: "Type",
    size: 100,
  },
  {
    accessorKey: "creationdate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

export function SFNPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<StateMachineListItem | null>(null);

  const client = createClient(SFNService, transport);
  const { queryKey } = useListKey("sfn");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listStateMachines({}),
    refetchInterval: 30_000,
  });

  const items: StateMachineListItem[] = dropEmpty(data?.statemachines ?? [], "name");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔀</span>
          <h1>Step Functions</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔀</span>
          <h1>Step Functions</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔀</span>
        <h1>Step Functions</h1>
        <span className="resource-count">{items.length} state machines</span>
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
