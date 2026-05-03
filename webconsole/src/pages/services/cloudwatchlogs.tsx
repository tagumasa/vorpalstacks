import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { CloudWatchLogsService, type LogGroupSummary } from "@/gen/cloudwatchlogs_pb";
import { CreateLogGroupRequestSchema } from "@/gen/cloudwatchlogs_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";

const columns: ColumnDef<LogGroupSummary, any>[] = [
  {
    accessorKey: "loggroupname",
    header: "Log Group",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "loggrouparn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "creationtime",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as number | bigint;
      if (!v) return "\u2014";
      try { return new Date(Number(v)).toLocaleString(); } catch { return String(v); }
    },
  },
  {
    accessorKey: "storedbytes",
    header: "Stored Bytes",
    size: 110,
  },
];

export function CloudWatchLogsPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<LogGroupSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");

  const client = createClient(CloudWatchLogsService, transport);
  const { queryKey } = useListKey("cloudwatchlogs");
  const invalidate = useListInvalidator();

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listLogGroups({}),
    refetchInterval: 30_000,
  });

  const items: LogGroupSummary[] = dropEmpty(data?.loggroups ?? [], "loggroupname");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createLogGroup(
        create(CreateLogGroupRequestSchema, { loggroupname: formName }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (logGroupName: string) =>
      client.deleteLogGroup({ loggroupname: logGroupName }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📜</span>
          <h1>CloudWatch Logs</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📜</span>
          <h1>CloudWatch Logs</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📜</span>
        <h1>CloudWatch Logs</h1>
        <span className="resource-count">{items.length} log groups</span>
        <div className="page-actions">
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            Create log group
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.loggroupname}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.loggroupname}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.loggroupname}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <JsonViewer data={selectedItem} />
          </div>
        )}
      </div>

      <Modal open={showCreate} onClose={() => setShowCreate(false)}>
        <h2>Create log group</h2>
        {createMutation.error && (
          <div className="modal-error">{String(createMutation.error)}</div>
        )}
        <label>
          Log Group Name *
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder="/aws/lambda/my-function"
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCreate(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!formName || createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            {createMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <ConfirmDialog
        open={showDelete && !!selectedItem}
        title="Delete log group"
        message={
          <>
            {t("confirm.confirmDelete", { name: "" })}
            <strong>{selectedItem?.loggroupname}</strong>?
          </>
        }
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.loggroupname)}
        onClose={() => setShowDelete(false)}
      />
    </div>
  );
}
