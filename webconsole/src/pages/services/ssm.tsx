import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { SSMService } from "@/gen/ssm_pb";
import {
  PutParameterRequestSchema,
  DeleteParameterRequestSchema,
  type ParameterMetadata,
  ParameterType,
  ParameterTier,
} from "@/gen/ssm_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";

const columns: ColumnDef<ParameterMetadata, any>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "type",
    header: "Type",
    cell: ({ getValue }) => {
      const t = getValue() as ParameterType;
      const labels: Record<number, string> = {
        [ParameterType.SECURE_STRING]: "SecureString",
        [ParameterType.STRING_LIST]: "StringList",
        [ParameterType.STRING]: "String",
      };
      return <span className="badge">{labels[t] ?? String(t)}</span>;
    },
  },
  {
    accessorKey: "version",
    header: "Ver",
    size: 50,
  },
  {
    accessorKey: "tier",
    header: "Tier",
    cell: ({ getValue }) => {
      const t = getValue() as ParameterTier;
      const labels: Record<number, string> = {
        [ParameterTier.STANDARD]: "Standard",
        [ParameterTier.ADVANCED]: "Advanced",
        [ParameterTier.INTELLIGENT_TIERING]: "Intelligent",
      };
      return <span>{labels[t] ?? String(t)}</span>;
    },
    size: 90,
  },
  {
    accessorKey: "lastmodifieddate",
    header: "Last Modified",
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

const PARAM_TYPES = [
  { value: ParameterType.STRING, label: "String" },
  { value: ParameterType.SECURE_STRING, label: "SecureString" },
  { value: ParameterType.STRING_LIST, label: "StringList" },
];

export function SSMPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<ParameterMetadata | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formValue, setFormValue] = useState("");
  const [formType, setFormType] = useState(ParameterType.STRING);
  const [formDesc, setFormDesc] = useState("");

  const client = createClient(SSMService, transport);
  const { queryKey } = useListKey("ssm");
  const invalidate = useListInvalidator();

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.describeParameters({}),
    refetchInterval: 30_000,
  });

  const items: ParameterMetadata[] = dropEmpty(data?.parameters ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.putParameter(
        create(PutParameterRequestSchema, {
          name: formName,
          value: formValue,
          type: formType,
          description: formDesc,
          tier: ParameterTier.STANDARD,
          overwrite: false,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormValue("");
      setFormDesc("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) =>
      client.deleteParameter(
        create(DeleteParameterRequestSchema, { name }),
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
          <span className="page-icon">📋</span>
          <h1>Systems Manager Parameter Store</h1>
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
          <h1>Systems Manager Parameter Store</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📋</span>
        <h1>Systems Manager Parameter Store</h1>
        <span className="resource-count">{t("common.resources", { count: items.length })}</span>
        <div className="page-actions">
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            Create parameter
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
        <h2>Create parameter</h2>
        {createMutation.error && (
          <div className="modal-error">{String(createMutation.error)}</div>
        )}
        <label>
          Name
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder="/my-app/config-key"
            className="modal-input"
          />
        </label>
        <label>
          Value
          <textarea
            value={formValue}
            onChange={(e) => setFormValue(e.target.value)}
            placeholder="parameter value"
            className="modal-textarea"
            rows={3}
          />
        </label>
        <label>
          Type
          <select
            value={formType}
            onChange={(e) => setFormType(Number(e.target.value))}
            className="modal-select"
          >
            {PARAM_TYPES.map((pt) => (
              <option key={pt.value} value={pt.value}>
                {pt.label}
              </option>
            ))}
          </select>
        </label>
        <label>
          Description
          <input
            value={formDesc}
            onChange={(e) => setFormDesc(e.target.value)}
            placeholder="optional description"
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCreate(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!formName || !formValue || createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            {createMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <ConfirmDialog
        open={showDelete && !!selectedItem}
        title={t("confirm.deleteParameter")}
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
