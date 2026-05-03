import { useTranslation } from "react-i18next";
import { useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { SQSService } from "@/gen/sqs_pb";
import { CreateQueueRequestSchema } from "@/gen/sqs_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";

interface TableRow {
  url: string;
  name: string;
}

function queueNameFromUrl(url: string): string {
  const idx = url.lastIndexOf("/");
  return idx >= 0 ? url.slice(idx + 1) : url;
}

const columns: ColumnDef<TableRow, any>[] = [
  {
    accessorKey: "name",
    header: "Queue Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "url",
    header: "URL",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
];

export function SQSPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");

  const client = createClient(SQSService, transport);
  const { queryKey } = useListKey("sqs");
  const invalidate = useListInvalidator();

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listQueues({}),
    refetchInterval: 30_000,
  });

  const items: TableRow[] = dropEmpty((data?.queueurls ?? []).map((url) => ({
    url,
    name: queueNameFromUrl(url),
  })), "url");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createQueue(
        create(CreateQueueRequestSchema, { queuename: formName }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (queueUrl: string) =>
      client.deleteQueue({ queueurl: queueUrl }),
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
          <span className="page-icon">📋</span>
          <h1>SQS</h1>
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
          <h1>SQS</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📋</span>
        <h1>SQS</h1>
        <span className="resource-count">{items.length} queues</span>
        <div className="page-actions">
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            Create queue
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
            getRowId={(row) => row.url}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.url}
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

      <Modal open={showCreate} onClose={() => setShowCreate(false)}>
        <h2>Create queue</h2>
        {createMutation.error && (
          <div className="modal-error">{String(createMutation.error)}</div>
        )}
        <label>
          Queue Name
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder="my-queue"
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
        title="Delete queue"
        message={
          <>
            {t("confirm.confirmDelete", { name: "" })}
            <strong>{selectedItem?.name}</strong>?
          </>
        }
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.url)}
        onClose={() => setShowDelete(false)}
      />
    </div>
  );
}
