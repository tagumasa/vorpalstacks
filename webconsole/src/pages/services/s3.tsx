import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { S3Service } from "@/gen/s3_pb";
import {
  CreateBucketRequestSchema,
  DeleteBucketRequestSchema,
  type Bucket,
} from "@/gen/s3_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";

const columns: ColumnDef<Bucket, any>[] = [
  {
    accessorKey: "name",
    header: "Bucket Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "creationDate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try {
        return new Date(v).toLocaleString();
      } catch {
        return v;
      }
    },
  },
];

export function S3Page() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<Bucket | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formBucket, setFormBucket] = useState("");

  const client = createClient(S3Service, transport);
  const { queryKey } = useListKey("s3");
  const invalidate = useListInvalidator();

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listBuckets({}),
    refetchInterval: 30_000,
  });

  const items: Bucket[] = dropEmpty(data?.buckets ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createBucket(
        create(CreateBucketRequestSchema, { bucket: formBucket }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormBucket("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (bucket: string) =>
      client.deleteBucket(
        create(DeleteBucketRequestSchema, { bucket }),
      ),
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
          <span className="page-icon">📦</span>
          <h1>S3 Buckets</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📦</span>
          <h1>S3 Buckets</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📦</span>
        <h1>S3 Buckets</h1>
        <span className="resource-count">{t("common.resources", { count: items.length })}</span>
        <div className="page-actions">
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            Create bucket
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

      <Modal open={showCreate} onClose={() => setShowCreate(false)}>
        <h2>Create bucket</h2>
        {createMutation.error && (
          <div className="modal-error">{String(createMutation.error)}</div>
        )}
        <label>
          Bucket name
          <input
            value={formBucket}
            onChange={(e) => setFormBucket(e.target.value)}
            placeholder="my-bucket-name"
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCreate(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!formBucket || createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            {createMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <ConfirmDialog
        open={showDelete && !!selectedItem}
        title={t("confirm.deleteBucket")}
        message={
          <>
            {t("confirm.confirmDelete", { name: "" })}
            <strong>{selectedItem?.name}</strong>?
          </>
        }
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </div>
  );
}
