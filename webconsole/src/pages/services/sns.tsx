/**
 * SNS service page. Lists topics via ListTopics RPC.
 */
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { SNSService, type Topic } from "@/gen/sns_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<Topic, any>[] = [
  {
    accessorKey: "topicarn",
    header: "Topic ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
];

function topicNameFromArn(arn: string): string {
  const parts = arn.split(":");
  return parts[parts.length - 1] || arn;
}

export function SNSPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<Topic | null>(null);

  const client = createClient(SNSService, transport);
  const { queryKey } = useListKey("sns");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listTopics({});
      return dropEmpty(resp.topics ?? [], "topicarn");
    },
    refetchInterval: 30_000,
  });

  const items: Topic[] = data ?? [];

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📨</span>
          <h1>SNS</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📨</span>
          <h1>SNS</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📨</span>
        <h1>SNS</h1>
        <span className="resource-count">{items.length} topics</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.topicarn}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.topicarn}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{topicNameFromArn(selectedItem.topicarn)}</h2>
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
